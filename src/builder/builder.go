package builder

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
)

type Builder struct {
	config    Config
	watcher   *fsnotify.Watcher
	files     []FileDefinition
	generator *FileGenerator
	processor *SourceProcessor
}

func NewBuilder(config Config) (*Builder, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	files := config.Files
	files = append(files, FileDefinition{
		Name:  "template",
		Path:  config.TemplatePath,
		Title: "Template",
	})

	if err := ensurePublicDir(); err != nil {
		return nil, err
	}

	return &Builder{
		config:  config,
		watcher: watcher,
		files:   files,
	}, nil
}

func (b *Builder) Initialize(watchMode bool) error {
	b.processor = NewSourceProcessor(b.files, b.watcher)
	b.generator = NewFileGenerator(
		b.config.TemplatePath,
		"public/index.html",
		watchMode,
		b.files,
	)

	return nil
}

func (b *Builder) Build() error {
	for content := range b.processor.Process() {
		b.generator.SetContent(content)
	}
	return b.generator.Generate()
}

func (b *Builder) Watch(onChange func(FileContent)) error {
	contentsChannel := make(chan FileContent, 2)
	b.processor.Watch(contentsChannel)
	defer b.processor.Stop()

	wh := &WebHandler{}
	wh.StartServer()

	for content := range contentsChannel {
		fmt.Printf("Generating HTML file with type %s\n", content.Name)
		b.generator.SetContent(content)
		if err := b.generator.Generate(); err != nil {
			return fmt.Errorf("failed to generate file: %w", err)
		}

		if onChange != nil {
			onChange(content)
		}
	}

	return nil
}

func (b *Builder) Close() error {
	return b.watcher.Close()
}

func ensurePublicDir() error {
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		if err := os.Mkdir("public", 0755); err != nil {
			return fmt.Errorf("failed to create public directory: %w", err)
		}
	}
	return nil
}
