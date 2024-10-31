package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/undine-project/undine/src/builder"
	"github.com/undine-project/undine/src/support"
	"log"
	"os"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate HTML file and exit",
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
			false,
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
		fmt.Println("HTML generated without watching.")
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
