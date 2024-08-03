package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rocco-gossmann/go_utils"
)

func ZipExists(target string) bool {

	zipfile, err := zip.OpenReader(target)

	if err != nil {
		return false
	} else {
		zipfile.Close()
		return true
	}
}

func ExtractZip(inFile string, dstPath string) (didStuff bool, err error) {

	reader, err := zip.OpenReader(inFile)
	go_utils.Err(err)
	defer reader.Close()

	for _, file := range reader.File {

		var item string = fmt.Sprintf("%s/%s", dstPath, file.Name)

		_, err := os.Stat(item)
		if !os.IsNotExist(err) {
			continue
		}

		fmt.Print(CLEAR_CMD_LINE, "Extract: ", file.FileInfo().Name())

		if file.FileInfo().IsDir() {
			go_utils.Err(go_utils.MkDir(item))

		} else {

			dir := filepath.Dir(item)
			go_utils.Err(go_utils.MkDir(dir))

			dstFile, err := os.OpenFile(item, os.O_CREATE|os.O_WRONLY, file.Mode())
			go_utils.Err(err)
			defer dstFile.Close()

			zipFile, err := file.Open()
			go_utils.Err(err)
			defer zipFile.Close()

			_, err = io.Copy(dstFile, zipFile)
			go_utils.Err(err)

			didStuff = true
		}
	}

	if didStuff {
		fmt.Print(CLEAR_CMD_LINE, "Extract: Done !!!\n")
	}

	return

}
