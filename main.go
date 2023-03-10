package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dpanic/fuzzer/src/fuzzer"
)

func main() {
	maxTime := flag.Int("maxTime", 0, "maximum execution time")
	maxReqSec := flag.Int("maxReqSec", 0, "maximum requests per second, default unlimited")
	method := flag.String("X", "GET", "GET, POST, HEAD, OPTIONS, PUT ...")
	filterCodes := flag.String("fc", "", "403,404")
	filterLines := flag.String("fl", "", "123,321")
	filterWords := flag.String("fw", "", "1,2,3")
	filterSize := flag.String("fs", "", "300,200")
	userAgent := flag.String("ua", "", "custom user agent")
	pseudoRandomUserAgent := flag.Bool("prua", false, "pseudo random user agent")

	outFile := flag.String("o", "", "/tmp/outFile.json")
	wordList := flag.String("w", "", "wordlists/big.txt")
	url := flag.String("u", "", "https://www.google.com/FUZZ")
	proxyURL := flag.String("p", "", "http://127.0.0.1:9000")

	flag.Parse()

	f, err := fuzzer.New(&fuzzer.Config{
		URL:                   *url,
		Method:                *method,
		ProxyURL:              *proxyURL,
		OutFile:               *outFile,
		WordList:              *wordList,
		UserAgent:             *userAgent,
		PseudoRandomUserAgent: *pseudoRandomUserAgent,
		MaxTime:               time.Duration(*maxTime) * time.Second,
		MaxReqSec:             *maxReqSec,
		Filters: fuzzer.Filters{
			StatusCodes: fuzzer.GetUniqueNumbers(*filterCodes, ","),
			Words:       fuzzer.GetUniqueNumbers(*filterWords, ","),
			Lines:       fuzzer.GetUniqueNumbers(*filterLines, ","),
			Size:        fuzzer.GetUniqueNumbers(*filterSize, ","),
		},
	})

	if err != nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	go f.Start()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		for {
			s := <-signalChannel
			f.Log.Warn(fmt.Sprintf("received signal %s", s.String()))

			switch s {
			case syscall.SIGHUP:
			case syscall.SIGINT:
				f.Done <- "SIGINT"
				return
			case syscall.SIGTERM:
				f.Done <- "SIGTERM"
				return
			case syscall.SIGQUIT:
				f.Done <- "SIGQUIT"
				return
			}
		}
	}()

	<-f.Done
	f.Stop()
}
