package common

import (
	dt "MikeAustin71/datetimeopsgo/datetime"
	"fmt"
	"time"
	"strings"
)

/*  'logjobgroupconfig.go' is located in source code
		repository:

				https://github.com/MikeAustin71/logopsgo.git

 */

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
	AppVersion              string
	AppName                 string
	AppPathFileNameExt      FileMgr
	AppStartTimeDto         dt.TimeZoneDto // Time Zone Dto Containing App StartUp time
	AppDateTimeFormat       string
	AppErrPathFileNameExt   FileMgr
	BaseStartDir            DirMgr
	CurrentDirPath          DirMgr
	CmdPathFileNameExt      FileMgr
	BatchStartTimeDto       dt.TimeZoneDto
	IanaTimeZone            string
	LogPathFileNameExt      FileMgr
	NoOfJobs                int
	LogFileRetentionInDays  int
	KillAllJobsOnFirstError bool
	StartTime               dt.TimeZoneDto
	EndTime                 dt.TimeZoneDto
	Duration                dt.DurationTriad
	Dtfmt                   * dt.FormatDateTimeUtility
	Banner1                 string
	Banner2                 string
	Banner3                 string
	Banner4                 string
	Banner5                 string
	Banner6                	string
	Banner7                	string
	BannerLen               int
	LeftTab                 string
	AppNoOfJobs             int
	NoOfJobGroupMsgs        int
	NoOfJobsCompleted       int
	PurgeFileInfo           DirectoryDeleteFileInfo
	NoOfLogFilesPurged      int
}

// New - Initializes key
// elements of a Logging Configuration
func (logOps *LogJobGroup) New(parent []OpsMsgContextInfo) OpsMsgDto {


	om := logOps.baseLogErrConfig(parent, "New()")

	// Assumes CreateAllFormatsInMemory() has
	// already been called.

	logOps.StartTime = logOps.BatchStartTimeDto.CopyOut()
	logOps.AppName = logOps.AppPathFileNameExt.FileName

	om = logOps.purgeOldLogFiles(om.GetNewParentHistory())

	if om.IsFatalError() {
		return om
	}

	om = logOps.InitializeLogFile(om.GetNewParentHistory())

	if om.IsFatalError() {
		return om
	}


	logOps.BannerLen = 80

	logOps.Banner1 = strings.Repeat("#", logOps.BannerLen) + "\n"

	logOps.Banner2 = strings.Repeat("=", logOps.BannerLen) + "\n"

	logOps.Banner3 = strings.Repeat("*", logOps.BannerLen) + "\n"

	logOps.Banner4 = strings.Repeat("-", logOps.BannerLen) + "\n"

	logOps.Banner5 = strings.Repeat("!", logOps.BannerLen) + "\n"

	logOps.Banner6 = strings.Repeat("&", logOps.BannerLen) + "\n"

	logOps.Banner7 = strings.Repeat("+", logOps.BannerLen) + "\n"

	logOps.LeftTab = strings.Repeat(" ", 2)



	return logOps.writeJobGroupHeaderToLog(om.GetNewParentHistory())

}

// InitializeLogFile - Creates the log directory and
// a new log file. Opens the Log File. New(..) MUST
// be called before this method!
func (logOps *LogJobGroup) InitializeLogFile(parent []OpsMsgContextInfo) OpsMsgDto {

	var err error
	om := logOps.baseLogErrConfig(parent, "InitializeLogFile()")



	if !logOps.LogPathFileNameExt.DMgr.IsInitialized {
		s := "logOps.LogPath DirMgr ERROR!"
		err = fmt.Errorf("ERROR - logOps.LogPath DirMgr NOT Initialized! DirectoryPath='%v'", logOps.LogPathFileNameExt.DMgr.AbsolutePath)
		om.SetFatalError(s, err, 201)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	err = logOps.LogPathFileNameExt.DMgr.IsDirMgrValid("")

	if err != nil {
		s:= fmt.Sprintf("logOps.LogPath DirMgr Object is INVALID! DirectoryPath='%v'", logOps.LogPathFileNameExt.DMgr.AbsolutePath)
		om.SetFatalError(s, err, 203)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	err = logOps.LogPathFileNameExt.DMgr.MakeDir()

	if err != nil {
		s := fmt.Sprintf("MakeAbsolutePath Failed for LogPath '%v'", logOps.LogPathFileNameExt.DMgr.AbsolutePath)
		om.SetFatalError(s, err, 205)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}


	logOps.LogPathFileNameExt.DMgr.DoesDirectoryExist()

	if !logOps.LogPathFileNameExt.DMgr.PathDoesExist || !logOps.LogPathFileNameExt.DMgr.AbsolutePathDoesExist {

			s := "Log Directory Does NOT Exist"
			err = fmt.Errorf("After Make Directory Attempt, Log Directory DOES NOT EXIST! LogPath='%v'", logOps.LogPathFileNameExt.DMgr.AbsolutePath)
			om.SetFatalError(s, err, 207)
			logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
			return om

	}

	err = logOps.LogPathFileNameExt.CreateDirAndFile()

	if err != nil {
		s := fmt.Sprintf("Error returned from LogPathFileNameExt.CreateDirAndFile(). Log File='%v'\n", logOps.LogPathFileNameExt.AbsolutePathFileName )
		om.SetFatalError(s, err, 209)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	om.SetSuccessfulCompletionMessage("Finished InitializeLogFile", 299)

	return om
}


// purgeOldLogFiles - Deletes log files which are older than
// logOps.LogFileRetentionInDays
func (logOps *LogJobGroup) purgeOldLogFiles(parent []OpsMsgContextInfo) OpsMsgDto {

	om := logOps.baseLogErrConfig(parent, "purgeOldLogFiles()")

	if logOps.LogFileRetentionInDays < 0 {
		om.SetSuccessfulCompletionMessage("Finished purgeOldLogFiles - No Log Files To Delete.", 1801)
		return om
	}

	var err error

	if !logOps.LogPathFileNameExt.DMgr.IsInitialized {
		err = fmt.Errorf("Error: LogPathFileNameExt Directory Manager is NOT INITIALIZED!")
		om.SetFatalError("", err, 1801)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om

	}

	err = logOps.LogPathFileNameExt.IsFileMgrValid("")

	if err != nil {
		s := fmt.Sprintf("Error returned by LogPathFileNameExt.IsFileMgrValid(). Error: %v", err.Error())
		om.SetFatalError(s, err, 1803)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}


	if !logOps.LogPathFileNameExt.DMgr.DoesDirMgrAbsolutePathExist() {
		om.SetSuccessfulCompletionMessage("Finished purgeOldLogFiles - Log Directory Does NOT Exist.", 1805)
		return om
	}


	logDur := time.Duration(logOps.LogFileRetentionInDays*24*-1) * time.Hour
	du, err := dt.DurationTriad{}.NewStartTimeDurationCalcTz(
								time.Now().Local(),
								logDur,
								dt.TDurCalcTypeSTDYEARMTH,
								logOps.IanaTimeZone,
								logOps.AppDateTimeFormat)

	thresholdTime := du.LocalTime.EndTimeDateTz.DateTime

	if thresholdTime.IsZero() {
		om.SetSuccessfulCompletionMessage("Finished purgeOldLogFiles - thresholdTime is Zero.", 1807)
		return om
	}

	searchPattern := "*.log"
	filesOlderThan := thresholdTime
	filesNewerThan := time.Time{}

	fsc := FileSelectionCriteria{}

	fsc.FileNamePatterns = []string{searchPattern}
	fsc.FilesOlderThan = filesOlderThan
	fsc.FilesNewerThan = filesNewerThan
	fsc.SelectCriterionMode = ANDFILESELECTCRITERION

	logOps.PurgeFileInfo, err = logOps.LogPathFileNameExt.DMgr.DeleteWalkDirFiles(fsc)

	if err != nil {
		s := fmt.Sprintf("Error returned by LogPathFileNameExt.DMgr.DeleteWalkDirFiles(fsc). Error: %v", err.Error())
		om.SetFatalError(s, err, 1809)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.NoOfLogFilesPurged = len(logOps.PurgeFileInfo.DeletedFiles.FMgrs)
	s:= fmt.Sprintf("Finished purgeOldLogFiles - Successfully Purged %v old Log Files.", logOps.NoOfLogFilesPurged)
	om.SetSuccessfulCompletionMessage(s,1899 )
	return om

}

// writeJobGroupHeaderToLog - Writes the Job Group
// Header info to the log file. This is a one-time
// operation for each log file. Method InitializeLogFile(..)
// MUST be called before first use of this method.
// The Header is always the first element in the Log.
func (logOps *LogJobGroup) writeJobGroupHeaderToLog(parent []OpsMsgContextInfo) OpsMsgDto {
	var err error
	var str, stx string

	om := logOps.baseLogErrConfig(parent, "writeJobGroupHeaderToLog()")

	thisParentInfo := om.GetNewParentHistory()

	if logOps.LogPathFileNameExt.FilePtr == nil {
		s := "logOps.LogPathFileNameExt.FilePtr was not correctly initialized! logOps.LogPathFileNameExt.FilePtr *os.File pointer is nil!"
		om.SetFatalErrorMessage(s, 301)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	su := StringUtility{}
	stx = fmt.Sprintf("App Name: %v  AppVersion: %v \n", logOps.AppName, logOps.AppVersion)
	str, err = su.StrCenterInStrLeft(stx, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStrLeft threw error on AppName AppVersion"
		om.SetFatalError(s, err, 302)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(str, thisParentInfo)

	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	str = fmt.Sprintf("Execution Job Group Command File: %v \n", logOps.CmdPathFileNameExt.FileNameExt)

	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command File Name"
		om.SetFatalError(s, err, 303)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	str = fmt.Sprintf("Starting Execution of %v Jobs. \n", logOps.AppNoOfJobs)

	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Number of Jobs Executed"
		om.SetFatalError(s, err, 304)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)


	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("  Job Group Start Time UTC: %v \n",
		logOps.StartTime.TimeUTC.GetDateTimeYMDAbbrvDowNano())

	logOps.writeTabFileStr(str, 0, parent)

	str = fmt.Sprintf("Job Group Start Time Local: %v \n",
		logOps.StartTime.TimeLocal.GetDateTimeYMDAbbrvDowNano())

	logOps.writeTabFileStr(str, 0, parent)

	localZone := logOps.StartTime.TimeLocal.TimeZone.ZoneName

	str = fmt.Sprintf(" Job Group Local Time Zone: %v - %v \n", logOps.IanaTimeZone, localZone)

	logOps.writeTabFileStr(str, 0, parent)

	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	str = "This Application Executable File:\n"
	logOps.writeTabFileStr(str, 1, thisParentInfo)
	str = logOps.AppPathFileNameExt.AbsolutePathFileName + "\n\n"
	logOps.writeTabFileStr(str, 2, thisParentInfo)


	str = "This Log File:\n"
	logOps.writeTabFileStr(str, 1, thisParentInfo)
	str = logOps.LogPathFileNameExt.AbsolutePathFileName + "\n\n"
	logOps.writeTabFileStr(str, 2, thisParentInfo)

	str = "Command File:\n"
	logOps.writeTabFileStr(str, 1, thisParentInfo)
	str = logOps.CmdPathFileNameExt.AbsolutePathFileName + "\n\n"
	logOps.writeTabFileStr(str, 2, thisParentInfo)

	str = "Initial Application Path: \n"
	logOps.writeTabFileStr(str, 1, thisParentInfo)
	str = logOps.AppPathFileNameExt.DMgr.AbsolutePath + "\n\n"
	logOps.writeTabFileStr(str, 2, thisParentInfo)

	str = "Base Start Directory:\n"
	logOps.writeTabFileStr(str, 1, thisParentInfo)
	str = logOps.BaseStartDir.AbsolutePath + "\n\n"
	logOps.writeTabFileStr(str, 2, thisParentInfo)

	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	str = fmt.Sprintf("Number Of Old Log Files Deleted: %v \n", logOps.NoOfLogFilesPurged)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	if logOps.NoOfLogFilesPurged > 0 {

		str = "!!!!!! Log Files Deleted !!!!!!\n"
		stx, _ = su.StrCenterInStrLeft(str, logOps.BannerLen)

		logOps.writeTabFileStr(stx, 1, thisParentInfo)
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
		for i := 0; i < logOps.NoOfLogFilesPurged; i++ {
			fMgr := logOps.PurgeFileInfo.DeletedFiles.FMgrs[i]

			str = fmt.Sprintf("%v. File Date: %v   File Name: %v \n",
				i+1, fMgr.ActualFileInfo.ModTime().Format(dt.FmtDateTimeSecText),  fMgr.ActualFileInfo.Name())

			logOps.writeTabFileStr(str, 2, thisParentInfo)
		}

		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	}

	logOps.PurgeFileInfo = DirectoryDeleteFileInfo{}
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	// Signal Successful Completion
	om.SetSuccessfulCompletionMessage("Finished writeJobGroupHeaderToLog", 309)
	return om
}

// WriteJobGroupFooterToLog - Writes the trailing Job Group data
// to the Log File. This is the last entry in the Log. The File
// pointer is closed here.
func (logOps *LogJobGroup) WriteJobGroupFooterToLog(cmds CommandBatch, parent []OpsMsgContextInfo) OpsMsgDto {

	var err error
	var str, stx string

	om := logOps.baseLogErrConfig(parent, "WriteJobGroupFooterToLog()")

	thisParentInfo := om.GetNewParentHistory()

	if logOps.LogPathFileNameExt.FilePtr == nil {
		s := "logOps.LogPathFileNameExt.FilePtr was not correctly initialized! logOps.LogPathFileNameExt.FilePtr *os.File pointer is nil!"
		om.SetFatalErrorMessage(s, 8001)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	defer logOps.LogPathFileNameExt.CloseFile()

	logOps.writeFileStr("\n\n", thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	su := StringUtility{}

	str = "Completed Job Group Execution\n"
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Job Group Execution Title"
		om.SetFatalError(s, err, 8003)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("App Name: %v  AppVersion: %v \n", logOps.AppName, logOps.AppVersion)
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on AppName AppVersion"
		om.SetFatalError(s, err, 8005)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("Execution Job Group Command File: %v \n", logOps.CmdPathFileNameExt.FileNameExt)

	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command File Name"
		om.SetFatalError(s, err, 8007)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	stx = fmt.Sprintf("  Number of Jobs Executed: %v \n", logOps.NoOfJobsCompleted)

	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	stx =  fmt.Sprintf("Number of Messages Logged: %v \n", logOps.NoOfJobGroupMsgs)

	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	stx = "Job Group Execution Times:\n"
	logOps.writeTabFileStr(stx, 1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	str = logOps.StartTime.TimeUTC.GetDateTimeYMDAbbrvDowNano()
	stx = fmt.Sprintf("JobGroup   Start Time UTC: %v \n", str)
	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	om2 := cmds.SetBatchEndTime(thisParentInfo)

	if om2.IsFatalError() {
		logOps.AppErrPathFileNameExt.WriteStrToFile(om2.Error())
		return om2
	}

	logOps.EndTime = cmds.CmdJobsHdr.CmdBatchEndTime


	stx = fmt.Sprintf("JobGroup     End Time UTC: %v \n",
				logOps.EndTime.TimeUTC.GetDateTimeYMDAbbrvDowNano())

	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)


	stx = fmt.Sprintf("JobGroup Start Time Local: %v \n",
					logOps.StartTime.TimeLocal.GetDateTimeYMDAbbrvDowNano())

	logOps.writeTabFileStr(stx, 1, thisParentInfo)


	stx = fmt.Sprintf("JobGroup   End Time Local: %v \n",
						logOps.EndTime.TimeLocal.GetDateTimeYMDAbbrvDowNano())

	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	stx = fmt.Sprintf("JobGroup  Local Time Zone: %v - %v\n", logOps.IanaTimeZone,
							logOps.EndTime.TimeLocal.TimeZone.ZoneName)

	logOps.writeTabFileStr(stx, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	stx = "JobGroup Elapsed Time:\n"
	logOps.writeTabFileStr(stx, 1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner6, thisParentInfo)
	stx = cmds.CmdJobsHdr.CmdBatchDuration.UTCTime.GetElapsedTimeStr() + "\n"
	logOps.writeTabFileStr(stx, 1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner6, thisParentInfo)

	logOps.writeFileStr("\n\n", thisParentInfo)

	str = "End of Job Group Execution\n"
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on End of Job Group Execution"
		om.SetFatalError(s, err, 8009)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	logOps.writeFileStr(logOps.Banner1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner1, thisParentInfo)

	// Signal Successful Completion
	om.SetSuccessfulCompletionMessage("Finished WriteJobGroupFooterToLog", 8099)

	return om
}

func (logOps *LogJobGroup) WriteCmdJobHeaderToLog(job *CmdJob, parent []OpsMsgContextInfo) OpsMsgDto {

	om := logOps.baseLogErrConfig(parent, "WriteCmdJobHeaderToLog()")

	thisParentInfo := om.GetNewParentHistory()

	su := StringUtility{}

	str := "\n\n"

	logOps.writeFileStr(str, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)
	str = "Starting Command Job Execution\n"
	stx, err := su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStrLeft threw error on Starting Command Job Execution"
		om.SetFatalError(s, err, 2501)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	str = fmt.Sprintf("Command Job Name: %v\n", job.CmdDisplayName)
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Starting Command Job Display Name"
		om.SetFatalError(s, err, 2502)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	str = fmt.Sprintf("Command Job Number: %v\n", job.CmdJobNo)
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStrLeft threw error on Starting Command Job Number"
		om.SetFatalError(s, err, 2503)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("   Command Job Description: %v\n", job.CmdDescription)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	str = fmt.Sprintf("               Command Type: %v\n", job.CmdDescription)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("   Execute Cmd In Directory: %v\n", job.ExeCmdInDir)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("    Delay Cmd Start Seconds: %v\n", job.DelayCmdStartSeconds)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf(" Delay Cmd To Date Time: %v\n", job.DelayStartCmdDateTime)

	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("    Cmd Time Out in Seconds: %v\n", job.CommandTimeOutInSeconds)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("          Execution Command: %v\n", job.ExeCommand)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	largs := len(job.CmdArguments.CmdArgs)

	if largs < 1 {
    str = "          Command Arguments: NO COMMAND ARGUMENTS PROVIDED\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)

	}	else if len(job.CombinedArguments) < 50 {
		str = fmt.Sprintf("          Command Arguments: %v\n", job.CombinedArguments)
		logOps.writeTabFileStr(str, 1, thisParentInfo)
	} else {
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
		str = "Command Arguments:\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
		for i:= 0; i < largs; i++ {
		 	str = job.CmdArguments.CmdArgs[i] + "\n"
			logOps.writeTabFileStr(str, 3, thisParentInfo)
		}

		logOps.writeFileStr(logOps.Banner4 + "\n", thisParentInfo)
	}

	largs = len(job.CmdInputs.InputArgs)

	if largs < 1 {
		str = "            Input Arguments: NO INPUT ARGUMENTS PROVIDED\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)

	} else if len(job.CombinedInputArguments) < 50 {
		str = fmt.Sprintf("            Input Arguments: %v\n", job.CombinedInputArguments)
		logOps.writeTabFileStr(str, 1, thisParentInfo)

	} else {

		str = "Input Arguments:\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)

		for i:= 0; i < largs; i++ {
			str = job.CmdInputs.InputArgs[i] + "\n"
			logOps.writeTabFileStr(str, 3, thisParentInfo)
		}

		logOps.writeTabFileStr("\n", 1, thisParentInfo)
	}


	logOps.writeFileStr(logOps.Banner3, thisParentInfo)


	str = "Command Job Start Times\n"
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command Job Start Times"
		om.SetFatalError(s, err, 2504)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner3, thisParentInfo)

	str = fmt.Sprintf("     Cmd Job Start Time UTC: %v\n",
		job.CmdJobStartTimeValue.TimeUTC.DateTime.Format(logOps.AppDateTimeFormat))
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("   Cmd Job Start Time Local: %v\n",
		job.CmdJobStartTimeValue.TimeLocal.DateTime.Format(logOps.AppDateTimeFormat))
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	tzLocal := job.CmdJobStartTimeValue.TimeLocal.TimeZone.ZoneName
	str = fmt.Sprintf("    Cmd Job Local Time Zone: %v - %v\n", job.IanaTimeZone, tzLocal)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner2, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)


	om.SetSuccessfulCompletionMessage("Finished WriteCmdJobHeaderToLog", 2509)
	return om
}

func (logOps *LogJobGroup) WriteCmdJobFooterToLog(job *CmdJob, parent []OpsMsgContextInfo) OpsMsgDto {

	om := logOps.baseLogErrConfig(parent, "WriteCmdJobHeaderToLog()")
	thisParentInfo := om.GetNewParentHistory()
	su := StringUtility{}

	logOps.writeFileStr("\n\n", thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)
	logOps.writeFileStr(logOps.Banner2, thisParentInfo)
	str := "Completed Command Job Execution\n"
	stx, err := su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Completed Command Job Execution"
		om.SetFatalError(s, err, 2601)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner5, thisParentInfo)
	str = "Job Execution Results\n"
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Completed Command Job Execution"
		om.SetFatalError(s, err, 2602)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	str = fmt.Sprintf("Command Job Name: %v\n", job.CmdDisplayName)
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Completing Command Job Display Name"
		om.SetFatalError(s, err, 2603)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	str = fmt.Sprintf("Command Job Number: %v\n", job.CmdJobNo)
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Completing Command Job Number"
		om.SetFatalError(s, err, 2604)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	if job.CmdJobExecutionStatus=="" {
		job.CmdJobExecutionStatus = "Successful Completion"
	}

	logOps.writeFileStr(logOps.Banner5, thisParentInfo)
	str = fmt.Sprintf("Command Execution Status: %v\n", job.CmdJobExecutionStatus)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner2, thisParentInfo)


	str = fmt.Sprintf("Cmd Job Description: %v\n", job.CmdDescription)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("  Execution Command: %v\n", job.ExeCommand)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	largs := len(job.CmdArguments.CmdArgs)

	if largs < 1 {
		str = "          Command Arguments: NO COMMAND ARGUMENTS PROVIDED\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)

	}	else if len(job.CombinedArguments) < 50 {
		str = fmt.Sprintf("          Command Arguments: %v\n", job.CombinedArguments)
		logOps.writeTabFileStr(str, 1, thisParentInfo)
	} else {
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
		str = "Command Arguments:\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)
		for i:= 0; i < largs; i++ {
			str = job.CmdArguments.CmdArgs[i] + "\n"
			logOps.writeTabFileStr(str, 3, thisParentInfo)
		}

		logOps.writeFileStr(logOps.Banner4 + "\n", thisParentInfo)
	}

	largs = len(job.CmdInputs.InputArgs)

	if largs < 1 {
		str = "            Input Arguments: NO INPUT ARGUMENTS PROVIDED\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)

	} else if len(job.CombinedInputArguments) < 50 {
		str = fmt.Sprintf("            Input Arguments: %v\n", job.CombinedInputArguments)
		logOps.writeTabFileStr(str, 1, thisParentInfo)

	} else {

		str = "Input Arguments:\n"
		logOps.writeTabFileStr(str, 1, thisParentInfo)
		logOps.writeFileStr(logOps.Banner4, thisParentInfo)

		for i:= 0; i < largs; i++ {
			str = job.CmdInputs.InputArgs[i] + "\n"
			logOps.writeTabFileStr(str, 3, thisParentInfo)
		}

		logOps.writeTabFileStr("\n", 1, thisParentInfo)
	}

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	str = "Command Job Execution Times\n"
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on Command Job Execution Times"
		om.SetFatalError(s, err, 2605)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)
	logOps.writeFileStr(logOps.Banner4, thisParentInfo)
	logOps.writeTabFileStr("UTC Start\\End Times:\n", 1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	str = fmt.Sprintf(" Cmd Job Start Time UTC: %v\n",
			job.CmdJobStartTimeValue.TimeUTC.DateTime.Format(logOps.AppDateTimeFormat))
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("   Cmd Job End Time UTC: %v\n",
				job.CmdJobEndTimeValue.TimeUTC.DateTime.Format(logOps.AppDateTimeFormat))
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	logOps.writeTabFileStr("Local Start\\End Times:\n", 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	str = fmt.Sprintf("Cmd Job Start Time Local: %v\n",
		job.CmdJobStartTimeValue.TimeLocal.DateTime.Format(logOps.AppDateTimeFormat))
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("  Cmd Job End Time Local: %v\n",
			job.CmdJobEndTimeValue.TimeLocal.DateTime.Format(logOps.AppDateTimeFormat))

	logOps.writeTabFileStr(str, 1, thisParentInfo)

	tzLocal := job.CmdJobEndTimeValue.TimeLocal.TimeZone.ZoneName
	str = fmt.Sprintf(" Cmd Job Local Time Zone: %v - %v\n", job.IanaTimeZone, tzLocal)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner4, thisParentInfo)

	str = "Cmd Job Execution Elapsed Time:\n"
	logOps.writeTabFileStr(str, 1, thisParentInfo)
	logOps.writeFileStr(logOps.Banner6, thisParentInfo)

	str = fmt.Sprintf("%v\n", job.CmdJobDuration.UTCTime.GetElapsedTimeStr())
	logOps.writeTabFileStr(str, 2, thisParentInfo)
	logOps.writeFileStr(logOps.Banner6 + "\n", thisParentInfo)

	str = fmt.Sprintf("      Number of Messages: %v\n", job.CmdJobNoOfMsgs)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	str = fmt.Sprintf("Number of Error Messages: %v\n", job.CmdJobNoOfErrorMsgs)
	logOps.writeTabFileStr(str, 1, thisParentInfo)

	logOps.writeFileStr(logOps.Banner2, thisParentInfo)

	str = fmt.Sprintf("End of Command Job Number: %v\n", job.CmdJobNo)
	stx, err = su.StrCenterInStrLeft(str, logOps.BannerLen)
	if err != nil {
		s := "StrCenterInStr threw error on End of Command Job Number"
		om.SetFatalError(s, err, 2606)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		return om
	}

	logOps.writeFileStr(stx, thisParentInfo)

	logOps.writeFileStr(logOps.Banner2, thisParentInfo)
	logOps.writeFileStr("\n\n", thisParentInfo)

	job.CmdJobIsCompleted = true

	om.SetSuccessfulCompletionMessage("Finished WriteCmdJobFooterToLog", 2609)
	return om

}

func (logOps *LogJobGroup) WriteOpsMsgToLog(opsMsgToLog OpsMsgDto, job *CmdJob, parent []OpsMsgContextInfo) OpsMsgDto {


	om := logOps.baseLogErrConfig(parent, "WriteOpsMsgToLog()")

	thisParentInfo := om.GetNewParentHistory()

	logOps.writeTabFileStr(opsMsgToLog.String(), 2, thisParentInfo)

	om.SetSuccessfulCompletionMessage("Finished WriteOpsMsgToLog", 2709)

	return om
}

func (logOps *LogJobGroup) writeTabFileStr(s string, noOfTabs int, parent []OpsMsgContextInfo) {

	om := logOps.baseLogErrConfig(parent, "writeTabFileStr()")

	stx := ""

	for i := 0; i < noOfTabs; i++ {
		stx += logOps.LeftTab
	}

	stx += s

	_, err := logOps.LogPathFileNameExt.WriteStrToFile(stx)

	if err != nil {
		s := fmt.Sprintf("LogFilePtr.WriteString threw error on string: '%v'", stx)
		om.SetFatalError(s, err, 1001)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		panic(om)
	}

}

func (logOps *LogJobGroup) writeFileStr(s string, parent []OpsMsgContextInfo) {

	_, err := logOps.LogPathFileNameExt.WriteStrToFile(s)

	om := logOps.baseLogErrConfig(parent, "writeFileStr()")

	if err != nil {
		s := fmt.Sprintf("LogFilePtr.WriteString threw error on string: '%v'", s)
		om.SetFatalError(s, err, 901)
		logOps.AppErrPathFileNameExt.WriteStrToFile(om.Error())
		panic(om)
	}

}

// baseLogErrConfig - Used internally by LogJobGroup
// methods to set up error messages.
func (logOps *LogJobGroup) baseLogErrConfig(parent []OpsMsgContextInfo, funcName string) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
							SourceFileName: LogGroupConfigSrcFile,
							ParentObjectName: "LogJobGroup",
							FuncName: funcName,
							BaseMessageId: LogGroupConfigErrBlockNo,
						}

	return OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)
}
