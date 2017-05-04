package common

import (
	"errors"
	"testing"
)

func TestErrorUtility(t *testing.T) {

	ex1 := "errprefix"
	ex2 := "errutil_test.go"
	ex3 := "TestErrorUtility"
	ex4 := int64(334)
	ex99 := "Test Error #1"

	err := errors.New(ex99)

	var se SpecErr

	x := se.New(ex1, err, true, ex2, ex3, ex4)

	if x.PrefixMsg != ex1 {
		t.Error("Expected 'errprefix' got", x.PrefixMsg)
	}

	if x.ErrMsg != ex99 {
		t.Error("Expected 'Test Error #1' got", x.ErrMsg)
	}

	if x.SrcFile != ex2 {
		t.Error("Expected 'errutil_test.go' got", x.SrcFile)
	}

	if x.FuncName != ex3 {
		t.Error("Expected 'TestErrorUtility' got", x.FuncName)
	}

	if x.ErrNo != ex4 {
		t.Error("Expected '334' got", x.ErrNo)
	}

}

func TestNoSetError(t *testing.T) {
	var se SpecErr

	s := se.SetNoError()

	if s.IsErr == true {
		t.Error("SetNoError() Expected false got", s.IsErr)
	}
}

func TestIsSpecErrNo(t *testing.T) {
	var se SpecErr
	s := se.SetNoError()

	isErr := CheckIsSpecErr(s)

	if isErr {
		t.Error("Expected CheckIsSpecErr() to return false, go", isErr)
	}

}

func TestIsSpecErrYes(t *testing.T) {
	ex1 := "errprefix"
	ex2 := "errutil_test.go"
	ex3 := "TestErrorUtility"
	ex4 := int64(334)
	ex99 := "Test Error #1"

	err := errors.New(ex99)

	var se SpecErr

	x := se.New(ex1, err, true, ex2, ex3, ex4)

	isErr := CheckIsSpecErr(x)

	if !isErr {
		t.Error("Expected CheckIsSpecErr() to return true, go", isErr)
	}

}
