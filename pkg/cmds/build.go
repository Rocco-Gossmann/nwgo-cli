package cmds

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rocco-gossmann/go_fileutils"
	"github.com/rocco-gossmann/go_utils"
	"github.com/rocco-gossmann/nwgo-cli/pkg"
	"github.com/spf13/cobra"
)

func packIntoZip(zipArchive *zip.Writer, file, targetInZip string) (err error) {
	var fOut io.Writer
	if fOut, err = zipArchive.Create(targetInZip); err != nil {
		return err
	}

	var fIn *os.File
	if fIn, err = os.Open(file); err != nil {
		return err
	}
	defer fIn.Close()

	_, err = go_utils.CopyWithProgress(fIn, fOut, func(bytesCopied int) {
		fmt.Print(go_utils.CLEAR_CMD_LINE_SEQ, "packing ", targetInZip, ": ", bytesCopied)
	})

	return
}

var buildCommand cobra.Command = cobra.Command{
	Use:   "build projectPath",
	Args:  cobra.ExactArgs(1),
	Short: "builds a runnable executable ob your project",
	Run: func(cmd *cobra.Command, args []string) {

		var platform = pkg.SetupPlatform(pkg.SetupPlatformOpts{Production: true})

		var fullBasePath, err = filepath.Abs(args[0])
		go_utils.Err(err)

		fullBasePath += string(os.PathSeparator)
		var fullBuildPath = fullBasePath + "build" + string(os.PathSeparator)

		// Build Backend
		gocmd := exec.Command("go", "build", "-o", fullBasePath+platform.BackendBinary)
		gocmd.Dir = args[0]
		go_utils.Err(gocmd.Run())

		go_utils.Err(os.RemoveAll(fullBuildPath))
		go_utils.Err(go_utils.MkDir(fullBuildPath))

		progressChan := go_fileutils.CopyRecursive(platform.Extract_build_target, fullBuildPath, platform.BuildCutPath)

	copy_runtime:
		for {
			var progress go_fileutils.BatchProgress = <-progressChan

			switch progress.State {
			case go_fileutils.STATE_START_FILE:
				fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "copy platform: ", progress.CurrentTarget)

			case go_fileutils.STATE_COPY:
				fmt.Print(go_utils.CLEAR_CMD_LINE_SEQ, progress.BytesCopied, " of ", progress.BytesTotal)

			case go_fileutils.STATE_FINISHED:
				fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "copy platform: done !!!")
				break copy_runtime

			case go_fileutils.STATE_ERROR:
				panic(progress.Error)

			}
		}

		tmpFileName := os.TempDir() + string(os.PathSeparator) + "nwgoapp_build.zip"
		err = os.Remove(tmpFileName)
		if err != nil && !os.IsNotExist(err) {
			panic(err)
		}

		zipFile, err := os.Create(tmpFileName)
		go_utils.Err(err)
		defer zipFile.Close()

		zipArchive := zip.NewWriter(zipFile)
		defer zipArchive.Close()

		os.Chmod(fullBasePath+platform.BackendBinary, 0755)

		go_utils.Err(packIntoZip(zipArchive, fullBasePath+platform.BackendBinary, platform.BackendBinary))
		go_utils.Err(packIntoZip(zipArchive, fullBasePath+"index.html", "index.html"))
		go_utils.Err(packIntoZip(zipArchive, fullBasePath+"package.json", "package.json"))

		fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "packaging: done!!!")

		go_utils.Err(zipArchive.Close())

		zipFile.Seek(0, 0)

		if platform.PostBuild != nil {
			//DONE: Finalize Platform specific build
			platform.PostBuild(platform, zipFile, fullBuildPath)
		}
	},
}
