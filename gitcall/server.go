package gitcall

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"sync"
	"time"
)

type UsercodeFunc func(ctx context.Context, data map[string]interface{}) error

type Usercode struct {
	handler UsercodeFunc
	context context.Context
}

// nolint: gocritic
func (h *Usercode) Run(request map[string]interface{}, response *map[string]interface{}) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("usercode:%v", p)
		}
	}()

	d, ok := request["data"]
	if !ok {
		return fmt.Errorf("'data' was not found")
	}

	data, ok := d.(map[string]interface{})
	if !ok {
		return fmt.Errorf("'data' is invalid type")
	}

	if err := h.handler(h.context, data); err != nil {
		return fmt.Errorf("usercode:%w", err)
	}

	(*response)["data"] = data

	return nil
}

type Server struct {
	server      *rpc.Server
	listener    net.Listener
	connections map[int]net.Conn
	context     context.Context
	stopped     bool
	lock        sync.Mutex
}

func NewServer(ctx context.Context, uri string, usercode UsercodeFunc) (*Server, error) {
	handler := &Usercode{
		handler: usercode,
		context: ctx,
	}

	server := rpc.NewServer()
	err := server.Register(handler)
	if err != nil {
		return nil, fmt.Errorf("register handler failed: %w", err)
	}

	listener, err := net.Listen("tcp", uri)
	if err != nil {
		return nil, fmt.Errorf("listen failed: %v", err)
	}

	s := &Server{
		server:      server,
		listener:    listener,
		connections: make(map[int]net.Conn),
		context:     ctx,
	}

	return s, nil
}

func (s *Server) Run() {
	wg := sync.WaitGroup{}
	index := 0

	for {
		conn, err := s.listener.Accept()

		s.lock.Lock()
		stopped := s.stopped
		s.lock.Unlock()

		if stopped {
			break
		}

		if err != nil {
			fmt.Printf("connection failed: %v\n", err)

			continue
		}

		index++

		s.lock.Lock()
		s.connections[index] = conn
		s.lock.Unlock()

		wg.Add(1)
		go func(index int) {
			defer wg.Done()

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

	if err := s.listener.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		fmt.Printf("close listener failed: %v\n", err)
	}

	for _, conn := range s.connections {
		if err := conn.SetReadDeadline(time.Now()); err != nil {
			fmt.Printf("set read deadline failed: %v\n", err)
		}
	}
}
