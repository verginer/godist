package godist

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func LoadTransactionsFromCSV(transactionsPath string, csvInfo TransactionCSVInfo) TransactionsCollection {

	nTransactions, _ := lineCounter(transactionsPath)
	transactions := make(TransactionsCollection, nTransactions)

	csvFile, err := os.Open(transactionsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	r.Comma = csvInfo.Separator
	r.LazyQuotes = true

	for i := 0; true; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		quantityString := record[csvInfo.Quantity]
		quantityString = strings.Split(quantityString, ".")[0]
		quantity, err := strconv.ParseInt(quantityString, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		date := NewTransactionDate(record[csvInfo.Date], csvInfo.DateFormat)

		transactions[i] = Transaction{
			sendingId:    record[csvInfo.SendingId],
			sendingAct:   record[csvInfo.SendingAct],
			receivingId:  record[csvInfo.ReceivingId],
			receivingAct: record[csvInfo.ReceivingAct],
			date:         date,
			quantity:     int(quantity),
		}
	}

	sort.Sort(transactions)

	return transactions
}

type TransactionCSVInfo struct {
	SendingId    int
	SendingAct   int
	ReceivingId  int
	ReceivingAct int
	Quantity     int
	Date         int
	DateFormat   string
	Separator    rune
}

func NewTransactionDate(date, dateFormat string) time.Time {
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Fatal(err)
	}
	return t
}
