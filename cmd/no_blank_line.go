package cmd

import (
	"bufio"
	"os"
	"strings"
)

func determineOutputDir() {
	if outputDir == "" {
		outputDir = "./tmp"
	} else {
		outputDir = strings.TrimRight(outputDir, "/")
	}
}

func noBlankLine() error {
	determineOutputDir()
	f, err := os.OpenFile(outputDir+"/joined.txt", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	scnr := bufio.NewScanner(f)
	var lines []string
	for scnr.Scan() {
		txt := scnr.Text()
		if txt != "" {
			lines = append(lines, txt)
		}
	}

	f.Truncate(0)
	for _, line := range lines {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
