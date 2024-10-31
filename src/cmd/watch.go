package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/undine-project/undine/src/builder"
	"github.com/undine-project/undine/src/support"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes in the diagram.md file",
	Run: func(cmd *cobra.Command, args []string) {
		c := support.LoadConfig()

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
			true,
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
		startWatching(sp, fg)
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
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
