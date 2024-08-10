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
		Download_target_prod: "nwjs.mac.prod.zip",
		Extract_sdk_target:   "extract.mac.sdk",
		Extract_build_target: "extract.mac.build",
		Launch_file:          "nwjs-sdk-v0.90.0-osx-arm64/nwjs.app/Contents/MacOS/nwjs",
		PostSetup: func(pe PlatformEnv) {
			exec.Command("xattr", "-cr", "nwjs-sdk-v0.90.0-osx-arm64/nwjs.app").Output()
		},
		Extractor:     ExtractZip,
		DirSeparator:  "/",
		BackendBinary: "backend",
		BuildCutPath:  "nwjs-v0.90.0-osx-arm64/",
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

	platform.Download_target = fmt.Sprintf("%s%c%s", base, os.PathSeparator, platform.Download_target)
	platform.Download_target_prod = fmt.Sprintf("%s%c%s", base, os.PathSeparator, platform.Download_target_prod)
	platform.Extract_sdk_target = fmt.Sprintf("%s%c%s", base, os.PathSeparator, platform.Extract_sdk_target)
	platform.Extract_build_target = fmt.Sprintf("%s%c%s", base, os.PathSeparator, platform.Extract_build_target)
	platform.Extract_sdk_target = fmt.Sprintf("%s%c%s", base, os.PathSeparator, platform.Extract_sdk_target)
	platform.Extract_build_target = fmt.Sprintf("%s%c%s", base, os.PathSeparator, platform.Extract_build_target)
	platform.Launch_file = fmt.Sprintf("%s%c%s", platform.Extract_sdk_target, os.PathSeparator, platform.Launch_file)

	return platform
}

type SetupPlatformOpts struct {
	Production bool
}

var currentSDKPlatform, currentProdPlatform PlatformEnv
var setSDKPlatform, setProdPlatform bool

func SetupPlatform(opts SetupPlatformOpts) PlatformEnv {

	config := GetPlatformConfig()

	var dlURL, dlTarget, dlLabel, dlExtract string
	var currentPlatform *PlatformEnv
	var setCurrentPlatform *bool

	if opts.Production {
		if setProdPlatform {
			return currentProdPlatform
		}
		dlURL = config.Download_build
		dlTarget = config.Download_target_prod
		dlExtract = config.Extract_build_target
		dlLabel = "NWJS - Runtime"

		currentPlatform = &currentProdPlatform
		setCurrentPlatform = &setProdPlatform
	} else {
		if setSDKPlatform {
			return currentSDKPlatform
		}
		dlURL = config.Download_sdk
		dlTarget = config.Download_target
		dlLabel = "NWJS - SDK"
		dlExtract = config.Extract_sdk_target
		currentPlatform = &currentSDKPlatform
		setCurrentPlatform = &setSDKPlatform
	}

	if !ZipExists(dlTarget) {
		err := DownloadFile(dlURL, dlTarget, dlLabel)
		go_utils.Err(err)
	}

	var didStuff bool
	var err error

	didStuff, err = config.Extractor(dlTarget, dlExtract)

	go_utils.Err(err)

	if didStuff {
		config.PostSetup(config)
	}

	*currentPlatform = config
	*setCurrentPlatform = true
	return config
}
