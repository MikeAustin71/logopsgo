package main

import (
	"MikeAustin71/logopsgo/common"
	"fmt"
	"time"
)

const (
	srcFileNameLogOpsMain = "main.go"
	errBlockNoLogOpsMain  = int64(1000000)
	cmdPathFileName = "D:/go/work/src/MikeAustin71/logopsgo/app/cmdrXCmds.xml"
	appPathFileName = "D:/go/work/src/MikeAustin71/logopsgo/app/cmdrX.exe"
	appLogPathOnly = "D:/go/work/src/MikeAustin71/logopsgo/app/cmdrX/alog.log"
)

func main() {

	fh := common.FileHelper{}
	s, _ := fh.GetCurrentDir()

	fmt.Println("Current Directory: ", s)

	s, _ = fh.GetExecutablePathFileName()

	fmt.Println("Executable Directory: ", s)

	lg := common.LogJobGroup{}
	parent := common.OpsMsgContextInfo{
		SourceFileName: srcFileNameLogOpsMain,
		ParentObjectName: "",
		FuncName: "main",
		BaseMessageId: errBlockNoLogOpsMain,
	}

	parentHistory := []common.OpsMsgContextInfo{parent}

	cmds, om := startUp(&lg, parentHistory)

	if om.IsFatalError() {
		panic(om)
	}

	for _, cmdJob := range cmds.CmdJobs.CmdJobArray {

		om2 := executeJob(&cmdJob, &lg, parentHistory)

		if cmdJob.CmdJobIsCompleted && !om2.IsError() {
			lg.NoOfJobsCompleted++
		}

		lg.NoOfJobGroupMsgs += cmdJob.CmdJobNoOfMsgs

		if om2.IsError() && lg.KillAllJobsOnFirstError {
			doLogWrapUp(&lg, cmds, parentHistory)
			panic(om2)
		}
	}

	doLogWrapUp(&lg, cmds, parentHistory)
}

func doLogWrapUp(lg *common.LogJobGroup, cmds common.CommandBatch, parent []common.OpsMsgContextInfo) common.OpsMsgDto {

	msgCtx := common.OpsMsgContextInfo{
							SourceFileName:"main.go",
							ParentObjectName: "main()",
							FuncName:" doLogWrapUp",
							BaseMessageId: errBlockNoLogOpsMain,
						}

	om1 := common.OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	om2 := lg.WriteJobGroupFooterToLog(cmds, om1.GetNewParentHistory())


	fmt.Println("AppLogPathFileName", lg.AppLogPathFileName)

	fmt.Println("AppLogBanner1", lg.Banner1)

	return om2
}

func startUp(lg *common.LogJobGroup,parent []common.OpsMsgContextInfo) (common.CommandBatch, common.OpsMsgDto) {

	om := baseLogMsgConfigMain(parent, "startUp")
	parentHistory := om.GetNewParentHistory()
	parms := common.StartupParameters{}

	appFileParms, om2 := parms.AssembleAppPath(appPathFileName, parentHistory)

	if om2.IsFatalError() {
		panic(om2)
	}

	cmdFileParms, om3 := parms.AssembleCmdPath(cmdPathFileName, parentHistory)

	if om3.IsFatalError() {
		panic(om3)
	}

	appLogPathParms, om4 := parms.AssembleLogPath(appLogPathOnly, parentHistory)

	if om4.IsFatalError() {
		panic(om4)
	}

	cmds, om5 := common.ParseXML(cmdFileParms.AbsolutePathFileName, parentHistory)

	if om5.IsFatalError() {
		panic(om5)
	}

	cmds.FormatCmdParameters()
	om6 := cmds.SetBatchStartTime(parentHistory)
	if om6.IsFatalError() {
		panic(om6)
	}

	parms.IanaTimeZone = 	cmds.CmdJobsHdr.IanaTimeZone
	parms.KillAllJobsOnFirstError = cmds.CmdJobsHdr.KillAllJobsOnFirstError
	parms.StartTimeUTC = cmds.CmdJobsHdr.CmdBatchStartUTC
	parms.StartTime = cmds.CmdJobsHdr.CmdBatchStartTime
	dtf := common.DateTimeFormatUtility{}
	dtf.CreateAllFormatsInMemory()


	parms.AppVersion = "2.0.0"
	parms.LogMode = common.LogVERBOSE
	parms.AppLogPath = appLogPathParms.AbsolutePath
	parms.AppName = appFileParms.FileName
	parms.AppExeFileName = appFileParms.FileNameExt
	parms.AppPath = appFileParms.AbsolutePath
	parms.BaseStartDir = appFileParms.AbsolutePath
	parms.CommandFileName = cmdFileParms.FileNameExt
	parms.NoOfJobs = cmds.CmdJobsHdr.NoOfCmdJobs
	parms.Dtfmt = &dtf

	om7 := lg.New(parms, parentHistory)

	if om7.IsFatalError() {
		panic(om7)
	}

	dur := common.DurationUtility{}

	time.Sleep(dur.GetDurationFromSeconds(10))

	om8 := cmds.SetBatchEndTime(parentHistory)

	if om8.IsFatalError() {
		panic(om8)
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
		return om2
	}

	time.Sleep(time.Duration(5) * time.Second)

	s := fmt.Sprintf("Completed Job: %v. No Errors!", job.CmdDisplayName)
	om.SetInfoMessage(s, 612)

	logOps.WriteOpsMsgToLog(om, job, thisParentInfo)

	job.SetCmdJobActualEndTime(thisParentInfo)

	logOps.WriteCmdJobFooterToLog(job, thisParentInfo)

	return om.SignalSuccessfulCompletion(619)
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

