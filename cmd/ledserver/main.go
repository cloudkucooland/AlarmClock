package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/log"
	"github.com/cloudkucooland/AlarmClock/ledserver"
)

const fulldir = "/var/local/ledserver"
const pin = "56383990"

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

	hk := NewLedServer(accessory.Info{
		Name:         "Birdhouse",
		SerialNumber: "1997-07-16",
		Manufacturer: "I❤️Jen",
		Model:        "💒🏩🐥",
		Firmware:     "69nice",
	}, led)

	ctx, cancel := context.WithCancel(context.Background())
	s, err := hap.NewServer(hap.NewFsStore(fulldir), hk.A)
	if err != nil {
		log.Info.Panic(err)
	}
	s.Pin = pin

	// await context cancel
	go s.ListenAndServe(ctx)

	sigch := make(chan os.Signal, 3)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	sig := <-sigch

	fmt.Printf("shutdown requested by signal: %s", sig)
	cancel()

	// wait for cancel...
}
