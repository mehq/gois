// +build !test

package main

import (
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	var bytesToDisk int64 = 0
	var downloadStartedAt time.Time

	btdUpdateMu := &sync.Mutex{}
	consumerWG := &sync.WaitGroup{}
	failCount := 0
	fcUpdateMu := &sync.Mutex{}
	opts := MakeOptions()
	producerClient := MakeHTTPClient()
	producerWG := &sync.WaitGroup{}
	progressWG := &sync.WaitGroup{}
	rcUpdateMu := &sync.Mutex{}
	resultChannel := make(chan string, 4096)
	resultCount := 0
	scUpdateMu := &sync.Mutex{}
	successCount := 0

	// Start progress hook
	progressWG.Add(1)
	go ProgressBar(progressWG, func() bool {
		return (failCount + successCount) == resultCount
	}, func() {
		out := MakeProgressBarOutput(&downloadStartedAt, bytesToDisk, successCount, failCount, resultCount)
		fmt.Printf("%s\r", out)
	})

	// Start producers

	// Bing
	producerWG.Add(1)
	go Produce(resultChannel, producerWG, func() []string {
		bing := Bing{
			client: producerClient,
			opts:   opts,
		}

		return bing.Scrape()
	}, func(ic int) {
		rcUpdateMu.Lock()
		resultCount += ic
		rcUpdateMu.Unlock()
	})

	// Wait for producers to finish their jobs first
	producerWG.Wait()

	// Close channel after producers are done
	close(resultChannel)

	downloadStartedAt = time.Now()

	for i := 0; i < runtime.NumCPU(); i++ {
		consumerWG.Add(1)
		go Consume(i, resultChannel, consumerWG, func(success bool, bw int64) {
			if success {
				scUpdateMu.Lock()
				successCount++
				scUpdateMu.Unlock()

				btdUpdateMu.Lock()
				bytesToDisk += bw
				btdUpdateMu.Unlock()
			} else {
				fcUpdateMu.Lock()
				failCount++
				fcUpdateMu.Unlock()
			}
		})
	}

	// Wait for consumers to finish
	consumerWG.Wait()

	// Finally, wait for progress bar to finish
	progressWG.Wait()

	fmt.Println("")
}

// Consume is a consumer which reads from ch and downloads corresponding
// image file.
func Consume(cpuNum int, ch chan string, wg *sync.WaitGroup, updateCounters func(success bool, bw int64)) {
	defer wg.Done()

	client := MakeHTTPClient()
	consumeCount := 0

	for item := range ch {
		outFilePath := fmt.Sprintf("%d-%d.jpg", cpuNum, consumeCount)
		success, bytesWritten := Download(client, item, outFilePath)
		updateCounters(success, bytesWritten)
	}
}

// MakeOptions parses command line arguments into Options
func MakeOptions() *Options {
	//help := flag.Bool("help", false, "help text")
	safe := flag.Bool("safe", true, "a bool")
	gif := flag.Bool("gif", false, "a bool")
	gray := flag.Bool("gray", false, "a bool")
	height := flag.Int("height", 0, "a bool")
	width := flag.Int("width", 0, "a bool")

	flag.Parse()

	tail := flag.Args()

	//if *help {
	//	flag.PrintDefaults()
	//	os.Exit(0)
	//}

	tailLen := len(tail)

	if tailLen < 1 {
		panic("Missing query")
	}

	if tailLen > 1 {
		panic("Too many queries")
	}

	opts := &Options{
		query:  tail[0],
		safe:   *safe,
		gif:    *gif,
		gray:   *gray,
		height: *height,
		width:  *width,
	}

	return opts
}

// Produce is a producer which scrapes data a source and writes to ch.
func Produce(ch chan string, wg *sync.WaitGroup, getItems func() []string, updateCounters func(ic int)) {
	defer wg.Done()

	items := getItems()

	for _, item := range items {
		ch <- item
	}

	updateCounters(len(items))
}

// ProgressBar is responsible to output relevant information regarding
// this program.
func ProgressBar(wg *sync.WaitGroup, shouldContinue func() bool, hook func()) {
	defer wg.Done()

	for shouldContinue() {
		time.Sleep(1500 * time.Millisecond)
		hook()
	}
}
