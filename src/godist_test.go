package godist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var csvInfo = TransactionCSVInfo{
	SendingId:    0,
	SendingAct:   0,
	ReceivingId:  1,
	ReceivingAct: 1,
	Date:         3,
	Quantity:     2,
	DateFormat:   "01022006",
	Separator:    ',',
}

var testDatPath = "./testdata/synth_transactions.tsv"

func TestReplay(t *testing.T) {
	opioidTransactions := LoadTransactionsFromCSV(testDatPath, csvInfo)

	distSystem := NewSupplySystem()
	distSystem.ReplayTransactions(opioidTransactions)

	assert.Equal(t, distSystem.TotalManufactured(), distSystem.TotalInStock())

	traces := distSystem.ExtractTraces()
	assert.Equal(t, distSystem.TotalManufactured(), traces.Sum())

	traces.ToJson("./testdata/synth-traces.json")
}

func ReplayDifferentSizes(transactions TransactionsCollection, nIterations int) {
	distSystem := NewSupplySystem()
	for i := 0; i < nIterations; i++ {
		distSystem.ReplayTransactions(transactions)
	}
	distSystem.ExtractTraces()
}

func BenchmarkSupplySystem_ReplayTransactions1(b *testing.B) {
	opioidTransactions := LoadTransactionsFromCSV(testDatPath, csvInfo)
	for i := 0; i < b.N; i++ {
		ReplayDifferentSizes(opioidTransactions, 1)
	}
}

func BenchmarkSupplySystem_ReplayTransactions10(b *testing.B) {
	opioidTransactions := LoadTransactionsFromCSV(testDatPath, csvInfo)
	for i := 0; i < b.N; i++ {
		ReplayDifferentSizes(opioidTransactions, 10)
	}
}

func BenchmarkSupplySystem_ReplayTransactions100(b *testing.B) {
	opioidTransactions := LoadTransactionsFromCSV(testDatPath, csvInfo)
	for i := 0; i < b.N; i++ {
		ReplayDifferentSizes(opioidTransactions, 100)
	}
}

func BenchmarkSupplySystem_ReplayTransactions1000(b *testing.B) {
	opioidTransactions := LoadTransactionsFromCSV(testDatPath, csvInfo)
	for i := 0; i < b.N; i++ {
		ReplayDifferentSizes(opioidTransactions, 1000)
	}
}
