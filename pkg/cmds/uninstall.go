package cmds

import (
	"fmt"
	"os"

	"github.com/rocco-gossmann/go_utils"
	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

func ifExists(path string, action func(string) error) {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		go_utils.Err(action(path))
	}
}

var uninstallCommand cobra.Command = cobra.Command{
	Use:   "uninstall",
	Short: "removes every system file created through nwgo/nwgo-cli",
	Run: func(cmd *cobra.Command, args []string) {

		rmFile := func(s string) error {
			fmt.Println("rm file ", s)
			return os.Remove(s)
		}

		ifExists(absoluteBinPath, rmFile)
		ifExists(BIN_FILE, rmFile)

		ifExists(pkg.NWGO_BASEPATH, func(s string) error {
			fmt.Println("rm dir ", s)
			return os.RemoveAll(s)
		})

		fmt.Println("Done !!!\n")
	},
}
