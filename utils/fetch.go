package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func Fetch(wg *sync.WaitGroup, ctx context.Context, url string, idx uint32, outputDir string) error {
	defer wg.Done()
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
	// log.Println(cmdYTCommentDownloader.Args)

	var outBuf, errBuf bytes.Buffer
	cmdYTCommentDownloader.Stdout = io.MultiWriter(os.Stdout, &outBuf)
	cmdYTCommentDownloader.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	err := cmdYTCommentDownloader.Start()
	if err != nil {
		return err
	}

	// fmt.Printf("Downloading idx: %d\n", idx)

	err = cmdYTCommentDownloader.Wait()
	if err != nil {
		return err
	}

	// defer fmt.Println(outBuf.String())

	return nil
}
