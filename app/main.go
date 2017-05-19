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
		parent := common.ErrBaseInfo{}.GetNewParentInfo(thisSrcFileName, "ini2", thisErrBlockNo)

	log := common.LogJobGroupConfig{}
	sErr := log.NewLogGroupConfig("CmdrX Ver 2.0", "CmdrX.xml",time.Now(), parent)

	if common.CheckIsSpecErr(sErr){
		fmt.Println(sErr.Error())
		return
	}

	fmt.Println("AppLogPathFileName", log.AppLogPathFileName)

	fmt.Println("AppLogBanner1", log.Banner1)
}
