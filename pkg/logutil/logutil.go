package logutil

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func Fatal(logger log.Logger, err error, msg string, extra ...interface{}) {
	if err == nil {
		level.Debug(logger).Log("err", err, "msg", msg)
		return
	}
	level.Info(logger).Log("err", err, "msg", msg)
	os.Exit(1)
}

func New() *log.SwapLogger {
	base := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	base = log.With(base, "ts", log.DefaultTimestampUTC)
	base = level.NewInjector(base, level.InfoValue())
	lev := level.AllowInfo()

	var swapLogger log.SwapLogger
	swapLogger.Swap(level.NewFilter(base, lev))

	go swapLevelHandler(base, &swapLogger, false)
	return &swapLogger
}

type key int

const loggerKey key = 0

func NewContext(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) log.Logger {
	v, ok := ctx.Value(loggerKey).(log.Logger)
	if !ok {
		return log.NewNopLogger()
	}

	return v
}

func swapLevelHandler(base log.Logger, swapLogger *log.SwapLogger, debug bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGUSR2)
	for {
		<-sigChan
		if debug {
			newLogger := level.NewFilter(base, level.AllowInfo())
			swapLogger.Swap(newLogger)
		} else {
			newLogger := level.NewFilter(base, level.AllowDebug())
			swapLogger.Swap(newLogger)
		}
		level.Info(swapLogger).Log("msg", "swapping level", "debug", !debug)
		debug = !debug
	}
}
