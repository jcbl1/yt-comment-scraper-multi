package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"

	"github.com/jcbl1/yt-comment-scraper-multi/utils"
	"github.com/spf13/cobra"
)

var maxThreads = 8

var rootCmdName = "yt-comment-scraper-multi"
var urlListFilename, outputDir string
var debug bool
var quiet bool
var convertOpt bool
var joinOpt bool

func init() {
	rootCmd.Flags().StringVarP(&urlListFilename, "url-list-filename", "f", "", "URL list")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", "", "Specify the output directory to store results (the default one is ./tmp/)")

	rootCmd.Flags().BoolVarP(&convertOpt, "convert", "c", false, "Append converting operation")
	rootCmd.Flags().BoolVarP(&joinOpt, "join", "j", false, "Append joining operation")

	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Amit all outputs")
	rootCmd.Flags().IntVar(&maxThreads, "max-threads", 8, "Max threads scraping comments")
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

		if convertOpt {
			convertCmd.Run(cmd, args)
			if joinOpt {
				joinCmd.Run(cmd, args)
			}
		}

	},

	PostRun: func(cmd *cobra.Command, args []string) {
		if !quiet {
			fmt.Println("Success")
		}
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
	threads := 0
	done := make(chan struct{}, maxThreads)
	defer close(done)
	for scnr.Scan() {
		text := scnr.Text()
		if text != "" {
			wg.Add(1)
			if threads < maxThreads {
				threads++
			} else {
				<-done
			}
			go func() {
				atomic.AddUint32(&idx, 1)
				err := utils.Fetch(&wg, &done, text, idx, outputDir)
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
