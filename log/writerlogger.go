package log

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"
)

//DefaultLogger keeps track of a writer and allows logging to that writer
type DefaultLogger struct {
	Writer      io.Writer
	writerMutex sync.Mutex
	Level       Priority
}

//Log simply logs a log message to the log. loggy log
func (logger *DefaultLogger) Log(priority Priority, message string, data map[string]string) {
	if priority > logger.Level {
		// Not high enough log level
		return
	}
	logMessage := time.Now().Format(time.RFC3339)
	logMessage += " " + PriorityString[priority]

	var dataKeys []string
	if data != nil {
		for key := range data {
			dataKeys = append(dataKeys, key)
		}
	}
	sort.Strings(dataKeys)

	for _, key := range dataKeys {
		logMessage += fmt.Sprintf(" %s=%s", key, data[key])
	}

	logger.writerMutex.Lock()
	defer logger.writerMutex.Unlock()
	fmt.Fprintf(logger.Writer, "%s %s\n", logMessage, message)
}

//NewDefaultLogger creates a new default logger with the defaultiest defaults there are
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		Writer: os.Stderr,
		Level:  Info,
	}
}
