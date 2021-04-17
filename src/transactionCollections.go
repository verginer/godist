package godist

import (
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

//func (t TransactionsCollection) Less(i, j int) bool {
//    return t[i].date.Before(t[j].date)
//}

func (t TransactionsCollection) Less(i, j int) bool {
	if t[i].date.Before(t[j].date) {
		return true
	} else if t[i].date.After(t[j].date) {
		return false
	}

	if t[i].quantity > t[j].quantity {
		return true
	} else if t[i].quantity < t[j].quantity {
		return false
	}

	return t[i].sendingId < t[j].sendingId
}

func (t TransactionsCollection) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
