package gitcall

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	defer fmt.Println("server stopped")

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
			fmt.Println("signal caught. stopping server")
			cancel()

			return
		}

		fmt.Printf("unknown signal caught: %d\n", sig)
	}
}
