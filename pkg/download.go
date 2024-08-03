package pkg

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/rocco-gossmann/go_utils"
)

func DownloadFile(url string, outputPath string, label string) (err error) {

	rawFile, err := os.Create(outputPath)
	go_utils.Err(err)
	defer rawFile.Close()

	res, err := http.Get(url)
	go_utils.Err(err)
	defer res.Body.Close()

	var totalBytes int64
	totalBytes, err = strconv.ParseInt(res.Header.Get("Content-Length"), 10, 64)
	go_utils.Err(err)

	fmt.Print("Downloading ", label, ":\nstart")

	go_utils.CopyWithProgress(res.Body, rawFile, func(readBytes int) {
		fmt.Printf("%s--- %d of %d bytes read --- ( %.2f %% )", CLEAR_CMD_LINE, readBytes, totalBytes, float64(readBytes)/float64(totalBytes)*100.0)
	})

	fmt.Println(CLEAR_CMD_LINE, "Done !!!")

	return
}
