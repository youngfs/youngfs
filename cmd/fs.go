package main

import (
	"github.com/spf13/cobra"
	"github.com/youngfs/youngfs/pkg/fs/routers"
	"github.com/youngfs/youngfs/pkg/fs/server"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/vars"
)

var fsCmd = &cobra.Command{
	Use:   "fs",
	Short: "youngfs virtual file system",
	Long:  "youngfs virtual file system",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.InitLogger()
		defer log.Sync()
		server.InitServer()
		routers.InitRouter()
		routers.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fsCmd)
	fsCmd.Flags().StringVarP(&vars.Port, "Port", "p", "9876", "fs server port")
	fsCmd.Flags().StringVar(&vars.RedisSocket, "RedisSocket", "localhost:6379", "redis ip and port")
	fsCmd.Flags().StringVar(&vars.RedisPassword, "RedisPassword", "", "redis password")
	fsCmd.Flags().IntVar(&vars.RedisDatabase, "RedisDatabase", 0, "redis database")
	fsCmd.Flags().StringVar(&vars.SeaweedFSMaster, "SeaweedFSMaster", "", "seaweedFS master")
	fsCmd.Flags().BoolVar(&vars.Debug, "Debug", false, "debug mode")
	fsCmd.Flags().BoolVar(&vars.InfoLog, "InfoLog", true, "print info log")
	_ = fsCmd.MarkFlagRequired("SeaweedFSMaster")
}
