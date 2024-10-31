package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/undine-project/undine/src/builder"
	"log"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Generate HTML file and exit",
	Run: func(cmd *cobra.Command, args []string) {
		config := builder.LoadConfig()

		b, err := builder.NewBuilder(config)
		if err != nil {
			log.Fatal(err)
		}
		defer b.Close()

		if err := b.Initialize(false); err != nil {
			log.Fatal(err)
		}

		if err := b.Build(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("HTML generated without watching.")
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
