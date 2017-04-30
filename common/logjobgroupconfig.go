package common

import (
	fp "path/filepath"
	"strings"
	"time"
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
	AppLogBanner1      string
	AppLogBanner2      string
	AppLogBanner3      string
	AppLogBanner4      string
	AppExeDir          string
	AppNoOfJobs        int
}

// NewLogGroupConfig - Initializes key
// elements of a Logging Configuration
func (logCfg *LogJobGroupConfig) NewLogGroupConfig(t time.Time) {

	var err error
	logCfg.LogMode = LogVERBOSE
	logCfg.StartTime = t
	logCfg.AppLogFileName = "CmdrX" + "_" + GetDateTimeStr(logCfg.StartTime) + ".log"
	logCfg.AppLogDir, err = MakeAbsolutePath("./")
	if err != nil {
		panic(err)
	}

	logCfg.AppLogPathFileName = fp.Join(logCfg.AppLogDir, logCfg.AppLogFileName)
	logCfg.AppLogBanner1 = strings.Repeat("=", 79)
	logCfg.AppLogBanner2 = strings.Repeat("#", 79)
	logCfg.AppLogBanner3 = strings.Repeat("*", 79)
	logCfg.AppLogBanner4 = strings.Repeat("-", 79)

	logCfg.AppExeDir = getExeDir()

}

func getExeDir() string {

	ex, err1 := GetExecutablePathFileName()

	if err1 != nil {
		panic(err1)
	}

	p, err2 := GetAbsPathFromFilePath(ex)

	if err2 != nil {
		panic(err2)
	}

	return p
}
