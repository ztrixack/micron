package micron

import (
	"context"
	"os"
	"path"

	"github.com/ztrixack/micron/log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Boot struct {
	Id         string
	ConfigPath string
}

// NewBoot create a bootstrapper.
func NewBoot() *Boot {
	boot := &Boot{
		Id: uuid.New().String(),
	}
	log.Init(log.Config{
		Development: true,
	})
	AppCtx.rawConfig = boot.readYaml()
	return boot
}

func (boot *Boot) WaitForTerminateSignal(ctx context.Context) {
	log.Event(LOG, "Service Started")
	AppCtx.WaitForTerminateSignal()
}

// readYaml read YAML file
func (boot *Boot) readYaml() []byte {
	if len(boot.ConfigPath) < 1 {
		boot.ConfigPath = "boot.yaml"
	}

	if !path.IsAbs(boot.ConfigPath) {
		wd, _ := os.Getwd()
		boot.ConfigPath = path.Join(wd, boot.ConfigPath)
	}

	log.Info(LOG, zap.String("configPath", boot.ConfigPath))
	res, err := os.ReadFile(boot.ConfigPath)
	if err != nil {
		log.Error(LOG, err, zap.String("configPath", boot.ConfigPath))
	}

	return res
}
