package pkg

type PlatformEnv struct {
	DirSeparator         string
	Download_sdk         string
	Download_build       string
	Download_target      string
	Download_target_prod string
	Extract_sdk_target   string
	Extract_build_target string
	Launch_file          string
	BackendBinary        string
	BackendBinarySlash   string
	BuildCutPath         string
	BuildPreRunJS        string

	PostSetup   func(PlatformEnv)
	PostExtract func(PlatformEnv)
	PostBuild   func(pl PlatformEnv, buildPath string)
	Extractor   func(string, string) (bool, error)
}
