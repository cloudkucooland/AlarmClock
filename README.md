# Birdhouse Alarm Clock

This is part of a Chrsitmas present for my wife.

The only use I think others will have for it is as a source for examples.

# Parts

[Raspberry Pi 5](https://www.raspberrypi.com/products/raspberry-pi-5/) Raspberry Pi 5, absurd overkill, I wanted power for future expansion. I have ideas...

[HiFiBerry Amp4](https://www.hifiberry.com/shop/boards/hifiberry-amp4/) 35W is plenty of power for this. I like the ability to have a single power supply for both the amp and the Pi.

[Power Supply](https://www.pishop.us/product/18v-power-supply-with-power-cable/) Recommended by PiShop.us for use with the Amp4

[5 inch touch display](https://www.waveshare.com/5inch-dsi-lcd.htm) My first DSI display. Seems good.

[Dayton Audio PC105-8 4-inch full-range driver](https://www.microcenter.com/product/633680/PC105-8_4%22_Full-Range_Poly_Cone_Driver) Selected because Microcenter had it in stock and the size was right. Sound quality is adequate but not great. Maybe it'll be better once the enclosure is finished.

I used a small microsd card, there is no need in my usecase for a big one.

# Need to add

RTC clock battery, obvious reasons.

# Glitches & Hitches

Once, after disassembling to test fit, the touch on the display quit working. After a few dis/re-assembling and moving from DSI-1 to DSI-0 it started working again. I think it was just underpowered and needed a bit to warm up. Or, it needed the speaker connected, I don't know. I've not gone looking to recreate the issue.

# Setup

use raspi-config to switch to Wayland / lwc

install go
go get github.com/cloudkucooland/AlarmClock/cmd/ac

~/config.../ whatever I did to make it start the clock on boot and nothing else


