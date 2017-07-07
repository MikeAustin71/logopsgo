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
	parent := common.ErrBaseInfo{}.GetNewParentInfo(srcFileNameLogOpsMain, "main", errBlockNoLogOpsMain)
	cmds, sea := startUp(&lg, parent)

	if sea.IsErr {
		panic(sea)
	}

	for _, cmdJob := range cmds.CmdJobs.CmdJobArray {

		se2 := executeJob(&cmdJob, &lg, parent)

		if se2.IsErr && lg.KillAllJobsOnFirstError {
			panic(se2)
		}
	}

}

func startUp(lg *common.LogJobGroup,parent []common.ErrBaseInfo) (common.CommandBatch, common.SpecErr) {

	se := baseLogErrConfigMain(parent, "startUp")

	parms := common.StartupParameters{}

	appFileParms, sea := parms.AssembleAppPath(appPathFileName, se.AddBaseToParentInfo())

	if sea.IsErr {
		panic(sea)
	}

	cmdFileParms, sea := parms.AssembleCmdPath(cmdPathFileName, se.AddBaseToParentInfo())

	if sea.IsErr {
		panic(sea)
	}

	appLogPathParms, sea := parms.AssembleLogPath(appLogPathOnly, se.AddBaseToParentInfo())

	cmds, sea := common.ParseXML(cmdFileParms.AbsolutePathFileName, parent)

	if sea.IsErr {
		panic(sea)
	}

	cmds.FormatCmdParameters()
	sea = cmds.SetBatchStartTime(se.AddBaseToParentInfo())
	if sea.IsErr {
		panic(sea)
	}

	parms.IanaTimeZone = cmds.CmdJobsHdr.IanaTimeZone
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

	sea = lg.New(parms, parent)

	if sea.IsErr {
		panic(sea)
	}

	dur := common.DurationUtility{}

	time.Sleep(dur.GetDurationFromSeconds(10))

	sea = cmds.SetBatchEndTime(se.AddBaseToParentInfo())

	if sea.IsErr {
		panic(sea)
	}

	lg.WriteJobGroupFooterToLog(cmds, parent)

	fmt.Println("AppLogPathFileName", lg.AppLogPathFileName)

	fmt.Println("AppLogBanner1", lg.Banner1)

	return cmds, se.SignalNoErrors()
}


func executeJob(job *common.CmdJob, logOps *common.LogJobGroup, parent []common.ErrBaseInfo) common.SpecErr {

	se := baseLogErrConfigMain(parent, "executeJob")
	opsMsg := common.OpsMsgDto{}
	thisParentInfo := se.AddBaseToParentInfo()
	job.SetCmdJobActualStartTime(thisParentInfo)

	sea := logOps.WriteCmdJobHeaderToLog(job, thisParentInfo)

	if sea.IsErr {
		om := opsMsg.NewSpecErr(sea)
		logOps.WriteOpsMsgToLog(om, job, thisParentInfo)
		return sea
	}

	time.Sleep(time.Duration(5) * time.Second)

	s := fmt.Sprintf("Completed Job: %v. No Errors!", job.CmdDisplayName)
	om := opsMsg.NewInfoMsg(s)

	logOps.WriteOpsMsgToLog(om, job, thisParentInfo)

	job.SetCmdJobActualEndTime(thisParentInfo)


	logOps.WriteCmdJobFooterToLog(job, thisParentInfo)

	return se.SignalNoErrors()
}


// baseLogErrConfig - Used internally by LogJobGroup
// methods to set up error messages.
func baseLogErrConfigMain(parent []common.ErrBaseInfo, funcName string) common.SpecErr {

	bi := common.ErrBaseInfo{}.New(srcFileNameLogOpsMain, funcName, errBlockNoLogOpsMain)

	return common.SpecErr{}.InitializeBaseInfo(parent, bi)
}

