package cmds

import (
	"fmt"
	"go/build"
	"os"

	"github.com/spf13/cobra"
)

const BIN_FILE = "./nwgo-cli"
const BIN_DEST = "nwgo"

var absoluteBinPath = BIN_DEST

var cobraHead = cobra.Command{

	Use: "nwgo command",

	RunE: func(cmd *cobra.Command, args []string) error {
		//		if cmd.Flag("version").Value.String() == "true" {
		//			fmt.Println(env.Version)
		//			return nil
		//		}

		return fmt.Errorf("nothing to do")
	},

	DisableFlagsInUseLine: true,
}

func init() {

	path := os.Getenv("GOPATH")
	if path == "" {
		path = build.Default.GOPATH
	}

	absoluteBinPath = fmt.Sprintf("%s/bin/%s", path, BIN_DEST)

	cobraHead.AddCommand(
		&installCommand,
		&initCommand,
		&runCommand,
		&uninstallCommand,
	)
}

func LetsGo() {
	cobraHead.Execute()
}
