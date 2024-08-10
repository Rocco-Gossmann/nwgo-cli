package cmds

import (
	"archive/zip"
	"fmt"
	"os"
	"os/exec"

	"github.com/rocco-gossmann/go_fileutils"
	"github.com/rocco-gossmann/go_utils"
	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

var buildCommand cobra.Command = cobra.Command{
	Use:   "build projectPath",
	Args:  cobra.ExactArgs(1),
	Short: "builds a runnable executable ob your project",
	Run: func(cmd *cobra.Command, args []string) {

		var platform = pkg.SetupPlatform(pkg.SetupPlatformOpts{Production: true})

		gocmd := exec.Command("go", "build", "-o", platform.BackendBinary)
		gocmd.Dir = args[0]
		go_utils.Err(gocmd.Run())

		// DONE: get full Buildpath
		var fullBuildPath = args[0] + "/build"

		// done: empty build dir (in case it existed beforehand)
		go_utils.Err(os.RemoveAll(fullBuildPath))

		// DONE: create output build dir
		go_utils.Err(go_utils.MkDir(fullBuildPath))

		// DONE: copy Platform Runtime to builddir

		progressChan := go_fileutils.CopyRecursive(platform.Extract_build_target, fullBuildPath, platform.BuildCutPath)
		//copyRecursive(platform.Extract_build_target, fullBuildPath, platform.BuildCutPath)

	copy_runtime:
		for {
			var progress go_fileutils.BatchProgress = <-progressChan

			switch progress.State {
			case go_fileutils.STATE_START_FILE:
				fmt.Print(go_utils.CLEAR_CMD_LINE_SEQ, "copy platform: ", progress.CurrentTarget)

			case go_fileutils.STATE_FINISHED:
				fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "copy platform: done !!!")
				break copy_runtime

			case go_fileutils.STATE_ERROR:
				panic(progress.Error)

			}
		}

		// TODO: package everything into a zip
		zipFile, err := os.Create(fmt.Sprintf("%s%cout.zip", fullBuildPath, os.PathSeparator))
		go_utils.Err(err)
		defer zipFile.Close()

		zipArchive := zip.NewWriter(zipFile)
		defer zipArchive.Close()

		fOut, err := zipArchive.Create(platform.BackendBinary)
		go_utils.Err(err)

		fIn, err := os.Open(platform.BackendBinary)
		go_utils.Err(err)
		defer fIn.Close()

		go_utils.CopyWithProgress(fIn, fOut, func(bytesCopied int) {
			fmt.Print(go_utils.CLEAR_CMD_LINE_SEQ, "packing backend: ", bytesCopied)
		})

		fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "packaging backend: done!!!")

		// TODO:
		// - backend
		// - index.html
		// - package.json
		// - every none *.go or *.tmpl file

		// TODO: rename .zip to .nw

		// TODO: Put .nw file into build directory

		// TODO: link nw and runtime into one executable

	},
}
