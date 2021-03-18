package config

type configs struct{}

var config configs

// Below vars are exposed to the build-layer (Makefile) so that they be overridden at build time.
var (
	Version         = "unknown"
	CotterAPIKey    = "unknown"
	BrevAPIEndpoint = "https://app.brev.dev"
)

func Init() {
	config = configs{}
}

func GetVersion() string {
	return Version
}

func GetCotterAPIKey() string {
	return CotterAPIKey
}

func GetBrevAPIEndpoint() string {
	return BrevAPIEndpoint
}
