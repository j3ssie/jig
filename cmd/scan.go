package cmd

import (
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/j3ssie/jig/core"
	"github.com/j3ssie/jig/utils"
	"github.com/panjf2000/ants"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"sync"
)

func init() {
	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Generate input from list of URLs",
		Long:  core.Banner(),
		RunE:  runScan,
	}

	scanCmd.Flags().StringP("url", "u", "", "URL of target")
	scanCmd.Flags().StringP("urls", "U", "", "URLs file of target")
	scanCmd.Flags().StringP("raw", "r", "", "Raw request from Burp for origin")

	scanCmd.Flags().StringP("otype", "I", "", "Output type")
	RootCmd.AddCommand(scanCmd)
}

func runScan(cmd *cobra.Command, _ []string) error {
	var urls []string
	// parse URL input here
	urlFile, _ := cmd.Flags().GetString("urls")
	urlInput, _ := cmd.Flags().GetString("url")
	options.OutputType, _ = cmd.Flags().GetString("otype")
	if urlInput != "" {
		urls = append(urls, urlInput)
	}
	// input as a file
	if urlFile != "" {
		URLs := utils.ReadingLines(urlFile)
		for _, url := range URLs {
			urls = append(urls, url)
		}
	}

	// input as stdin
	if len(urls) == 0 {
		stat, _ := os.Stdin.Stat()
		// detect if anything came from std
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			sc := bufio.NewScanner(os.Stdin)
			for sc.Scan() {
				url := strings.TrimSpace(sc.Text())
				if err := sc.Err(); err == nil && url != "" {
					urls = append(urls, url)
				}
			}
		}
	}

	if len(urls) == 0 {
		fmt.Fprintf(os.Stderr, "[Error] No input loaded")
		os.Exit(1)
	}

	/* ---- Really start do something ---- */

	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(options.Concurrency, func(i interface{}) {
		startJob(i)
		wg.Done()
	}, ants.WithPreAlloc(true))
	defer p.Release()

	// Submit tasks one by one.
	for _, url := range urls {
		wg.Add(1)
		//job := libs.Job{URL: url, Sign: sign}
		_ = p.Invoke(url)
	}

	wg.Wait()
	return nil
}

func startJob(j interface{}) {
	job := j.(string)
	SendRequest(job, options)
}

func SendRequest(url string, options core.Options) {
	req, res := core.SendGET(url, options)
	data := GenOutput(req, res, options)
	fmt.Println(data)
	if data != "" && !options.NoOutput {
		utils.AppendToContent(options.Output, data)
	}
	//spew.Dump(res)
}

func GenOutput(req core.Request, res core.Response, options core.Options) string {
	data := make(map[string]string)
	switch options.OutputType {
	case "location":
		data["BaseURL"] = res.Location
	}
	if options.Debug {
		spew.Dump(data)

	}
	return core.ConvertToJson(data)
}
