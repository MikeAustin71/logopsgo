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
	appVersion						= "2.0.0"
)



type AppStartUpInfo struct {
	AppVersion						string
	AppPathFileNameExt    common.FileMgr		// Path File Name Extension for Application Executable
	AppPath               common.DirMgr			// Path to the Application Executable
	CmdPathFileNameExt    common.FileMgr		// Path File Name Extension for XML Command File
	CmdPath               common.DirMgr			// Path to the XML Command File
	LogPathFileNameExt		common.FileMgr		// Path, File Name and Extension for Log File
	LogPath               common.DirMgr			// Path to the Log Directory
	AppErrPathFileNameExt common.FileMgr		// Path File Name Extension for the Error File
	CurrentDirPath        common.DirMgr		  // Path for the current Directory
	AppStartTimeTzu       common.TimeZoneUtility 	// Time Zone Utility Containing App StartUp time
																								// for UTC and 'Local' TimeZone

}


func main() {


	parent := common.OpsMsgContextInfo{
		SourceFileName: srcFileNameLogOpsMain,
		ParentObjectName: "",
		FuncName: "main",
		BaseMessageId: errBlockNoLogOpsMain,
	}

	om := common.OpsMsgDto{}.InitializeWithMessageContext(parent)

	isRunningInDebugMode := false

	var appDirs AppStartUpInfo
	var om2 common.OpsMsgDto

	if isRunningInDebugMode {

		appDirs, om2 = assembleDebugAppPathFiles(om.GetNewParentHistory())

	} else {

		appDirs, om2 = assembleAppPathFiles(om.GetNewParentHistory())

	}

	if om2.IsFatalError() {
		om2.PrintToConsole()
		return
	}


	su := common.StringUtility{}
	lBanner1 := len(appBanner1)
	fmt.Println("\n\n"+ appBanner1)
	strx, _ := su.StrCenterInStr("CmdrX.exe", lBanner1 - 2)
	fmt.Println("=" + strx + "=")
	fmt.Println(appBanner1)
	fmt.Println("Current Directory: ", appDirs.CurrentDirPath.AbsolutePath)
	fmt.Println("Executable Directory: ", appDirs.AppPathFileNameExt.DMgr.AbsolutePath)

	parentHistory := []common.OpsMsgContextInfo{parent}

	lg := common.LogJobGroup{}

	cmds, om3 := startUp(&lg, appDirs, parentHistory)

	if om3.IsFatalError() {
		panic(om3)
	}

	runCmdJobs(&lg, cmds, parentHistory)

	doLogWrapUp(&lg, cmds, parentHistory)
}

func assembleDebugAppPathFiles(parent []common.OpsMsgContextInfo) (AppStartUpInfo, common.OpsMsgDto) {

	const (
		debugAppPathFileName  = "D:/go/work/src/MikeAustin71/logopsgo/app/cmdrX.exe"
		debugCmdPathFileName  = "D:/go/work/src/MikeAustin71/logopsgo/app/cmdrXCmds.xml"
		debugAppLogPathOnly   = "D:/go/work/src/MikeAustin71/logopsgo/app/cmdrXLog"
	)

	var err error

	appDirs := AppStartUpInfo{}

	msgCtx := common.OpsMsgContextInfo{
		SourceFileName:"main.go",
		ParentObjectName: "main()",
		FuncName:" assembleDebugAppPathFiles",
		BaseMessageId: errBlockNoLogOpsMain,
	}

	om := common.OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	// Set Application Version
	appDirs.AppVersion = appVersion

	// Compute App Start Time
	utcTime := time.Now().UTC()
	appDirs.AppStartTimeTzu, err = common.TimeZoneUtility{}.ConvertTz(utcTime, "Local")

	if err != nil {
		s:= fmt.Sprintf("TimeZoneUtility{}.ConvertTz Error - Failed to convert UTC to local Time Zone. Start UTC: %v. Iana Time Zone: %v", utcTime, "Local")
		om.SetFatalError(s, err, 2701)
		return AppStartUpInfo{},  om
	}

	// App Path Exe File
	appDirs.AppPathFileNameExt, err  = common.FileMgr{}.New(debugAppPathFileName)

	if err != nil {
		s:= fmt.Sprintf("appDirs.AppPathFileNameExt Error returned from FileMgr{}.New(debugAppPathFileName) debugAppPathFileName='%v'\n", debugAppPathFileName)
		om.SetFatalError(s, err, 2703)
		return AppStartUpInfo{},  om
	}

	// App Path
	appDirs.AppPath = appDirs.AppPathFileNameExt.DMgr.CopyOut()

	// Cmd Path Cmd File Name Ext
	appDirs.CmdPathFileNameExt, err = common.FileMgr{}.New(debugCmdPathFileName)

	if err != nil {
		s:= fmt.Sprintf("appDirs.AppPathFileNameExt Error returned from FileMgr{}.New(debugCmdPathFileName) debugCmdPathFileName='%v'\n", debugCmdPathFileName)
		om.SetFatalError(s, err, 2705)
		return AppStartUpInfo{},  om
	}

	// Cmd Path
	appDirs.CmdPath = appDirs.CmdPathFileNameExt.DMgr.CopyOut()

	// App LogPath
	appDirs.LogPath, err = common.DirMgr{}.New(debugAppLogPathOnly)

	if err != nil {
		s:= fmt.Sprintf("appDirs.AppPathFileNameExt Error returned from DirMgr{}.New(debugAppLogPathOnly) debugAppLogPathOnly='%v'\n", debugAppLogPathOnly)
		om.SetFatalError(s, err, 2707)
		return AppStartUpInfo{},  om
	}

	// App Log File Name Ext
	dt := common.DateTimeUtility{}
	logFileNameExt := appDirs.AppPathFileNameExt.FileName + "_" + dt.GetDateTimeStr(appDirs.AppStartTimeTzu.TimeOut) + ".log"
	appDirs.LogPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(appDirs.LogPath, logFileNameExt)

	if err != nil {
		s:= fmt.Sprintf("appDirs.AppPathFileNameExt Error returned from FileMgr{}.NewFromDirMgrFileNameExt(appDirs.LogPath, logFileNameExt) appDirs.LogPath='%v' logFileNameExt='%v' \n", appDirs.LogPath, logFileNameExt)
		om.SetFatalError(s, err, 2709)
		return AppStartUpInfo{},  om
	}

	// Error File
	appErrFileName := appDirs.AppPathFileNameExt.FileName + "_Errors_" + dt.GetDateTimeStr(appDirs.AppStartTimeTzu.TimeOut) + ".txt"

	// AppErrPath File Name Ext
	appDirs.AppErrPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(appDirs.AppPathFileNameExt.DMgr, appErrFileName)

	if err != nil {
		s:= fmt.Sprintf("appDirs.AppPathFileNameExt Error returned from FileMgr{}.NewFromDirMgrFileNameExt(appDirs.AppPathFileNameExt.DMgr, appErrFileName) appDirs.AppPathFileNameExt.DMgr.Path='%v' logFileNameExt='%v'\n", appDirs.AppPathFileNameExt.DMgr.Path, appErrFileName)
		om.SetFatalError(s, err, 2711)
		return AppStartUpInfo{},  om
	}

	// Current Directory
	fh := common.FileHelper{}

	currDirPath, err := fh.GetCurrentDir()

	appDirs.CurrentDirPath, err = common.DirMgr{}.New(currDirPath)

	if err != nil {
		s:= fmt.Sprintf("appDirs.AppPathFileNameExt Error returned from FileMgr{}.NewFromDirMgrFileNameExt(appDirs.AppPathFileNameExt.DMgr, appErrFileName) appDirs.AppPathFileNameExt.DMgr.Path='%v' logFileNameExt='%v'\n", appDirs.AppPathFileNameExt.DMgr.Path, appErrFileName)
		om.SetFatalError(s, err, 2715)
		return AppStartUpInfo{},  om
	}

	return appDirs, om.SignalNoErrors(2799)
}

func assembleAppPathFiles(parent []common.OpsMsgContextInfo) (AppStartUpInfo, common.OpsMsgDto) {

	const (
		cmdPathFileName  = "cmdrXCmds.xml"
		appLogPathOnly   = "cmdrXLog"
		appName 				 = "cmdrX"
	)

	fh := common.FileHelper{}

	appDirs := AppStartUpInfo{}

	var err error

	msgCtx := common.OpsMsgContextInfo{
		SourceFileName:"main.go",
		ParentObjectName: "main()",
		FuncName:" assembleAppPathFiles",
		BaseMessageId: errBlockNoLogOpsMain,
	}

	om := common.OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	// Set Application Version
	appDirs.AppVersion = appVersion

	// Compute App Start Time
	utcTime := time.Now().UTC()
	appDirs.AppStartTimeTzu, err = common.TimeZoneUtility{}.ConvertTz(utcTime, "Local")

	if err != nil {
		s:= fmt.Sprintf("TimeZoneUtility{}.ConvertTz Error - Failed to convert UTC to local Time Zone. Start UTC: %v. Iana Time Zone: %v", utcTime, "Local")
		om.SetFatalError(s, err, 3701)
		return AppStartUpInfo{},  om
	}

	// App Path File Name Ext
	appExePathFileNameStr, err := fh.GetExecutablePathFileName()

	if err!=nil {
		om.SetFatalError("fh.GetExecutablePathFileName() FAILED!\n", err, 3703)
		return AppStartUpInfo{}, om
	}

	appDirs.AppPathFileNameExt, err = common.FileMgr{}.New(appExePathFileNameStr)

	if err!=nil {
		s:= fmt.Sprintf("FileMgr{}.New(appExePathFileNameStr) appExePathFileNameStr='%v' FAILED!\n", appExePathFileNameStr)
		om.SetFatalError(s, err, 3705)
		return AppStartUpInfo{}, om
	}

	if !appDirs.AppPathFileNameExt.AbsolutePathFileNameDoesExist {
		om.SetFatalError("Error: App File Name Extension Does NOT Exist!\n", fmt.Errorf("File Does NOT Exist! AppPathFileName='%v' ", appDirs.AppPathFileNameExt.AbsolutePathFileName), 3707)
		return AppStartUpInfo{}, om
	}

	// App Path
	appDirs.AppPath = appDirs.AppErrPathFileNameExt.DMgr.CopyOut()

	// Command Path
	appDirs.CmdPath = appDirs.AppPath.CopyOut()

	// Command Path File Name Ext
	appDirs.CmdPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(appDirs.CmdPath,cmdPathFileName)

	if err!=nil {
		s:= fmt.Sprintf("Error returned by FileMgr{}.NewFromDirMgrFileNameExt(appDirs.CmdPath,cmdPathFileName ) appDirs.CmdPath.Path='%v' cmdPathFileName='%v'\n", appDirs.CmdPath.Path, cmdPathFileName)
		om.SetFatalError(s, err, 3709)
		return AppStartUpInfo{}, om
	}

	// Log Path
	logPath := appDirs.AppPath.GetPathWithSeparator() + appLogPathOnly
	appDirs.LogPath, err = common.DirMgr{}.New(logPath)

	if err!=nil {
		s:= fmt.Sprintf("Error returned by DirMgr{}.New(logPath) logPath='%v'\n", logPath)
		om.SetFatalError(s, err, 3711)
		return AppStartUpInfo{}, om
	}

	// Log Path File Name Ext
	dt := common.DateTimeUtility{}
	dateTimeStamp := dt.GetDateTimeStr(appDirs.AppStartTimeTzu.TimeOut)
	logFileNameExt := appName + "_" + dateTimeStamp + ".log"
	appDirs.LogPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(appDirs.LogPath, logFileNameExt)

	if err != nil {
		s:= fmt.Sprintf("appDirs.LogPathFileNameExt Error returned from FileMgr{}.NewFromDirMgrFileNameExt(appDirs.LogPath, logFileNameExt) appDirs.LogPath='%v' logFileNameExt='%v' \n", appDirs.LogPath.Path, logFileNameExt)
		om.SetFatalError(s, err, 3715)
		return AppStartUpInfo{},  om
	}

	// App Error Path File Name Ext
	appErrFileNameExt := appName + "_Errors_" + dateTimeStamp + ".txt"

	appDirs.AppErrPathFileNameExt, err = common.FileMgr{}.NewFromDirMgrFileNameExt(appDirs.LogPath,appErrFileNameExt)

	if err != nil {
		s:= fmt.Sprintf("appDirs.AppErrPathFileNameExt Error returned from FileMgr{}.NewFromDirMgrFileNameExt(appDirs.LogPath,appErrFileNameExt) appDirs.LogPath.Path='%v' appErrFileNameExt='%v' \n", appDirs.LogPath, appErrFileNameExt)
		om.SetFatalError(s, err, 3717)
		return AppStartUpInfo{},  om
	}

	// Set Current Directory Path
	currDirPath, err := fh.GetCurrentDir()

	if err!=nil {
		om.SetFatalError("fh.GetCurrentDir() FAILED!\n", err, 3719)
		return AppStartUpInfo{}, om
	}

	appDirs.CurrentDirPath, err = common.DirMgr{}.New(currDirPath)

	if err!=nil {
		s := fmt.Sprintf("Error returned by DirMgr{}.New(currDirPath). currDirPath='%v'\n",currDirPath)
		om.SetFatalError(s, err,3721)
		return AppStartUpInfo{}, om
	}


	return appDirs, om.SignalNoErrors(3799)
}

func doLogWrapUp(lg *common.LogJobGroup, cmds common.CommandBatch, parent []common.OpsMsgContextInfo) common.OpsMsgDto {

	msgCtx := common.OpsMsgContextInfo{
							SourceFileName:"main.go",
							ParentObjectName: "main()",
							FuncName:" doLogWrapUp",
							BaseMessageId: errBlockNoLogOpsMain,
						}

	om1 := common.OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	// Closes lg.LogFilePtr
	om2 := lg.WriteJobGroupFooterToLog(cmds, om1.GetNewParentHistory())

	if om2.IsFatalError(){
		om2.PrintToConsole()
		return om2
	}

	fmt.Println(appBanner1)
	fmt.Println("See Log File:")
	fmt.Println(lg.AppLogPathFileName)

	fmt.Println(appBanner1)

	return om2
}

func startUp(lg *common.LogJobGroup,appStartDirs AppStartUpInfo,
		parent []common.OpsMsgContextInfo) (common.CommandBatch, common.OpsMsgDto) {

	om := baseLogMsgConfigMain(parent, "startUp")
	parentHistory := om.GetNewParentHistory()
	parms := common.StartupParameters{}

	cmds, om5 := common.ParseXML(appStartDirs.CmdPathFileNameExt.AbsolutePathFileName, parentHistory)

	if om5.IsFatalError() {
		panic(om5)
	}

	om6 := cmds.FormatCmdParameters(parentHistory)

	if om6.IsFatalError() {
		panic(om6)
	}


	om7 := cmds.SetBatchStartTime(appStartDirs.AppStartTimeTzu.TimeUTC, parentHistory)
	if om7.IsFatalError() {
		panic(om7)
	}

	parms.IanaTimeZone = 	cmds.CmdJobsHdr.IanaTimeZone
	parms.KillAllJobsOnFirstError = cmds.CmdJobsHdr.KillAllJobsOnFirstError
	parms.StartTimeUTC = cmds.CmdJobsHdr.CmdBatchStartUTC
	parms.StartTime = cmds.CmdJobsHdr.CmdBatchStartTime
	dtf := common.DateTimeFormatUtility{}
	dtf.CreateAllFormatsInMemory()


	parms.AppVersion = appStartDirs.AppVersion
	parms.LogMode = common.LogVERBOSE
	parms.AppLogPath = appStartDirs.LogPath.AbsolutePath
	parms.AppName = appStartDirs.AppPathFileNameExt.FileName
	parms.AppExeFileName = appStartDirs.AppPathFileNameExt.FileNameExt
	parms.AppPath = appStartDirs.AppPath.AbsolutePath
	parms.BaseStartDir = appStartDirs.AppPath.AbsolutePath
	parms.CommandFileName = appStartDirs.CmdPathFileNameExt.FileNameExt
	parms.NoOfJobs = cmds.CmdJobsHdr.NoOfCmdJobs
	parms.Dtfmt = &dtf

	om8 := lg.New(parms, parentHistory)

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
		om1.LogMsgToFile(logOps.LogFilePtr)
		panic(om1)
	}

	om2 := job.SetCmdJobActualEndTime(om.GetNewParentHistory())

	if om2.IsFatalError() {
		om2.LogMsgToFile(logOps.LogFilePtr)
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
			om4.LogMsgToFile(lg.LogFilePtr)
			doLogWrapUp(lg, cmds, parentHistory)
			panic(om4)
		}
	}

	om.SetSuccessfulCompletionMessage("Finished runCmdJobs", 679)

	return om
}
