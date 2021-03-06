package common

/*  'logLevel.go' is located in source code
		repository:

		https://github.com/MikeAustin71/logopsgo.git

 */

// LogLevel - Holds the Message level indicating the relative importance of a specific log Message.
type LogLevel int

func (level LogLevel) String() string {
	return LogLevelNames[level]
}

const (
	// LogDEBUG - Message
	LogDEBUG LogLevel = iota
	// LogOPERROR - 1 Message is an Error Message
	LogOPERROR
	// LogFATAL - 2 Message is a Fatal Error Message
	LogFATAL
	// LogINFO - 3 Message is an Informational Message
	LogINFO
	// LogWARN - 4 Message is a warning Message
	LogWARN
)

// LogLevelNames - string array containing names of Log Levels
var LogLevelNames = [...]string{"DEBUG", "OPERROR", "FATAL", "INFO", "WARN"}
