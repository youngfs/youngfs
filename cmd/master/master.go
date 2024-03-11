package master

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/youngfs/youngfs/pkg/chunk/master"
	"github.com/youngfs/youngfs/pkg/kv/bbolt"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/log/zap"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"
)

const (
	port        = "PORT"
	logLevel    = "LOG_LEVEL"
	logAge      = "LOG_AGE"
	logFileSize = "LOG_FILE_SIZE"
	dir         = "DIR"
)

const (
	cmdPort        = "port"
	cmdLogLevel    = "logLevel"
	cmdLogAge      = "logAge"
	cmdLogFileSize = "logFileSize"
	cmdDir         = "dir"
)

var cmdM = map[string]string{
	cmdPort:        port,
	cmdLogLevel:    logLevel,
	cmdLogAge:      logAge,
	cmdLogFileSize: logFileSize,
	cmdDir:         dir,
}

var markRequired = []string{}

var Cmd = &cobra.Command{
	Use:   "master",
	Short: "youngfs chunk master server",
	Long:  "youngfs chunk master server",
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
		logger := zap.New("master", logOptions...)
		syncer = append(syncer, logger)

		// kv
		mkv, err := bbolt.New(path.Join(viper.GetString(dir), "master.db"), []byte("master"))
		if err != nil {
			return err
		}
		closers = append(closers, mkv)

		svr := master.New(mkv, logger)

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
	Cmd.Flags().Int(cmdPort, 9433, "port")
	Cmd.Flags().String(cmdLogLevel, "info", "log level (debug, info, warn, error, dpanic, panic)")
	Cmd.Flags().Uint64(cmdLogAge, 31, "log age (day)")
	Cmd.Flags().Uint64(cmdLogFileSize, 32, "log file max size (MiB)")

	Cmd.Flags().String(cmdDir, ".", "data dir")

	Cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(cmdM[flag.Name], Cmd.PersistentFlags().Lookup(flag.Name))
	})
	Cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(cmdM[flag.Name], Cmd.Flags().Lookup(flag.Name))
	})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}
