package common

// LogLevel - Holds the message level indicating the relative importance of a specific log message.
type LogLevel int

func (level LogLevel) String() string {
	return LogLevelNames[level]
}

const (
	// DEBUG - message
	DEBUG LogLevel = iota
	// OPERROR - 1 Message is an Error Message
	OPERROR
	// FATAL - 2 Message is a Fatal Error Message
	FATAL
	// INFO - 3 Message is an Informational Message
	INFO
	// WARN - 4 Message is a warning Message
	WARN
)

// LogLevelNames - string array containing names of Log Levels
var LogLevelNames = [...]string{"DEBUG", "OPERROR", "FATAL", "INFO", "WARN"}
