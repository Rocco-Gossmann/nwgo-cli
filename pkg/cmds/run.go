package cmds

import (
	"fmt"
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

		gocmd := exec.Command("go", "build", "-o", "backend")
		gocmd.Dir = args[0]

		out, err := gocmd.CombinedOutput()
		if err == nil {
			out, err = exec.Command(platform.Launch_file, args[0]).CombinedOutput()
			if err != nil {
				fmt.Println("nwjs launch error:\n--------------------------")
				fmt.Println(string(out))
				fmt.Println()
			}
		} else {
			fmt.Println("backend compilation error:\n--------------------------")
			fmt.Println(string(out))
			fmt.Println()
		}
	},
}
