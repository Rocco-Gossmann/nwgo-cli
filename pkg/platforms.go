package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/rocco-gossmann/go_utils"
)

var platformMap = map[string]PlatformEnv{

	"darwin_arm64": {
		Download_build:       "https://dl.nwjs.io/v0.90.0/nwjs-v0.90.0-osx-arm64.zip",
		Download_sdk:         "https://dl.nwjs.io/v0.90.0/nwjs-sdk-v0.90.0-osx-arm64.zip",
		Download_target:      "nwjs.mac.zip",
		Extract_sdk_target:   "extract.mac.sdk",
		Extract_build_target: "extract.mac.build",
		Launch_file:          "nwjs-sdk-v0.90.0-osx-arm64/nwjs.app/Contents/MacOS/nwjs",
		PostSetup: func(pe PlatformEnv) {
			exec.Command("xattr", "-cr", "nwjs-sdk-v0.90.0-osx-arm64/nwjs.app").Output()
		},
		Extractor:    ExtractZip,
		DirSeparator: "/",
	},

	"linux_amd64": {
		Download_build:       "https://dl.nwjs.io/v0.90.0/nwjs-v0.90.0-linux-x64.tar.gz",
		Download_sdk:         "https://dl.nwjs.io/v0.90.0/nwjs-sdk-v0.90.0-linux-x64.tar.gz",
		Download_target:      "nwjs.linux_x86.tar.gz",
		Extract_sdk_target:   "extract.linux_x86.sdk",
		Extract_build_target: "extract.linux_x86.build",
		Launch_file:          "nwjs-sdk-v0.90.0-linux-x64/nw",
		PostSetup: func(pe PlatformEnv) {

		},
		Extractor:    ExtractTarGZ,
		DirSeparator: "/",
	},

	"windows_amd64": {
		Download_build:       "https://dl.nwjs.io/v0.90.0/nwjs-v0.90.0-win-x64.zip",
		Download_sdk:         "https://dl.nwjs.io/v0.90.0/nwjs-sdk-v0.90.0-win-x64.zip",
		Download_target:      "nwjs.win.zip",
		Extract_sdk_target:   "extract.win_x64.sdk",
		Extract_build_target: "extract.win.build",
		Launch_file:          "nwjs-sdk-v0.90.0-win-x64\\nw.exe",
		PostSetup: func(pe PlatformEnv) {

		},
		Extractor:    ExtractZip,
		DirSeparator: "\\",
	},
}

func GetPlatformConfig() PlatformEnv {

	var (
		platform_key = fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
		platform     = PlatformEnv{}
		ok           bool
	)

	platform, ok = platformMap[platform_key]
	if !ok {
		fmt.Println("Platform:", platform_key, "is not supported yet")
		os.Exit(1)
	}

	var base = NWGO_BASEPATH
	go_utils.MkDir(base)

	platform.Download_target = fmt.Sprintf("%s%s%s", base, platform.DirSeparator, platform.Download_target)
	platform.Extract_sdk_target = fmt.Sprintf("%s%s%s", base, platform.DirSeparator, platform.Extract_sdk_target)
	platform.Extract_build_target = fmt.Sprintf("%s%s%s", base, platform.DirSeparator, platform.Extract_build_target)
	platform.Launch_file = fmt.Sprintf("%s%s%s", platform.Extract_sdk_target, platform.DirSeparator, platform.Launch_file)

	return platform
}

func SetupPlatform() PlatformEnv {

	config := GetPlatformConfig()

	if !ZipExists(config.Download_target) {
		err := DownloadFile(
			config.Download_sdk,
			config.Download_target,
			"NWJS - SDK",
		)
		go_utils.Err(err)
	}

	didStuff, err := config.Extractor(
		config.Download_target,
		config.Extract_sdk_target,
	)
	go_utils.Err(err)

	if didStuff {
		config.PostSetup(config)
	}

	return config
}
