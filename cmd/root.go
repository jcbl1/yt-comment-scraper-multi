package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/jcbl1/yt-comment-scraper-multi/utils"
	"github.com/spf13/cobra"
)

const (
	sUCCESS_OUTPUT = "Success."
)

var rootCmdName = "yt-comment-scraper-multi"
var urlListFilename, outputDir string
var debug = false

func init() {
	rootCmd.Flags().StringVarP(&urlListFilename, "url-list-filename", "f", "", "URL list")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", "", "Specify the output directory to store results (the default one is ./tmp/)")
}

var rootCmd = &cobra.Command{
	Use:                   rootCmdName,
	Short:                 "A tool to scrape youtube comments",
	DisableFlagsInUseLine: true,
	Args:                  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetFlags(log.Llongfile)
		}
		if urlListFilename == "" {
			cmd.Usage()
			return
		}
		err := process()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(sUCCESS_OUTPUT)
	},
}

func process() error {
	f, err := os.Open(urlListFilename)
	if err != nil {
		return err
	}
	defer f.Close()
	scnr := bufio.NewScanner(f)
	wg := sync.WaitGroup{}
	var idx uint32
	for scnr.Scan() {
		text := scnr.Text()
		if text != "" {
			wg.Add(1)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				atomic.AddUint32(&idx, 1)
				err := utils.Fetch(&wg, ctx, text, idx, outputDir)
				if err != nil {
					log.Fatalln(err)
				}
			}()
		}
	}
	wg.Wait()
	return nil
}

func Execute() {
	rootCmd.Execute()
}
