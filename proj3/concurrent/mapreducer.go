package concurrent

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

func mapper(config Config, fileNo int, ch chan map[string][]MapReducer) {
	filePath := "../data/effects" + fmt.Sprint(fileNo) + ".txt"

	file, _ := os.Open(filePath)
	fileReader := json.NewDecoder(file)
	resultMap := make(map[string][]MapReducer)

	for fileReader.More() {
		var currReq MapReducer
		err := fileReader.Decode(&currReq)
		if err != nil {
			panic(err)
		}

		for _, dir := range strings.Split(config.DataDirs, "+") {
			currReq.dataDir = dir
			resultMap[currReq.Region] = append(resultMap[currReq.Region], currReq)
		}
	}

	ch <- resultMap

}

func reducer(imgArr []MapReducer, wg *sync.WaitGroup) {
	for _, imgTask := range imgArr {
		processImage(Request{imgTask.InPath, imgTask.OutPath, imgTask.Effects, imgTask.dataDir}, imgTask.dataDir)
	}
	// wg.Done()
}

// func shuffler(intermediateMap []map[string][]MapReducer) {
// 	shuffledResults := make(map[string][]MapReducer)

// 	for i := 0; i < len(intermediateMap); i += 1 {
// 		for key, value := range intermediateMap[i] {
// 			shuffledResults[key] = append(shuffledResults[key], value...)
// 		}
// 	}

// 	waitGroup := &sync.WaitGroup{}

// 	// create a channel to send the shuffled results to
// 	// create a wait group to wait for all the goroutines to finish
// 	// iterate through the shuffled results

// 	shuffledResultsChannel := make(chan []MapReducer, len(shuffledResults))
// 	for sr := range shuffledResults {
// 		go func() {
// 			shuffledResultsChannel <- shuffledResults[sr]
// 		}
// 	}

// 	for sr := range shuffledResultsChannel {
// 		waitGroup.Add(1)
// 		go reducer(sr, waitGroup)
// 	}

// 	waitGroup.Wait()

// }

func shuffler(intermediateMap []map[string][]MapReducer, config Config) {
	// Steps:
	// 1. Create a map to store the shuffled results by key
	// 2. Create a channel to send the shuffled results to
	// 3. Each of the config.ThreadCount goroutines should remove one entry from the channel and process it
	// 4. Each goroutine should wait for all the goroutines to finish
	// 5. Each goroutine should call the reducer function

	shuffledResults := make(map[string][]MapReducer)

	for i := 0; i < len(intermediateMap); i += 1 {
		for key, value := range intermediateMap[i] {
			shuffledResults[key] = append(shuffledResults[key], value...)
		}
	}

	ch := make(chan []MapReducer, len(shuffledResults))
	wg := &sync.WaitGroup{}

	for _, value := range shuffledResults {
		ch <- value
	}

	// fmt.Println("Shuffled results channel length:", len(ch))
	for i := 0; i < config.ThreadCount; i++ {
		wg.Add(1)
	
		go func(threadID int) {
			// remove one entry from the channel and process it until the channel is empty
			defer wg.Done() // Ensure that wg.Done() is called when the goroutine exits
	
			for {
				select {
				case key1, ok := <-ch:
					if !ok {
						// Channel is closed, no more entries
						return
					}
					// fmt.Println("Thread", threadID, "is processing", key1)
					reducer(key1, wg)
				default:
					// Channel is empty, exit the goroutine
					return
				}
			}
		}(i)
	}
	

	wg.Wait()
}