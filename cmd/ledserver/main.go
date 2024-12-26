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
	rpc.Register(led)
	rpc.HandleHTTP()

	l, err := net.Listen("unix", ledserver.Pipefile)
	if err != nil {
		fmt.Println("listen error:", err)
		panic(err)
	}

	go http.Serve(l, nil)

	// ctx, cancel := context.WithCancel(context.Background())

	sigch := make(chan os.Signal, 3)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	sig := <-sigch

	l.Close()

	fmt.Printf("shutdown requested by signal: %s", sig)
	// cancel()
}
