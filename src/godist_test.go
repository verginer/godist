package godist

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestReplay(t *testing.T) {

	//csvInfo := TransactionCSVInfo{
	//	SendingId:    0,
	//	SendingAct:   0,
	//	ReceivingId:  1,
	//	ReceivingAct: 1,
	//	Date:         3,
	//	Quantity:     2,
	//	DateFormat:   "01022006",
	//	Separator: ',',
	//}

    //csvInfo := TransactionCSVInfo{
    //   SendingId:    0,
    //   SendingAct:   1,
    //   ReceivingId:  2,
    //   ReceivingAct: 3,
    //   Date:         5,
    //   Quantity:     4,
    //   DateFormat:   "2006-01-02",
    //}


    csvInfo := TransactionCSVInfo{
        SendingId:    0,
        SendingAct:   1,
        ReceivingId:  10,
        ReceivingAct: 11,
        Quantity:     24,
        Date:         30,
        DateFormat:   "01022006",
        Separator:    '\t',
    }


    //testDatPath := "/Users/lucaverginer/Downloads/presorted_transactions_example.csv"
	//testDatPath := "/Users/lucaverginer/Downloads/random_transactions.csv"
    //testDatPath := "./testdata/synth_transactions.tsv"
    testDatPath := "../input2/004091255"
	opioidTransactions := LoadTransactionsFromCSV(testDatPath, csvInfo)

	distSystem := NewSupplySystem()
	distSystem.ReplayTransactions(opioidTransactions)

	assert.Equal(t, distSystem.TotalManufactured(), distSystem.TotalInStock())

	traces := distSystem.ExtractTraces()
	assert.Equal(t, distSystem.TotalManufactured(), traces.Sum())

	traces.ToJson("../tmp/gotraces.json")
}
