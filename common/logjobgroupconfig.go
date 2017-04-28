package common

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// LogJobGroupConfig - holds logging configuration for the
// current group of jobs
type LogJobGroupConfig struct {
	LogMode            LoggingMode
	StartTime          time.Time
	EndTime            time.Time
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
func (logCfg LogJobGroupConfig) NewLogGroupConfig(t time.Time) *LogJobGroupConfig {

	var err error
	logCfg.LogMode = LogVERBOSE
	logCfg.StartTime = t
	logCfg.AppLogFileName = "CmdrX" + "_" + GetDateTimeStr(logCfg.StartTime) + ".log"
	logCfg.AppLogDir, err = filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	logCfg.AppLogPathFileName = logCfg.AppLogDir + logCfg.AppLogFileName
	logCfg.AppLogBanner1 = strings.Repeat("=", 79)
	logCfg.AppLogBanner2 = strings.Repeat("#", 79)
	logCfg.AppLogBanner3 = strings.Repeat("*", 79)
	logCfg.AppLogBanner4 = strings.Repeat("-", 79)
	ex, err2 := os.Executable()
	if err2 != nil {
		panic(err2)
	}

	logCfg.AppExeDir = path.Dir(ex)

	return &logCfg
}
