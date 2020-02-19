package runner

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

type UsercodeFunc func(ctx context.Context, data map[string]interface{}) error

type Usercode struct {
	handler UsercodeFunc
	context context.Context
}

func (h *Usercode) Run(request map[string]interface{}, response *map[string]interface{}) error {
	d, ok := request["data"]
	if !ok {
		return fmt.Errorf("'data' was not found")
	}

	data, ok := d.(map[string]interface{})
	if !ok {
		return fmt.Errorf("'data' is invalid type")
	}

	if err := h.handler(h.context, data); err != nil {
		return err
	}

	(*response)["data"] = data

	return nil
}

type Server struct {
	server *rpc.Server
	listener net.Listener
	connections map[int]net.Conn
	context context.Context
	stopped bool
	lock sync.Mutex
}

func NewServer(ctx context.Context, socket string, usercode UsercodeFunc) (*Server, error) {
	handler := &Usercode{
		handler: usercode,
		context: ctx,
	}

	server := rpc.NewServer()
	err := server.Register(handler)
	if err != nil {
		return nil, fmt.Errorf("register handler failed: %w", err)
	}

	listener, err := net.Listen("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("listen failed: %v", err)
	}

	s := &Server{
		server:   server,
		listener: listener,
		connections: make(map[int]net.Conn, 0),
		context: ctx,
	}

	return s, nil
}

func (s *Server) Run() {
	wg := sync.WaitGroup{}
	index := 0

	BreakFor:for {
		conn, err := s.listener.Accept()

		s.lock.Lock()
		stopped := s.stopped
		s.lock.Unlock()

		if stopped {
			break BreakFor
		}

		if err != nil {
			log.Printf("connection failed: %v", err)

			continue
		}

		index++

		s.lock.Lock()
		s.connections[index] = conn
		s.lock.Unlock()

		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			defer log.Printf("connection handler closed: %d", index)

			s.server.ServeCodec(jsonrpc.NewServerCodec(conn))

			s.lock.Lock()
			delete(s.connections, index)
			s.lock.Unlock()
		}(index)
	}

	wg.Wait()
}

func (s *Server) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.stopped = true

	log.Print("close listener")
	if err := s.listener.Close(); err != nil {
		log.Printf("close listener failed: %v", err)
	}

	for i, conn := range s.connections {
		log.Printf("close connection: %d", i)
		if err := conn.Close(); err != nil {
			log.Printf("close connection failed: idx:%d err: %v", i, err)
		}
	}
}

