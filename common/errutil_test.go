package common

import (
	"errors"
	"fmt"
	"testing"
	"strings"
)

func TestErrorUtility(t *testing.T) {

	ex1 := "errutil_test.go"
	ex2 := "TestErrorUtility"
	ex3 := int64(10000)

	ex4 := "errprefix"
	ex5 := int64(334)
	ex6 := "Test Error #1"
	ex7 := ex5 + ex3

	err := errors.New(ex6)
	a := make([]ErrBaseInfo, 0, 10)
	bi := ErrBaseInfo{}.New(ex1, ex2, ex3)
	se := SpecErr{}.InitializeBaseInfo(a, bi)
	x := se.New(ex4, err, true, ex5)

	if x.PrefixMsg != ex4 {
		t.Error(fmt.Sprintf("Expected '%v' got", ex4), x.PrefixMsg)
	}

	if x.ErrMsg != ex6 {
		t.Error(fmt.Sprintf("Expected '%v' got", ex6), x.ErrMsg)
	}

	if x.BaseInfo.SourceFileName != ex1 {
		t.Error(fmt.Sprintf("Expected '%v' got", ex1), x.BaseInfo.SourceFileName)
	}

	if x.BaseInfo.FuncName != ex2 {
		t.Error(fmt.Sprintf("Expected '%v' got", ex2), x.BaseInfo.FuncName)
	}

	if x.ErrNo != ex7 {
		t.Error(fmt.Sprintf("Expected '%v' got", ex7), x.ErrNo)
	}

}

func TestUninitializedBaseInfo(t *testing.T) {
	var se SpecErr

	if se.BaseInfo.SourceFileName != "" {
		t.Error("String SourceFileName was uninitialized. Was expecting empty string, got", se.BaseInfo.SourceFileName)
	}

	if se.BaseInfo.FuncName != "" {
		t.Error("String FuncName was uninitialized. Was expecting empty string, got", se.BaseInfo.FuncName)
	}

	if se.BaseInfo.BaseErrorID != 0 {
		t.Error("Int64 BaseErrorID was uninitialized. Was expecting value of zero, got", se.BaseInfo.BaseErrorID)
	}

}

func TestInitializeParentInfo(t *testing.T) {

	bi := ErrBaseInfo{}

	x := bi.New("TestSourceFileName", "TestFuncName", 9000)
	y := bi.New("TestSrcFileName2", "TestFuncName2", 14000)
	z := bi.New("TestSrcFileName3", "TestFuncName3", 15000)

	se := SpecErr{}
	se.ParentInfo = append(se.ParentInfo, x)
	se.ParentInfo = append(se.ParentInfo, y)
	se.ParentInfo = append(se.ParentInfo, z)

	l := len(se.ParentInfo)

	if l != 3 {
		t.Error("Expected ParentInfo Length of 3, got", l)
	}

	if se.ParentInfo[1].FuncName != "TestFuncName2" {
		t.Error("Expected 2nd Element 'TestFuncName2', got", se.ParentInfo[1].FuncName)
	}

}

func TestAddSlicesParentInfo(t *testing.T) {
	var bi ErrBaseInfo
	x := bi.New("TestSourceFileName", "TestFuncName", 9000)
	y := bi.New("TestSrcFileName2", "TestFuncName2", 14000)
	z := bi.New("TestSrcFileName3", "TestFuncName3", 15000)

	a := make([]ErrBaseInfo, 0, 30)

	a = append(a, x, y, z)

	var se SpecErr

	se.ParentInfo = a

	l := len(se.ParentInfo)

	if l != 3 {
		t.Error("Expected ParentInfo Length of 3, got", l)
	}

	if se.ParentInfo[1].FuncName != "TestFuncName2" {
		t.Error("Expected 2nd Element 'TestFuncName2', got", se.ParentInfo[1].FuncName)
	}

}

func TestSetParentInfo(t *testing.T) {
	var bi ErrBaseInfo
	x := bi.New("TestSourceFileName", "TestFuncName", 9000)
	y := bi.New("TestSrcFileName2", "TestFuncName2", 14000)
	z := bi.New("TestSrcFileName3", "TestFuncName3", 15000)

	a := make([]ErrBaseInfo, 0, 30)

	a = append(a, x, y, z)

	var se SpecErr

	se.ParentInfo = se.DeepCopyParentInfo(a)

	l := len(se.ParentInfo)

	if l != 3 {
		t.Error("Expected ParentInfo length of 3, go length of ", l)
	}

	if se.ParentInfo[1].FuncName != "TestFuncName2" {
		t.Error("Expected 2nd Element 'TestFuncName2', got", se.ParentInfo[1].FuncName)
	}
}

func TestSetErrDetail(t *testing.T) {

	ex1 := "errutil_test.go"
	ex2 := "TestErrorUtility"
	ex3 := int64(10000)

	ex4 := "errprefix"
	ex5 := int64(338)
	ex6 := "Test Error #21"
	ex7 := ex5 + ex3

	err := errors.New(ex6)
	a := make([]ErrBaseInfo, 0, 10)
	bi := ErrBaseInfo{}.New(ex1, ex2, ex3)

	x := SpecErr{}.InitializeBaseInfo(a, bi).New(ex4, err, true, ex5)

	if x.ErrNo != ex7 {
		t.Error(fmt.Sprintf("Expected Err No '%v', got", ex7), x.ErrNo)
	}

	if x.BaseInfo.SourceFileName != ex1 {
		t.Error(fmt.Sprintf("Expected Source File: '%v',got", ex1), x.BaseInfo.SourceFileName)
	}

	if x.BaseInfo.FuncName != ex2 {
		t.Error(fmt.Sprintf("Expected FuncName: '%v', got", ex2), x.BaseInfo.FuncName)
	}

}

func TestIsSpecErrNo(t *testing.T) {

	s := SpecErr{}.SignalNoErrors()

	isErr := s.CheckIsSpecErr()

	if isErr {
		t.Errorf("Expected CheckIsSpecErr() to return false. Instead it retrned %v", isErr)
	}

}

func TestIsSpecErrYes(t *testing.T) {
	ex1 := "errprefix"
	ex4 := int64(334)
	ex99 := "Test Error #1"

	err := errors.New(ex99)

	x := SpecErr{}.New(ex1, err, true, ex4)

	isErr := x.CheckIsSpecErr()

	if !isErr {
		t.Errorf("Expected x.CheckIsSpecErr() to return 'true'. Instead it returned %v", isErr)
	}

}

func TestSetNoErr(t *testing.T) {
	x := SpecErr{}.SignalNoErrors()

	if x.IsErr {
		t.Errorf("Expected IsErr= 'false'. Instead IsErr= '%v'", x.IsErr)
	}
}

func TestQuickInitialize(t *testing.T) {

	ex1 := "prefixMsg"
	ex2 := "Error Msg X"
	err := errors.New(ex2)
	ex4 := int64(499)

	x := SpecErr{}.New(ex1, err, false, ex4)

	if x.ErrNo != ex4 {
		t.Error(fmt.Sprintf("Expected ErrNo: '%v', got", ex4), x.ErrNo)
	}

	if x.PrefixMsg != ex1 {
		t.Error(fmt.Sprintf("Expected PrefixMsg == '%v', got", ex1), x.PrefixMsg)
	}

	if x.ErrMsg != ex2 {
		t.Error(fmt.Sprintf("Expected Error Msg '%v', got,", ex2), x.ErrMsg)
	}

	if x.IsPanic == true {
		t.Errorf("Expected IsPanic == '%v'. Instead IsPanic== '%v' ", false, x.IsPanic)
	}
}

func TestFullInitialize(t *testing.T) {

	bi := ErrBaseInfo{}

	f := bi.New("TestSourceFileName", "TestFuncName", 9000)
	g := bi.New("TestSrcFileName2", "TestFuncName2", 14000)
	h := bi.New("TestSrcFileName3", "TestFuncName3", 15000)

	ex1 := make([]ErrBaseInfo, 0, 10)
	ex1 = append(ex1, f, g, h)

	ex21 := "TestSrcFileName99"
	ex22 := "TestFuncName99"
	ex23 := int64(16000)
	ex2 := bi.New(ex21, ex22, ex23)

	ex3 := "prefixString"
	ex4 := "Error Msg 99"
	err := errors.New(ex4)
	ex5 := false
	ex6 := int64(22)
	ex7 := int64(16022)

	x := SpecErr{}.Initialize(ex1, ex2, ex3, err, false, ex6)

	pl := len(x.ParentInfo)

	if pl != 3 {
		t.Error("Expected ParentInfo length == 3, got", pl)
	}

	if x.PrefixMsg != ex3 {
		t.Error(fmt.Sprintf("Expected PrefixMsg '%v', got", ex3), x.PrefixMsg)
	}

	if x.ErrMsg != ex4 {
		t.Error(fmt.Sprintf("Expected ErrMsg '%v', got", ex4), x.ErrMsg)
	}

	if x.IsErr == false {
		t.Error(fmt.Sprintf("Expected IsErr == '%v', got", true), x.IsErr)
	}

	if x.IsPanic != ex5 {
		t.Errorf("Expected IsPanic == '%v'. Instead IsPannic=='%v'", ex5, x.IsPanic)
	}

	if x.ErrNo != ex7 {
		t.Error(fmt.Sprintf("Expected ErrNo '%v', got", ex7), x.ErrNo)
	}

	if x.BaseInfo.SourceFileName != ex21 {
		t.Error(fmt.Sprintf("Expected SourceFileName '%v', got", ex21), x.BaseInfo.SourceFileName)
	}

	if x.BaseInfo.FuncName != ex22 {
		t.Error(fmt.Sprintf("Expected Function Name '%v', got", ex22), x.BaseInfo.FuncName)
	}

	if x.BaseInfo.BaseErrorID != ex23 {
		t.Error(fmt.Sprintf("Expected Base Error ID '%v', got", ex23), x.BaseInfo.BaseErrorID)
	}

}

func TestBlankInitialize(t *testing.T) {

	ex1 := "prefixString"
	ex2 := "Error Msg 99"
	err := errors.New(ex2)
	ex3 := false
	ex4 := int64(22)

	x := SpecErr{}.Initialize(blankParentInfo, blankErrBaseInfo, ex1, err, ex3, ex4)

	if x.PrefixMsg != ex1 {
		t.Error(fmt.Sprintf("Expected PrefixMsg '%v', got", ex1), x.PrefixMsg)
	}

	if x.ErrMsg != ex2 {
		t.Error(fmt.Sprintf("Expected ErrMsg '%v', got", ex2), x.ErrMsg)
	}

	if x.IsErr == false {
		t.Error(fmt.Sprintf("Expected IsErr == '%v', got", true), x.IsErr)
	}

	if x.IsPanic != ex3 {
		t.Error(fmt.Sprintf("Expected IsPanic == '%v', got", ex3), x.IsPanic)
	}

	if x.ErrNo != ex4 {
		t.Error(fmt.Sprintf("Expected ErrNo '%v', got", ex4), x.ErrNo)
	}

	if x.BaseInfo.SourceFileName != "" {
		t.Error("Expected BaseInfo Source File Name == '', got", x.BaseInfo.SourceFileName)
	}

	if x.BaseInfo.FuncName != "" {
		t.Error("Expected BaseInfo FuncName == '', got", x.BaseInfo.FuncName)
	}

	if x.BaseInfo.BaseErrorID != int64(0) {
		t.Error("Expected BaseInfo BaseErrorID == 'Zero', got", x.BaseInfo.BaseErrorID)
	}

}

func TestAddBaseInfoToParent(t *testing.T) {

	bi := ErrBaseInfo{}

	f := bi.New("TestSourceFileName", "TestFuncName", 9000)
	g := bi.New("TestSrcFileName2", "TestFuncName2", 14000)
	h := bi.New("TestSrcFileName3", "TestFuncName3", 15000)

	ex1 := make([]ErrBaseInfo, 0, 10)
	ex1 = append(ex1, f, g, h)

	ex21 := "TestSrcFileName99"
	ex22 := "TestFuncName99"
	ex23 := int64(16000)
	ex2 := bi.New(ex21, ex22, ex23)

	ex3 := "prefixString"
	ex4 := "Error Msg 99"
	err := errors.New(ex4)
	ex5 := false
	ex6 := int64(22)

	x := SpecErr{}.Initialize(ex1, ex2, ex3, err, ex5, ex6)

	pl := len(x.ParentInfo)

	if pl != 3 {
		t.Error("Expected Parent Info to contain 3-elements. Actual number of elements was ", pl)
	}

	p2 := x.AddBaseToParentInfo()

	pl2 := len(p2)

	if pl2 != 4 {
		t.Error("Expected New Parent Info to contain 4-elements. Actual number of elements was ", pl2)
	}

}

func TestSpecErr_CheckIsSpecErr(t *testing.T) {

	bi := ErrBaseInfo{}

	f := bi.New("TestSourceFileName", "TestFuncName", 9000)
	g := bi.New("TestSrcFileName2", "TestFuncName2", 14000)
	h := bi.New("TestSrcFileName3", "TestFuncName3", 15000)

	ex1 := make([]ErrBaseInfo, 0, 10)
	ex1 = append(ex1, f, g, h)

	ex21 := "TestSrcFileName99"
	ex22 := "TestFuncName99"
	ex23 := int64(16000)
	ex2 := bi.New(ex21, ex22, ex23)

	ex3 := "prefixString"
	ex4 := "Error Msg 99"
	err := errors.New(ex4)
	ex6 := int64(22)

	x := SpecErr{}.Initialize(ex1, ex2, ex3, err, false, ex6)

	result := x.CheckIsSpecErr()

	if result == false {
		t.Errorf("Expected x.CheckIsSpecErr()== 'true'. Instead, x.CheckIsSpecErr()== '%v'", result)
	}

}

func TestSpecErr_CheckIsSpecErrPanic(t *testing.T) {
	bi := ErrBaseInfo{}

	f := bi.New("TestSourceFileName", "TestFuncName", 9000)
	g := bi.New("TestSrcFileName2", "TestFuncName2", 14000)
	h := bi.New("TestSrcFileName3", "TestFuncName3", 15000)

	ex1 := make([]ErrBaseInfo, 0, 10)
	ex1 = append(ex1, f, g, h)

	ex21 := "TestSrcFileName99"
	ex22 := "TestFuncName99"
	ex23 := int64(16000)
	ex2 := bi.New(ex21, ex22, ex23)

	ex3 := "prefixString"
	ex4 := "Error Msg 99"
	err := errors.New(ex4)
	ex6 := int64(22)

	x := SpecErr{}.Initialize(ex1, ex2, ex3, err, false, ex6)

	x.IsPanic = true

	result := x.CheckIsSpecErrPanic()

	if result != true {
		t.Errorf("Expected x.CheckIsSpecErrPanic()== 'true' .  Instead, x.CheckIsSpecErrPanic()== '%v'", result)
	}

}

func TestSpecErr_NewErrorMsgString_01(t *testing.T) {
	prefixString := "prefixString"
	errMsg := "This is the Error Msg!"
	isPanic := false
	errNo := int64(22)

	s :=  SpecErr{}.NewErrorMsgString(prefixString, errMsg, isPanic, errNo )

	if s.ErrMsg != errMsg {
		t.Errorf("Expected Error Message= '%v'.  Instead, Error Message = '%v'", errMsg, s.ErrMsg)
	}


}

func TestSpecErr_NewErrorMsgString_02(t *testing.T) {
	prefixString := "prefixString"
	errMsg := "This is the Error Msg!"
	isPanic := false
	errNo := int64(0)

	s :=  SpecErr{}.NewErrorMsgString(prefixString, errMsg, isPanic, errNo )

	if s.ErrMsg != errMsg {
		t.Errorf("Expected Error Message= '%v'.  Instead, Error Message = '%v'", errMsg, s.ErrMsg)
	}

	errResult := s.Error()

	hasErrNo := strings.Contains(errResult, "ErrNo:")

	if true == hasErrNo {
		t.Error("Due to Error Number=Zero, expected error message WITHOUT Error Number. Instead, Error Number was included in Error Message")
	}
}

func TestSpecErr_SetErrorMessageLabel(t *testing.T) {

	prefixString := "prefixString"
	errMsg := "This is the Error Msg!"
	isPanic := false
	errNo := int64(0)

	s :=  SpecErr{}.NewErrorMsgString(prefixString, errMsg, isPanic, errNo )

	s.SetErrorMessageLabel("StdOut Error")

	expectedErrMsg := "StdOut Error: " + errMsg
	actualErrMsg := s.Error()

	if strings.Contains(actualErrMsg, expectedErrMsg) == false {
		t.Errorf("Expected Error Message= '%v'.  Instead, Actual Error Message = '%v'", expectedErrMsg, actualErrMsg)
	}

}
