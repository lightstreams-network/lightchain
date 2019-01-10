package node


type Config struct {
	DataDir string
}

func NewConfig(homeDir string) Config {
	return Config {
		homeDir,
	}
}
