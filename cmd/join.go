package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

func init() {
	rootCmd.AddCommand(joinCmd)
}

var joinCmd = &cobra.Command{
	Use:                   "join",
	Short:                 "Join sheets",
	Args:                  cobra.ExactArgs(0),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetFlags(log.Llongfile)
		}
		err := joinProcess()
		if err != nil {
			log.Fatalln(err)
		}

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if !quiet {
			fmt.Println("Join success")
		}
	},
}

func joinProcess() error {
	if !quiet {
		fmt.Println("Joining...")
	}

	if outputDir == "" {
		outputDir = "./tmp"
	}
	joinedFile, err := os.OpenFile(strings.TrimRight(outputDir, "/")+"/joined.xlsx", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer joinedFile.Close()

	mother := excelize.NewFile()

	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		return err
	}

	for _, v := range files {
		if !strings.HasSuffix(v.Name(), ".xlsx") || v.Name() == "joined.xlsx" {
			continue
		}
		sheetIdx, err := mother.NewSheet(v.Name())
		if err != nil {
			return fmt.Errorf("Error creating new sheet: %s", err)
		}
		child, err := excelize.OpenFile(strings.TrimRight(outputDir, "/") + "/" + v.Name())
		if err != nil {
			return fmt.Errorf("Error open child xlsx file: %s", err)
		}

		maxEmpty := 20
		emptyCount := 0

		for colNum := 1; ; colNum++ {
			for rowNum := 1; ; rowNum++ {
				colName, err := excelize.ColumnNumberToName(colNum)
				if err != nil {
					return err
				}
				cell, err := child.GetCellValue("Sheet1", fmt.Sprintf("%s%d", colName, rowNum))
				if err != nil {
					return err
				}

				if cell == "" {
					emptyCount++
					break
				} else {
					emptyCount = 0
				}

				mother.SetCellValue(
					mother.GetSheetName(sheetIdx),
					fmt.Sprintf("%s%d", colName, rowNum),
					cell,
				)
			}
			if emptyCount > maxEmpty {
				break
			}
		}
		child.Close()
	}

	if _, err := mother.WriteTo(joinedFile); err != nil {
		return err
	}
	return nil
}
