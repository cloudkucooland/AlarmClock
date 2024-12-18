# Birdhouse Alarm Clock

This is part of a Chrsitmas present for my wife.  The only use I think others will have for it is as a source for examples.

She had a Grace Digital Mondo that was a fantastic alarm clock for her. They quit supporting it and most of the features were shut down. She had a very detailed set of features she wants in an alarmclock. There was nothing in her list that I didn't think I couldn't code myself with a little effort.

Back in the last 1990's and early 2000's she made a lot of [folder icons](https://totoro.org/jen/bluecat/icons/) and [desktop wallpaper](https://totoro.org/jen/bluecat/desktop.shtml). I thought it would be cool to repurpose these for this project.

While trying to think of how to build the case, I had the idea of a birdhouse, using her bird icons for the GUI.

# Parts

[Raspberry Pi 5](https://www.raspberrypi.com/products/raspberry-pi-5/) Raspberry Pi 5, absurd overkill, I wanted power for future expansion. I have ideas...

[HiFiBerry Amp4](https://www.hifiberry.com/shop/boards/hifiberry-amp4/) 35W is plenty of power for this. I like the ability to have a single power supply for both the amp and the Pi.

[Power Supply](https://www.pishop.us/product/18v-power-supply-with-power-cable/) Recommended by PiShop.us for use with the Amp4

[5 inch touch display](https://www.waveshare.com/5inch-dsi-lcd.htm) My first DSI display. Seems good.

[Dayton Audio PC105-8 4-inch full-range driver](https://www.microcenter.com/product/633680/PC105-8_4%22_Full-Range_Poly_Cone_Driver) Selected because Microcenter had it in stock and the size was right. Sound quality is adequate but not great. Maybe it'll be better once the enclosure is finished.

I used a small microsd card, there is no need in my usecase for a big one. I bought a 128G Microcenter branded one with the Pi and speaker, but it was defective. I've not returned it yet because I'm lazy and it was only $10. I've never had a problem with Microcenter cards in the past, so I'm assuming this is a fluke. Why am I telling you this? I don't know. I don't know!

# Need to add

RTC clock battery, obvious reasons.

# Glitches & Hitches

Once, after disassembling to test fit, the touch on the display quit working. After a few dis/re-assembling and moving from DSI-1 to DSI-0 it started working again. I think it was just underpowered and needed a bit to warm up. Or, it needed the speaker connected, I don't know. I've not gone looking to recreate the issue.

# Setup

use raspi-config to switch to Wayland / lwc

[install go](https://go.dev/doc/install)

Install the code
```go get github.com/cloudkucooland/AlarmClock/cmd/ac```

~/config.../ whatever I did to make it start the clock on boot and nothing else

[set up for mono](https://askubuntu.com/questions/1439652/how-can-i-downmix-stereo-audio-output-to-mono-in-pipewire-on-22-10)

```pactl set-default-sink 37```

# Code notes

This is my first project using [ebiten](https://ebitengine.org) or using a game engine at all. It has been fun. I'm learning a lot and now have the urge to write and 8-bit-looking RPG...


