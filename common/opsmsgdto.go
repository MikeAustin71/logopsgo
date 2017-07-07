package common

// OpsMsgDto - Data Transfer Object
// containing information about an
// operations message
type OpsMsgDto struct {
	message  []string
	msgType  LogMsgType
	msgLevel LogLevel
	specErr  SpecErr
}

func (opsMsg *OpsMsgDto) NewSpecErr(se SpecErr) OpsMsgDto {

	om := OpsMsgDto{}

	if se.IsPanic {
		om.msgType = LogERRORMSGTYPE
		om.msgLevel = LogFATAL
		opsMsg.specErr = se

	} else if se.IsErr {
		om.msgType = LogERRORMSGTYPE
		om.msgLevel = LogDEBUG
		opsMsg.specErr = se

	} else {
		om.msgType = LogINFOMSGTYPE
		om.msgLevel = LogINFO
		opsMsg.specErr = se
	}

	opsMsg.message = append(opsMsg.message, se.Error())

	return om
}

func(opsMsg *OpsMsgDto) NewInfoMsg(msg string) OpsMsgDto {

	om := OpsMsgDto{}
	om.msgType = LogINFOMSGTYPE
	om.msgLevel = LogINFO

	opsMsg.message = append(opsMsg.message, msg)

	return om
}