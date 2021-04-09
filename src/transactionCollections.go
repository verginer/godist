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

type Transaction struct {
	sendingId    string
	sendingAct   string
	receivingId  string
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
    if t[i].date.Before(t[j].date) {
        return true
    } else if t[i].date.After(t[j].date) {
        return false
    }
    // Otherwise the dates are equal

    if t[i].sendingId < t[j].sendingId {
        return true
    } else if  t[i].sendingId > t[j].sendingId {
        return false
    }
    // Otherwise they are equal

    return t[i].receivingId < t[j].receivingId
}


//func (t TransactionsCollection) Less(i, j int) bool {
//	return t[i].date.Before(t[j].date)
//}

func (t TransactionsCollection) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
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
