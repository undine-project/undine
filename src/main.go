package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/undine-project/undine/src/builder"
	"github.com/undine-project/undine/src/cmd"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	TemplatePath string                   `yaml:"templatePath"`
	Files        []builder.FileDefinition `yaml:"files"`
}

func main() {
	cmd.Execute()
	c := loadConfig()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)

		return
	}

	files := c.Files
	files = append(files, builder.FileDefinition{
		Name:  "template",
		Path:  c.TemplatePath,
		Title: "Template",
	})
	sp := builder.NewSourceProcessor(files, watcher)

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		err := os.Mkdir("public", 0755)
		if err != nil {
			panic(err)
		}
	}

	fg := builder.NewFileGenerator(
		c.TemplatePath,
		"public/index.html",
		cmd.WatchMode,
		files,
	)
	for content := range sp.Process() {
		fg.SetContent(content)
	}
	err = fg.Generate()
	if err != nil {
		log.Fatal(err)

		return
	}

	if cmd.WatchMode {
		startWatching(sp, fg)
	} else {
		fmt.Println("HTML generated without watching.")
	}
}

func loadConfig() config {
	configFile := "docs-config.yaml"

	yamlFile, err := os.Open(configFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func(yamlFile *os.File) {
		_ = yamlFile.Close()
	}(yamlFile)

	c := config{}
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

func startWatching(sourceProcessor *builder.SourceProcessor, generator *builder.FileGenerator) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal %s, exiting...", sig)
		os.Exit(0)
	}()

	go func() {
		contentsChannel := make(chan builder.FileContent, 2)
		sourceProcessor.Watch(contentsChannel)
		defer sourceProcessor.Stop()

		wh := &builder.WebHandler{}
		wh.StartServer()

		for content := range contentsChannel {
			fmt.Printf("Generating HTML file with type %s\n", content.Name)
			generator.SetContent(content)
			err := generator.Generate()
			if err != nil {
				log.Fatal(err)

				return
			}

			fmt.Println("Sending content...")
			wh.SendContent(content)
		}
	}()

	// Keep the main goroutine running to prevent the program from exiting
	select {}
}
