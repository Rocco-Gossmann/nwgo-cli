package pkg

type PlatformEnv struct {
	Download_sdk         string
	Download_build       string
	Download_target      string
	Extract_sdk_target   string
	Extract_build_target string
	Launch_file          string

	PostSetup func(PlatformEnv)
	Extractor func(string, string) (bool, error)
}
