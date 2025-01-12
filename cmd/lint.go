package main

import (
	"fmt"

	jyoro "github.com/llamerada-jp/jyoro/internal"
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use: "lint",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		_, err := jyoro.LoadConfig(configPath)
		if err != nil {
			return err
		}
		fmt.Println("config is valid")
		return nil
	},
}
