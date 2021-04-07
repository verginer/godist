package godist

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type Transaction struct {
	sendingID    string
	sendingAct   string
	receivingID  string
	receivingAct string
	date         time.Time
	quantity     int
}

type TransactionsCollection []Transaction

// Implement the interface to use sort.Sort
func (t TransactionsCollection) Len() int {
	return len(t)
}

func (t TransactionsCollection) Less(i, j int) bool {
	return t[i].date.Before(t[j].date)
}

func (t TransactionsCollection) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func LoadTransactionsFromCSV(transactionsPath string, sortByDate bool) TransactionsCollection {
	nTransactions, _ := lineCounter(transactionsPath)
	transactions := make(TransactionsCollection, nTransactions)

	csvFile, err := os.Open(transactionsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	r.Comma = '\t'
	r.LazyQuotes = true

	for i := 0; true; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		quantity, err := strconv.ParseFloat(record[24], 64)
		if err != nil {
			log.Fatal(err)
		}

		date := NewTransactionDate(record[30])

		transactions[i] = Transaction{
			sendingID:    record[0],
			sendingAct:   record[1],
			receivingID:  record[10],
			receivingAct: record[11],
			date:         date,
			quantity:     int(quantity),
		}
	}
	if sortByDate {
		sort.Sort(transactions)
	}
	return transactions
}
