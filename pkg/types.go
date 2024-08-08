package pkg

type PlatformEnv struct {
	DirSeparator         string
	Download_sdk         string
	Download_build       string
	Download_target      string
	Extract_sdk_target   string
	Extract_build_target string
	Launch_file          string
	BackendBinary        string
	BackendBinarySlash   string

	PostSetup func(PlatformEnv)
	Extractor func(string, string) (bool, error)
}
