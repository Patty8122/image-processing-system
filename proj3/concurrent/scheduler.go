package concurrent

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	DataDirs string 
	Mode     string 
	ThreadCount int
}

type MapReducer struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
	dataDir string
	Region    string `json:"region"`
}

func RunWorkStealing(config Config) {
	numThreads := config.ThreadCount
	pathToFile := "../data/effects.txt"
	file, _ := os.Open(pathToFile)
	reader := json.NewDecoder(file)
	ws := NewWorkStealingExecutor(numThreads, 10)

	for reader.More() {
		var req Request
		err := reader.Decode(&req)
		if err != nil {
			panic(err)
		}
	
		for _, dir := range strings.Split(config.DataDirs, "+") {
			req.dataDir = dir
			ws.Submit(req)
		}
	}
	ws.Shutdown()
}

func RunMapReduce(config Config) {
	resultChannel := make(chan map[string][]MapReducer, 2)
	go mapper(config, 1, resultChannel)
	go mapper(config, 2, resultChannel)

	mapped := make([]map[string][]MapReducer, 2)
	mapped[0] = <-resultChannel
	mapped[1] = <-resultChannel

	shuffler(mapped, config)
}

func Schedule(config Config) {
	if config.Mode == "ws" {
		RunWorkStealing(config)
	} else if config.Mode == "mr" {
		RunMapReduce(config)
	} else {
		RunSequential(config)
	}
}
