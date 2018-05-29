package main

import (
	//"fmt"
	"github.com/hoarfw/gmsm/catool/ca"
	//"github.com/hoarfw/gmsm/catool/common"
)

func main() {

	// setup system-wide logging backend based on settings from core.yaml
	//logging.InitBackend(logging.SetFormat(viper.GetString("logging.format")), logOutput)
	ca.RootCmd()
	//logger.Info("Exiting.....")
}
