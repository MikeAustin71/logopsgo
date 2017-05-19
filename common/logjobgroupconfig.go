package common

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	// LogGroupConfigSrcFile - Log Group Cfg source
	// code file name.
	LogGroupConfigSrcFile = "logjobgroupconfig.go"

	// LogGroupConfigErrBlockNo - The block of numbers
	// assigned as error codes for this source file.
	LogGroupConfigErrBlockNo = int64(647000000)
)

// LogJobGroupConfig - holds logging configuration for the
// current group of jobs
type LogJobGroupConfig struct {
	LogMode            LoggingMode
	StartTime          time.Time
	EndTime            time.Time
	Duration           time.Time
	CommandFileName    string
	AppNameVersion     string
	AppLogFileName     string
	AppLogDir          string
	AppLogPathFileName string
	FilePtr            *os.File
	Banner1            string
	Banner2            string
	Banner3            string
	Banner4            string
	BannerLen          int
	BaseStartDir       string
	AppNoOfJobs        int
}

// NewLogGroupConfig - Initializes key
// elements of a Logging Configuration
func (logCfg *LogJobGroupConfig) NewLogGroupConfig(appNameVer string, commandFileName string, t time.Time, parent []ErrBaseInfo) SpecErr {
	var err error
	var isPanic bool
	se := logCfg.BaseLogErrConfig(parent, "NewLogGroupConfig()")

	logCfg.LogMode = LogVERBOSE
	logCfg.StartTime = t
	logCfg.CommandFileName = commandFileName
	logCfg.AppNameVersion = appNameVer
	dt := DateTimeUtility{}
	logCfg.AppLogFileName = "CmdrX" + "_" + dt.GetDateTimeStr(logCfg.StartTime) + ".log"

	fh := FileHelper{}
	logCfg.AppLogDir, err = fh.MakeAbsolutePath("./CmdrX")

	if err != nil {
		s := "MakeAbsolutePath Failed for './CmdrX'"
		isPanic = true
		return se.New(s, err, isPanic, 101)
	}

	logCfg.AppLogPathFileName = fh.JoinPathsAdjustSeparators(logCfg.AppLogDir, logCfg.AppLogFileName)

	logCfg.BannerLen = 79

	logCfg.Banner1 = strings.Repeat("#", logCfg.BannerLen) + "\n"

	logCfg.Banner2 = strings.Repeat("=", logCfg.BannerLen) + "\n"

	logCfg.Banner3 = strings.Repeat("*", logCfg.BannerLen) + "\n"

	logCfg.Banner4 = strings.Repeat("-", logCfg.BannerLen) + "\n"

	logCfg.BaseStartDir, err = fh.GetAbsCurrDir()

	if err != nil {
		s := "GetAbsCurrDir() Failed!"
		isPanic = true
		return se.New(s, err, isPanic, 102)
	}

	// Signal Successful Completion
	return SpecErr{}.SignalNoErrors()
}

// InitializeLogFile - Initializes the log directory and
// a new log file.
func (logCfg *LogJobGroupConfig) InitializeLogFile(parent []ErrBaseInfo) SpecErr {

	var err error
	var isPanic bool
	se := logCfg.BaseLogErrConfig(parent, "InitializeLogFile()")

	fh := FileHelper{}

	if !fh.DoesFileExist(logCfg.AppLogDir) {

		err = fh.MakeDirAll(logCfg.AppLogDir)

		if err != nil {
			s := fmt.Sprintf("MakeDirAll() Failed on Dir: %v", logCfg.AppLogDir)
			isPanic = true
			return se.New(s, err, isPanic, 201)
		}

	}

	// At this point, logCfg.AppLogDir exists.
	// Check to determine if the log file exists.
	// If log file already exists, delete it.
	if fh.DoesFileExist(logCfg.AppLogPathFileName) {

		err = fh.DeleteDirFile(logCfg.AppLogPathFileName)

		if err != nil {
			s := fmt.Sprintf("DeleteDirFile() Failed on File: %v", logCfg.AppLogPathFileName)
			isPanic = true
			return se.New(s, err, isPanic, 202)
		}
	}

	// Create a new log file
	logCfg.FilePtr, err = fh.CreateFile(logCfg.AppLogPathFileName)

	if err != nil {
		s := fmt.Sprintf("CreateFile() Failed on File: %v", logCfg.AppLogPathFileName)
		isPanic = true
		return se.New(s, err, isPanic, 203)
	}

	// Signal Successful Completion
	return SpecErr{}.SignalNoErrors()
}

// WriteFileGroupHeaderToLog - Writes the Job Group
// Header info to the log file. This is a one-time
// operation for each log file. BE SURE that
// logCfg.AppNoOfJobs is previously initialized!
func (logCfg *LogJobGroupConfig) WriteFileGroupHeaderToLog(parent []ErrBaseInfo, exeVersionNo string) SpecErr {
	var err error
	var isPanic bool
	var str, stx string

	se := logCfg.BaseLogErrConfig(parent, "WriteFileGroupHeaderToLog()")

	if logCfg.FilePtr == nil {
		s := "logCfg.FilePtr was not correctly initialized! logCfg.FilePtr *os.File pointer is nil!"
		isPanic = true
		return se.New(s, err, isPanic, 301)
	}

	logCfg.FilePtr.WriteString(logCfg.Banner1)
	logCfg.FilePtr.WriteString(logCfg.Banner1)

	su := StringUtility{}

	str, err = su.StrCenterInStr(logCfg.AppNameVersion, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr through error on AppNameVersion"
		isPanic = true
		return se.New(s, err, isPanic, 302)
	}

	logCfg.FilePtr.WriteString(str)
	logCfg.FilePtr.WriteString(logCfg.Banner1)
	str = fmt.Sprintf("Execution Job Group: %v Commands", logCfg.CommandFileName)

	stx, err = su.StrCenterInStr(str, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command File Name"
		isPanic = true
		return se.New(s, err, isPanic, 303)
	}

	logCfg.FilePtr.WriteString(stx)

	str = fmt.Sprintf("Starting Execution of %v Jobs.", logCfg.AppNoOfJobs)

	stx, err = su.StrCenterInStr(str, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Number of Jobs Executed"
		isPanic = true
		return se.New(s, err, isPanic, 304)
	}

	// Signal Successful Completion
	return SpecErr{}.SignalNoErrors()
}

// BaseLogErrConfig - Used internally by LogJobGroupConfig
// methods to set up error messages.
func (logCfg *LogJobGroupConfig) BaseLogErrConfig(parent []ErrBaseInfo, funcName string) SpecErr {

	bi := ErrBaseInfo{}.New(LogGroupConfigSrcFile, funcName, LogGroupConfigErrBlockNo)

	return SpecErr{}.InitializeBaseInfo(parent, bi)
}
