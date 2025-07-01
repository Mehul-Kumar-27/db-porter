package core

import (
	"os"

	"gopkg.in/yaml.v2"
)

type AdapterConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

type PipelineConfig struct {
	SourceConfig      AdapterConfig `yaml:"source"`
	DestinationConfig AdapterConfig `yaml:"destination"`
}

func LoadPipelineConfig(path string) (*PipelineConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg PipelineConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
