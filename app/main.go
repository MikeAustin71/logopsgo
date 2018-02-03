package main

import (
	"MikeAustin71/logopsgo/common"
	"fmt"
	"time"
)

const (
	srcFileNameLogOpsMain = "main.go"
	errBlockNoLogOpsMain  = int64(10000000)
	appBanner1            = "===================================================================="
	applicationVersion		= "2.0.0"
	debugAppPath			= "D:/go/work/src/MikeAustin71/logopsgo/app"
	debugAppName 			= "cmdrX"
)


func main() {


	parent := common.OpsMsgContextInfo{
		SourceFileName: srcFileNameLogOpsMain,
		ParentObjectName: "",
		FuncName: "main",
		BaseMessageId: errBlockNoLogOpsMain,
	}

	om := common.OpsMsgDto{}.InitializeWithMessageContext(parent)

	isRunningInDebugMode := true

	lg := common.LogJobGroup{}

	om2 := assembleAppPathFiles(isRunningInDebugMode, &lg, om.GetNewParentHistory())

	if om2.IsFatalError() {
		om2.PrintToConsole()
		return
	}


	su := common.StringUtility{}
	lBanner1 := len(appBanner1)
	fmt.Println("\n\n"+ appBanner1)
	strx, _ := su.StrCenterInStr(lg.AppPathFileNameExt.FileNameExt, lBanner1 - 2)
	fmt.Println("=" + strx + "=")
	fmt.Println(appBanner1)
	fmt.Println("Current Directory: ", lg.CurrentDirPath.AbsolutePath)
	fmt.Println("Executable Directory: ", lg.AppPathFileNameExt.DMgr.AbsolutePath)

	parentHistory := []common.OpsMsgContextInfo{parent}



	cmds, om3 := startUp(&lg, parentHistory)

	if om3.IsFatalError() {
		panic(om3)
	}

	runCmdJobs(&lg, cmds, parentHistory)

	doLogWrapUp(&lg, cmds, parentHistory)
}


func assembleAppPathFiles(isDebugMode bool, lg *common.LogJobGroup, parent []common.OpsMsgContextInfo) common.OpsMsgDto {

	fh := common.FileHelper{}
	var err error

	msgCtx := common.OpsMsgContextInfo{
		SourceFileName:"main.go",
		ParentObjectName: "main()",
		FuncName:" assembleAppPathFiles",
		BaseMessageId: errBlockNoLogOpsMain,
	}

	om := common.OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	// Set Application Version
	lg.AppVersion = applicationVersion

	// Compute App Start Time
	lg.AppStartTimeTzu, err = common.TimeZoneUtility{}.ConvertTz(time.Now().UTC(), "Local")

	if err != nil {
		s:= fmt.Sprintf("TimeZoneUtility{}.ConvertTz Error - Failed to convert UTC to local Time Zone. Start UTC: %v. Time Zone 'Local': %v", time.Now().UTC(), "Local")
		om.SetFatalError(s, err, 3701)
		return om
	}

	// Setup App Path
	if isDebugMode {
		appPath, err  := common.DirMgr{}.New(debugAppPath)

		if err!=nil {
			s:= fmt.Sprintf("DEBUG MODE-Error returned from DirMgr{}.New(debugAppPath) debugAppPath='%v'\n", debugAppPath)
			om.SetFatalError(s, err, 3703)
			return om
		}

		lg.AppPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(appPath, debugAppName + ".exe" )

		if err!=nil {
			s:= fmt.Sprintf("DEBUG MODE-Error returned from FileMgr{}.NewFromDirMgrFileNameExt(appPath, appName + \".exe\" ) (lg.AppPath.Path='%v' appFileNameExt='%v' \n", appPath.Path, debugAppName + ".exe")
			om.SetFatalError(s, err, 3705)
			return om
		}

	} else {
		// NOT DEBUG BUG - ACTUAL MODE!
		// App Path File Name Ext

		appExePathFileNameStr, err := fh.GetExecutablePathFileName()

		if err!=nil {
			om.SetFatalError("ACTUAL MODE-fh.GetExecutablePathFileName() FAILED!\n", err, 3707)
			return om
		}

		lg.AppPathFileNameExt, err = common.FileMgr{}.New(appExePathFileNameStr)

		if err!=nil {
			s:= fmt.Sprintf("ACTUAL MODE- Error returned from FileMgr{}.New(appExePathFileNameStr) appExePathFileNameStr='%v'!\n", appExePathFileNameStr)
			om.SetFatalError(s, err, 3709)
			return om
		}

	}

	err = lg.AppPathFileNameExt.IsFileMgrValid("")

	if err!=nil {
		s:= fmt.Sprintf("Error: lg.AppPathFileNameExt is INVALID! lg.AppPathFileNameExt='%v' \n", lg.AppPathFileNameExt.AbsolutePathFileName)
		om.SetFatalError(s, err, 3711)
		return om
	}

	if !lg.AppPathFileNameExt.AbsolutePathFileNameDoesExist {
		s:= "Error: lg.AppPathFileNameExt DOES NOT EXIST!  \n"
		om.SetFatalError(s, fmt.Errorf("File Does NOT Exist! lg.AppPathFileNameExt='%v'\n", lg.AppPathFileNameExt.AbsolutePathFileName), 3715)
		return om
	}

	// Command Path

	cmdPath := lg.AppPathFileNameExt.DMgr.CopyOut()
	cmdFileNameExt := lg.AppPathFileNameExt.FileName + "Cmds.xml"
	
	// Command Path File Name Ext
	lg.CmdPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(cmdPath,cmdFileNameExt)

	if err!=nil {
		s:= fmt.Sprintf("Error returned by FileMgr{}.NewFromDirMgrFileNameExt(cmdPath,cmdPathFileName ) cmdPath.Path='%v' cmdPathFileName='%v'\n", cmdPath.Path, cmdFileNameExt)
		om.SetFatalError(s, err, 3717)
		return om
	}

	if !lg.CmdPathFileNameExt.AbsolutePathFileNameDoesExist {
		s:= "Error: XML Commands File Does NOT EXIST!"		
		err = fmt.Errorf("XML Command File DOES NOT EXIST! CmdPathFileNameExt='%v'", lg.CmdPathFileNameExt.AbsolutePathFileName)
		om.SetFatalError(s, err, 3719)
		return om
	}
	
	// Log Path
	appLogDirectory := lg.AppPathFileNameExt.FileName + "Log"
	logPath := lg.AppPathFileNameExt.DMgr.GetPathWithSeparator() + appLogDirectory

	logPathDMgr, err := common.DirMgr{}.New(logPath)

	if err!=nil {
		s:= fmt.Sprintf("Error returned by DirMgr{}.New(logPath) logPath='%v'\n", logPath)
		om.SetFatalError(s, err, 3721)
		return om
	}

	// Log Path File Name Ext
	dt := common.DateTimeUtility{}
	dateTimeStamp := dt.GetDateTimeStr(lg.AppStartTimeTzu.TimeOut)
	logFileNameExt :=  lg.AppPathFileNameExt.FileName + "_" + dateTimeStamp + ".log"
	lg.LogPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(logPathDMgr, logFileNameExt)

	if err != nil {
		s:= fmt.Sprintf("lg.LogPathFileNameExt Error returned from FileMgr{}.NewFromDirMgrFileNameExt(logPathDMgr, logFileNameExt) logPathDMgr.Path='%v' logFileNameExt='%v' \n", logPathDMgr.Path, logFileNameExt)
		om.SetFatalError(s, err, 3725)
		return om
	}

	// App Error Path File Name Ext
	appErrFileNameExt := lg.AppPathFileNameExt.FileName + "_Errors_" + dateTimeStamp + ".txt"

	lg.AppErrPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(lg.LogPathFileNameExt.DMgr,appErrFileNameExt)

	if err != nil {
		s:= fmt.Sprintf("lg.AppErrPathFileNameExt - Error returned from FileMgr{}.NewFromDirMgrFileNameExt(lg.LogPathFileNameExt.DMgr,appErrFileNameExt) lg.LogPath.Path='%v' appErrFileNameExt='%v' \n", lg.LogPathFileNameExt.DMgr, appErrFileNameExt)
		om.SetFatalError(s, err, 3727)
		return om
	}

	// Set Current Directory Path
	currDirPath, err := fh.GetCurrentDir()

	if err!=nil {
		om.SetFatalError("fh.GetCurrentDir() FAILED!\n", err, 3731)
		return om
	}

	lg.CurrentDirPath, err = common.DirMgr{}.New(currDirPath)

	if err!=nil {
		s := fmt.Sprintf("Error returned by DirMgr{}.New(currDirPath). currDirPath='%v'\n",currDirPath)
		om.SetFatalError(s, err,3733)
		return om
	}


	return om.SignalNoErrors(3799)
}

func doLogWrapUp(lg *common.LogJobGroup, cmds common.CommandBatch, parent []common.OpsMsgContextInfo) common.OpsMsgDto {

	msgCtx := common.OpsMsgContextInfo{
							SourceFileName:"main.go",
							ParentObjectName: "main()",
							FuncName:" doLogWrapUp",
							BaseMessageId: errBlockNoLogOpsMain,
						}

	om1 := common.OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	defer lg.LogPathFileNameExt.CloseFile()

	// Closes lg.LogFilePtr
	om2 := lg.WriteJobGroupFooterToLog(cmds, om1.GetNewParentHistory())

	if om2.IsFatalError(){
		om2.PrintToConsole()
		return om2
	}

	fmt.Println(appBanner1)
	fmt.Println("See Log File:")
	fmt.Println(lg.LogPathFileNameExt.AbsolutePathFileName)

	fmt.Println(appBanner1)

	return om2
}

func startUp(lg *common.LogJobGroup, parent []common.OpsMsgContextInfo) (common.CommandBatch, common.OpsMsgDto) {

	om := baseLogMsgConfigMain(parent, "startUp")
	parentHistory := om.GetNewParentHistory()


	cmds, om5 := common.ParseXML(lg.CmdPathFileNameExt, parentHistory)

	if om5.IsFatalError() {
		panic(om5)
	}

	om6 := cmds.FormatCmdParameters(parentHistory)

	if om6.IsFatalError() {
		panic(om6)
	}


	om7 := cmds.SetBatchStartTime(lg.AppStartTimeTzu, parentHistory)
	if om7.IsFatalError() {
		return cmds, om7
	}

	tzxu := common.TimeZoneUtility{}

	isValidTz, _, _ := tzxu.IsValidTimeZone(cmds.CmdJobsHdr.IanaTimeZone)

	if !isValidTz {
		cmds.CmdJobsHdr.IanaTimeZone = "Local"
	}

	lg.IanaTimeZone = 	cmds.CmdJobsHdr.IanaTimeZone
	lg.KillAllJobsOnFirstError = cmds.CmdJobsHdr.KillAllJobsOnFirstError
	lg.LogFileRetentionInDays = cmds.CmdJobsHdr.LogFileRetentionInDays


	var err error

	lg.BatchStartTimeTzu, err = common.TimeZoneUtility{}.New(cmds.CmdJobsHdr.CmdBatchStartUTC, lg.IanaTimeZone)

	if err != nil {

		om.SetFatalError("Error from TimeZoneUtility{}.New(cmds.CmdJobsHdr.CmdBatchStartUTC, lg.IanaTimeZone)\n", err, 4701)
		return cmds, om
	}

	dtf := common.DateTimeFormatUtility{}
	dtf.CreateAllFormatsInMemory()
	lg.LogMode = common.LogVERBOSE
	lg.BaseStartDir = lg.AppPathFileNameExt.DMgr.CopyOut()
	lg.NoOfJobs = cmds.CmdJobsHdr.NoOfCmdJobs
	lg.Dtfmt = &dtf

	om8 := lg.New(parentHistory)

	if om8.IsFatalError() {
		panic(om8)
	}

	dur := common.DurationUtility{}

	time.Sleep(dur.GetDurationFromSeconds(10))

	om9 := cmds.SetBatchEndTime(parentHistory)

	if om9.IsFatalError() {
		panic(om9)
	}

	om.SetSuccessfulCompletionMessage("Finished startUp()", 79)
	return cmds, om
}


func executeJob(job *common.CmdJob, logOps *common.LogJobGroup, parent []common.OpsMsgContextInfo) common.OpsMsgDto {

	om := baseLogMsgConfigMain(parent, "executeJob")

	thisParentInfo := om.GetNewParentHistory()

	job.SetCmdJobActualStartTime(thisParentInfo)

	om2 := logOps.WriteCmdJobHeaderToLog(job, thisParentInfo)

	if om2.IsFatalError() {
		logOps.WriteOpsMsgToLog(om2, job, thisParentInfo)
		om2.PrintToConsole()
		return om2
	}

	executeJobCommand(job,logOps,om.GetNewParentHistory())

	executeJobCommand(job,logOps, om.GetNewParentHistory())

	logOps.WriteCmdJobFooterToLog(job, thisParentInfo)

	return om.SignalSuccessfulCompletion(619)
}

func executeJobCommand(job *common.CmdJob, logOps *common.LogJobGroup, parent []common.OpsMsgContextInfo) common.OpsMsgDto {

	om := baseLogMsgConfigMain(parent, "executeJobCommand")

	time.Sleep(time.Duration(5) * time.Second)

	s := fmt.Sprintf("Completed Job: %v. No Errors!", job.CmdDisplayName)
	opsMsg := common.OpsMsgDto{}

	job.CmdJobNoOfMsgs++

	opsMsg.SetTimeZone(job.IanaTimeZone)

	opsMsg.SetInfoMessage(s, int64((job.CmdJobNo * 10000000) + job.CmdJobNoOfMsgs))

	if opsMsg.IsError() {
		job.CmdJobExecutionStatus = "Errors - See Error Messages"
		job.CmdJobNoOfErrorMsgs++
	}

	om1 := logOps.WriteOpsMsgToLog(opsMsg, job, om.GetNewParentHistory())

	if om1.IsFatalError() {
		om1.LogMsgToFile(logOps.LogPathFileNameExt.FilePtr)
		panic(om1)
	}

	om2 := job.SetCmdJobActualEndTime(om.GetNewParentHistory())

	if om2.IsFatalError() {
		om2.LogMsgToFile(logOps.LogPathFileNameExt.FilePtr)
		panic(om2)
	}

	om.SetSuccessfulCompletionMessage("Completed executeJobCommand", 629)

	return om
}

// baseLogErrConfig - Used internally by LogJobGroup
// methods to set up error messages.
func baseLogMsgConfigMain(parent []common.OpsMsgContextInfo, funcName string) common.OpsMsgDto {

	opsContext := common.OpsMsgContextInfo{
									SourceFileName: srcFileNameLogOpsMain,
									ParentObjectName: "",
									FuncName: funcName,
									BaseMessageId: errBlockNoLogOpsMain,
	}


	return common.OpsMsgDto{}.InitializeAllContextInfo(parent, opsContext)
}

func runCmdJobs(lg *common.LogJobGroup, cmds common.CommandBatch, parentHistory []common.OpsMsgContextInfo) common.OpsMsgDto {

	om := baseLogMsgConfigMain(parentHistory, "runCmdJobs")

	for _, cmdJob := range cmds.CmdJobs.CmdJobArray {

		om4 := executeJob(&cmdJob, lg, parentHistory)

		if cmdJob.CmdJobIsCompleted && !om4.IsError() {
			lg.NoOfJobsCompleted++
		}


		lg.NoOfJobGroupMsgs += cmdJob.CmdJobNoOfMsgs

		if om4.IsError() && lg.KillAllJobsOnFirstError {
			om4.LogMsgToFile(lg.LogPathFileNameExt.FilePtr)
			doLogWrapUp(lg, cmds, parentHistory)
			panic(om4)
		}
	}

	om.SetSuccessfulCompletionMessage("Finished runCmdJobs", 679)

	return om
}
