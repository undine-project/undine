package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/undine-project/undine/src/builder"
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
		config := builder.LoadConfig()

		b, err := builder.NewBuilder(config)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := b.Close(); err != nil {
				log.Printf("Error closing builder: %v", err)
				if err == nil {
					os.Exit(1)
				}
			}
		}()

		if err := b.Initialize(true); err != nil {
			log.Fatal(err)
		}

		// Initial build
		if err := b.Build(); err != nil {
			log.Fatal(err)
		}

		// Setup signal handling
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
		go func() {
			sig := <-sigChan
			log.Printf("Received signal %s, exiting...", sig)
			os.Exit(0)
		}()

		// Start watching
		go func() {
			if err := b.Watch(func(content builder.FileContent) {
				fmt.Println("Sending content...")
				// Note: WebHandler should be refactored to be passed in or managed by Builder
				(&builder.WebHandler{}).SendContent(content)
			}); err != nil {
				log.Fatal(err)
			}
		}()

		// Keep the main goroutine running
		select {}
	},
}

func init() {
	RootCmd.AddCommand(watchCmd)
}
