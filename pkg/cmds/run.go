package cmds

import (
	"os/exec"

	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

var runCommand cobra.Command = cobra.Command{
	Use:   "run projectPath [args...]",
	Args:  cobra.ExactArgs(1),
	Short: "runs a NW.JS Project",
	Run: func(cmd *cobra.Command, args []string) {
		var platform = pkg.SetupPlatform()
		exec.Command(platform.Launch_file, args[0]).Output()
	},
}
