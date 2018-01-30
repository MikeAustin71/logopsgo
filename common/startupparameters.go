package common

import (
	"time"
)

/*  'startupparameters.go' is located in source code
		repository:

		https://github.com/MikeAustin71/logopsgo.git

*/


const (
	srcFileNameStartUpParms = "startupparameters.go"
	errBlockNoStartUpParms  = int64(351230000)
)

// StartupParameters - This data structure is
// used by the calling method to initialize a
// LogJobGroup using the 'New(...)' method.
type StartupParameters struct {
	AppName                 string
	AppExeFileName          string
	AppPath									string
	BaseStartDir						string
	AppVersion              string
	CommandFileName         string
	StartTimeUTC						time.Time
	StartTime               time.Time
	AppLogPath              string
	AppLogFileName          string
	NoOfJobs                int
	LogFileRetentionInDays  int
	CmdExeDir               string
	KillAllJobsOnFirstError bool
	IanaTimeZone            string
	LogMode                 LoggingMode
	Dtfmt										* DateTimeFormatUtility
}


// baseLogErrConfig - Used internally by LogJobGroup
// methods to set up error messages.
func (sUp *StartupParameters) startUpBaseErrConfig(parent []OpsMsgContextInfo, funcName string) OpsMsgDto {

	// bi := ErrBaseInfo{}.New(srcFileNameStartUpParms, funcName, errBlockNoStartUpParms)
	msgCtx := OpsMsgContextInfo{
							SourceFileName: srcFileNameStartUpParms,
							ParentObjectName: "StartupParameters",
							FuncName: funcName,
							BaseMessageId: errBlockNoStartUpParms,
						}

	return OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)
}
