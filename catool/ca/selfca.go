package ca

import (
	//	"errors"

	"github.com/spf13/cobra"
)

const selfcaCmdDescription = "Generate ca."

func selfcaCmd() *cobra.Command {
	// Set the flags on the channel start command.
	certCmd := &cobra.Command{
		Use:   "selfca",
		Short: selfcaCmdDescription,
		Long:  selfcaCmdDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return selfca(cmd, args)
		},
	}

	// flagList := []string{
	// 	"reqfile",
	// 	"cakey",
	// }
	//attachFlags(certCmd, flagList)

	return certCmd
}

func executeSelfCa() (err error) {

	logger.Info("Successfully sign the request")
	return nil
}

func selfca(cmd *cobra.Command, args []string) error {

	return executeSelfCa()
}
