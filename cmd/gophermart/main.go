package main

import (
	"context"
	"errors"
	"gophermart-points/internal/repo/pgsql"
	"gophermart-points/internal/srv"
	"gophermart-points/internal/srv/external"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Address        string `mapstructure:"RUN_ADDRESS"`
	DBConn         string `mapstructure:"DATABASE_URI"`
	AccrualAddress string `mapstructure:"ACCRUAL_SYSTEM_ADDRESS"`
	AuthKey        string `mapstructure:"AUTH_KEY"`
}

func init() {
	pflag.StringP("RUN_ADDRESS", "a", "localhost:8080", "Address to run server on")
	pflag.StringP("DATABASE_URI", "d", "", "DB connection string")
	pflag.StringP("ACCRUAL_SYSTEM_ADDRESS", "r", "localhost:8081", "Address of accrual service")
	pflag.StringP("AUTH_KEY", "k", "", "Authentication key")
}

func LoadConfig(path string) (cfg Config, err error) {
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.AddConfigPath(path)
	viper.SetConfigName("conf")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func initLogger(mode string) *zap.Logger {
	var core zapcore.Core
	if mode == "prod" {
		f, err := os.OpenFile("server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}

		prodcfg := zap.NewProductionEncoderConfig()
		fileEncoder := zapcore.NewJSONEncoder(prodcfg)
		sync := zapcore.AddSync(f)
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, sync, zapcore.InfoLevel),
		)
	} else {
		std := zapcore.AddSync(os.Stdout)

		devcfg := zap.NewDevelopmentEncoderConfig()
		devcfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(devcfg)
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, std, zapcore.InfoLevel),
		)
	}

	l := zap.New(core)
	defer l.Sync()

	return l
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	var err error
	logger := initLogger("dev").Sugar()
	cfg, err := LoadConfig("./")
	if err != nil {
		logger.Warn("failed to load config, cause: %s:", err.Error())
	}

	var db *pgsql.PgSQL
	db, err = pgsql.Init(ctx, cfg.DBConn, logger)
	if err != nil {
		logger.Fatalf("failed to init DB, cause: %s:", err.Error())
	}

	srv := srv.New(cfg.Address, db, logger)
	orderPoolProc := external.NewOrderProcessing(ctx, db, cfg.AccrualAddress, 36, logger)
	srv.MountHandlers(cfg.AuthKey, orderPoolProc)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		orderPoolProc.RunPool(4)
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		db.Close()

		go func() error {
			<-timeoutCtx.Done()
			if timeoutCtx.Err() == context.DeadlineExceeded {
				return errors.New("timed out performing graceful shutdown")
			}

			return nil
		}()

		return srv.Shutdown(timeoutCtx)
	})

	logger.Infow("server was started", "addr:", cfg.Address)
	if err = g.Wait(); err != nil {
		logger.Errorf("exit reason: %s", err)
	}
}
