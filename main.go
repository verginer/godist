package main

import (
	godist "github.com/verginer/godist/src"
	"log"
	"os"
)

func main() {
	pathToTransactions := os.Args[1]
	pathToJson := os.Args[2]

	//pathToTransactions := "/Users/lucaverginer/Research/Projects/distpath/data/interim/transactions_by_ndc/621750490"
	//pathToJson := "out.json"
	log.Print("Reading:", pathToTransactions)
	supplySystem := godist.ReplayTransactionsFromFile(pathToTransactions)
	traces := supplySystem.ExtractTraces()
	godist.WriteTracesToJson(traces, pathToJson)
}
