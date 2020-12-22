package log

//Priority simply limits an integer to a set of values acceptable by the logger
type Priority int

//GlobalLogger is used for all your logging needs
var GlobalLogger Logger = NewDefaultLogger()

const (
	//Emergency means everything is chaos and we should give up
	Emergency Priority = iota
	//Alert means we should give up and try later
	Alert
	//Critical means a serious issue has appeared and we should give up
	Critical
	//Error indicates an unexpected condition has occured that we can not recover from
	Error
	//Warning indicates an unexpected condition has occured that we can recover from
	Warning
	//Notice indicates that an unusual decision was made
	Notice
	//Info simply informes about execution flow
	Info
	//Debug contains everything
	Debug
)

//PriorityString converts a numerical value to the log level string representation
var PriorityString = map[Priority]string{
	Emergency: "EMERGENCY",
	Alert:     "ALERT",
	Critical:  "CRITICAL",
	Error:     "ERROR",
	Warning:   "WARNING",
	Notice:    "NOTICE",
	Info:      "INFO",
	Debug:     "DEBUG",
}

//Logger defines the interface required for logging systems
type Logger interface {
	Log(priority Priority, message string, data map[string]string)
}

//Log logs something
func Log(priority Priority, message string, data map[string]string) {
	GlobalLogger.Log(priority, message, data)
}
