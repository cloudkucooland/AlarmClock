package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
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

	hk := ledserver.NewLedServer(accessory.Info{
		Name:         "Birdhouse",
		SerialNumber: "1997-07-16",
		Manufacturer: "I❤️Jen",
		Model:        "💒🏩🐥",
		Firmware:     "69nice",
	}, led)

	ctx, cancel := context.WithCancel(context.Background())
	s, err := hap.NewServer(hap.NewFsStore(fulldir), hk.A)
	if err != nil {
		panic(err)
	}
	s.Pin = pin

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()
		s.ListenAndServe(ctx)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()

		t := time.Tick(time.Minute)

		for {
			select {
			case <-t:
				max, err := thermal()
				if err != nil {
					return
				}
				if max > threshold {
					ledserver.ThermalHigh()
				} else {
					ledserver.ThermalNormal()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP, os.Interrupt)
	sig := <-sigch

	fmt.Printf("shutdown requested by signal: %s\n", sig)
	cancel()
	fmt.Printf("waiting for goroutines\n")
	wg.Wait()
}
