package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func Fetch(wg *sync.WaitGroup, done *chan struct{}, url string, idx uint32, outputDir string) error {
	defer wg.Done()
	defer func() {
		*done <- struct{}{}
	}()
	var cmdYTCommentDownloader = exec.Command("youtube-comment-downloader")
	cmdYTCommentDownloader.Args = append(
		cmdYTCommentDownloader.Args,
		"--url",
		url,
		"--output",
	)
	if outputDir == "" {
		cmdYTCommentDownloader.Args = append(cmdYTCommentDownloader.Args, fmt.Sprintf("./tmp/%d.json", idx))
	} else {
		cmdYTCommentDownloader.Args = append(cmdYTCommentDownloader.Args, fmt.Sprintf("%s/%d.json", strings.TrimRight(outputDir, "/"), idx))
	}

	var outBuf, errBuf bytes.Buffer
	cmdYTCommentDownloader.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	cmdYTCommentDownloader.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	err := cmdYTCommentDownloader.Start()
	if err != nil {
		return err
	}

	err = cmdYTCommentDownloader.Wait()
	if err != nil {
		return err
	}

	defer fmt.Println(outBuf.String())

	return nil
}
