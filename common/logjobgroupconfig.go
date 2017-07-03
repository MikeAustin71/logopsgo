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

// LogJobGroup - holds logging configuration for the
// current group of jobs
type LogJobGroup struct {
	LogMode                 LoggingMode
	StartTimeUTC            time.Time
	StartTime               time.Time
	EndTimeUTC              time.Time
	EndTime                 time.Time
	Duration                time.Time
	Dtfmt                   *DateTimeFormatUtility
	CommandFileName         string
	AppName                 string
	AppExeFileName          string
	AppVersion              string
	AppPath                 string
	AppLogFileName          string
	AppLogDir               string
	AppLogPathFileName      string
	AppLogDirWalkInfo       DirWalkInfo
	FilePtr                 *os.File
	Banner1                 string
	Banner2                 string
	Banner3                 string
	Banner4                 string
	BannerLen               int
	LeftTab                 string
	BaseStartDir            string
	AppNoOfJobs             int
	NoOfJobGroupMsgs        int
	NoOfJobMsgs             int
	NoOfJobsCompleted       int
	LogFileRetentionInDays  int
	NoOfLogFilesPurged      int
	CmdExeDir               string
	KillAllJobsOnFirstError bool
	IanaTimeZone            string
}

// New - Initializes key
// elements of a Logging Configuration
func (logOps *LogJobGroup) New(startParams StartupParameters,
	parent []ErrBaseInfo) SpecErr {

	var err error
	var isPanic bool

	se := logOps.baseLogErrConfig(parent, "New()")

	// Assumes CreateAllFormatsInMemory() has
	// already been called.
	logOps.Dtfmt = startParams.Dtfmt

	logOps.LogMode = startParams.LogMode
	logOps.StartTimeUTC = startParams.StartTimeUTC
	logOps.StartTime = startParams.StartTime
	logOps.CommandFileName = startParams.CommandFileName
	logOps.AppName = startParams.AppName
	logOps.AppExeFileName = startParams.AppExeFileName
	logOps.AppPath = startParams.AppPath
	logOps.BaseStartDir = startParams.BaseStartDir
	logOps.AppVersion = startParams.AppVersion
	logOps.AppNoOfJobs = startParams.NoOfJobs
	logOps.LogFileRetentionInDays = startParams.LogFileRetentionInDays
	logOps.CmdExeDir = startParams.CmdExeDir
	logOps.KillAllJobsOnFirstError = startParams.KillAllJobsOnFirstError
	logOps.IanaTimeZone = startParams.IanaTimeZone

	fh := FileHelper{}

	if startParams.AppLogFileName == "" {
		logOps.AppLogFileName = logOps.AppName
	}

	dt := DateTimeUtility{}
	logOps.AppLogFileName = logOps.AppLogFileName + "_" + dt.GetDateTimeStr(logOps.StartTime) + ".log"

	logOps.AppLogDir, err = fh.MakeAbsolutePath(startParams.AppLogPath)

	if err != nil {
		s := fmt.Sprintf("MakeAbsolutePath Failed for AppLogPath '%v'", startParams.AppLogPath)
		isPanic = true
		return se.New(s, err, isPanic, 101)
	}

	logOps.AppLogPathFileName = fh.JoinPathsAdjustSeparators(logOps.AppLogDir, logOps.AppLogFileName)

	logOps.BannerLen = 80

	logOps.Banner1 = strings.Repeat("#", logOps.BannerLen) + "\n"

	logOps.Banner2 = strings.Repeat("=", logOps.BannerLen) + "\n"

	logOps.Banner3 = strings.Repeat("*", logOps.BannerLen) + "\n"

	logOps.Banner4 = strings.Repeat("-", logOps.BannerLen) + "\n"

	logOps.LeftTab = strings.Repeat(" ", 2)

	logOps.BaseStartDir, err = fh.GetAbsCurrDir()

	if err != nil {
		s := "GetAbsCurrDir() Failed!"
		isPanic = true
		return se.New(s, err, isPanic, 102)
	}

	return logOps.InitializeLogFile(se.AddBaseToParentInfo())

}

// InitializeLogFile - Creates the log directory and
// a new log file. Opens the Log File. New(..) MUST
// be called before this method!
func (logOps *LogJobGroup) InitializeLogFile(parent []ErrBaseInfo) SpecErr {

	var err error
	var isPanic bool
	se := logOps.baseLogErrConfig(parent, "InitializeLogFile()")

	fh := FileHelper{}

	if !fh.DoesFileExist(logOps.AppLogDir) {

		err = fh.MakeDirAll(logOps.AppLogDir)

		if err != nil {
			s := fmt.Sprintf("MakeDirAll() Failed on Dir: %v", logOps.AppLogDir)
			isPanic = true
			return se.New(s, err, isPanic, 201)
		}

	}

	// At this point, logOps.AppLogPath exists.
	// Check to determine if the log file exists.
	// If log file already exists, delete it.
	if fh.DoesFileExist(logOps.AppLogPathFileName) {

		err = fh.DeleteDirFile(logOps.AppLogPathFileName)

		if err != nil {
			s := fmt.Sprintf("DeleteDirFile() Failed on File: %v", logOps.AppLogPathFileName)
			isPanic = true
			return se.New(s, err, isPanic, 202)
		}
	}

	// Create a new log file
	logOps.FilePtr, err = fh.CreateFile(logOps.AppLogPathFileName)

	if err != nil {
		s := fmt.Sprintf("CreateFile() Failed on File: %v Error: %v", logOps.AppLogPathFileName, err.Error())
		isPanic = true
		return se.New(s, err, isPanic, 203)
	}

	si := logOps.purgeOldLogFiles(se.AddBaseToParentInfo())

	if si.IsErr {
		si.IsPanic = true
		return si
	}

	return logOps.writeFileGroupHeaderToLog(se.AddBaseToParentInfo())
}

// purgeOldLogFiles - Deletes log files which are older than
// logOps.LogFileRetentionInDays
func (logOps *LogJobGroup) purgeOldLogFiles(parent []ErrBaseInfo) SpecErr {
	se := logOps.baseLogErrConfig(parent, "purgeOldLogFiles()")

	fh := FileHelper{}

	if logOps.LogFileRetentionInDays < 0 {
		return se.SignalNoErrors()
	}

	logDur := time.Duration(logOps.LogFileRetentionInDays*24*-1) * time.Hour
	du := DurationUtility{}
	du.SetStartTimeDuration(time.Now(), logDur)
	thresholdTime := du.StartDateTime

	if thresholdTime.IsZero() {
		return se.SignalNoErrors()
	}

	logOps.AppLogDirWalkInfo = DirWalkInfo{}

	logOps.AppLogDirWalkInfo.StartPath = logOps.AppLogDir
	logOps.AppLogDirWalkInfo.PatternMatch = "*.log"
	logOps.AppLogDirWalkInfo.DeleteFilesOlderThan = thresholdTime
	err := fh.WalkDirPurgeFilesOlderThan(&logOps.AppLogDirWalkInfo)
	if err != nil {
		s := fmt.Sprintf("WalkDirPurgeFilesOlderThan() Failed on File: %v Error: %v", logOps.AppLogPathFileName, err.Error())
		return se.New(s, err, true, 1801)
	}

	logOps.NoOfLogFilesPurged = len(logOps.AppLogDirWalkInfo.DeletedFiles)

	return se.SignalNoErrors()

}

// writeFileGroupHeaderToLog - Writes the Job Group
// Header info to the log file. This is a one-time
// operation for each log file. Method InitializeLogFile(..)
// MUST be called before first use of this method.
// The Header is always the first element in the Log.
func (logOps *LogJobGroup) writeFileGroupHeaderToLog(parent []ErrBaseInfo) SpecErr {
	var err error
	var isPanic bool
	var str, stx string

	se := logOps.baseLogErrConfig(parent, "writeFileGroupHeaderToLog()")

	thisParentInfo := se.AddBaseToParentInfo()

	if logOps.FilePtr == nil {
		s := "logOps.FilePtr was not correctly initialized! logOps.FilePtr *os.File pointer is nil!"
		isPanic = true
		return se.New(s, err, isPanic, 301)
	}

	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	su := StringUtility{}
	stx = fmt.Sprintf("App Name: %v  AppVersion: %v \n", logOps.AppName, logOps.AppVersion)
	str, err = su.StrCenterInStr(stx, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on AppName AppVersion"
		isPanic = true
		return se.New(s, err, isPanic, 302)
	}

	logOps.writeFileStr(str, thisParentInfo)

	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	str = fmt.Sprintf("Execution Job Group Command File: %v \n", logOps.CommandFileName)

	stx, err = su.StrCenterInStr(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command File Name"
		isPanic = true
		return se.New(s, err, isPanic, 303)
	}

	logOps.writeFileStr(stx, thisParentInfo)

	str = fmt.Sprintf("Starting Execution of %v Jobs. \n", logOps.AppNoOfJobs)

	stx, err = su.StrCenterInStr(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Number of Jobs Executed"
		isPanic = true
		return se.New(s, err, isPanic, 304)
	}

	logOps.writeFileStr(stx, thisParentInfo)

	dt := DateTimeUtility{}

	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("  Job Execution Start Time UTC: %v \n", dt.GetDateTimeEverything(logOps.StartTimeUTC))

	logOps.writeTabFileStr(str, 0, parent)

	str = fmt.Sprintf("Job Execution Start Time Local: %v \n", dt.GetDateTimeEverything(logOps.StartTime))

	logOps.writeTabFileStr(str, 0, parent)

	str = fmt.Sprintf(" Job Execution Local Time Zone: %v \n", logOps.IanaTimeZone)

	logOps.writeTabFileStr(str, 0, parent)

	logOps.writeFileStr(logOps.Banner2, thisParentInfo)
	logOps.writeFileStr("\n", thisParentInfo)
	logOps.writeFileStr(logOps.Banner3, thisParentInfo)
	logOps.writeFileStr("\n", thisParentInfo)

	str = "Initial Application Path: \n"

	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = logOps.AppPath + "\n"

	logOps.writeTabFileStr(str, 2, thisParentInfo)

	logOps.writeFileStr(logOps.Banner3, thisParentInfo)
	logOps.writeFileStr("\n", thisParentInfo)

	str = "This Log File:\n"

	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = logOps.AppLogPathFileName + "\n"

	logOps.writeTabFileStr(str, 2, thisParentInfo)

	logOps.writeFileStr(logOps.Banner3, thisParentInfo)
	logOps.writeFileStr("\n", thisParentInfo)

	str = "Base Start Directory:\n"
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = logOps.BaseStartDir + "\n"
	logOps.writeTabFileStr(str, 2, thisParentInfo)
	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	logOps.writeFileStr("\n", thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("Number Of Old Log Files Deleted: %v \n", logOps.NoOfLogFilesPurged)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	if logOps.NoOfLogFilesPurged > 0 {

		str = "!!!!!! Log Files Deleted !!!!!!\n"
		stx, _ = su.StrCenterInStr(str, logOps.BannerLen)

		logOps.writeTabFileStr(stx, 1, thisParentInfo)
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
		for i := 0; i < logOps.NoOfLogFilesPurged; i++ {
			fi := logOps.AppLogDirWalkInfo.DeletedFiles[i]

			str = fmt.Sprintf("%v. File Date: %v   File Name: %v \n",
				i+1, fi.Info.ModTime().Format(FmtDateTimeSecText), fi.Info.Name())

			logOps.writeTabFileStr(str, 2, thisParentInfo)
		}

		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	}

	logOps.AppLogDirWalkInfo = DirWalkInfo{}
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	// Signal Successful Completion
	return se.SignalNoErrors()
}

// WriteJobGroupFooterToLog - Writes the trailing Job Group data
// to the Log File. This is the last entry in the Log. The File
// pointer is closed here.
func (logOps *LogJobGroup) WriteJobGroupFooterToLog(parent []ErrBaseInfo) SpecErr {

	var err error
	var isPanic bool
	var str, stx string

	se := logOps.baseLogErrConfig(parent, "WriteJobGroupFooterToLog()")

	thisParentInfo := se.AddBaseToParentInfo()

	if logOps.FilePtr == nil {
		s := "logOps.FilePtr was not correctly initialized! logOps.FilePtr *os.File pointer is nil!"
		isPanic = true
		return se.New(s, err, isPanic, 801)
	}

	defer logOps.FilePtr.Close()

	logOps.writeFileStr("\n\n", thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	su := StringUtility{}

	str = "Completed Job Group Execution\n"
	stx, err = su.StrCenterInStr(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Job Group Execution Title"
		isPanic = true
		return se.New(s, err, isPanic, 802)
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("App Name: %v  AppVersion: %v \n", logOps.AppName, logOps.AppVersion)
	stx, err = su.StrCenterInStr(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on AppName AppVersion"
		isPanic = true
		return se.New(s, err, isPanic, 803)
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("Execution Job Group Command File: %v \n", logOps.CommandFileName)

	stx, err = su.StrCenterInStr(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command File Name"
		isPanic = true
		return se.New(s, err, isPanic, 804)
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	stx = fmt.Sprintf("  Number of Jobs Executed: %v \n", logOps.NoOfJobsCompleted)

	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	stx = logOps.LeftTab + fmt.Sprintf("Number of Messages Logged: %v \n", logOps.NoOfJobGroupMsgs)

	logOps.writeFileStr(stx, thisParentInfo)

	logOps.writeFileStr(logOps.Banner3, thisParentInfo)
	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	dt := DateTimeUtility{}
	dur := DurationUtility{}
	tz := TimeZoneUtility{}
	str = dt.GetDateTimeNanoSecText(logOps.StartTimeUTC)
	stx = fmt.Sprintf("JobGroup   Start Time UTC: %v \n", str)
	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	logOps.EndTimeUTC = time.Now().UTC()
	tzLocal, _ := tz.ConvertTz(logOps.EndTimeUTC, logOps.IanaTimeZone)

	logOps.EndTime = tzLocal.TimeOut
	str = dt.GetDateTimeNanoSecText(logOps.EndTimeUTC)
	stx = fmt.Sprintf("JobGroup     End Time UTC: %v \n", str)
	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	str = dt.GetDateTimeNanoSecText(logOps.StartTime)
	stx = fmt.Sprintf("JobGroup Start Time Local: %v \n", str)
	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	str = dt.GetDateTimeNanoSecText(logOps.EndTime)
	stx = fmt.Sprintf("JobGroup   End Time Local: %v \n", str)
	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	dur.SetStartEndTimes(logOps.StartTimeUTC, logOps.EndTimeUTC)
	ed, err2 := dur.GetYearsMthsWeeksTimeAbbrv()

	if err2 != nil {
		s := fmt.Sprintf("DateTimeUtility:GetElapsedTime threw error on Start '%v' and End Time '%v'", logOps.StartTime, logOps.EndTime)
		isPanic = true
		return se.New(s, err, isPanic, 805)
	}

	logOps.writeFileStr(logOps.Banner3, thisParentInfo)
	stx = "JobGroup Elapsed Time:\n"
	logOps.writeTabFileStr(stx, 1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	stx = ed.DisplayStr + "\n"
	logOps.writeTabFileStr(stx, 2, thisParentInfo)
	logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	logOps.writeFileStr(logOps.Banner3, thisParentInfo)
	logOps.writeFileStr("\n", thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	// Signal Successful Completion
	return se.SignalNoErrors()
}

func (logOps *LogJobGroup) writeTabFileStr(s string, noOfTabs int, parent []ErrBaseInfo) {

	stx := ""

	for i := 0; i < noOfTabs; i++ {
		stx += logOps.LeftTab
	}

	stx += s

	_, err := logOps.FilePtr.WriteString(stx)

	if err != nil {
		s := fmt.Sprintf("FilePtr.WriteString threw error on string: '%v'", stx)
		se :=
			logOps.baseLogErrConfig(parent, "writeTabFileStr()").
				New(s, err, true, 1001)
		panic(se)
	}

}

func (logOps *LogJobGroup) writeFileStr(s string, parent []ErrBaseInfo) {

	_, err := logOps.FilePtr.WriteString(s)

	if err != nil {
		s := fmt.Sprintf("FilePtr.WriteString threw error on string: '%v'", s)
		se :=
			logOps.baseLogErrConfig(parent, "writeFileStr()").
				New(s, err, true, 901)
		panic(se)
	}

}

// baseLogErrConfig - Used internally by LogJobGroup
// methods to set up error messages.
func (logOps *LogJobGroup) baseLogErrConfig(parent []ErrBaseInfo, funcName string) SpecErr {

	bi := ErrBaseInfo{}.New(LogGroupConfigSrcFile, funcName, LogGroupConfigErrBlockNo)

	return SpecErr{}.InitializeBaseInfo(parent, bi)
}
