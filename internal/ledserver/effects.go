package ledserver

import (
	"context"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

var colors = map[string]color.RGBA{
	"red":    {0xff, 0x00, 0x00, 0x00},
	"orange": {0xff, 0xa5, 0x00, 0x00},
	"yellow": {0xff, 0xff, 0x00, 0x00},
	"green":  {0x00, 0xff, 0x00, 0x00},
	"blue":   {0x00, 0x00, 0xff, 0x00},
	"indigo": {0x4b, 0x00, 0x82, 0x00},
	"violet": {0x7f, 0x00, 0xff, 0x00},
}

func (l *LED) rainbow(ctx context.Context, done chan struct{}) {
	l.off(true)
	inorder := []string{"red", "orange", "yellow", "green", "blue", "indigo", "violet", "indigo", "blue", "green", "yellow", "orange"}

	go func() {
		defer close(done)
		ticker := time.NewTicker(25 * time.Millisecond)
		defer ticker.Stop()

		for {
			for _, cname := range inorder {
				cc := colors[cname]
				for i := 0; i < l.bufsize; i += channels {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						l.mu.Lock()
						l.buf[i] = cc.R
						l.buf[i+1] = cc.G
						l.buf[i+2] = cc.B
						_, _ = l.leds.Write(l.buf)
						l.mu.Unlock()
					}
				}
			}
		}
	}()
}

func (l *LED) startup_test(ctx context.Context, done chan struct{}) {
	l.off(true)

	go func() {
		defer close(done)
		ticker := time.NewTicker(25 * time.Millisecond)
		defer ticker.Stop()

		for i := 0; i < l.bufsize; i += channels {
			for j := 0; j < channels; j++ {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					l.mu.Lock()
					l.buf[i+j] = 0x1f
					_, _ = l.leds.Write(l.buf)
					l.mu.Unlock()
				}

				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					l.mu.Lock()
					l.buf[i+j] = 0x00
					_, _ = l.leds.Write(l.buf)
					l.mu.Unlock()
				}
			}
		}
	}()
}

// sunrise smoothly ramps from deep red -> warm orange -> soft warm daylight over duration.
func (l *LED) sunrise(ctx context.Context, done chan struct{}, duration time.Duration) {
	go func() {
		defer close(done)

		const steps = 180
		stepInterval := duration / steps
		if stepInterval < 50*time.Millisecond {
			stepInterval = 50 * time.Millisecond
		}

		// Color stops in LAB space for smooth perceptual transitions
		cStart := colorful.Color{R: 0.08, G: 0.005, B: 0.0} // faint ember red
		cMid := colorful.Color{R: 1.0, G: 0.28, B: 0.02}    // warm sunrise orange
		cEnd := colorful.Color{R: 1.0, G: 0.82, B: 0.55}    // soft warm morning white

		ticker := time.NewTicker(stepInterval)
		defer ticker.Stop()

		for step := 0; step <= steps; step++ {
			t := float64(step) / float64(steps)
			var cf colorful.Color
			if t < 0.5 {
				cf = cStart.BlendLab(cMid, t*2.0)
			} else {
				cf = cMid.BlendLab(cEnd, (t-0.5)*2.0)
			}
			r, g, b := cf.RGB255()
			col := color.RGBA{R: r, G: g, B: b}

			// Update HomeKit every ~15 steps so we don't spam HAP notifications
			l.staticColor(col, step%15 == 0 || step == steps)

			if step == steps {
				return
			}

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
}

// candle reproduces the warm, jittered lantern flicker from the lantern ESP32 project.
func (l *LED) candle(ctx context.Context, done chan struct{}, brightness uint8) {
	if brightness == 0 {
		brightness = 4 // matches 0..63 * 4 -> 0..252 max red
	}
	go func() {
		defer close(done)

		for {
			l.mu.Lock()
			hot := runningHot.Load()
			var lastRnd uint32
			for i := 0; i < l.bufsize; i += channels {
				if hot && i%(4*channels) != 0 {
					l.buf[i], l.buf[i+1], l.buf[i+2] = 0, 0, 0
					continue
				}
				rnd := rand.Uint32()
				lastRnd = rnd

				rRaw := (rnd >> 8) & 0x3F
				rLow := rRaw & 0x0F
				gRaw := min(rnd&0x0F, rLow)
				// Keep blue very low on 3-channel RGB so low-level flicker stays amber/orange instead of lavender/white
				bRaw := min((rnd>>16)&0x07, rLow/2)

				// Subtle warm ember/wick glint (matches the 1-bit W flicker without washing out to blue-white)
				var sparkR, sparkG uint32
				if (rnd>>24)&0x03 == 0 {
					sparkR = 4
					sparkG = 1
				}

				l.buf[i] = uint8(min(255, (rRaw+sparkR)*uint32(brightness)))
				l.buf[i+1] = uint8(min(255, (gRaw+sparkG)*uint32(brightness)))
				l.buf[i+2] = uint8(min(255, bRaw*uint32(brightness)))
			}
			_, _ = l.leds.Write(l.buf)
			l.mu.Unlock()

			delay := time.Duration(50+(lastRnd>>20&0x7F)) * time.Millisecond
			select {
			case <-ctx.Done():
				return
			case <-time.After(delay):
			}
		}
	}()
}

type firefly struct {
	pixel int
	phase float64
	speed float64
	tint  color.RGBA
}

// fireflies renders a dim night-sky backdrop with a few drifting warm bioluminescent glows.
func (l *LED) fireflies(ctx context.Context, done chan struct{}) {
	palette := []color.RGBA{
		{R: 180, G: 220, B: 20}, // classic warm yellow-green firefly
		{R: 255, G: 190, B: 30}, // warm amber-gold
		{R: 140, G: 210, B: 45}, // soft chartreuse
	}

	const numFlies = 4
	flies := make([]firefly, numFlies)
	for f := range flies {
		flies[f] = firefly{
			pixel: rand.Intn(numpixels),
			phase: float64(f) / float64(numFlies), // stagger initial phases so they don't pulse in unison
			speed: 0.012 + rand.Float64()*0.015,   // ~1.5–3.5 seconds per full glow cycle at 30ms ticks
			tint:  palette[rand.Intn(len(palette))],
		}
	}

	go func() {
		defer close(done)
		ticker := time.NewTicker(30 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.mu.Lock()
				hot := runningHot.Load()

				// Faint sparse moonlight base (or pure black when running hot)
				for i := 0; i < l.bufsize; i += channels {
					if hot || (i/channels)%3 != 0 {
						l.buf[i], l.buf[i+1], l.buf[i+2] = 0, 0, 0
					} else {
						l.buf[i], l.buf[i+1], l.buf[i+2] = 0, 0, 1
					}
				}

				for idx := range flies {
					f := &flies[idx]
					f.phase += f.speed
					if f.phase >= 1.0 {
						// Respawn on a new pixel with a slightly varied tempo and hue
						f.phase = 0
						f.pixel = rand.Intn(numpixels)
						if hot {
							f.pixel = (f.pixel / 4) * 4
						}
						f.speed = 0.012 + rand.Float64()*0.015
						f.tint = palette[rand.Intn(len(palette))]
						continue
					}

					// Smooth bell curve: 0 -> 1 -> 0
					intensity := 0.5 * (1.0 - math.Cos(2.0*math.Pi*f.phase))
					// Sharpen peak slightly so the swell feels more bioluminescent
					intensity *= intensity

					base := f.pixel * channels
					if base+2 < l.bufsize {
						r := int(l.buf[base]) + int(float64(f.tint.R)*intensity)
						g := int(l.buf[base+1]) + int(float64(f.tint.G)*intensity)
						b := int(l.buf[base+2]) + int(float64(f.tint.B)*intensity)

						l.buf[base] = uint8(min(255, r))
						l.buf[base+1] = uint8(min(255, g))
						l.buf[base+2] = uint8(min(255, b))
					}
				}

				_, _ = l.leds.Write(l.buf)
				l.mu.Unlock()
			}
		}
	}()
}

// alertPulse smoothly breathes all pixels between a dim floor and full peak color.
// If col is black (zero), it defaults to warm coral-gold.
// If period <= 0, it defaults to 1.5 seconds per cycle (pass ~4s for a gentle snooze breathe).
func (l *LED) alertPulse(ctx context.Context, done chan struct{}, col color.RGBA, period time.Duration) {
	if col.R == 0 && col.G == 0 && col.B == 0 {
		col = color.RGBA{R: 255, G: 110, B: 20} // warm coral-gold default
	}
	if period <= 0 {
		period = 1500 * time.Millisecond
	}

	go func() {
		defer close(done)

		const tickInterval = 25 * time.Millisecond
		ticker := time.NewTicker(tickInterval)
		defer ticker.Stop()

		start := time.Now()
		const minFloor = 0.08 // keep a faint glow at the trough so it's clearly active

		for {
			select {
			case <-ctx.Done():
				return
			case now := <-ticker.C:
				elapsed := now.Sub(start)
				phase := float64(elapsed%period) / float64(period)

				// Smooth breathing curve: minFloor -> 1.0 -> minFloor
				wave := 0.5 * (1.0 - math.Cos(2.0*math.Pi*phase))
				scale := minFloor + (1.0-minFloor)*(wave*wave)

				r := uint8(float64(col.R) * scale)
				g := uint8(float64(col.G) * scale)
				b := uint8(float64(col.B) * scale)

				l.mu.Lock()
				hot := runningHot.Load()
				for i := 0; i < l.bufsize; i += channels {
					if hot && i%(4*channels) != 0 {
						l.buf[i], l.buf[i+1], l.buf[i+2] = 0, 0, 0
						continue
					}
					l.buf[i] = r
					l.buf[i+1] = g
					l.buf[i+2] = b
				}
				_, _ = l.leds.Write(l.buf)
				l.mu.Unlock()
			}
		}
	}()
}
