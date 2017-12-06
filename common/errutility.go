package common

import (
	"fmt"
	"errors"
)

/*
		The source code for errutility.go is located in source
		code repository:

						https://github.com/MikeAustin71/ErrHandlerGo.git
*/



// ErrBaseInfo is intended for use with
// the SpecErr Structure. It sets up base
// error info to be used repeatedly.
type ErrBaseInfo struct {
	SourceFileName string
	FuncName       string
	BaseErrorID    int64
}

// New - returns a new, populated ErrBaseInfo Structure
func (b ErrBaseInfo) New(srcFile, funcName string, baseErrID int64) ErrBaseInfo {
	return ErrBaseInfo{SourceFileName: srcFile, FuncName: funcName, BaseErrorID: baseErrID}
}

// NewFunc - Returns a New ErrBaseInfo structure with a different Func Name
func (b ErrBaseInfo) NewFunc(funcName string) ErrBaseInfo {
	return ErrBaseInfo{SourceFileName: b.SourceFileName, FuncName: funcName, BaseErrorID: b.BaseErrorID}
}

// NewBaseInfo - returns a deep copy of the current
// ErrBaseInfo structure.
func (b ErrBaseInfo) NewBaseInfo() ErrBaseInfo {
	return ErrBaseInfo{SourceFileName: b.SourceFileName, FuncName: b.FuncName, BaseErrorID: b.BaseErrorID}
}

// DeepCopyBaseInfo - Same as NewBaseInfo()
func (b ErrBaseInfo) DeepCopyBaseInfo() ErrBaseInfo {
	return ErrBaseInfo{SourceFileName: b.SourceFileName, FuncName: b.FuncName, BaseErrorID: b.BaseErrorID}
}

// GetBaseSpecErr - Returns an empty
// SpecErr structure populated with
// Base Error Info
func (b ErrBaseInfo) GetBaseSpecErr() SpecErr {

	return SpecErr{BaseInfo: b.NewBaseInfo()}
}

// GetNewParentInfo - Returns a slice of ErrBaseInfo
// structures with the first element initialized to a
// new ErrBaseInfo structure.
func (b ErrBaseInfo) GetNewParentInfo(srcFile, funcName string, baseErrID int64) []ErrBaseInfo {
	var parent []ErrBaseInfo

	bi := b.New(srcFile, funcName, baseErrID)

	return append(parent, bi)
}

// SpecErr - A data structure used
// to hold custom error information
type SpecErr struct {
	ParentInfo []ErrBaseInfo
	BaseInfo   ErrBaseInfo
	IsErr      bool
	IsPanic    bool
	PrefixMsg  string
	ErrMsgLabel string
	ErrMsg     string
	ErrNo      int64
}


// AddParentInfo - Adds ParentInfo elements to
// the current SpecErr ParentInfo slice
func (s *SpecErr) AddParentInfo(parent []ErrBaseInfo) {
	if len(parent) == 0 {
		return
	}

	x := s.DeepCopyParentInfo(parent)

	for _, bi := range x {
		s.ParentInfo = append(s.ParentInfo, bi.NewBaseInfo())
	}

	return

}

// AddBaseToParentInfo - Adds the structure's
// ErrBaseInfo data to ParentInfo and returns a
// new ParentInfo Array
func (s *SpecErr) AddBaseToParentInfo() []ErrBaseInfo {

	a := s.DeepCopyParentInfo(s.ParentInfo)
	return append(a, s.BaseInfo.DeepCopyBaseInfo())
}

// CheckErrPanic - Checks for error and then
// executes 'panic'
func (s *SpecErr) CheckErrPanic(e error) {
	if e != nil {
		panic(e)
	}
}

// CheckIsSpecErr - If error is present,
// returns 'true'.  If NO Error, returns
// 'false'.
func (s *SpecErr) CheckIsSpecErr() bool {

	if s.IsErr {
		return true
	}

	return false

}

// CheckIsSpecErrPanic - Returns 'true' if
// SpecErr object is configured as a panic
// error.
func (s *SpecErr) CheckIsSpecErrPanic() bool {

	return s.IsPanic
}

// DeepCopyBaseInfo - Returns a deep copy of the
// current BaseInfo structure.
func (s SpecErr) DeepCopyBaseInfo() ErrBaseInfo {
	return s.BaseInfo.DeepCopyBaseInfo()
}

// DeepCopyParentInfo - Receives an array of slices
// type ErrBaseInfo and appends deep copies
// of those slices to the SpecErr ParentInfo
// field.
func (s *SpecErr) DeepCopyParentInfo(pi []ErrBaseInfo) []ErrBaseInfo {

	if len(pi) == 0 {
		return pi
	}

	a := make([]ErrBaseInfo, 0, len(pi)+10)
	for _, bi := range pi {
		a = append(a, bi.NewBaseInfo())
	}

	return a
}

// Error - Implements Error Interface
func (s SpecErr) Error() string {

	m := "\nError Message:"
	m += "\n---------------------"
	if s.PrefixMsg != "" {
		m += "\n"
		m += s.PrefixMsg
	}

	m+= "\n"

	if s.ErrMsgLabel != "" {
		m+= s.ErrMsgLabel + ": "
	}

	m += s.ErrMsg
	m += "\n---------------------"

	if s.BaseInfo.SourceFileName != "" {
		m += "\nSourceFile: " + s.BaseInfo.SourceFileName
	}

	if s.BaseInfo.FuncName != "" {
		m += "\nFuncName: " + s.BaseInfo.FuncName
	}

	if s.ErrNo != 0 {
		m += fmt.Sprintf("\nErrNo: %v", s.ErrNo)
	}

	m += fmt.Sprintf("\nIsErr: %v", s.IsErr)
	m += fmt.Sprintf("\nIsPanic: %v", s.IsPanic)

	// If parent Function Info Exists
	// Print it out.
	if len(s.ParentInfo) > 0 {
		m += "\n---------------------"
		m += "\n  Parent Func Info"
		m += "\n---------------------"

		for _, bi := range s.ParentInfo {
			m += "\n" + bi.SourceFileName + "-" + bi.FuncName
			if bi.BaseErrorID != 0 {
				m += fmt.Sprintf(" ErrorID: %v", bi.BaseErrorID)
			}
		}
	}

	return m
}

// InitializeBaseInfo - Initializes a SpecErr Structure
// from a ParentInfo array and a ErrBaseInfo
// structure
func (s SpecErr) InitializeBaseInfo(parent []ErrBaseInfo, bi ErrBaseInfo) SpecErr {

	return SpecErr{
		ParentInfo: s.DeepCopyParentInfo(parent),
		BaseInfo:   bi.DeepCopyBaseInfo()}
}

// Initialize - Initializes all elements of
// the SpecErr structure
func (s SpecErr) Initialize(parent []ErrBaseInfo, bi ErrBaseInfo, prefix string, err error, isPanic bool, errNo int64) SpecErr {
	return s.InitializeBaseInfo(parent, bi).New(prefix, err, isPanic, errNo)

}

// New - Creates new SpecErr Type. Uses existing
// Parent and ErrBaseInfo data. The error is based on
// a parameter of type 'error' passed to the method.
//
// Note: If you set errNo == zero, no error number will be displayed in the
// in the error message.
func (s SpecErr) New(prefix string, err error, isPanic bool, errNo int64) SpecErr {

	x := SpecErr{
		ParentInfo: s.DeepCopyParentInfo(s.ParentInfo),
		BaseInfo:   s.BaseInfo.DeepCopyBaseInfo(),
		PrefixMsg:  prefix,
		IsPanic:    isPanic}

	if errNo != 0 {
		x.ErrNo = errNo + x.BaseInfo.BaseErrorID
	}


	if err != nil {
		x.ErrMsg = err.Error()
		x.IsErr = true
	} else {
		x.ErrMsg = ""
		x.IsErr = false
		x.IsPanic = false
	}

	return x
}

// NewErrorMsgString - Creates a new error message
// based on an error message string.
//
// Note: If you set errNo == zero, no error number will be displayed in the
// in the error message.
func (s SpecErr) NewErrorMsgString(prefix string, errMsg string, isPanic bool, errNo int64) SpecErr {
		er := errors.New(errMsg)

		return s.New(prefix, er, isPanic, errNo)
}

// Panic - Executes 'panic' command
// if IsPanic == 'true'
func (s *SpecErr) Panic() {
	if s.IsPanic {
		panic(s)
	}
}

// PanicOnSpecErr - Issues a 'panic'
// command if SpecErr IsPanic flag is set
func (s *SpecErr) PanicOnSpecErr(eSpec SpecErr) {

	if s.IsPanic {
		panic(s)
	}
}

// SignalNoErrors - Returns a SpecErr
// structure with IsErr set to false.
func (s SpecErr) SignalNoErrors() SpecErr {
	return SpecErr{IsErr: false, IsPanic: false}
}

// SetBaseInfo - Sets the SpecErr ErrBaseInfo internal
// structure. This data is used for creating repetitive
// error information.
func (s *SpecErr) SetBaseInfo(bi ErrBaseInfo) {
	s.BaseInfo = bi.NewBaseInfo()
}

// SetErrorLabel - If an Error Message Label is needed
// the Error message, set the value Error Message Label
// here.  This method merely sets the SpecErr string field,
// SpecErr.ErrMsgLabel. Of course this field can also be
// set directly with the use of this method.
//
// If the SpecErr.ErrMsgLabel is set to "StdOut Err", the
// error message will be formatted as :
// 						"StdOut Err: Your Error Message"
func (s *SpecErr) SetErrorMessageLabel(errorMsgLabel string) {
	s.ErrMsgLabel = errorMsgLabel
}
// SetParentInfo - Sets the ParentInfo Slice for
// the current SpecErr structure
func (s *SpecErr) SetParentInfo(parent []ErrBaseInfo) {
	if len(parent) == 0 {
		return
	}

	s.ParentInfo = s.DeepCopyParentInfo(parent)
}

var blankErrBaseInfo = ErrBaseInfo{}
var blankParentInfo = make([]ErrBaseInfo, 0, 10)

