package builder

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	TemplatePath string           `yaml:"templatePath"`
	Files        []FileDefinition `yaml:"files"`
}

// LoadConfig Temporary, until #1 implemented
func LoadConfig() Config {
	configFile := "docs-config.yaml"

	yamlFile, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(yamlFile *os.File) {
		_ = yamlFile.Close()
	}(yamlFile)

	c := Config{}
	yamlDecoder := yaml.NewDecoder(yamlFile)
	err = yamlDecoder.Decode(&c)
	if err != nil {
		log.Fatal(err)
	}

	if c.TemplatePath == "" {
		panic("template path is empty")
	}

	_, err = os.Stat(c.TemplatePath)
	if err != nil {
		panic(fmt.Sprintf("template file %s doesn't exist", c.TemplatePath))
	}

	return c
}
