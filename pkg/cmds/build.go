package cmds

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gfu "github.com/rocco-gossmann/go_fileutils"
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

func submitErr(err error, progressChan chan gfu.BatchProgress) error {
	progressChan <- gfu.BatchProgress{
		Error: err,
		State: gfu.STATE_ERROR,
	}

	return err
}
func packFolderContent(zipArchive *zip.Writer, srcFolder, targetFolder string, onProgress func(progress gfu.BatchProgress)) error {

	var progressChan = make(chan gfu.BatchProgress)

	root, err := filepath.Abs(srcFolder)
	go_utils.Err(err)

	go func() {

		progress := gfu.BatchProgress{}

		filepath.WalkDir(srcFolder, func(path string, d fs.DirEntry, err error) error {

			if d.IsDir() {
				return nil
			}

			cutPath, _ := strings.CutPrefix(path, root)
			dstPath := fmt.Sprintf("%s%s", targetFolder, cutPath)

			progress.CurrentSource = path
			progress.CurrentTarget = dstPath

			progress.State = gfu.STATE_START_FILE
			progressChan <- progress

			if err := packIntoZip(zipArchive, path, dstPath); err != nil {
				return submitErr(err, progressChan)
			}

			progress.State = gfu.STATE_END_FILE
			progressChan <- progress

			return nil

		})

		progress.State = gfu.STATE_FINISHED
		progressChan <- progress

	}()

	for {
		progress := <-progressChan
		if progress.State == gfu.STATE_FINISHED {
			break

		} else if progress.State == gfu.STATE_ERROR {
			return progress.Error

		} else {
			onProgress(progress)

		}
	}

	return nil
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

		progressChan := gfu.CopyRecursive(platform.Extract_build_target, fullBuildPath, platform.BuildCutPath)

	copy_runtime:
		for {
			var progress gfu.BatchProgress = <-progressChan

			switch progress.State {
			case gfu.STATE_START_FILE:
				fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "copy platform: ", progress.CurrentTarget)

			case gfu.STATE_COPY:
				fmt.Print(go_utils.CLEAR_CMD_LINE_SEQ, progress.BytesCopied, " of ", progress.BytesTotal)

			case gfu.STATE_FINISHED:
				fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "copy platform: done !!!")
				break copy_runtime

			case gfu.STATE_ERROR:
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

		var staticDir = fullBasePath + "static"
		_, err = os.Stat(staticDir)

		if !os.IsNotExist(err) {
			fmt.Println()
			go_utils.Err(packFolderContent(zipArchive, staticDir, "static", func(progress gfu.BatchProgress) {
				switch progress.State {
				case gfu.STATE_START_FILE:
					fmt.Print("packing: ", progress.CurrentTarget)

				case gfu.STATE_END_FILE:
					fmt.Println(" => done !!!")

				}
			}))
		}

		fmt.Println(go_utils.CLEAR_CMD_LINE_SEQ, "packaging: done!!!")

		go_utils.Err(zipArchive.Close())

		zipFile.Seek(0, 0)

		if platform.PostBuild != nil {
			//DONE: Finalize Platform specific build
			platform.PostBuild(platform, zipFile, fullBuildPath)
		}
	},
}
