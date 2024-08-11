package pkg

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rocco-gossmann/go_utils"
)

func ZipExists(target string) bool {

	_, err := os.Stat(target)

	if os.IsNotExist(err) {
		return false
	}

	go_utils.Err(err)

	return true
}

func ExtractTarGZ(inFile string, dstPath string) (didStuff bool, err error) {

	file, err := os.Open(inFile)
	go_utils.Err(err)
	defer file.Close()

	gz, err := gzip.NewReader(file)
	go_utils.Err(err)

	tarReader := tar.NewReader(gz)

	var item string

	for {
		tarItem, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		go_utils.Err(err)

		item = fmt.Sprintf("%s%c%s", dstPath, os.PathSeparator, tarItem.Name)

		switch tarItem.Typeflag {
		case tar.TypeDir:
			go_utils.Err(go_utils.MkDir(item))

		case tar.TypeReg:
			dir := filepath.Dir(item)
			go_utils.Err(go_utils.MkDir(dir))

			func() {
				_, err := os.Stat(item)
				if !os.IsNotExist(err) {
					return
				}

				fmt.Print(CLEAR_CMD_LINE, "Extract: ", tarItem.Name)
				dstFile, err := os.OpenFile(item, os.O_CREATE|os.O_WRONLY, 0777)
				go_utils.Err(err)
				defer dstFile.Close()

				_, err = io.Copy(dstFile, tarReader)
				go_utils.Err(err)

				didStuff = true
			}()
		}
	}

	if didStuff {
		fmt.Print(CLEAR_CMD_LINE, "Extract: done !!!\n")
	}
	return
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
