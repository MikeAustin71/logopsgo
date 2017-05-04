package common

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	srcFileThis      = "logjobgroupconfig.go"
	loggrouperrblock = int64(647000000)
)

// LogJobGroupConfig - holds logging configuration for the
// current group of jobs
type LogJobGroupConfig struct {
	LogMode            LoggingMode
	StartTime          time.Time
	EndTime            time.Time
	Duration           time.Time
	AppLogFileName     string
	AppLogDir          string
	AppLogPathFileName string
	AppLogFile         *os.File
	AppLogBanner1      string
	AppLogBanner2      string
	AppLogBanner3      string
	AppLogBanner4      string
	BaseStartDir       string
	AppNoOfJobs        int
}

// NewLogGroupConfig - Initializes key
// elements of a Logging Configuration
func (logCfg *LogJobGroupConfig) NewLogGroupConfig(t time.Time) SpecErr {

	funcName := "NewLogGroupConfig()"
	var err error
	logCfg.LogMode = LogVERBOSE
	logCfg.StartTime = t
	logCfg.AppLogFileName = "CmdrX" + "_" + GetDateTimeStr(logCfg.StartTime) + ".log"
	logCfg.AppLogDir, err = MakeAbsolutePath("./CmdrX")
	if err != nil {
		s := "MakeAbsolutePathFailed for './CmdrX'"
		return setError(true, s, err, funcName, loggrouperrblock+101)
	}

	logCfg.AppLogPathFileName = JoinPathsAdjustSeparators(logCfg.AppLogDir, logCfg.AppLogFileName)
	logCfg.AppLogBanner1 = strings.Repeat("=", 79)
	logCfg.AppLogBanner2 = strings.Repeat("#", 79)
	logCfg.AppLogBanner3 = strings.Repeat("*", 79)
	logCfg.AppLogBanner4 = strings.Repeat("-", 79)

	logCfg.BaseStartDir, err = GetAbsCurrDir()

	if err != nil {
		s := "GetAbsCurrDir() Failed!"
		return setError(true, s, err, funcName, loggrouperrblock+102)
	}

	return SpecErr{IsErr: false}
}

// InitializeLogFile - Initializes the log directory and
// a new log file.
func (logCfg *LogJobGroupConfig) InitializeLogFile() SpecErr {
	var err error
	funcName := "(logCfg *LogJobGroupConfig) InitializeLogFile()"

	if !DoesFileExist(logCfg.AppLogDir) {

		err = MakeDirAll(logCfg.AppLogDir)

		if err != nil {
			s := fmt.Sprintf("MakeDirAll() Failed on Dir: %v", logCfg.AppLogDir)
			return setError(true, s, err, funcName, loggrouperrblock+201)
		}

	}

	// At this point, logCfg.AppLogDir exists.
	// Check to determine if the log file exists.
	// If log file already exists, delete it.
	if DoesFileExist(logCfg.AppLogPathFileName) {
		err = DeleteDirFile(logCfg.AppLogPathFileName)
		if err != nil {
			s := fmt.Sprintf("DeleteDirFile() Failed on File: %v", logCfg.AppLogPathFileName)
			return setError(true, s, err, funcName, loggrouperrblock+202)
		}
	}

	// Create a new log file
	logCfg.AppLogFile, err = CreateFile(logCfg.AppLogPathFileName)

	if err != nil {
		s := fmt.Sprintf("CreateFile() Failed on File: %v", logCfg.AppLogPathFileName)
		return setError(true, s, err, funcName, loggrouperrblock+203)
	}

	return SpecErr{IsErr: false}
}

func setError(isPanic bool, prefix string, err error, funcName string, errNo int64) SpecErr {
	var s SpecErr

	return s.New(prefix, err, isPanic, srcFileThis, funcName, errNo)
}

func writeFileGroupHeaderToLog(logCfg *LogJobGroupConfig) SpecErr {
	funcName := "writeFileGroupHeaderToLog()"

	if logCfg.AppLogFile == nil {
		s := "logCfg.AppLogFile was not correctly initialized!"
		err := errors.New(" logCfg.AppLogFile *os.File pointer is nil")
		return setError(true, s, err, funcName, loggrouperrblock+301)
	}

	return SpecErr{IsErr: false}
}
