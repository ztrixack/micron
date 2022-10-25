package micron

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Boot struct {
	configPath string
}

func NewBoot() *Boot {
	return NewBootWithPath("boot")
}

func NewBootWithPath(bootfile string) *Boot {
	boot := &Boot{configPath: bootfile}
	AppCtx.config = boot.loadConfig()

	return boot
}

func (boot *Boot) RootConfig(config interface{}) *Boot {
	AppCtx.GetRootConfig(config)

	return boot
}

func (boot *Boot) WaitForTerminateSignal() {
	AppCtx.WaitForTerminateSignal()
}

func (boot *Boot) loadConfig() *viper.Viper {
	conf := viper.New()
	conf.AddConfigPath(".")
	conf.SetConfigName(boot.configPath)
	conf.SetConfigType("yaml")
	conf.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	conf.AutomaticEnv()

	err := conf.ReadInConfig()
	if err != nil {
		log.Fatalf("fatal error config file: %s \n", err)
	}

	return conf
}
