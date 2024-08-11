package pkg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rocco-gossmann/go_fileutils"
	"github.com/rocco-gossmann/go_utils"
)

var platformMap = map[string]PlatformEnv{

	"darwin_arm64": {
		Download_build: "https://dl.nwjs.io/v0.90.0/nwjs-v0.90.0-osx-arm64.zip",
		BuildCutPath:   "nwjs-v0.90.0-osx-arm64/",
		//Download_build:       "https://dl.nwjs.io/v0.90.0/nwjs-sdk-v0.90.0-osx-arm64.zip",
		//BuildCutPath:         "nwjs-sdk-v0.90.0-osx-arm64/",
		Download_sdk:         "https://dl.nwjs.io/v0.90.0/nwjs-sdk-v0.90.0-osx-arm64.zip",
		Download_target:      "nwjs.mac.zip",
		Download_target_prod: "nwjs.mac.prod.zip",
		Extract_sdk_target:   "extract.mac.sdk",
		Extract_build_target: "extract.mac.build",
		Launch_file:          "nwjs-sdk-v0.90.0-osx-arm64/nwjs.app/Contents/MacOS/nwjs",
		BuildPreRunJS:        `spawn('chmod', [ '+x', './backend' ])`,
		PostSetup:            MacPostSetup,
		Extractor:            ExtractZip,
		DirSeparator:         "/",
		BackendBinary:        "backend",
		PostBuild:            moveAppToPackage("nwjs.app/Contents/Resources/app.nw"),
	},

	"linux_amd64": {
		Download_build:       "https://dl.nwjs.io/v0.90.0/nwjs-v0.90.0-linux-x64.tar.gz",
		Download_sdk:         "https://dl.nwjs.io/v0.90.0/nwjs-sdk-v0.90.0-linux-x64.tar.gz",
		Download_target:      "nwjs.linux_x86.tar.gz",
		Download_target_prod: "nwjs.linux_x86.prod.tar.gz",
		Extract_sdk_target:   "extract.linux_x86.sdk",
		Extract_build_target: "extract.linux_x86.build",
		Launch_file:          "nwjs-sdk-v0.90.0-linux-x64/nw",
		Extractor:            ExtractTarGZ,
		DirSeparator:         "/",
		BackendBinary:        "backend",
		BuildPreRunJS:        `spawn('chmod', [ '+x', './backend' ])`,
		PostBuild:            moveAppToPackage("nwjs-v0.90.0-linux-x64/package.nw"),
	},

	"windows_amd64": {
		Download_build:       "https://dl.nwjs.io/v0.90.0/nwjs-v0.90.0-win-x64.zip",
		Download_sdk:         "https://dl.nwjs.io/v0.90.0/nwjs-sdk-v0.90.0-win-x64.zip",
		Download_target:      "nwjs.win.zip",
		Download_target_prod: "nwjs.win.prod.zip",
		Extract_sdk_target:   "extract.win_x64.sdk",
		Extract_build_target: "extract.win.build",
		Launch_file:          "nwjs-sdk-v0.90.0-win-x64\\nw.exe",
		Extractor:            ExtractZip,
		DirSeparator:         "\\",
		BackendBinary:        "backend.exe",
		BackendBinarySlash:   "\\",
		PostBuild:            moveAppToPackage("nwjs-v0.90.0-win-x64\\package.nw"),
	},
}

func moveAppToPackage(targetFile string) func(PlatformEnv, string) {
	return func(pl PlatformEnv, buildPath string) {
		progChan := go_fileutils.CopyFile(
			buildPath+"app.zip",
			buildPath+targetFile,
		)

	copy_loop:
		for {
			prog := <-progChan
			switch prog.State {
			case go_fileutils.STATE_START_FILE:
				fmt.Println("copying ", prog.BytesTotal, " bytes")

			case go_fileutils.STATE_COPY:
				fmt.Print(go_utils.CLEAR_CMD_LINE_SEQ, "copy => ", prog.BytesCopied, " to go")

			case go_fileutils.STATE_ERROR:
				panic(prog.Error)

			case go_fileutils.STATE_END_FILE:
				break copy_loop
			}
		}

		os.Remove(buildPath + "app.zip")
	}
}

// Mac Specific functions
// ====================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================================
func MacPostSetup(pe PlatformEnv) {
	exec.Command("xattr", "-cr", "nwjs-sdk-v0.90.0-osx-arm64/nwjs.app").Output()
}
