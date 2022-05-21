package command

import (
	"flag"
	"icesos/command/vars"
)

func InitCommand() {
	flag.StringVar(&vars.Port, "Port", "9876", "server port")
	flag.StringVar(&vars.RedisSocket, "RedisSocket", "localhost:6379", "redis ip and port")
	flag.StringVar(&vars.RedisPassword, "RedisPassword", "", "redis password")
	flag.IntVar(&vars.RedisDatabase, "RedisDatabase", 0, "redis database")
	flag.StringVar(&vars.SeaweedFSMaster, "SeaweedFSMaster", "", "seaweedFS master")
	flag.BoolVar(&vars.Debug, "Debug", false, "debug mode")
	flag.Parse()
}
