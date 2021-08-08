package MTGoLogger

import "go.uber.org/zap"

type AssetLog struct {
	sugarLogger *zap.SugaredLogger
	conf        *LoggerConfig
}

type LoggerConfigs struct {
	Configs map[string]*LoggerConfig
}

type LoggerConfig struct {
	Level       string   `yaml:"level"`
	OutputPaths []string `yaml:"output_paths"`
	Appends     []string `yaml:"appends"`
}
