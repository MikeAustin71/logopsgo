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

	parms := common.LogStartupParameters{}

	parms.StartTime = time.Now()
	parms.AppVersion = "2.0.0"
	parms.LogMode = common.LogVERBOSE
	parms.AppLogDir = "./CmdrX"
	parms.AppName = "CmdrX"
	parms.AppExeFileName = "CmdrX.exe"
	parms.NoOfJobs = 37
	parms.CommandFileName = "CmdrX.xml"

	parent := common.ErrBaseInfo{}.GetNewParentInfo(thisSrcFileName, "ini2", thisErrBlockNo)

	lg := common.LogJobGroupConfig{}

	se := lg.New(parms, parent)

	if common.CheckIsSpecErr(se) {
		fmt.Println(se.Error())
		return
	}

	fmt.Println("AppLogPathFileName", lg.AppLogPathFileName)

	fmt.Println("AppLogBanner1", lg.Banner1)
}
