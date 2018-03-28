package common

import (
	dt "MikeAustin71/datetimeopsgo/datetime"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/*  'loggingmode.go' is located in source code
		repository:

		https://github.com/MikeAustin71/logopsgo.git

 */


const (
	srcFileNameXMLCmdJobsData = "xmlcmdjobsdata.go"
	errBlockNoXMLCmdJobsData  = int64(9230610000)
	logOpsTimeFmt = "2006-01-02 Mon 15:04:05.000000000 -0700 MST"
)

// CommandBatch - Xml Root and Parent Element
type CommandBatch struct {
	CmdJobsHdr CommandJobsHdr  `xml:"CommandJobsHeader"`
	CmdJobs    CommandJobArray `xml:"CommandJobs"`
}

func (cmdBatch *CommandBatch) FormatCmdParameters(parentHistory []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
		SourceFileName: srcFileNameXMLCmdJobsData,
		ParentObjectName: "CommandBatch",
		FuncName: "FormatCmdParameters",
		BaseMessageId: errBlockNoXMLCmdJobsData,
	}

	om := OpsMsgDto{}.InitializeAllContextInfo(parentHistory, msgCtx)

	om2 := cmdBatch.assembleTimeFormats(om.GetNewParentHistory())

	if om2.IsError() {
		return om2
	}

	om3 := cmdBatch.assembleCmdElements(om.GetNewParentHistory())

	if om3.IsError() {
		return om3
	}

	om.SetSuccessfulCompletionMessage("Completed FormatCmdParameters", 109)

	return om
}

func (cmdBatch *CommandBatch) assembleTimeFormats(parentHistory []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
		SourceFileName: srcFileNameXMLCmdJobsData,
		ParentObjectName: "CommandBatch",
		FuncName: "assembleTimeFormats",
		BaseMessageId: errBlockNoXMLCmdJobsData,
	}

	om := OpsMsgDto{}.InitializeAllContextInfo(parentHistory, msgCtx)

	if cmdBatch.CmdJobsHdr.IanaTimeZone == "" {
		cmdBatch.CmdJobsHdr.IanaTimeZone = "Local"
	}

	tzu := dt.TimeZoneDto{}
	isValidTz, _, _ := tzu.IsValidTimeZone(cmdBatch.CmdJobsHdr.IanaTimeZone)

	if !isValidTz {
		cmdBatch.CmdJobsHdr.IanaTimeZone = "Local"
	}

	cmdBatch.CmdJobsHdr.StdTimeFormat = logOpsTimeFmt

	om.SetSuccessfulCompletionMessage("Completed assembleTimeFormats", 209)

	return om
}

// assembleCmdElements - Assembles
// Command and Command Arguments.
// These command elements are then
// stored in struct CombinedExeCommand
func (cmdBatch *CommandBatch) assembleCmdElements(parentHistory []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
		SourceFileName: srcFileNameXMLCmdJobsData,
		ParentObjectName: "CommandBatch",
		FuncName: "assembleTimeFormats",
		BaseMessageId: errBlockNoXMLCmdJobsData,
	}

	om := OpsMsgDto{}.InitializeAllContextInfo(parentHistory, msgCtx)

	cmdBatch.CmdJobsHdr.NoOfCmdJobs = len(cmdBatch.CmdJobs.CmdJobArray)

	lJobs := len(cmdBatch.CmdJobs.CmdJobArray)

	for i := 0; i < lJobs; i++ {
		job := &cmdBatch.CmdJobs.CmdJobArray[i]

		// Set Job No.
		job.CmdJobNo = i + 1

		// sync time zones
		job.IanaTimeZone = cmdBatch.CmdJobsHdr.IanaTimeZone

		// sync time formats
		job.CmdJobTimeFormat = cmdBatch.CmdJobsHdr.StdTimeFormat

		omx := cmdBatch.assembleCmdArgs(job, om.GetNewParentHistory())

		if omx.IsError() {
			return omx
		}

		if job.ExeCmdInDir == "" {
			job.ExeCmdInDir = "Current CmdrX.exe Directory"
		}

		omy := cmdBatch.assembleInputArgs(job, om.GetNewParentHistory())

		if omy.IsError() {
			return omy
		}

		job.ExeCommand = strings.TrimLeft(strings.TrimRight(job.ExeCommand, " "), " ")

		if len(job.CombinedArguments) > 0 {
			job.CombinedExeCommand =
				job.ExeCommand + " " + job.CombinedArguments
		} else {
			job.CombinedExeCommand = job.ExeCommand
		}

	}

	om.SetSuccessfulCompletionMessage("Finished assembleCmdElements", 309)

	return om
}

// assembleCmdArgs - Processes Command Arguments located in job.CmdArguments.CmdArgs
func (cmdBatch *CommandBatch) assembleCmdArgs(job *CmdJob, parentHistory []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
		SourceFileName: srcFileNameXMLCmdJobsData,
		ParentObjectName: "CommandBatch",
		FuncName: "assembleCmdArgs",
		BaseMessageId: errBlockNoXMLCmdJobsData,
	}

	om := OpsMsgDto{}.InitializeAllContextInfo(parentHistory, msgCtx)

	job.CombinedArguments = ""

	lCmdArgs := len(job.CmdArguments.CmdArgs)

	if lCmdArgs == 0 {
		om.SetSuccessfulCompletionMessage("Finished assembleCmdArgs", 408)
		return om
	}

	cmdArgs := ""
	doCmdArgsExist := false

	for i := 0; i < lCmdArgs; i++ {
		job.CmdArguments.CmdArgs[i] = strings.TrimRight(strings.TrimLeft(job.CmdArguments.CmdArgs[i], " "), " ")

		if job.CmdArguments.CmdArgs[i] != "" {

			if i==0 {
				cmdArgs += job.CmdArguments.CmdArgs[i]
			} else {
				cmdArgs += " " + job.CmdArguments.CmdArgs[i]
			}
			doCmdArgsExist = true
		}
	}

	if doCmdArgsExist {
		job.CombinedArguments = cmdArgs
	} else {
		job.CmdArguments.CmdArgs = make([]string, 0, 0)
	}

	om.SetSuccessfulCompletionMessage("Finished assembleCmdArgs", 409)
	return om
}

// assembleInputArgs -
func (cmdBatch *CommandBatch) assembleInputArgs(job *CmdJob, parentHistory []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
		SourceFileName: srcFileNameXMLCmdJobsData,
		ParentObjectName: "CommandBatch",
		FuncName: "assembleInputArgs",
		BaseMessageId: errBlockNoXMLCmdJobsData,
	}

	om := OpsMsgDto{}.InitializeAllContextInfo(parentHistory, msgCtx)

	job.CombinedInputArguments = ""

	lCmdInputs := len(job.CmdInputs.InputArgs)

	if lCmdInputs == 0 {
		om.SetSuccessfulCompletionMessage("Finished assembleInputArgs", 508)
		return om
	}

	cmdInputs := ""
	doCmdInputsExist := false

	for i:=0; i < lCmdInputs; i++ {
		job.CmdInputs.InputArgs[i] = strings.TrimRight(strings.TrimLeft(job.CmdInputs.InputArgs[i], " "), " ")

		if job.CmdInputs.InputArgs[i] != "" {

			if i==0 {
				cmdInputs += job.CmdInputs.InputArgs[i]
			} else {
				cmdInputs += " " + job.CmdInputs.InputArgs[i]
			}

			doCmdInputsExist = true
		}

	}

	if doCmdInputsExist {
		job.CombinedInputArguments = cmdInputs
	} else {
		job.CmdInputs.InputArgs = make([]string, 0, 0)
	}

	om.SetSuccessfulCompletionMessage("Finished assembleInputArgs", 509)
	return om
}

// SetBatchStartTime - Sets the time at which jobs in this
// Command Batch began processing.
func (cmdBatch *CommandBatch) SetBatchStartTime(appStartTime dt.TimeZoneDto, parent []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
							SourceFileName: srcFileNameXMLCmdJobsData,
							ParentObjectName: "CommandBatch",
							FuncName: "SetBatchStartTime",
							BaseMessageId: errBlockNoXMLCmdJobsData,
						}

	om := OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	var err error

	tzxu := dt.TimeZoneDto{}

	isValidTz, _, _ :=tzxu.IsValidTimeZone(cmdBatch.CmdJobsHdr.IanaTimeZone)

	if !isValidTz {
		cmdBatch.CmdJobsHdr.IanaTimeZone = "Local"
	}

	cmdBatch.CmdJobsHdr.CmdBatchStartTime, err  =
		dt.TimeZoneDto{}.New(appStartTime.TimeLocal.DateTime,
													cmdBatch.CmdJobsHdr.IanaTimeZone,
														dt.FmtDateTimeYMDAbbrvDowNano)

	if err != nil {
		s:= fmt.Sprintf("Error returned by dt.TimeZoneDto{}.New( " +
			"appStartTime.TimeLocal.DateTime, " +
			"cmdBatch.CmdJobsHdr.IanaTimeZone, fmt) " +
			"appStartTime.TimeLocal.DateTime: %v. " +
			"cmdBatch.CmdJobsHdr.IanaTimeZone : %v \n",
				appStartTime.TimeLocal.DateTime.Format(dt.FmtDateTimeYrMDayFmtStr),
					cmdBatch.CmdJobsHdr.IanaTimeZone)

		om.SetFatalError(s, err, 601)
		return  om
	}

	om.SetSuccessfulCompletionMessage("Finished SetBatchStartTime", 609)

	return om
}

// SetBatchEndTime - Sets the time at which all jobs
// in this Command Batch ended and, in addition,
// computes the elapsed time to complete all jobs in
// this Command Batch.
func (cmdBatch *CommandBatch) SetBatchEndTime(parent []OpsMsgContextInfo) OpsMsgDto {


	msgCtx := OpsMsgContextInfo{
							SourceFileName: srcFileNameXMLCmdJobsData,
							ParentObjectName: "CommandBatch",
							FuncName: "SetBatchEndTime",
							BaseMessageId: errBlockNoXMLCmdJobsData,
	}

	om:= OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	var err error



	cmdBatch.CmdJobsHdr.CmdBatchEndTime, err = dt.TimeZoneDto{}.New(
																							time.Now().Local(),
																							cmdBatch.CmdJobsHdr.IanaTimeZone,
																							dt.FmtDateTimeYMDAbbrvDowNano )

	if err != nil {
		s:= fmt.Sprintf("Error returned by dt.TimeZoneDto{}.New(time.Now().Local(), " +
			"cmdBatch.CmdJobsHdr.IanaTimeZone, dt.FmtDateTimeYMDAbbrvDowNano ) " +
			"time.Now().Local(): %v. Iana Time Zone: %v \n",
			time.Now().Local().Format(dt.FmtDateTimeYMDAbbrvDowNano),
			cmdBatch.CmdJobsHdr.IanaTimeZone)

		om.SetFatalError(s, err, 701)
		return  om
	}



	cmdBatch.CmdJobsHdr.CmdBatchDuration, err =
		dt.DurationTriad{}.NewStartEndTimesCalcTz(
			cmdBatch.CmdJobsHdr.CmdBatchStartTime.TimeLocal.DateTime,
			cmdBatch.CmdJobsHdr.CmdBatchEndTime.TimeLocal.DateTime,
			dt.TDurCalcTypeSTDYEARMTH,
			cmdBatch.CmdJobsHdr.IanaTimeZone,
			dt.FmtDateTimeYMDAbbrvDowNano)

	if err != nil {
		s:= fmt.Sprintf("Error returned by dt.DurationTriad{}.NewStartEndTimesCalcTz() " +
			"cmdBatch.CmdJobsHdr.CmdBatchStartTime.TimeLocal.DateTime: %v. " + "" +
			"cmdBatch.CmdJobsHdr.CmdBatchEndTime.TimeLocal.DateTime: %v " +
			"cmdBatch.CmdJobsHdr.IanaTimeZone='%v' \n",
			cmdBatch.CmdJobsHdr.CmdBatchStartTime.TimeLocal.DateTime.Format(dt.FmtDateTimeYMDAbbrvDowNano),
			cmdBatch.CmdJobsHdr.CmdBatchEndTime.TimeLocal.DateTime.Format(cmdBatch.CmdJobsHdr.IanaTimeZone),
			cmdBatch.CmdJobsHdr.IanaTimeZone)

		om.SetFatalError(s, err, 702)
		return  om
	}

	om.SetSuccessfulCompletionMessage("Finished SetBatchEndTime", 709)
	return om

}

// CommandJobsHdr - Holds base info related to
// command jobs
type CommandJobsHdr struct {
	CmdFileVersion          string `xml:"CmdFileVersion"`
	LogFileRetentionInDays  int    `xml:"LogFileRetentionInDays"`
	KillAllJobsOnFirstError bool   `xml:"KillAllJobsOnFirstError"`
	IanaTimeZone            string `xml:"IanaTimeZone"`
	NoOfCmdJobs             int
	StdTimeFormat           string
	CmdBatchStartTime       dt.TimeZoneDto
	CmdBatchEndTime         dt.TimeZoneDto
	CmdBatchDuration        dt.DurationTriad
}

// CommandJobArray - Holds individual
// CommandJob structs
type CommandJobArray struct {
	CmdJobArray []CmdJob `xml:"CommandJob"`
}

// CmdJob - Command Job information
type CmdJob struct {
	CmdDisplayName           string `xml:"CommandDisplayName"`
	CmdDescription           string `xml:"CommandDescription"`
	CmdJobNo                 int
	CmdType                  string `xml:"CommandType"`
	ExeCmdInDir              string `xml:"ExecuteCmdInDir"`
	DelayCmdStartSeconds     string `xml:"DelayCmdStartSeconds"`
	DelayCmdStartDuration    dt.DurationTriad
	DelayStartCmdDateTime    string `xml:"DelayStartCmdDateTime"`
	CommandTimeOutInSeconds  float64 `xml:"CmdTimeOutInSeconds"`
	CommandTimeOutDuration   dt.DurationTriad
	CombinedExeCommand       string
	ExeCommand               string `xml:"ExeCommand"`
	CombinedArguments        string
	CmdArguments             CommandArgumentsArray `xml:"CmdArguments"`
	CmdInputs                CommandInputsArray	`xml:"CmdInputs"`
	CombinedInputArguments   string
	IanaTimeZone             string
	CmdJobTimeFormat         string
	CmdJobStartTimeValue     dt.TimeZoneDto
	CmdJobEndTimeValue       dt.TimeZoneDto
	CmdJobDuration           dt.DurationTriad
	CmdJobNoOfMsgs					 int
	CmdJobNoOfErrorMsgs			 int
	CmdJobIsCompleted				 bool
	CmdJobExecutionStatus		 string
}

// SetDelayCmdStartTime - Sets the date time at which the command
// job can begin execution. User has the option to input a delay
// factor expressed in seconds or to input a specific time at which
// the command will begin execution. If no start time is signaled
// by the user, the command will begin execution immediately.
func (job *CmdJob) SetDelayCmdStartTime(dtf *dt.FormatDateTimeUtility, parentHistory []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo{
							SourceFileName: srcFileNameXMLCmdJobsData,
							ParentObjectName: "CmdJob",
							FuncName: "SetDelayCmdStartTime",
							BaseMessageId: errBlockNoXMLCmdJobsData,
						}

	var err error


	om:= OpsMsgDto{}.InitializeAllContextInfo(parentHistory, msgCtx)

	// Set a Default Delay as Current Time - 3-years. In other words
	tDto := dt.TimeDto{Years: 3}

	job.DelayCmdStartDuration, err = dt.DurationTriad{}.NewEndTimeMinusTimeDtoCalcTz(
												time.Now().Local(),
												tDto,
												dt.TDurCalcTypeSTDYEARMTH,
												job.IanaTimeZone,
												dt.FmtDateTimeYMDAbbrvDowNano)


	if err != nil {
		s:= fmt.Sprintf("Error job.DelayCmdStartDuration = " +
			"dt.DurationTriad{}.NewEndTimeMinusTimeDtoCalcTz() " +
			" Iana Time Zone='%v' \n",
			job.IanaTimeZone)

		om.SetFatalError(s, err, 1201)

		return om
	}

	if job.DelayCmdStartSeconds != "" {

		dSecs, err := strconv.ParseInt(job.DelayCmdStartSeconds, 10, 64)

		if err == nil {
			// Delay seconds is populated. Re-Calculate Delay Command Execution
			// times.
			dur := job.DelayCmdStartDuration.GetDurationFromSeconds(dSecs)

			job.DelayCmdStartDuration, err =
				dt.DurationTriad{}.NewStartTimeDurationCalcTz(
													time.Now().Local(),
													dur,
													dt.TDurCalcTypeSTDYEARMTH,
													job.IanaTimeZone,
													dt.FmtDateTimeYMDAbbrvDowNano)

			if err != nil {
				s := fmt.Sprintf("Error job.DelayCmdStartDuration = "+
					"dt.DurationTriad{}.NewEndTimeMinusTimeDtoCalcTz() "+
					" Iana Time Zone='%v' \n",
					job.IanaTimeZone)

				om.SetFatalError(s, err, 1203)
			}

		}
	}

	if job.DelayStartCmdDateTime != "" {
		// The Delay Start Execution Date Time is populated! Set Execution for
		// a specific datetime.
		tStart, err := dtf.ParseDateTimeString(job.DelayStartCmdDateTime, "")

		if err != nil {
			s:= fmt.Sprintf("ParseDateTimeString Error. job.DelayStartCmdDateTime: %v", job.DelayStartCmdDateTime)
			om.SetFatalError(s, err, 1205)
			return  om
		}

		startDateTz, err := dt.DateTzDto{}.NewTz(
														tStart,
														job.IanaTimeZone,
														dt.FmtDateTimeYMDAbbrvDowNano)

		if err != nil {
			s:= fmt.Sprintf("Error returned from tStart = dt.DateTzDto{}.NewTz(). " +
				"tStart='%v' Iana Time Zone='%v'",
				tStart.Format(dt.FmtDateTimeYMDAbbrvDowNano),
					job.IanaTimeZone)

			om.SetFatalError(s, err, 1207)
			return  om
		}

		nowDateTz, err := dt.DateTzDto{}.NewTz(
														time.Now().Local(),
														job.IanaTimeZone,
														dt.FmtDateTimeYMDAbbrvDowNano)

		if err != nil {
			s:= fmt.Sprintf("Error returned from nowDateTz = dt.DateTzDto{}.NewTz(). " +
				"Now ='%v' Iana Time Zone='%v'\n",
				time.Now().Local().Format(dt.FmtDateTimeYMDAbbrvDowNano),
				job.IanaTimeZone)

			om.SetFatalError(s, err, 1209)
			return  om
		}

		if nowDateTz.DateTime.Before(startDateTz.DateTime) {
			job.DelayCmdStartDuration, err =
				dt.DurationTriad{}.NewStartEndTimesCalcTz(
					nowDateTz.DateTime,
					startDateTz.DateTime,
					dt.TDurCalcTypeSTDYEARMTH,
					job.IanaTimeZone,
					dt.FmtDateTimeYMDAbbrvDowNano)

			if err != nil {
				s := fmt.Sprintf("Error Target Start Execution Date job.DelayCmdStartDuration = " +
					"dt.DurationTriad{}.NewStartEndTimesCalcTz() "+
					" Iana Time Zone='%v' \n",
					job.IanaTimeZone)

				om.SetFatalError(s, err, 1212)
				return om
			}
		}
	}

	om.SetSuccessfulCompletionMessage("Finished SetDelayCmdStartTime", 1299)
	return om
}

// SetCmdJobActualStartTime - Computes and saves the date time
// at which job execution commenced.
func (job *CmdJob) SetCmdJobActualStartTime(parent []OpsMsgContextInfo) OpsMsgDto {

	msgCtx := OpsMsgContextInfo {
								SourceFileName: srcFileNameXMLCmdJobsData,
								ParentObjectName: "CmdJob",
								FuncName: "SetCmdJobActualStartTime",
								BaseMessageId: errBlockNoXMLCmdJobsData,
						}

	om := OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	var err error

	job.CmdJobStartTimeValue, err =
						dt.TimeZoneDto{}.New(time.Now().UTC(),
																	job.IanaTimeZone,
																	dt.FmtDateTimeYMDAbbrvDowNano)

	if err != nil {
		s:= fmt.Sprintf("Error returned by dt.TimeZoneDto{}.New(time.Now().UTC()...) " +
			" Start Time %v. Iana Time Zone %v. \n",
			time.Now().UTC().Format(dt.FmtDateTimeYMDAbbrvDowNano),
			job.IanaTimeZone)

		om.SetFatalError(s, err, 1401)

		return  om
	}

	om.SetSuccessfulCompletionMessage("Finished SetCmdJobActualStartTime", 1409)
	return om
}

// SetCmdJobActualEndTime - Sets the time at which this job completed processing.
// In addition, this method also computes the elapsed time required to
// complete processing of this command job.

func (job *CmdJob) SetCmdJobActualEndTime(parent []OpsMsgContextInfo) OpsMsgDto {

	//bi := ErrBaseInfo{}.New(srcFileNameXMLCmdJobsData, "CmdJob.SetCmdJobActualEndTime", errBlockNoXMLCmdJobsData)
	msgCtx := OpsMsgContextInfo{
							SourceFileName: srcFileNameXMLCmdJobsData,
							ParentObjectName: "CmdJob",
							FuncName: "SetCmdJobActualEndTime",
							BaseMessageId: errBlockNoXMLCmdJobsData,
					}

	//se:= SpecErr{}.InitializeBaseInfo(parent, bi)
	om := OpsMsgDto{}.InitializeAllContextInfo(parent, msgCtx)

	var err error

	job.CmdJobEndTimeValue, err =
		dt.TimeZoneDto{}.New( time.Now().UTC(),
													job.IanaTimeZone,
													dt.FmtDateTimeYMDAbbrvDowNano)

	if err != nil {
		s:= fmt.Sprintf("Error job.CmdJobEndTimeValue = " +
			"dt.TimeZoneDto{}.New(time.Now().UTC()...) " +
			"End Time %v. Local Time Zone %v.",
			time.Now().UTC().Format(dt.FmtDateTimeYMDAbbrvDowNano),
			job.IanaTimeZone)

		om.SetFatalError(s, err, 1401)
		return om
	}

	job.CmdJobDuration, err = dt.DurationTriad{}.NewStartEndTimesCalcTz(
								job.CmdJobStartTimeValue.TimeOut.DateTime,
								job.CmdJobEndTimeValue.TimeOut.DateTime,
								dt.TDurCalcTypeSTDYEARMTH,
								job.IanaTimeZone,
								dt.FmtDateTimeYMDAbbrvDowNano)

	if err != nil {
		s:= fmt.Sprintf("Error job.CmdJobDuration = " +
			"dt.DurationTriad{}.NewStartEndTimesCalcTz(). " +
			"Start Time %v. End Time %v. Iana Time Zone='%v'\n",
			job.CmdJobStartTimeValue.TimeOut.String(),
			job.CmdJobEndTimeValue.TimeOut.String(),
			job.IanaTimeZone)

		om.SetFatalError(s, err, 1403)
		return om
	}

	om.SetSuccessfulCompletionMessage("Finished InitializeAllContextInfo", 1409)
	return om
}

// CommandArgumentsArray - Holds CmdElement structures
type CommandArgumentsArray struct {
	CmdArgs []string `xml:"CmdArg"`
}

type CommandInputsArray struct {
	InputArgs [] string `xml:"InputArg"`
}
