package config

type configs struct{}

var config configs

// Below vars are exposed to the build-layer (Makefile) so that they be overridden at build time.
var (
	Version = "unknown"
)

func Init() {
	config = configs{}
}

func GetVersion() string {
	return Version
}
