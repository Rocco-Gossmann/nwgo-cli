package pkg

import (
	"fmt"
	"os"

	"github.com/rocco-gossmann/go_utils"
)

const (
	NWGO_BASEDIR   = ".local/state/nwgo"
	CLEAR_CMD_LINE = go_utils.CLEAR_CMD_LINE_SEQ
)

var NWGO_BASEPATH = NWGO_BASEDIR

func init() {
	home, err := os.UserHomeDir()
	go_utils.Err(err)

	NWGO_BASEPATH = fmt.Sprintf("%s/%s", home, NWGO_BASEDIR)
}
