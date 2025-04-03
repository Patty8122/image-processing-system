package concurrent

import (
	"os"
	"proj3/png"
	"strings"
	"encoding/json"
)

type Request struct {
	InPath  string   `json:"inPath"`
	OutPath string   `json:"outPath"`
	Effects []string `json:"effects"`
	dataDir string
}



func processImage(request Request, dataDir string) {
	fileInpath := "../data/in/" + dataDir + "/" + request.InPath
	fileOutpath := "../data/out/" + dataDir + "_" + request.OutPath
	pngImg, err := png.Load(fileInpath)
	if err != nil {
		panic(err)
	}
	pngImg.RunEffects(request.Effects)
	pngImg.Save(fileOutpath)
}

func RunSequential(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	effectsPathFile := "../data/effects.txt"

	effectsFile, _ := os.Open(effectsPathFile)
	reader := json.NewDecoder(effectsFile)

	for reader.More() {
		var request Request
		err := reader.Decode(&request)
		if err != nil {
			panic(err)
		}
		for _, dataDir := range dataDirs {
			processImage(request, dataDir)
		}
	}
}
