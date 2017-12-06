package common

import "time"

/*  'opsmsgdto.go' is located in source code
		repository:

		https://github.com/MikeAustin71/logopsgo.git

*/



// OpsMsgDto - Data Transfer Object
// containing information about an
// operations Message
type OpsMsgDto struct {
	Message  []string
	MsgType  LogMsgType
	MsgLevel LogLevel
	MsgTimeUTC time.Time
	MsgTimeLocal time.Time
	LocalTimeZone string
	ErrDto   SpecErr
}

func (opsMsg *OpsMsgDto) NewSpecErr(se SpecErr) OpsMsgDto {

	om := OpsMsgDto{}

	if se.IsPanic {
		om.MsgType = LogERRORMSGTYPE
		om.MsgLevel = LogFATAL
		opsMsg.ErrDto = se

	} else if se.IsErr {
		om.MsgType = LogERRORMSGTYPE
		om.MsgLevel = LogDEBUG
		opsMsg.ErrDto = se

	} else {
		om.MsgType = LogINFOMSGTYPE
		om.MsgLevel = LogINFO
		opsMsg.ErrDto = se
	}

	opsMsg.Message = append(opsMsg.Message, se.Error())

	return om
}

func(opsMsg *OpsMsgDto) NewInfoMsg(msg string) OpsMsgDto {

	om := OpsMsgDto{}
	om.MsgType = LogINFOMSGTYPE
	om.MsgLevel = LogINFO

	om.Message = append(om.Message, msg)

	return om
}

func(opsMsg *OpsMsgDto)SetTime(localTimeZone string){

	tz := TimeZoneUtility{}

	isValid, _, _ := tz.IsValidTimeZone(localTimeZone)

	if !isValid {
		localTimeZone = "Local"
	}

	opsMsg.MsgTimeUTC = time.Now().UTC()
	opsMsg.LocalTimeZone = localTimeZone

	tzLocal, _ := tz.ConvertTz(opsMsg.MsgTimeUTC, opsMsg.LocalTimeZone)
	opsMsg.MsgTimeLocal = tzLocal.TimeOut

}