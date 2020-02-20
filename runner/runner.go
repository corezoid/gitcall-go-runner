package runner

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(usercodeFunc UsercodeFunc) {
	udsPath := os.Getenv("DUNDERGITCALL_UDS")
	if udsPath == "" {
		log.Fatal("DUNDERGITCALL_UDS env is required but not set")
	}

	_ = os.Remove(udsPath)

	ctx, cancel := context.WithCancel(context.Background())
	go sigHandler(cancel)

	server, err := NewServer(ctx, udsPath, usercodeFunc)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	go func() {
		<-ctx.Done()

		server.Stop()
	}()

	defer log.Print("server stopped")

	server.Run()
}

func sigHandler(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
	)

	for sig := range signals {
		if sig == syscall.SIGTERM || sig == syscall.SIGQUIT || sig == syscall.SIGINT || sig == os.Kill {
			log.Print("signal caught. stopping app")
			go func() {
				<-time.After(time.Second * 10)

				log.Print("graceful shutdown timed out")
				os.Exit(1)
			}()
			cancel()

			return
		}

		log.Printf("unknown signal caught: %d", sig)
	}
}
