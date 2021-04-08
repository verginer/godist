package main

import (
    godist "github.com/verginer/godist/src"
    "log"
    "os"
    "path/filepath"
)

func worker(pathToTransactions, pathToJson string) {

    csvInfo := godist.TransactionCSVInfo{
        SendingId:    0,
        SendingAct:   1,
        ReceivingId:  10,
        ReceivingAct: 11,
        Quantity:     24,
        Date:         30,
        DateFormat:   "01022006",
        Separator:    '\t',
    }

	transactions := godist.LoadTransactionsFromCSV(pathToTransactions, csvInfo)

	supplySystem := godist.NewSupplySystem()
	supplySystem.ReplayTransactions(transactions)

	totalInStock := supplySystem.TotalInStock()
	totalManufactured := supplySystem.TotalManufactured()
	if totalInStock != totalManufactured {
		log.Fatal("The quantity manufactured is not equal to the quantity in stock. This cannot be.")
	} else {
		log.Printf("A total of %d pills were manufactured and shipped.", totalManufactured)
	}
	traces := supplySystem.ExtractTraces()
	log.Printf("Extracted %d traces", len(traces))
	traces.ToJson(pathToJson)
}

func main() {
	ndcDirectory := os.Args[1]
	jsonDirectory := os.Args[2]

	ndcFiles, err := os.ReadDir(ndcDirectory)

	if err != nil {
		log.Fatal(err)
	}

	//var wg sync.WaitGroup
	for _, ndcFile := range ndcFiles {

		drugName := ndcFile.Name()
		jsonPath := filepath.Join(jsonDirectory, drugName+".json")
		ndcPath := filepath.Join(ndcDirectory, drugName)

		worker(ndcPath, jsonPath)
		//wg.Add(1)
		//go func(ndcPath, jsonPath string) {
		//	log.Println("Extracting: ", drugName)
		//	worker(ndcPath, jsonPath)
		//	wg.Done()
		//}(ndcPath, jsonPath)
	}

	//wg.Wait()
}
