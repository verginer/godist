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

	totalInStock := supplySystem.TotalInStock()
	totalManufactured := supplySystem.TotalManufactured()
	if totalInStock != totalManufactured {
	    log.Fatal("The quantity manufactured is not equal to the quantity in stock. This cannot be.")
    } else {
        log.Printf("A total of %d pills were manufactured and shipped.", totalManufactured)
    }
	traces := supplySystem.ExtractTraces()
	log.Printf("Extracted %d traces", len(traces))
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
