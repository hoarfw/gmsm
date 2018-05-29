package ca

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// Program name
const ProgramName = "ca"

// Cmd returns the Cobra Command for Version
func versionCmd() *cobra.Command {
	return versionCommand
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print ca version.",
	Long:  `Print current version of the ca tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(GetInfo())
	},
}

// GetInfo returns version information for the peer
func GetInfo() string {

	ccinfo := fmt.Sprintf(" Version: %s\n", "1.0")

	return fmt.Sprintf("%s:\n  Go version: %s\n OS/Arch: %s\n"+
		" Tools:\n %s\n",
		ProgramName, runtime.Version(),
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		ccinfo)
}
