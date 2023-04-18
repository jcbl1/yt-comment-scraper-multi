package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().BoolVarP(&joinOpt, "join", "j", false, "Append joining operation")
}

var convertCmd = &cobra.Command{
	Use:                   "convert",
	DisableFlagsInUseLine: true,
	Example:               "convert --out-dir ./somepath/",
	Short:                 "Convert to Excel sheets",
	Args:                  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetFlags(log.Llongfile)
		}
		err := convertProcess()
		if err != nil {
			log.Fatalln(err)
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if !quiet {
			fmt.Println("Conversion success")
		}
	},
}

var waitConvert sync.WaitGroup

func convertProcess() error {
	if !quiet {
		fmt.Println("Converting...")
	}

	if outputDir == "" {
		outputDir = "./tmp"
	} else {
		outputDir = strings.TrimRight(outputDir, "/")
	}
	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		filename := outputDir + "/" + f.Name()
		waitConvert.Add(1)
		go func() {
			err := convert(filename)
			if err != nil {
				log.Fatalln(err)
			}
		}()
	}
	waitConvert.Wait()

	return nil
}

type Comment struct {
	CID        string  `json:"cid"`
	Text       string  `json:"text"`
	Time       string  `json:"time"`
	Author     string  `json:"author"`
	Channel    string  `json:"channel"`
	Votes      string  `json:"votes"`
	Photo      string  `json:"photo"`
	Heart      bool    `json:"heart"`
	Reply      bool    `json:"reply"`
	TimeParsed float64 `json:"time_parsed"`
}

func convert(filename string) error {
	defer waitConvert.Done()

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scnr := bufio.NewScanner(f)
	var comments []Comment
	for scnr.Scan() {
		line := scnr.Text()
		if line != "" {
			comment := Comment{}
			err := json.Unmarshal([]byte(line), &comment)
			if err != nil {
				return err
			}
			comments = append(comments, comment)
		}
	}

	// Writing to file
	outputFile, err := os.OpenFile(
		strings.TrimSuffix(filename, ".json")+".xlsx",
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	eFile := excelize.NewFile()
	if err := eFile.SetCellStr("Sheet1", "A1", "CID"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "B1", "Text"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "C1", "Time"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "D1", "Author"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "E1", "Channel"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "F1", "Votes"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "G1", "Photo"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "H1", "Heart"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "I1", "Reply"); err != nil {
		return err
	}
	if err := eFile.SetCellStr("Sheet1", "J1", "TimeParsed"); err != nil {
		return err
	}

	for i, v := range comments {
		if err := eFile.SetCellStr("Sheet1", fmt.Sprintf("A%d", i+2), v.CID); err != nil {
			return err
		}
		if err := eFile.SetCellStr("Sheet1", fmt.Sprintf("B%d", i+2), v.Text); err != nil {
			return err
		}
		if err := eFile.SetCellStr("Sheet1", fmt.Sprintf("C%d", i+2), v.Time); err != nil {
			return err
		}
		if err := eFile.SetCellStr("Sheet1", fmt.Sprintf("D%d", i+2), v.Author); err != nil {
			return err
		}
		if err := eFile.SetCellStr("Sheet1", fmt.Sprintf("E%d", i+2), v.Channel); err != nil {
			return err
		}
		if err := eFile.SetCellStr("Sheet1", fmt.Sprintf("F%d", i+2), v.Votes); err != nil {
			return err
		}
		if err := eFile.SetCellStr("Sheet1", fmt.Sprintf("G%d", i+2), v.Photo); err != nil {
			return err
		}
		if err := eFile.SetCellBool("Sheet1", fmt.Sprintf("H%d", i+2), v.Heart); err != nil {
			return err
		}
		if err := eFile.SetCellBool("Sheet1", fmt.Sprintf("I%d", i+2), v.Reply); err != nil {
			return err
		}
		if err := eFile.SetCellFloat("Sheet1", fmt.Sprintf("J%d", i+2), v.TimeParsed, -1, 64); err != nil {
			return err
		}
	}

	_, err = eFile.WriteTo(outputFile)
	if err != nil {
		return err
	}

	return nil
}
