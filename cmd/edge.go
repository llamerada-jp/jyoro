package main

import (
	"log"
	"time"

	jyoro "github.com/llamerada-jp/jyoro/internal"
	"github.com/spf13/cobra"
)

var edgeCmd = &cobra.Command{
	Use: "edge",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		config, err := jyoro.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}
		usb := jyoro.NewUSB()
		ticker := time.NewTicker(time.Second)

		log.Println("start edge")

		for range ticker.C {
			now := time.Now()
			for _, entry := range config.Entries {
				current, err := usb.IsON(entry.Location, entry.Port)
				if err != nil {
					log.Printf("failed to get usb status: %v", err)
					continue
				}
				if entry.Match(&now, config.Location) != current {
					err := usb.Power(entry.Location, entry.Port, !current)
					if err != nil {
						log.Printf("failed to switch usb: %v", err)
						continue
					}
					powerStatus := "off"
					if !current {
						powerStatus = "on"
					}
					log.Printf("switched usb %s:%d %s", entry.Location, entry.Port, powerStatus)
				}
			}
		}

		return nil
	},
}
