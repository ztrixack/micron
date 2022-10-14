package micron

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ztrixack/micron/log"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

const LOG = "BOOT"

type appContext struct {
	startTime       time.Time
	rawConfig       []byte
	terminateSignal chan os.Signal
	terminateFunc   []TerminateFunc
}

type TerminateFunc func(context.Context)

var (
	AppCtx = &appContext{
		startTime:       time.Now(),
		terminateSignal: make(chan os.Signal, 1),
		terminateFunc:   make([]TerminateFunc, 0),
	}
)

// Init global app context with bellow fields.
func init() {
	signal.Notify(
		AppCtx.terminateSignal,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
}

// WaitForTerminateSignal waits for shutdown signal.
func (ctx *appContext) GetConfig(out interface{}) {
	if err := yaml.Unmarshal(ctx.rawConfig, out); err != nil {
		log.Error(LOG, err, zap.String("rawConfig", string(ctx.rawConfig)))
	}
}

// WaitForTerminateSignal waits for shutdown signal.
func (ctx *appContext) WaitForTerminateSignal() {
	<-ctx.terminateSignal

	bgCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	for _, fn := range ctx.terminateFunc {
		fn(bgCtx)
	}
}

// GetTerminateSignal returns shutdown signal.
func (ctx *appContext) GetTerminateSignal() chan os.Signal {
	return ctx.terminateSignal
}

// WaitForTerminateSignal waits for shutdown signal.
func (ctx *appContext) AddTerminateFunc(fn TerminateFunc) {
	ctx.terminateFunc = append(ctx.terminateFunc, fn)
}
