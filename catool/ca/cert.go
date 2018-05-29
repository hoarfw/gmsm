package ca

import (
	//	"errors"
	"fmt"
	//	"io/ioutil"
	"github.com/hoarfw/gmsm/catool/logging"
	"github.com/spf13/cobra"
)

var logger = logging.MustGetLogger("cert")

const commandDescription = "sign a req with ca."

func certCmd() *cobra.Command {
	// Set the flags on the channel start command.
	certCmd := &cobra.Command{
		Use:   "cert",
		Short: commandDescription,
		Long:  commandDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cert(cmd, args)
		},
	}

	// flagList := []string{
	// 	"reqfile",
	// 	"cakey",
	// }


	return certCmd
}

type FileNotFoundErr string

func (e FileNotFoundErr) Error() string {
	return fmt.Sprintf(" file not found %s", string(e))
}

func executeCert() (err error) {

	logger.Info("Successfully sign the request")
	return nil
}

func cert(cmd *cobra.Command, args []string) error {

	return executeCert()
}
