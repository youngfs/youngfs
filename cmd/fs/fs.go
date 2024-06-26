package fs

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/youngfs/youngfs/pkg/fs/engine"
	"github.com/youngfs/youngfs/pkg/fs/engine/chunk"
	"github.com/youngfs/youngfs/pkg/fs/engine/simple"
	"github.com/youngfs/youngfs/pkg/fs/handler"
	"github.com/youngfs/youngfs/pkg/fs/meta/s3"
	"github.com/youngfs/youngfs/pkg/fs/router"
	"github.com/youngfs/youngfs/pkg/fs/server"
	"github.com/youngfs/youngfs/pkg/kv"
	"github.com/youngfs/youngfs/pkg/kv/badger"
	"github.com/youngfs/youngfs/pkg/kv/bbolt"
	"github.com/youngfs/youngfs/pkg/kv/leveldb"
	"github.com/youngfs/youngfs/pkg/kv/tikv"
	"github.com/youngfs/youngfs/pkg/log"
	"github.com/youngfs/youngfs/pkg/log/zap"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

const (
	port         = "PORT"
	logLevel     = "LOG_LEVEL"
	logAge       = "LOG_AGE"
	logFileSize  = "LOG_FILE_SIZE"
	dir          = "DIR"
	master       = "MASTER"
	meta         = "META"
	tikvEndpoins = "TIKV_ENDPOINTS"
)

const (
	cmdPort         = "port"
	cmdLogLevel     = "logLevel"
	cmdLogAge       = "logAge"
	cmdLogFileSize  = "logFileSize"
	cmdDir          = "dir"
	cmdMaster       = "master"
	cmdMeta         = "meta"
	cmdTikvEndpoins = "tikvEndpoints"
)

var cmdM = map[string]string{
	cmdPort:         port,
	cmdLogLevel:     logLevel,
	cmdLogAge:       logAge,
	cmdLogFileSize:  logFileSize,
	cmdDir:          dir,
	cmdMaster:       cmdMaster,
	cmdMeta:         meta,
	cmdTikvEndpoins: tikvEndpoins,
}

var markRequired = []string{}

var Cmd = &cobra.Command{
	Use:   "fs",
	Short: "youngfs virtual file system",
	Long:  "youngfs virtual file system",
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
		logger := zap.New("fs", logOptions...)
		syncer = append(syncer, logger)

		// kv
		var s3kv kv.TransactionStore
		var s3cnkv kv.TransactionStore
		switch strings.ToLower(viper.GetString(meta)) {
		case "badger":
			s3kv, err = badger.New(path.Join(viper.GetString(dir), "meta", "s3kv"))
			if err != nil {
				return err
			}
			s3cnkv, err = badger.New(path.Join(viper.GetString(dir), "meta", "s3continuekv"))
			if err != nil {
				return err
			}
		case "bboltdb":
			// need a file
			s3kv, err = bbolt.New(path.Join(viper.GetString(dir), "meta", "s3kv.db"), []byte("s3kv"))
			if err != nil {
				return err
			}
			s3cnkv, err = bbolt.New(path.Join(viper.GetString(dir), "meta", "s3continuekv.db"), []byte("s3continuekv"))
			if err != nil {
				return err
			}
		case "leveldb":
			s3kv, err = leveldb.New(path.Join(viper.GetString(dir), "meta", "s3kv"))
			if err != nil {
				return err
			}
			s3cnkv, err = leveldb.New(path.Join(viper.GetString(dir), "meta", "s3continuekv"))
			if err != nil {
				return err
			}
		case "tikv":
			endpoints := viper.GetStringSlice(tikvEndpoins)
			if len(endpoints) == 0 {
				return fmt.Errorf("tikv endpoints cannot be empty")
			}
			s3kv, err = tikv.New(endpoints, tikv.WithKeySpace("s3kv"))
			if err != nil {
				return err
			}
			s3cnkv, err = tikv.New(endpoints, tikv.WithKeySpace("s3continuekv"))
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported kv store type")
		}
		closers = append(closers, s3kv, s3cnkv)

		// engine
		var e engine.Engine
		if viper.GetString(master) != "" {
			e, err = chunk.New(viper.GetString(master))
			if err != nil {
				return err
			}
		} else {
			if err := os.MkdirAll(path.Join(viper.GetString(dir), "data"), 0755); err != nil {
				return err
			}
			e, err = simple.New(path.Join(viper.GetString(dir), "data"))
			if err != nil {
				return err
			}
		}
		// meta
		metaStore := s3.New(s3kv, s3cnkv)
		// server
		svr := server.NewServer(metaStore, e)
		// handler
		h := handler.New(svr, logger)
		// router
		r := router.New(h, router.WithMiddlewares(router.Logger(logger)))
		errChan := make(chan error, 1)
		go func(errChan chan<- error) {
			err := http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt(port)), r)
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
	Cmd.Flags().Int(cmdPort, 9432, "port")
	Cmd.Flags().String(cmdLogLevel, "info", "log level (debug, info, warn, error, dpanic, panic)")
	Cmd.Flags().Uint64(cmdLogAge, 31, "log age (day)")
	Cmd.Flags().Uint64(cmdLogFileSize, 32, "log file max size (MiB)")

	Cmd.Flags().String(cmdDir, ".", "data dir")
	Cmd.Flags().String(cmdMaster, "", "master endpoint")
	Cmd.Flags().String(cmdMeta, "badger", "kv store type [badger, bboltdb, leveldb, tikv]")
	Cmd.Flags().StringSlice(cmdTikvEndpoins, nil, "tikv endpoints")

	Cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(cmdM[flag.Name], Cmd.PersistentFlags().Lookup(flag.Name))
	})
	Cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		_ = viper.BindPFlag(cmdM[flag.Name], Cmd.Flags().Lookup(flag.Name))
	})

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
}
