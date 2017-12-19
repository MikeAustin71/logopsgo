package common

import (
	"testing"
	"strings"

)

func TestOpsMsgCollection_AddOpsMsg_01(t *testing.T) {

	om1 := testOpsMsgCollectionCreateFatalErrorMsg_01()

	om2 := testOpsMsgCollectionCreateStdErrorMsg_01()

	om3 := testOpsMsgCollectionCreateInfoMsg()

	om4 := testOpsMsgCollectionCreateDebugMsg()

	om5 := testOpsMsgCollectionCreateWarningMsg()

	expectedArrayLen := 5

	opMsgs := OpsMsgCollection{}

	opMsgs.AddOpsMsg(om1)
	opMsgs.AddOpsMsg(om2)
	opMsgs.AddOpsMsg(om3)
	opMsgs.AddOpsMsg(om4)
	opMsgs.AddOpsMsg(om5)

	actualArrayLen := len(opMsgs.OpsMessages)

	if expectedArrayLen != actualArrayLen {
		t.Errorf("Expected Message Array Length = '%v'. Actual Message Array Length = '%v'", expectedArrayLen, actualArrayLen)
	}

	xMsg := "This is Information Message for test object"
	msgId := int64(19)
	msgNo := int64(6019)
	msgType := OpsMsgTypeINFOMSG

	testParentHistory := testOpsMsgDtoCreateParentHistory()
	testMsgContext := testOpsMsgDtoCreateContextInfoObj()

	om := opMsgs.OpsMessages[2].CopyOut()

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgContext.Equal(&om.MsgContext) {
		t.Error("Expected testMsgContext to EQUAL om.MsgContext. It did NOT!")
	}

	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != false {
		t.Error("Expected Information Message to generate IsError='false'. It did NOT! IsError='true'.")
	}

	if om.IsFatalError() != false {
		t.Errorf("Expected Information to generate IsFatalError()='false'. It did NOT! IsFatalError()='%v'", om.IsFatalError())
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

}

func TestOpsMsgCollection_PopLastMsg_01(t *testing.T) {

	opMsgs := testOpsMsgCollectionCreateT01Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}


	om := opMsgs.PopLastMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()
	testMsgContext := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is FATAL Error Msg for test object"
	msgId := int64(152)
	msgNo := int64(6152)
	msgType := OpsMsgTypeFATALERRORMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgContext.Equal(&om.MsgContext) {
		t.Error("Expected testMsgContext to EQUAL om.MsgContext. It did NOT!")
	}


	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsError='true'. It did NOT! IsError='false'.")
	}

	if om.IsFatalError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsFatalError()='true'. It did NOT! IsFatalError()='%v'", om.IsFatalError())
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

	expectedLen--
	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PopLastMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PopLastMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

}

func TestOpsMsgCollection_PopLastMsg_02(t *testing.T) {

	opMsgs := testOpsMsgCollectionCreateT01Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

	om := opMsgs.PopLastMsg()
	om = opMsgs.PopLastMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()

	testMsgCtx := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is Standard Error Msg for test object"
	msgId := int64(429)
	msgNo := int64(6429)
	msgType := OpsMsgTypeERRORMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgCtx.Equal(&om.MsgContext) {
		t.Error("Expected testMsgCtx to EQUAL om.MsgContext. It did NOT!")
	}

	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != true {
		t.Error("Expected error msg to generate IsError='true'. It did NOT! IsError='false'.")
	}

	if om.IsFatalError() == true {
		t.Error("Expected standard error msg to generate IsFatalError()='false'. It did NOT! IsFatalError()='true'")
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

	expectedLen--
	expectedLen--
	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PopLastMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PopLastMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}



}

func TestOpsMsgCollection_PeekLastMsg_01(t *testing.T) {

	opMsgs := testOpsMsgCollectionCreateT01Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}


	om := opMsgs.PeekLastMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()
	testMsgContext := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is FATAL Error Msg for test object"
	msgId := int64(152)
	msgNo := int64(6152)
	msgType := OpsMsgTypeFATALERRORMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgContext.Equal(&om.MsgContext) {
		t.Error("Expected testMsgContext to EQUAL om.MsgContext. It did NOT!")
	}

	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsError='true'. It did NOT! IsError='false'.")
	}

	if om.IsFatalError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsFatalError()='true'. It did NOT! IsFatalError()='%v'", om.IsFatalError())
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}


	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PeekMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PeekMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

}

func TestOpsMsgCollection_PeekLastMsg_02(t *testing.T) {

	opMsgs := testOpsMsgCollectionCreateT01Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

	opMsgs.PopLastMsg()
	opMsgs.PopLastMsg()

	om := opMsgs.PeekLastMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()
	testMsgContext := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is Warning Message for test object."
	msgId := int64(67)
	msgNo := int64(6067)
	msgType := OpsMsgTypeWARNINGMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgContext.Equal(&om.MsgContext) {
		t.Error("Expected testMsgContext to EQUAL om.MsgContext. It did NOT!")
	}

	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != false {
		t.Error("Expected Warning Message to generate IsError='false'. It did NOT! IsError='true'.")
	}

	if om.IsFatalError() != false {
		t.Errorf("Expected Warning to generate IsFatalError()='false'. It did NOT! IsFatalError()='%v'", om.IsFatalError())
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

	expectedLen--
	expectedLen--
	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PeekMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PeekMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

}

func TestOpsMsgCollection_PopFirstMsg_01(t *testing.T) {

	opMsgs := testOpsMsgCollectionCreateT02Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}


	om := opMsgs.PopFirstMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()
	testMsgContext := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is FATAL Error Msg for test object"
	msgId := int64(152)
	msgNo := int64(6152)
	msgType := OpsMsgTypeFATALERRORMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgContext.Equal(&om.MsgContext) {
		t.Error("Expected testMsgContext to EQUAL om.MsgContext. It did NOT!")
	}


	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsError='true'. It did NOT! IsError='false'.")
	}

	if om.IsFatalError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsFatalError()='true'. It did NOT! IsFatalError()='%v'", om.IsFatalError())
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

	expectedLen--
	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PopLastMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PopLastMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

}

func TestOpsMsgCollection_PopFirstMsg_02(t *testing.T) {

	opMsgs := testOpsMsgCollectionCreateT02Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

	om := opMsgs.PopFirstMsg()

	om = opMsgs.PopFirstMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()

	testMsgCtx := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is Standard Error Msg for test object"
	msgId := int64(429)
	msgNo := int64(6429)
	msgType := OpsMsgTypeERRORMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgCtx.Equal(&om.MsgContext) {
		t.Error("Expected testMsgCtx to EQUAL om.MsgContext. It did NOT!")
	}

	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != true {
		t.Error("Expected error msg to generate IsError='true'. It did NOT! IsError='false'.")
	}

	if om.IsFatalError() == true {
		t.Error("Expected standard error msg to generate IsFatalError()='false'. It did NOT! IsFatalError()='true'")
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

	expectedLen--
	expectedLen--
	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PopLastMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PopLastMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

}

func TestOpsMsgCollection_PeekFirstMsg_01(t *testing.T) {
	opMsgs := testOpsMsgCollectionCreateT02Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}


	om := opMsgs.PeekFirstMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()
	testMsgContext := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is FATAL Error Msg for test object"
	msgId := int64(152)
	msgNo := int64(6152)
	msgType := OpsMsgTypeFATALERRORMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgContext.Equal(&om.MsgContext) {
		t.Error("Expected testMsgContext to EQUAL om.MsgContext. It did NOT!")
	}


	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsError='true'. It did NOT! IsError='false'.")
	}

	if om.IsFatalError() != true {
		t.Errorf("Expected Fatal Error Message to generate IsFatalError()='true'. It did NOT! IsFatalError()='%v'", om.IsFatalError())
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PeekMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PeekMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

}

func TestOpsMsgCollection_PeekFirstMsg_02(t *testing.T) {

	opMsgs := testOpsMsgCollectionCreateT02Collection()
	expectedLen := 7
	actualLen := len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("Expected opMsgs.OpsMessages Array Length = '%v'.  Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

	opMsgs.PopFirstMsg()
	opMsgs.PopFirstMsg()

	om := opMsgs.PeekFirstMsg()

	testParentHistory := testOpsMsgDtoCreateParentHistory()
	testMsgContext := testOpsMsgDtoCreateContextInfoObj()

	xMsg := "This is Information Message for test object"
	msgId := int64(19)
	msgNo := int64(6019)
	msgType := OpsMsgTypeINFOMSG

	l1 := len(testParentHistory)

	l2 := len(om.ParentContextHistory)

	if l1 != l2 {
		t.Errorf("Expected om.ParentContextHistory to equal length= '%v'. It did NOT! actual length= '%v'",l1, l2)
	}

	for i:=0; i<l1; i++ {
		if !testParentHistory[i].Equal(&om.ParentContextHistory[i]) {
			t.Errorf("Expected om.ParentContextHistory to Equal testParentContext History. It did NOT!. i= '%v'",i)
		}
	}

	if !testMsgContext.Equal(&om.MsgContext) {
		t.Error("Expected testMsgContext to EQUAL om.MsgContext. It did NOT!")
	}

	if om.MsgType != msgType {
		t.Errorf("Expected Messgage Type == '%v'. Instead, Message Type == '%v'.", msgType, om.MsgType)
	}

	if om.IsError() != false {
		t.Error("Expected Information Message to generate IsError='false'. It did NOT! IsError='true'.")
	}

	if om.IsFatalError() != false {
		t.Errorf("Expected Information to generate IsFatalError()='false'. It did NOT! IsFatalError()='%v'", om.IsFatalError())
	}

	mId := om.GetMessageId()

	if mId != msgId {
		t.Errorf("Expected message id = '%v'. Instead message id = '%v'.", msgId, mId)
	}

	mNo := om.GetMessageNumber()

	if msgNo != mNo {
		t.Errorf("Expected message number = '%v'. Instead message number = '%v'.", msgNo, mNo)
	}

	actMsg := om.GetFmtMessage()

	if !strings.Contains(actMsg, xMsg) {
		t.Errorf("Expected message to contain '%v'. It did NOT! Actual Message = '%v'",xMsg, actMsg)
	}

	if om.MsgTimeUTC.IsZero()  {
		t.Errorf("Error: om.MsgTimeUTC == Zero. om.MsgTimeUTC== '%v'", om.MsgTimeUTC)
	}

	if om.MsgTimeLocal.IsZero()  {
		t.Errorf("Error: om.MsgTimeLocal == Zero. om.MsgTimeLocal== '%v'",om.MsgTimeLocal)
	}

	if om.MsgLocalTimeZone != "Local" {
		t.Errorf("Error: om.MsgLocalTimeZone is NOT set to 'Local'. om.MsgLocalTimeZone== '%v' ", om.MsgLocalTimeZone)
	}

	expectedLen--
	expectedLen--
	actualLen = len(opMsgs.OpsMessages)

	if expectedLen != actualLen {
		t.Errorf("After PeekMsg - Expected opMsgs.OpsMessages Array Length = '%v'. After PeekMsg Actual opMsgs.OpsMessages Array Length = '%v'", expectedLen, actualLen)
	}

}

/*
=======================================================================================================
								Private Methods
=======================================================================================================
*/

func testOpsMsgCollectionCreateContextInfoObj() OpsMsgContextInfo {
	ci := OpsMsgContextInfo{}
	return ci.New("TSource06", "PObj06", "Func006", 6000)
}

func testOpsMsgCollectionCreateParentHistory() []OpsMsgContextInfo {
	ci := OpsMsgContextInfo{}

	x1 := ci.New("TSource01", "PObj01", "Func001", 1000)
	x2 := ci.New("TSource02", "PObj02", "Func002", 2000)
	x3 := ci.New("TSource03", "PObj03", "Func003", 3000)
	x4 := ci.New("TSource04", "PObj04", "Func004", 4000)
	x5 := ci.New("TSource05", "PObj05", "Func005", 5000)

	parent := make([]OpsMsgContextInfo,0,10)

	parent = append(parent, x1)
	parent = append(parent, x2)
	parent = append(parent, x3)
	parent = append(parent, x4)
	parent = append(parent, x5)

	return parent
}


func testOpsMsgCollectionCreateStdErrorMsg_01() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetStdErrorMessage("This is Standard Error Msg for test object", 429)
	return om
}

func testOpsMsgCollectionCreateStdErrorMsg_02() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetStdErrorMessage("This is Standard Error Msg #2 for test object", 430)
	return om
}

func testOpsMsgCollectionCreateFatalErrorMsg_01() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetFatalErrorMessage("This is FATAL Error Msg for test object", 152)
	return om
}

func testOpsMsgCollectionCreateFatalErrorMsg_02() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetFatalErrorMessage("This is FATAL Error Msg #2 for test object", 153)
	return om
}

func testOpsMsgCollectionCreateInfoMsg() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetInfoMessage("This is Information Message for test object", 19)
	return om
}

func testOpsMsgCollectionCreateWarningMsg() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetWarningMessage("This is Warning Message for test object.", 67)
	return om
}

func testOpsMsgCollectionCreateDebugMsg() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetDebugMessage("This is DEBUG Message for test object.", 238)
	return om
}

func testOpsMsgCollectionCreateSuccessfulCompletionMsg() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetSuccessfulCompletionMessage("", 64)
	return om
}

func testOpsMsgCollectionCreateNoErrorsNoMessagesMsg() OpsMsgDto {
	om := OpsMsgDto{}.InitializeAllContextInfo(testOpsMsgDtoCreateParentHistory(), testOpsMsgDtoCreateContextInfoObj())
	om.SetNoErrorsNoMessages("", 28)
	return om
}

func testOpsMsgCollectionCreateT01Collection() OpsMsgCollection {

	om1 := testOpsMsgCollectionCreateFatalErrorMsg_02()

	om2 := testOpsMsgCollectionCreateStdErrorMsg_02()

	om3 := testOpsMsgCollectionCreateInfoMsg()

	om4 := testOpsMsgCollectionCreateDebugMsg()

	om5 := testOpsMsgCollectionCreateWarningMsg()

	om6 := testOpsMsgCollectionCreateStdErrorMsg_01()

	om7 := testOpsMsgCollectionCreateFatalErrorMsg_01()

	opMsgs := OpsMsgCollection{}

	opMsgs.AddOpsMsg(om1)
	opMsgs.AddOpsMsg(om2)
	opMsgs.AddOpsMsg(om3)
	opMsgs.AddOpsMsg(om4)
	opMsgs.AddOpsMsg(om5)
	opMsgs.AddOpsMsg(om6)
	opMsgs.AddOpsMsg(om7)

	return opMsgs
}

func testOpsMsgCollectionCreateT02Collection() OpsMsgCollection {

	om1 := testOpsMsgCollectionCreateFatalErrorMsg_01()

	om2 := testOpsMsgCollectionCreateStdErrorMsg_01()

	om3 := testOpsMsgCollectionCreateInfoMsg()

	om4 := testOpsMsgCollectionCreateDebugMsg()

	om5 := testOpsMsgCollectionCreateWarningMsg()

	om6 := testOpsMsgCollectionCreateStdErrorMsg_02()

	om7 := testOpsMsgCollectionCreateFatalErrorMsg_02()

	opMsgs := OpsMsgCollection{}

	opMsgs.AddOpsMsg(om1)
	opMsgs.AddOpsMsg(om2)
	opMsgs.AddOpsMsg(om3)
	opMsgs.AddOpsMsg(om4)
	opMsgs.AddOpsMsg(om5)
	opMsgs.AddOpsMsg(om6)
	opMsgs.AddOpsMsg(om7)

	return opMsgs
}