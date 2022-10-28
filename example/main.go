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
	defer boot.WaitForTerminateSignal()

	err := log.New(conf.Log).Build().Err()
	boot.ErrorSignal(err)

	log.I("ENV", zap.String("LOG_ENV", conf.Log.Env))
	log.D("Debug")
	log.I("Info")
	log.W("Warn")
	log.E("Error")
	// log.Critical("Critical")
}
