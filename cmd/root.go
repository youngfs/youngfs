package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "youngfs",
	Short:        "youngfs",
	Long:         "youngfs, easy to use distributed file system",
	Version:      "1.1.0",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		_ = cmd.Usage()
		return err
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_ = fmt.Errorf("Error:%s\n", err.Error())
		os.Exit(1)
	}
}
