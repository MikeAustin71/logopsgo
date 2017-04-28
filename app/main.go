package main

import (
	common "MikeAustin71/logopsgo/common"
	"fmt"
	"time"
)

func main() {

	var log common.LogJobGroupConfig

	lg := log.NewLogGroupConfig(time.Now().Local())

	fmt.Println("AppLogPathFileName", lg.AppLogPathFileName)

	fmt.Println("AppExeDir", lg.AppExeDir)

	fmt.Println("AppLogBanner1", lg.AppLogBanner1)

}
