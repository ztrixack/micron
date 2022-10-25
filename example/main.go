package main

import (
	"github.com/ztrixack/micron"
	"github.com/ztrixack/micron/modules/log"
	"go.uber.org/zap"
)

type RootConfig struct {
	Log log.Config
}

func init() {
	// os.Setenv("LOG_ENV", "production")
}

func main() {
	conf := &RootConfig{}
	boot := micron.NewBoot().RootConfig(conf)

	log.New(conf.Log).Build()

	log.Info("ENV", zap.String("LOG_ENV", conf.Log.Env))

	log.Debug("Debug")
	log.Info("Info")
	log.Warn("Warn")
	log.Error("Error")
	// log.Critical("Critical")

	boot.WaitForTerminateSignal()
}
