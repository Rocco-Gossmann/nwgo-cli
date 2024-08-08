package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/rocco-gossmann/go_utils"
	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

func ifExists(path string, dirsep string, action func(string) error) {
	localPath := strings.ReplaceAll(path, "/", dirsep)
	_, err := os.Stat(localPath)
	if !os.IsNotExist(err) {
		go_utils.Err(action(localPath))
	}
}

var uninstallCommand cobra.Command = cobra.Command{
	Use:   "uninstall",
	Short: "removes every system file created through nwgo/nwgo-cli",
	Run: func(cmd *cobra.Command, args []string) {

		platform := pkg.GetPlatformConfig()

		rmFile := func(s string) error {
			fmt.Println("rm file ", s)
			return os.Remove(s)
		}

		ifExists(absoluteBinPath, platform.DirSeparator, rmFile)

		ifExists(pkg.NWGO_BASEPATH, platform.DirSeparator, func(s string) error {
			fmt.Println("rm dir ", s)
			return os.RemoveAll(s)
		})

		fmt.Print("Done !!!\n\n")
	},
}
