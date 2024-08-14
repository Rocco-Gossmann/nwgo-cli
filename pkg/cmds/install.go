package cmds

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rocco-gossmann/go_utils"
	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

var installCommand cobra.Command = cobra.Command{
	Use:   "install",
	Short: "Installs the nwgo commant into your $GOPATH ",
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		var platform = pkg.GetPlatformConfig()

		fmt.Println("compiling")
		_, err = exec.Command("go", "mod", "tidy").Output()
		go_utils.Err(err)

		_, err = exec.Command("go", "build", "-o", BIN_FILE, ".").Output()
		go_utils.Err(err)

		_, err = os.Stat(BIN_FILE)
		go_utils.Err(err)

		input, err := os.Open(BIN_FILE)
		go_utils.Err(err)
		defer input.Close()
		defer os.Remove(BIN_FILE)

		baseFile := strings.ReplaceAll(absoluteBinPath, "/", platform.DirSeparator)
		go_utils.Err(go_utils.MkDir(filepath.Dir(baseFile)))

		output, err := os.OpenFile(baseFile, os.O_WRONLY|os.O_CREATE, 0777)
		go_utils.Err(err)
		defer output.Close()

		go_utils.CopyWithProgress(input, output, func(_ int) {
			fmt.Print(".")
		})

		fmt.Print("\nDone !!!\nYou should now be able to call 'nwgo' from everywhere on your system")

	},
}
