package notion

var (
	ApiKey         string
	ApiVersion     = "2022-06-28"
	OutputDir      string
	RelativeImgDir string
)

func init() {
	ApiKey = ""
	OutputDir = ""
}

// Init 추가적인 인자의
func Init(apiKey, outputDir, relativeImgDir string) {
	ApiKey = apiKey
	OutputDir = outputDir
	RelativeImgDir = relativeImgDir
}
