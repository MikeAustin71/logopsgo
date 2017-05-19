package main

import (
	common "MikeAustin71/logopsgo/common"
	"fmt"
	"time"
)

const (
	mainSrcFileName = "main.go"
	mainErrBlockNo = int64(1000)
)
func main() {

	ini2()

}

func ini2() {
		parent := common.ErrBaseInfo{}.GetNewParentInfo(mainSrcFileName, "ini2", mainErrBlockNo )

	log := common.LogJobGroupConfig{}
	sErr := log.NewLogGroupConfig(time.Now().Local(), parent)

	if common.CheckIsSpecErr(sErr){
		fmt.Println(sErr.Error())
		return
	}

	fmt.Println("AppLogPathFileName", log.AppLogPathFileName)

	fmt.Println("AppLogBanner1", log.AppLogBanner1)
}
