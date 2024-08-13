package cmds

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/rocco-gossmann/go_utils"
	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

var runCommand cobra.Command = cobra.Command{
	Use:   "run projectPath [args...]",
	Short: "runs a NW.JS Project",
	Run: func(cmd *cobra.Command, args []string) {
		var platform = pkg.SetupPlatform(pkg.SetupPlatformOpts{})
		var path = "."

		if len(args) > 0 {
			path = args[0]
		}

		path, err := filepath.Abs(path)
		go_utils.Err(err)

		gocmd := exec.Command("go", "build", "-o", platform.BackendBinary)
		gocmd.Dir = path

		out, err := gocmd.CombinedOutput()
		if err == nil {
			out, err = exec.Command(platform.Launch_file, path).CombinedOutput()
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
