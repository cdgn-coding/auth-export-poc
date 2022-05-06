package env

import "os"

const (
	GoEnvironment          = "GO_ENVIRONMENT"
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "develop"

	Port      = "PORT"
	Scope     = "SCOPE"
	ConfigDir = "CONF_DIR"
	GoPath    = "GOPATH"
)

func GetEnvironment() string {
	return os.Getenv(GoEnvironment)
}

func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}

func GetPort() string {
	port := os.Getenv(Port)
	if port == "" {
		port = ":8080"
	}
	return port
}

func GetScope() string {
	return os.Getenv(Scope)
}

func GetConfigDir() string {
	return os.Getenv(ConfigDir)
}

func GetGoPath() string {
	return os.Getenv(GoPath)
}
