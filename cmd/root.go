package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vlaborie/remarkable-sync/remarkable"

	"github.com/bmaupin/go-epub"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "remarkable-sync",
		Short: "Sync tool for reMarkable paper tablet",
		Long: `Remarkable-sync is a Go applications for syncing external
services to reMarkable paper table, like Wallabag or Miniflux.`,
		Run: func(cmd *cobra.Command, args []string) {
			Remarkable := remarkable.New("/home/root/.local/share/remarkable/xochitl/")

			wallabagConfig, err := os.Open("/etc/remarkable-sync/wallabag.json")
			if err == nil {
				Remarkable.Wallabag(wallabagConfig)
			}

			minifluxConfig, err := os.Open("/etc/remarkable-sync/miniflux.json")
			if err == nil {
				Remarkable.Miniflux(minifluxConfig)
			}

			for _, RemarkableItem := range Remarkable.Items {
				if RemarkableItem.ContentType == "html" {
					RemarkableItem.ContentType = "epub"
					e := epub.NewEpub(RemarkableItem.VisibleName)
					e.AddSection(string(RemarkableItem.Content), "Section 1", "", "")
					e.Write(Remarkable.Dir + RemarkableItem.Id + "." + RemarkableItem.ContentType)
					fmt.Println("EPUB of " + RemarkableItem.Id + " writed")
				}

				j, _ := json.Marshal(RemarkableItem)
				_ = ioutil.WriteFile(Remarkable.Dir+RemarkableItem.Id+".metadata", j, 0644)
				fmt.Println("Metadata of " + RemarkableItem.Id + " updated")
			}

		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}