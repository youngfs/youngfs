package volume

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/youngfs/youngfs/pkg/chunk/volume"
	"github.com/youngfs/youngfs/pkg/chunk/volume/needle"
	"github.com/youngfs/youngfs/pkg/kv"
	"github.com/youngfs/youngfs/pkg/kv/badger"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/log/zap"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"
)

const (
	port           = "PORT"
	logLevel       = "LOG_LEVEL"
	logAge         = "LOG_AGE"
	logFileSize    = "LOG_FILE_SIZE"
	dir            = "DIR"
	masterEndpoint = "MASTER"
	localIP        = "LOCAL_IP"
)

const (
	cmdPort           = "port"
	cmdLogLevel       = "logLevel"
	cmdLogAge         = "logAge"
	cmdLogFileSize    = "logFileSize"
	cmdDir            = "dir"
	cmdMasterEndpoint = "master"
	cmdLocalIP        = "localIP"
)

var cmdM = map[string]string{
	cmdPort:           port,
	cmdLogLevel:       logLevel,
	cmdLogAge:         logAge,
	cmdLogFileSize:    logFileSize,
	cmdDir:            dir,
	cmdMasterEndpoint: masterEndpoint,
	cmdLocalIP:        localIP,
}

var markRequired = []string{masterEndpoint}

var Cmd = &cobra.Command{
	Use:   "volume",
	Short: "youngfs chunk volume server",
	Long:  "youngfs chunk volume server",
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = viper.ReadInConfig()
		viper.SetEnvPrefix("YOUNGFS")
		viper.AutomaticEnv()
		for _, flag := range markRequired {
			if !viper.IsSet(flag) {
				return fmt.Errorf("flag %s is not set", flag)
			}
		}

		var closers []io.Closer
		var syncer []interface{ Sync() error }

		// log
		var logOptions []zap.Option
		logOptions = append(logOptions, zap.WithLogPath(path.Join(viper.GetString(dir), "log")))
		logOptions = append(logOptions, zap.WithLogFileAge(int(viper.GetUint64(logAge))))
		logOptions = append(logOptions, zap.WithLogFileSize(int(viper.GetUint64(logFileSize))))
		level, err := log.ParserLogLevel(viper.GetString(logLevel))
		if err != nil {
			return err
		}
		logOptions = append(logOptions, zap.WithLevel(level))
		logger := zap.New("volume", logOptions...)
		syncer = append(syncer, logger)

		creator := needle.KvNeedleStoreCreator(func(path string) (kv.Store, error) {
			return badger.New(path)
		})
		var options []volume.Option
		if viper.GetString(localIP) != "" {
			options = append(options, volume.WithLocalIP(viper.GetString(localIP)))
		}
		svr := volume.New(viper.GetString(dir), viper.GetString(masterEndpoint), logger, creator, options...)

		errChan := make(chan error, 1)
		go func(errChan chan<- error) {
			err := svr.Run(viper.GetInt(port))
			if err != nil {
				errChan <- err
			}
		}(errChan)
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		select {
		case err = <-errChan:
			break
		case <-signals:
			break
		}
		for _, s := range syncer {
			if err := s.Sync(); err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStderr(), "sync failed: %s\n", err.Error())
			}
		}
		for _, c := range closers {
			if err := c.Close(); err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStderr(), "close failed: %s\n", err.Error())
			}
		}
		return nil
	},
}

func init() {
	Cmd.Flags().Int(cmdPort, 9434, "port")
	Cmd.Flags().String(cmdLogLevel, "info", "log level (debug, info, warn, error, dpanic, panic)")
	Cmd.Flags().Uint64(cmdLogAge, 31, "log age (day)")
	Cmd.Flags().Uint64(cmdLogFileSize, 32, "log file max size (MiB)")

	Cmd.Flags().String(cmdDir, ".", "data dir")

	Cmd.Flags().String(cmdMasterEndpoint, "", "master endpoint")
	Cmd.Flags().String(cmdLocalIP, "", "local ip")

	Cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(cmdM[flag.Name], Cmd.PersistentFlags().Lookup(flag.Name))
	})
	Cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(cmdM[flag.Name], Cmd.Flags().Lookup(flag.Name))
	})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}
