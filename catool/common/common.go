package common

import (
	//"fmt"
	//"os"
	"strings"
	"github.com/op/go-logging"
	"errors"	
	"github.com/spf13/viper"
)

// UndefinedParamValue defines what undefined parameters in the command line will initialise to
const UndefinedParamValue = ""

	

func init() {
}



// SetLogLevelFromViper sets the log level for 'module' logger to the value in
// core.yaml
func SetLogLevelFromViper(module string) error {
	var err error
	if module == "" {
		return errors.New("log level not set, no module name provided")
	}
	logLevelFromViper := viper.GetString("logging." + module)
	err = CheckLogLevel(logLevelFromViper)
	if err != nil {
		return err
	}
	// replace period in module name with forward slash to allow override
	// of logging submodules
	module = strings.Replace(module, ".", "/", -1)
	// only set logging modules that begin with the supplied module name here
	//_, err = logging.SetModuleLevel("^"+module, logLevelFromViper)
	return err
}

// CheckLogLevel checks that a given log level string is valid
func CheckLogLevel(level string) error {
	_, err := logging.LogLevel(level)
	return err
}
