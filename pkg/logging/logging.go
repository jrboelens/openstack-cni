package logging

import (
	"sync"

	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

var logger zerolog.Logger
var locker sync.RWMutex

func SetupLogging(name string, opts httplog.Options) zerolog.Logger {
	locker.Lock()
	defer locker.Unlock()
	logger = httplog.NewLogger(name, opts)
	return logger
}

func Log() *zerolog.Logger {
	locker.RLock()
	defer locker.RUnlock()
	return &logger
}

func Error(msg string, err error) {
	logger.Error().AnErr("err", err).Msg("failed to run server")
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
