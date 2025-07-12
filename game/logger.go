package game

import "sync"

type Logger struct {
	enabled bool
	size    int
	logs    []string
}

var loggerInstance *Logger
var loggerOnce sync.Once

func GetLoggerInstance() *Logger {
	loggerOnce.Do(func() {
		loggerInstance = &Logger{
			enabled: true,
			size:    100,
			logs:    make([]string, 0, 100),
		}
	})
	return loggerInstance
}

func (l *Logger) Log(log string) {
	if !l.enabled {
		return
	}
	if len(l.logs) >= l.size {
		l.logs = l.logs[1:]
	}
	l.logs = append(l.logs, log)
}

func (l *Logger) PrintLogs(maxLogs int) {
	if !l.enabled {
		return
	}
	var limit int
	if maxLogs > len(l.logs) {
		limit = len(l.logs)
	} else {
		limit = maxLogs
	}
	for i := len(l.logs) - limit; i < len(l.logs); i++ {
		println(l.logs[i])
	}
}
