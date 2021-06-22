// +build !test

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	buildNumber    string
	programVersion string
	programName    string
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	help := flag.Bool("help", false, "Show this help output.")
	safe := flag.Bool("safe", false, "Filter explicit images.")
	gif := flag.Bool("gif", false, "a bool")
	gray := flag.Bool("gray", false, "a bool")
	version := flag.Bool("version", false, fmt.Sprintf("Show the current %s version.", programName))
	height := flag.Int("height", 0, "a bool")
	width := flag.Int("width", 0, "a bool")
	wd := flag.String("chdir", "", "Switch to a different working directory before executing this program.")

	flag.Parse()

	tail := flag.Args()

	if *help {
		printUsage()
		return 0
	}

	if *version {
		printVersion()
		return 0
	}

	tailLen := len(tail)

	if tailLen < 1 {
		panic("Missing query")
	}

	if tailLen > 1 {
		panic("Too many queries")
	}

	if *wd != "" {
		err := os.Chdir(*wd)

		if err != nil {
			return 1
		}
	}

	opts := &Options{
		query:  tail[0],
		safe:   *safe,
		gif:    *gif,
		gray:   *gray,
		height: *height,
		width:  *width,
	}

	var bytesToDisk int64 = 0
	var downloadStartedAt time.Time

	btdUpdateMu := &sync.Mutex{}
	consumerWG := &sync.WaitGroup{}
	failCount := 0
	fcUpdateMu := &sync.Mutex{}
	producerClient := MakeHTTPClient()
	producerWG := &sync.WaitGroup{}
	progressWG := &sync.WaitGroup{}
	rcUpdateMu := &sync.Mutex{}
	resultChannel := make(chan string, 4096)
	resultCount := -2
	scUpdateMu := &sync.Mutex{}
	successCount := 0

	// Start progress hook
	progressWG.Add(1)
	go ProgressBar(progressWG, &failCount, &successCount, &resultCount, &bytesToDisk, &downloadStartedAt)

	// Start producers

	// Bing
	producerWG.Add(1)
	go Produce(resultChannel, producerWG, rcUpdateMu, &resultCount, func() []string {
		bing := Bing{
			client: producerClient,
			opts:   opts,
		}

		return bing.Scrape()
	})

	// Google
	producerWG.Add(1)
	go Produce(resultChannel, producerWG, rcUpdateMu, &resultCount, func() []string {
		google := Google{
			client: producerClient,
			opts:   opts,
		}

		return google.Scrape()
	})

	// Wait for producers to finish their jobs first
	producerWG.Wait()

	// Close channel after producers are done
	close(resultChannel)

	downloadStartedAt = time.Now()

	for i := 0; i < runtime.NumCPU(); i++ {
		consumerWG.Add(1)
		go Consume(i, resultChannel, consumerWG, scUpdateMu, fcUpdateMu, btdUpdateMu, &failCount, &successCount, &bytesToDisk)
	}

	// Wait for consumers to finish
	consumerWG.Wait()

	// Finally, wait for progress bar to finish
	progressWG.Wait()

	fmt.Println("")

	return 0
}

func printUsage() {
	fmt.Printf("Usage: %s [options] query \n\nOptions:\n", programName)
	flag.PrintDefaults()
}

func printVersion() {
	fmt.Printf("%s %s, build %s\n", programName, programVersion, buildNumber)
}

// Consume is a consumer which reads from ch and downloads corresponding
// image file.
func Consume(cpuNum int, ch chan string, wg *sync.WaitGroup, scUpdateMu, fcUpdateMu, btdUpdateMu *sync.Mutex, failCount, successCount *int, bytesToDisk *int64) {
	defer wg.Done()

	client := MakeHTTPClient()
	consumeCount := 0

	for item := range ch {
		outFilePath := fmt.Sprintf("%d-%d.jpg", cpuNum, consumeCount)
		success, bytesWritten := Download(client, item, outFilePath)

		if success {
			scUpdateMu.Lock()
			*successCount++
			scUpdateMu.Unlock()

			btdUpdateMu.Lock()
			*bytesToDisk += bytesWritten
			btdUpdateMu.Unlock()
		} else {
			fcUpdateMu.Lock()
			*failCount++
			fcUpdateMu.Unlock()
		}

		consumeCount++
	}
}

// Produce is a producer which scrapes data a source and writes to ch.
func Produce(ch chan string, wg *sync.WaitGroup, rcUpdateMu *sync.Mutex, resultCount *int, getItems func() []string) {
	defer wg.Done()

	items := getItems()

	for _, item := range items {
		ch <- item
	}

	rcUpdateMu.Lock()
	*resultCount += len(items) + 1
	rcUpdateMu.Unlock()
}

// ProgressBar is responsible to output relevant information regarding
// this program.
func ProgressBar(wg *sync.WaitGroup, failCount, successCount, resultCount *int, bytesToDisk *int64, downloadStartedAt *time.Time) {
	defer wg.Done()

	for (*failCount + *successCount) != *resultCount || *resultCount < 0 {
		time.Sleep(1500 * time.Millisecond)
		out := MakeProgressBarOutput(downloadStartedAt, *bytesToDisk, *successCount, *failCount, *resultCount)
		fmt.Printf("%s\r", out)
	}
}
