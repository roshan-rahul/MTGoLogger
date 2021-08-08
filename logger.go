package MTGoLogger

import "go.uber.org/zap"

type Config struct {
	level         string   `yaml:"level"`
	outputHandler string   `yaml:"output_paths"`
	appends       []string `yaml:"appends"`
}

type Configs struct {
	configs []*Config
}

type MTLogger struct {
	log     *zap.SugaredLogger
	configs *Configs
}

func getMyLogger() *MTLogger {
	return &MTLogger{}
}

func main() {

}
