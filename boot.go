package micron

import (
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
	go AppCtx.waitForErrorSignal()

	boot := &Boot{configPath: bootfile}
	AppCtx.config = boot.loadConfig()

	return boot
}

func (boot *Boot) RootConfig(config interface{}) *Boot {
	AppCtx.GetRootConfig(config)
	return boot
}

func (boot *Boot) WaitForTerminateSignal() {
	AppCtx.waitForTerminateSignal()
}

func (boot *Boot) ErrorSignal(err error) {
	if err != nil {
		AppCtx.ErrorSignal <- err
	}
}

func (boot *Boot) HandleErrorFunc(fn ErrorFunc) {
	AppCtx.errorFunc = fn
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
		AppCtx.ErrorSignal <- err
	}

	return conf
}
