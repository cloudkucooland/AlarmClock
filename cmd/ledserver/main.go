package main

import (
	"fmt"
	// "context"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudkucooland/AlarmClock/ledserver"
)

func main() {
	led := new(ledserver.LED)
	if err := led.Init(); err != nil {
		panic(err)
	}

	if err := rpc.Register(led); err != nil {
		panic(err)
	}
	// nosec G114
	rpc.HandleHTTP()

	l, err := net.Listen("unix", ledserver.Pipefile)
	if err != nil {
		fmt.Println("listen error:", err)
		panic(err)
	}

	// #nosec G114 -- this is a socket, no need for timeouts
	go http.Serve(l, nil)
	fmt.Println("ledserver running")

	// ctx, cancel := context.WithCancel(context.Background())

	sigch := make(chan os.Signal, 3)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	sig := <-sigch

	if err := l.Close(); err != nil {
		panic(err)
	}

	fmt.Printf("shutdown requested by signal: %s", sig)
	// cancel()
}
