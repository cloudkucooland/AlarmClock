package main

import (
	"context"
	"fmt"
	"image/color"
	"log"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/cloudkucooland/AlarmClock/ledserver"
)

func main() {
	client, err := rpc.DialHTTP("unix", ledserver.Pipefile)
	if err != nil {
		fmt.Printf("led server not connected: %s", err.Error())
		return
	}
	var res ledserver.Result

	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:  "sunrise",
				Usage: "run the sunrise ramp (optional duration, e.g. 30s or 15m)",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rpccmd := &ledserver.Command{
						Command: ledserver.Sunrise,
					}
					if arg := cmd.Args().Get(0); arg != "" {
						d, err := time.ParseDuration(arg)
						if err != nil {
							return err
						}
						rpccmd.Duration = d
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},

			{
				Name:  "fireflies",
				Usage: "run the fireflies effect",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rpccmd := &ledserver.Command{
						Command: ledserver.Fireflies,
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "alert",
				Usage: "run the alert pulse effect",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rpccmd := &ledserver.Command{
						Command: ledserver.AlertPulse,
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:  "candle",
				Usage: "run the candle effect",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rpccmd := &ledserver.Command{
						Command: ledserver.Candle,
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "rainbow",
				Aliases: []string{"r"},
				Usage:   "run the rainbow pattern",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rpccmd := &ledserver.Command{
						Command: ledserver.Rainbow,
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "off",
				Aliases: []string{"o"},
				Usage:   "turn all off",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rpccmd := &ledserver.Command{
						Command: ledserver.Off,
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "startup",
				Aliases: []string{"s"},
				Usage:   "run the startup test",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rpccmd := &ledserver.Command{
						Command: ledserver.Startup,
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "white",
				Aliases: []string{"w"},
				Usage:   "white",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rs := cmd.Args().Get(0)
					var r uint8
					if rs != "" {
						rx, err := strconv.ParseUint(rs, 16, 8)
						if err != nil {
							log.Fatal(err)
						}
						r = uint8(rx)
					}

					rpccmd := &ledserver.Command{
						Command: ledserver.AllOn,
						Color:   color.RGBA{R: r, G: r, B: r},
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
			{
				Name:    "color",
				Aliases: []string{"c"},
				Usage:   "run the color",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rs := cmd.Args().Get(0)
					var r uint8
					if rs != "" {
						rx, err := strconv.ParseUint(rs, 16, 8)
						if err != nil {
							log.Fatal(err)
						}
						r = uint8(rx)
					}

					gs := cmd.Args().Get(1)
					var g uint8
					if gs != "" {
						gx, err := strconv.ParseUint(gs, 16, 8)
						if err != nil {
							log.Fatal(err)
						}
						g = uint8(gx)
					}

					bs := cmd.Args().Get(2)
					var b uint8
					if bs != "" {
						bx, err := strconv.ParseUint(bs, 16, 8)
						if err != nil {
							log.Fatal(err)
						}
						b = uint8(bx)
					}

					rpccmd := &ledserver.Command{
						Command: ledserver.AllOn,
						Color:   color.RGBA{R: r, G: g, B: b},
					}

					if err := client.Call("LED.Set", rpccmd, &res); err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
