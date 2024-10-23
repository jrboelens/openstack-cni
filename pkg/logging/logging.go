package logging

import (
	"io"
	"sync"

	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger
var locker sync.RWMutex

func SetupLogging(name string, opts httplog.Options, output io.Writer) zerolog.Logger {
	locker.Lock()
	defer locker.Unlock()
	logger = httplog.NewLogger(name, opts)
	if output != nil {
		logger = logger.Output(zerolog.ConsoleWriter{Out: output, TimeFormat: opts.TimeFieldFormat})
	}
	return logger
}

func Log() *zerolog.Logger {
	locker.Lock()
	defer locker.Unlock()
	return &logger
}

func Error(msg string, err error) {
	logger.Error().AnErr("err", err).Msg(msg)
}

func AddStrings(event *zerolog.Event, strs [][]string) *zerolog.Event {
	lastEvent := event
	for _, pair := range strs {
		if len(pair) < 2 {
			continue
		}
		lastEvent = event.Str(pair[0], pair[1])
	}
	return lastEvent
}
