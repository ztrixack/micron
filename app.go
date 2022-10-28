package micron

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

type appContext struct {
	startTime       time.Time
	config          *viper.Viper
	ErrorSignal     chan error
	errorFunc       ErrorFunc
	terminateSignal chan os.Signal
	terminateFunc   []TerminateFunc
}

type TerminateFunc func(context.Context)

type ErrorFunc func(error)

var (
	AppCtx = &appContext{
		startTime:       time.Now(),
		ErrorSignal:     make(chan error, 1),
		errorFunc:       fatalHandler,
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
func (ctx *appContext) GetRootConfig(out interface{}) {
	if err := defaults.Set(out); err != nil {
		AppCtx.ErrorSignal <- err
	}

	if err := ctx.config.Unmarshal(out); err != nil {
		AppCtx.ErrorSignal <- err
	}
}

func (ctx *appContext) GetViper() *viper.Viper {
	return ctx.config
}

func (ctx *appContext) SetDefault(key string, value interface{}) {
	ctx.config.SetDefault(key, value)
}

func (ctx *appContext) waitForErrorSignal() {
	err := <-ctx.ErrorSignal
	ctx.errorFunc(err)
	os.Exit(1)
}

// waitForTerminateSignal waits for shutdown signal.
func (ctx *appContext) waitForTerminateSignal() {
	<-ctx.terminateSignal
	bgCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()

	for _, fn := range ctx.terminateFunc {
		fn(bgCtx)
	}
}

// SetErrorFunc waits for critical error signal.
func (ctx *appContext) SetErrorFunc(fn ErrorFunc) {
	ctx.errorFunc = fn
}

// WaitForTerminateSignal waits for shutdown signal.
func (ctx *appContext) AddTerminateFunc(fn TerminateFunc) {
	ctx.terminateFunc = append(ctx.terminateFunc, fn)
}

func fatalHandler(err error) {
	log.Printf("\tCRITICAL\t[micron]\t %v\n", err)
}
