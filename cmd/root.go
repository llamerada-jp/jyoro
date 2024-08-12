package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jyoro",
	Short: "jyoro is a tool for managing your USB power",
}

var (
	configPath string
)

func init() {
	pFlags := rootCmd.PersistentFlags()
	pFlags.StringVar(&configPath, "config", "", "path to the config file")
	rootCmd.MarkFlagRequired("config")

	rootCmd.AddCommand(edgeCmd)
	rootCmd.AddCommand(lintCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
