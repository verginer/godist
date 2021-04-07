package main

import (
	godist "github.com/verginer/godist/src"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func worker(pathToTransactions, pathToJson string) {
	supplySystem := godist.ReplayTransactionsFromFile(pathToTransactions)
	traces := supplySystem.ExtractTraces()
	godist.WriteTracesToJson(traces, pathToJson)
}

func main() {
	ndcDirectory := os.Args[1]
	jsonDirectory := os.Args[2]

	ndcFiles, err := os.ReadDir(ndcDirectory)

	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, ndcFile := range ndcFiles {

		drugName := ndcFile.Name()
		jsonPath := filepath.Join(jsonDirectory, drugName+".json")
		ndcPath := filepath.Join(ndcDirectory, drugName)

		wg.Add(1)
		go func(ndcPath, jsonPath string) {
			log.Println("Extracting: ", drugName)
			worker(ndcPath, jsonPath)
			wg.Done()
		}(ndcPath, jsonPath)
	}

	wg.Wait()
}
