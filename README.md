# jyoro

## system
TBD

## setup software

### OS
- [Raspberry Pi OS](https://www.raspberrypi.com/software/)
- [Enable SPI](https://www.raspberrypi-spy.co.uk/2014/08/enabling-the-spi-interface-on-the-raspberry-pi/)

### hub-ctrl

- [hub-ctrl](https://www.gniibe.org/development/ac-power-control-by-USB-hub/index.html)
- [an issue of combination raspberry pi and hub-ctrl](https://forums.raspberrypi.com/viewtopic.php?t=242059)

```sh
sudo apt-get install libusb-dev
wget http://www.gniibe.org/oitoite/ac-power-control-by-USB-hub/hub-ctrl.c
gcc -O2 hub-ctrl.c -o hub-ctrl-armhf-static -lusb -static
sudo cp hub-ctrl-armhf-static /usr/local/bin/hub-ctrl
```

## configuration

### example

```yaml
spi_path: "/dev/spidev0.0"
v_ref: 3.3
threshold_dry_probes: 2
spray_duration: "1m"
probes:
  0:
    name: "Dalmatie"
    a: 0.33
    b: -0.1
    threshold_dry: 0.5
  1:
    name: "Pastilliere"
  2:
    name: "Boujassotte Grise"
```