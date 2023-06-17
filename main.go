package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"sigs.k8s.io/yaml"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	host "periph.io/x/host/v3"
)

const Vref = 3.33

func main() {
	buf, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	if _, err := driverreg.Init(); err != nil {
		log.Fatal(err)
	}

	p, err := spireg.Open(config.Devices[0].Path)
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	c, err := p.Connect(10*physic.KiloHertz, spi.Mode3, 8)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
		for ch := 0; ch < 8; ch++ {
			write := []byte{
				0x06 | ((0x7 & byte(ch)) >> 2),
				(byte(ch) & 0x03) << 6,
				0x00,
			}

			read := make([]byte, len(write))
			if err := c.Tx(write, read); err != nil {
				log.Fatal(err)
			}

			readBin := (float64(read[1]&0x0f)*256 + float64(read[2])) / 4096.0
			fmt.Printf("%d %f\n", ch, Vref*readBin)
		}
		fmt.Println("----")
	}
}
