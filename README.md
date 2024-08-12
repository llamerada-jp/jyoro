# Jyoro🔌

## What's this?

RaspberryPiのUSBの給電を外部から制御するためのプログラムです。</br>
This program is for controlling the power supply of RaspberryPi USB from the outside.

## Use case

- USB給電のポンプを使った自動水やり</br>
  Automatic watering using USB-powered pumps.

## Setup device

### OS

- [Raspberry Pi OS](https://www.raspberrypi.com/software/)

### uhubctl

- [uhubctl](https://github.com/mvp/uhubctl)

```sh
# at Raspberry Pi
sudo apt-get install libusb-dev
git clone https://github.com/mvp/uhubctl
cd uhubctl
make
sudo make install
```

### Build

```sh
# at Raspberry Pi working directory
git clone https://github.com/llamerada-jp/jyoro.git
make -C jyoro
```

### Setup

Check your machines USB location & port like the followings
```sh
Current status for hub 1-1 👈️ location
  Port 1: 0503 power highspeed enable connect
  Port 2: 0000 off 👈️ target port
  Port 3: 0100 power
  Port 4: 0100 power
  Port 5: 0100 power
```

Prepare a configuration file like the following:

```json
{
	"location": "Asia/Tokyo", // 👈️  IANA Time Zone
  "entries": [
		{
      "location": "1-1", // 👈️ the location you looked up
      "port": 2, // 👈️ the port you looked up
		  "start_at": "08:00:00",
			"duration": "15m"
		}
	]
}

```

```sh
```

