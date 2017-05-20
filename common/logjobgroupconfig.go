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

// LogStartupParameters - This data structure is
// used by the calling method to initialize a
// LogJobGroupConfig using the 'New(...)' method.
type LogStartupParameters struct {
	AppName         string
	AppExeFileName  string
	AppVersion      string
	CommandFileName string
	StartTime       time.Time
	AppLogDir       string
	NoOfJobs        int
	LogMode         LoggingMode
}

// LogJobGroupConfig - holds logging configuration for the
// current group of jobs
type LogJobGroupConfig struct {
	LogMode            LoggingMode
	StartTime          time.Time
	EndTime            time.Time
	Duration           time.Time
	CommandFileName    string
	AppName            string
	AppExeFileName     string
	AppVersion         string
	AppLogFileName     string
	AppLogDir          string
	AppLogPathFileName string
	FilePtr            *os.File
	Banner1            string
	Banner2            string
	Banner3            string
	Banner4            string
	BannerLen          int
	LeftTab            string
	BaseStartDir       string
	AppNoOfJobs        int
	NoOfMessagesLogged int
	NoOfJobsCompleted  int
}

// New - Initializes key
// elements of a Logging Configuration
func (logCfg *LogJobGroupConfig) New(startParams LogStartupParameters,
	parent []ErrBaseInfo) SpecErr {

	var err error
	var isPanic bool
	se := logCfg.BaseLogErrConfig(parent, "New()")

	logCfg.LogMode = startParams.LogMode
	logCfg.StartTime = startParams.StartTime
	logCfg.CommandFileName = startParams.CommandFileName
	logCfg.AppName = startParams.AppName
	logCfg.AppExeFileName = startParams.AppExeFileName
	logCfg.AppVersion = startParams.AppVersion
	logCfg.AppNoOfJobs = startParams.NoOfJobs

	dt := DateTimeUtility{}
	logCfg.AppLogFileName = logCfg.AppName + "_" + dt.GetDateTimeStr(logCfg.StartTime) + ".log"

	fh := FileHelper{}
	logCfg.AppLogDir, err = fh.MakeAbsolutePath(startParams.AppLogDir)

	if err != nil {
		s := fmt.Sprintf("MakeAbsolutePath Failed for AppLogDir '%v'", startParams.AppLogDir)
		isPanic = true
		return se.New(s, err, isPanic, 101)
	}

	logCfg.AppLogPathFileName = fh.JoinPathsAdjustSeparators(logCfg.AppLogDir, logCfg.AppLogFileName)

	logCfg.BannerLen = 79

	logCfg.Banner1 = strings.Repeat("#", logCfg.BannerLen) + "\n"

	logCfg.Banner2 = strings.Repeat("=", logCfg.BannerLen) + "\n"

	logCfg.Banner3 = strings.Repeat("*", logCfg.BannerLen) + "\n"

	logCfg.Banner4 = strings.Repeat("-", logCfg.BannerLen) + "\n"

	logCfg.LeftTab = strings.Repeat(" ", 2)

	logCfg.BaseStartDir, err = fh.GetAbsCurrDir()

	if err != nil {
		s := "GetAbsCurrDir() Failed!"
		isPanic = true
		return se.New(s, err, isPanic, 102)
	}

	// Signal Successful Completion
	return SpecErr{}.SignalNoErrors()
}

// InitializeLogFile - Creates the log directory and
// a new log file. Opens the Log File. New(..) MUST
// be called before this method!
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
// operation for each log file. Method InitializeLogFile(..)
// MUST be called before first use of this method.
func (logCfg *LogJobGroupConfig) WriteFileGroupHeaderToLog(parent []ErrBaseInfo) SpecErr {
	var err error
	var isPanic bool
	var str, stx string

	se := logCfg.BaseLogErrConfig(parent, "WriteFileGroupHeaderToLog()")

	thisParentInfo := se.AddBaseToParentInfo()

	if logCfg.FilePtr == nil {
		s := "logCfg.FilePtr was not correctly initialized! logCfg.FilePtr *os.File pointer is nil!"
		isPanic = true
		return se.New(s, err, isPanic, 301)
	}

	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)

	su := StringUtility{}
	stx = fmt.Sprintf("App Name: %v  AppVersion: %v ", logCfg.AppName, logCfg.AppVersion)
	str, err = su.StrCenterInStr(stx, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on AppName AppVersion"
		isPanic = true
		return se.New(s, err, isPanic, 302)
	}

	logCfg.WriteFileString(str, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)

	str = fmt.Sprintf("Execution Job Group Command File: %v", logCfg.CommandFileName)

	stx, err = su.StrCenterInStr(str, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command File Name"
		isPanic = true
		return se.New(s, err, isPanic, 303)
	}

	logCfg.WriteFileString(stx, thisParentInfo)

	str = fmt.Sprintf("Starting Execution of %v Jobs.", logCfg.AppNoOfJobs)

	stx, err = su.StrCenterInStr(str, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Number of Jobs Executed"
		isPanic = true
		return se.New(s, err, isPanic, 304)
	}

	logCfg.WriteFileString(stx, thisParentInfo)

	dt := DateTimeUtility{}

	str = fmt.Sprintf("Job Execution Start Time: %v", dt.GetDateTimeEverything(logCfg.StartTime))

	stx, err = su.StrCenterInStr(str, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Number of Job Execution Start Time"
		isPanic = true
		return se.New(s, err, isPanic, 305)
	}

	logCfg.WriteFileString(stx, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)

	// Signal Successful Completion
	return SpecErr{}.SignalNoErrors()
}

func (logCfg *LogJobGroupConfig) WriteFileGroupFooterToLog(parent []ErrBaseInfo) SpecErr {

	var err error
	var isPanic bool
	var str, stx string

	se := logCfg.BaseLogErrConfig(parent, "WriteFileGroupFooterToLog()")

	thisParentInfo := se.AddBaseToParentInfo()

	if logCfg.FilePtr == nil {
		s := "logCfg.FilePtr was not correctly initialized! logCfg.FilePtr *os.File pointer is nil!"
		isPanic = true
		return se.New(s, err, isPanic, 801)
	}

	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)

	su := StringUtility{}
	str = fmt.Sprintf("App Name: %v  AppVersion: %v ", logCfg.AppName, logCfg.AppVersion)
	stx, err = su.StrCenterInStr(stx, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on AppName AppVersion"
		isPanic = true
		return se.New(s, err, isPanic, 802)
	}

	logCfg.WriteFileString(stx, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)

	str = fmt.Sprintf("Execution Job Group Command File: %v", logCfg.CommandFileName)

	stx, err = su.StrCenterInStr(str, logCfg.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command File Name"
		isPanic = true
		return se.New(s, err, isPanic, 803)
	}

	logCfg.WriteFileString(stx, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)

	stx = logCfg.LeftTab + fmt.Sprintf("  Number of Jobs Executed: %v ", logCfg.NoOfJobsCompleted)

	logCfg.WriteFileString(stx, thisParentInfo)

	stx = logCfg.LeftTab + fmt.Sprintf("Number of Messages Logged: %v ", logCfg.NoOfMessagesLogged)

	logCfg.WriteFileString(stx, thisParentInfo)

	logCfg.WriteFileString(logCfg.Banner3, thisParentInfo)

	dt := DateTimeUtility{}
	str = dt.GetDateTimeNanoSecText(logCfg.StartTime)
	stx = logCfg.LeftTab + fmt.Sprintf("JobGroup   Start Time: %v ", str)
	logCfg.WriteFileString(stx, thisParentInfo)

	str = dt.GetDateTimeNanoSecText(logCfg.EndTime)
	stx = logCfg.LeftTab + fmt.Sprintf("JobGroup     End Time: %v ", str)
	logCfg.WriteFileString(stx, thisParentInfo)

	ed, err2 := dt.GetElapsedTime(logCfg.StartTime, logCfg.EndTime)

	if err2 != nil {
		s := fmt.Sprintf("DateTimeUtility:GetElapsedTime threw error on Start '%v' and End Time '%v'", logCfg.StartTime, logCfg.EndTime)
		isPanic = true
		return se.New(s, err, isPanic, 804)

	}

	stx = logCfg.LeftTab + fmt.Sprintf("JobGroup Elapsed Time: %v ", ed.NanosecStr)
	logCfg.WriteFileString(stx, thisParentInfo)

	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)
	logCfg.WriteFileString(logCfg.Banner1, thisParentInfo)

	// Signal Successful Completion
	return SpecErr{}.SignalNoErrors()
}

func (logCfg *LogJobGroupConfig) WriteFileString(s string, parent []ErrBaseInfo) {

	_, err := logCfg.FilePtr.WriteString(s)

	if err != nil {
		s := fmt.Sprintf("FilePtr.WriteString threw error on string: '%v'", s)
		se :=
			logCfg.BaseLogErrConfig(parent, "WriteFileString()").
				New(s, err, true, 901)
		panic(se)
	}

}

// BaseLogErrConfig - Used internally by LogJobGroupConfig
// methods to set up error messages.
func (logCfg *LogJobGroupConfig) BaseLogErrConfig(parent []ErrBaseInfo, funcName string) SpecErr {

	bi := ErrBaseInfo{}.New(LogGroupConfigSrcFile, funcName, LogGroupConfigErrBlockNo)

	return SpecErr{}.InitializeBaseInfo(parent, bi)
}
