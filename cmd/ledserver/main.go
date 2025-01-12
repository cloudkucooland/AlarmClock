package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudkucooland/AlarmClock/ledserver"
)

func main() {
	led := &ledserver.LED{}
	if err := led.Init(); err != nil {
		panic(err)
	}
	defer led.Shutdown()

	if err := rpc.Register(led); err != nil {
		panic(err)
	}
	// nosec G114
	rpc.HandleHTTP()

	listener, err := net.Listen("unix", ledserver.Pipefile)
	if err != nil {
		fmt.Println("listen error:", err)
		panic(err)
	}

	if err := os.Chmod(ledserver.Pipefile, 0666); err != nil {
		panic(err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			panic(err)
		}
	}()

	// #nosec G114 -- this is a socket, no need for timeouts
	go http.Serve(listener, nil)

	// ctx, cancel := context.WithCancel(context.Background())

	sigch := make(chan os.Signal, 3)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	sig := <-sigch

	fmt.Printf("shutdown requested by signal: %s", sig)
	// cancel()
}
