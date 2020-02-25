package gitcall

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Handle(usercodeFunc UsercodeFunc) {
	uri := os.Getenv("DUNDERGITCALL_URI")
	if uri == "" {
		log.Fatal("DUNDERGITCALL_URI env is required but not set")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go sigHandler(cancel)

	server, err := NewServer(ctx, uri, usercodeFunc)
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
