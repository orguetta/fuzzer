package fuzzer

import (
	"encoding/json"
	"fuzzer/src/logger"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type result struct {
	RedirectLocation string `json:"redirectLocation"`
	URL              string `json:"url"`
	Size             int    `json:"size"`
	Lines            int    `json:"lines"`
	StatusCode       int    `json:"statusCode"`
	Words            int    `json:"words"`
}

// Results is worker which saves results one by one in jsonl format
func (f *Fuzzer) Results() {
	fd, err := os.OpenFile(f.OutFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Error("error in opening out file",
			zap.Error(err),
		)
		return
	}

	defer func() {
		f.mutex.Lock()
		defer f.mutex.Unlock()

		logger.Log.Debug("shutting down results worker",
			zap.Int("totalWorkers", f.totalWorkers),
		)

		f.totalWorkers--
	}()

	shouldWork := true

	// monitoring for control exit
	go func() {
		<-f.control

		f.mutex.Lock()
		defer f.mutex.Unlock()

		shouldWork = false
	}()

	for {
		f.mutex.Lock()
		if !shouldWork {
			f.mutex.Unlock()
			return
		}
		f.mutex.Unlock()

		var r result
		select {
		case r = <-f.results:
		case <-time.After(3 * time.Second):
			continue
		}

		raw, _ := json.Marshal(r)

		fd.WriteString(string(raw) + "\n")
	}
}

func GetUniqueNumbers(input, delimiter string) (res []int) {
	res = make([]int, 0)

	tmp := strings.Split(input, delimiter)
	unq := make(map[int]bool, 0)
	for _, t := range tmp {
		val, err := strconv.Atoi(t)
		if err != nil {
			continue
		}

		unq[val] = true
	}

	for key := range unq {
		res = append(res, key)
	}

	return
}
