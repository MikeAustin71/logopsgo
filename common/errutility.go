package common

import (
	"fmt"
)

// SpecErr - A data structure used
// to hold custom error information
type SpecErr struct {
	IsErr     bool
	IsPanic   bool
	PrefixMsg string
	ErrMsg    string
	SrcFile   string
	FuncName  string
	ErrNo     int64
}

// New - Creates new SpecErr Type
func (s SpecErr) New(prefix string, err error, isPanic bool, srcFile string, funcName string, errNo int64) SpecErr {

	x := SpecErr{PrefixMsg: prefix, IsPanic: isPanic, SrcFile: srcFile, FuncName: funcName, ErrNo: errNo}

	if err != nil {
		x.ErrMsg = err.Error()
		x.IsErr = true
	} else {
		x.ErrMsg = ""
		x.IsErr = false
	}

	return x
}

// SetNoError - Returns a SpecErr
// structure with IsErr set to false.
func (s SpecErr) SetNoError() SpecErr {
	return SpecErr {IsErr: false, IsPanic: false}
}

// Panic - Executes 'panic' command
// if IsPanic == 'true'
func (s SpecErr) Panic() {
	if s.IsPanic {
		panic(s)
	}
}

// Error - Implements Error Interface
func (s SpecErr) Error() string {
	m := s.PrefixMsg
	m += "\n" + s.ErrMsg

	if s.SrcFile != "" {
		m += "\nSourceFile: " + s.SrcFile
	}

	if s.FuncName != "" {
		m += "\nFuncName: " + s.FuncName
	}

	if s.ErrNo != 0 {
		m += fmt.Sprintf("\nErrNo: %v", s.ErrNo)
	}

	m += fmt.Sprintf("\nIsErr: %v", s.IsErr)
	m += fmt.Sprintf("\nIsPanic: %v", s.IsPanic)
	return m
}

// CheckErrPanic - Checks for error and then
// executes 'panic'
func CheckErrPanic(e error) {
	if e != nil {
		panic(e)
	}
}

// CheckIsSpecErr - If error is present,
// returns 'true'.  If NO Error, returns
// 'false'.
func CheckIsSpecErr(eSpec SpecErr) bool {

	if eSpec.IsErr {
		return true
	}

	return false

}

// CheckSpecErrPanic - Issues a 'panic'
// command if SpecErr IsPanic flag is set
func CheckSpecErrPanic(eSpec SpecErr) {

	if eSpec.IsPanic {
		panic(eSpec)
	}
}
