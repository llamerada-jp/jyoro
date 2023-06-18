package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"sigs.k8s.io/yaml"

	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	host "periph.io/x/host/v3"
)

const (
	defaultVRef         = 3.3
	defaultThresholdDry = 2.0
	defaultA            = 1.0
	defaultB            = 0.0
)

type Config struct {
	SPIPath            string         `json:"spi_path"`
	VRef               *float32       `json:"v_ref,omitempty"`
	Probes             map[uint]Probe `json:"probes"`
	ThresholdDryProbes float32        `json:"threshold_dry_probes"`
	SprayDuration      string         `json:"spray_duration"`
}

type Probe struct {
	Name         *string  `json:"name,omitempty"`
	A            *float32 `json:"a,omitempty"`
	B            *float32 `json:"b,omitempty"`
	ThresholdDry *float32 `json:"threshold_dry"`
}

func readConfig(file string) (*Config, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ADC MCP320x accessor via SPI
type MCP struct {
	port  spi.PortCloser
	conn  spi.Conn
	v_ref float32
}

func newProbeDevice(path string, v_ref *float32) (*MCP, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	if _, err := driverreg.Init(); err != nil {
		return nil, err
	}

	port, err := spireg.Open(path)
	if err != nil {
		return nil, err
	}

	conn, err := port.Connect(10*physic.KiloHertz, spi.Mode3, 8)
	if err != nil {
		log.Fatal(err)
	}

	vr := float32(defaultVRef)
	if v_ref != nil {
		vr = *v_ref
	}

	return &MCP{
		port:  port,
		conn:  conn,
		v_ref: vr,
	}, nil
}

func (mcp *MCP) read(ch uint) (float32, error) {
	w := []byte{
		0x06 | ((0x7 & byte(ch)) >> 2),
		(byte(ch) & 0x03) << 6,
		0x00,
	}

	r := make([]byte, len(w))
	if err := mcp.conn.Tx(w, r); err != nil {
		return 0, err
	}

	v := (float32(r[1]&0x0f)*256 + float32(r[2])) * mcp.v_ref / 4096.0
	return v, nil
}

func (mcp *MCP) close() {
	mcp.port.Close()
}

func countDryProbes(device *MCP, probes map[uint]Probe) int {
	count := 0

	for ch, probe := range probes {
		v, err := device.read(ch)
		if err != nil {
			log.Panicln("read error on probe ", ch, err)
		}

		a := float32(defaultA)
		if probe.A != nil {
			a = *probe.A
		}

		b := float32(defaultB)
		if probe.B != nil {
			b = *probe.B
		}

		threshold := float32(defaultThresholdDry)
		if probe.ThresholdDry != nil {
			threshold = *probe.ThresholdDry
		}

		if a*v+b > threshold {
			count += 1
		}
	}

	return count
}

func turnOnSpray() error {
	log.Println("turn on spray")

	for _, port := range []uint{5, 4, 3, 2} {
		if err := runHubCtrl(port, true); err != nil {
			return nil
		}
	}

	return nil
}

func turnOffSpray() error {
	log.Println("turn off spray")

	for _, port := range []uint{2, 3, 4, 5} {
		if err := runHubCtrl(port, false); err != nil {
			return nil
		}
	}

	return nil
}

func runHubCtrl(port uint, on bool) error {
	sw := "0"
	if on {
		sw = "1"
	}
	out, err := exec.Command("hub-ctrl",
		"-b", "1",
		"-d", "2",
		"-P", fmt.Sprintf("%d", port),
		"-p", sw).Output()
	if err != nil {
		log.Println(string(out))
		return err
	}
	return nil
}

func main() {
	defer turnOffSpray()

	configFile := flag.String("config", "config.yaml", "specify yaml formatted config file")
	flag.Parse()

	config, err := readConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	sprayDuration, err := time.ParseDuration(config.SprayDuration)
	if err != nil {
		log.Fatal(err)
	}

	device, err := newProbeDevice(config.SPIPath, config.VRef)
	if err != nil {
		log.Fatal(err)
	}
	defer device.close()

	spraying := false
	var sprayUntil time.Time
	turnOffSpray()

	for {
		time.Sleep(1 * time.Second)
		if time.Now().After(sprayUntil) {
			spraying = false
			turnOffSpray()
		}

		if countDryProbes(device, config.Probes) > int(config.ThresholdDryProbes) &&
			!spraying {
			spraying = true
			sprayUntil = time.Now().Add(sprayDuration)
			turnOnSpray()
		}
	}
}
