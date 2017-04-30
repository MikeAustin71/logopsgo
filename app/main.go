package main

import (
	common "MikeAustin71/logopsgo/common"
	"fmt"
	"time"
)

func main() {

	ini2()

}

func ini2() {
	log := common.LogJobGroupConfig{}
	log.NewLogGroupConfig(time.Now().Local())

	fmt.Println("AppLogPathFileName", log.AppLogPathFileName)

	fmt.Println("AppExeDir", log.AppExeDir)

	fmt.Println("AppLogBanner1", log.AppLogBanner1)
}
