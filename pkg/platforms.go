package pkg

import (
	"fmt"
	"os"
	"runtime"

	"github.com/rocco-gossmann/go_utils"
)

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

	if didStuff && config.PostSetup != nil {
		config.PostSetup(config)
	}

	*currentPlatform = config
	*setCurrentPlatform = true
	return config
}
