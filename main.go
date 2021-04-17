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
	err := supplySystem.ReplayTransactions(transactions)
	godist.LogFatalOnError(err)

	totalInStock := supplySystem.TotalInStock()
	totalManufactured := supplySystem.TotalManufactured()

	if totalInStock != totalManufactured {
		log.Fatal("The quantity manufactured is not equal to the quantity in stock. This cannot be.")
	}
	log.Printf("A total of %d pills were manufactured and shipped.", totalManufactured)

	traces := supplySystem.ExtractTraces()
	log.Printf("Extracted %d traces", len(traces))
	traces.ToJson(pathToJson)
}

func main() {
	ndcDirectory := os.Args[1]
	jsonDirectory := os.Args[2]

	ndcFiles, err := os.ReadDir(ndcDirectory)

	godist.LogFatalOnError(err)

	for _, ndcFile := range ndcFiles {
		drugName := ndcFile.Name()
		jsonPath := filepath.Join(jsonDirectory, drugName+".json")
		ndcPath := filepath.Join(ndcDirectory, drugName)

		worker(ndcPath, jsonPath)
	}
}
