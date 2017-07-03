package main

import (
	"MikeAustin71/logopsgo/common"
	"fmt"
	"time"
)

const (
	thisSrcFileName = "main.go"
	thisErrBlockNo  = int64(1000)
)

func main() {

	ini2()

}

func ini2() {

	tz := common.TimeZoneUtility{}
	parms := common.StartupParameters{}
	parms.IanaTimeZone = "Local"
	tzLocal, _ := tz.ConvertTz(time.Now().UTC(), parms.IanaTimeZone )
	dtf := common.DateTimeFormatUtility{}
	dtf.CreateAllFormatsInMemory()

	parms.StartTimeUTC = time.Now().UTC()
	parms.StartTime = tzLocal.TimeOut
	parms.AppVersion = "2.0.0"
	parms.LogMode = common.LogVERBOSE
	parms.AppLogPath = "./cmdrX"
	parms.AppName = "cmdrX"
	parms.AppExeFileName = "cmdrX.exe"
	parms.NoOfJobs = 37
	parms.CommandFileName = "cmdrX.xml"
	parms.Dtfmt = &dtf
	cmdPathFileName := "D:/go/work/src/MikeAustin71/logopsgo/common/cmdrXCmds.xml"
	fh, err := parms.AssembleAppPath(cmdPathFileName)

	if err != nil {
		panic(err)
	}

	if !fh.AbsolutePathIsPopulated {
		panic(fmt.Errorf("Failed to lcoate Command File %v", cmdPathFileName))
	}
	parms.AppPath = fh.AbsolutePath
	parms.BaseStartDir = fh.AbsolutePath

	parent := common.ErrBaseInfo{}.GetNewParentInfo(thisSrcFileName, "ini2", thisErrBlockNo)

	lg := common.LogJobGroup{}

	se := lg.New(parms, parent)

	if common.CheckIsSpecErr(se) {
		fmt.Println(se.Error())
		return
	}

	dur := common.DurationUtility{}


	time.Sleep(dur.GetDurationFromSeconds(10))

	lg.WriteJobGroupFooterToLog(parent)

	fmt.Println("AppLogPathFileName", lg.AppLogPathFileName)

	fmt.Println("AppLogBanner1", lg.Banner1)
}
