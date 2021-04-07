package godist

import "testing"

func TestSupplySystem_ShipStock(t *testing.T) {
	system := NewSupplySystem()
	transactionDate := NewTransactionDate("12122021")

	system.ShipStock("M1", "M", "D1", "D", 10, transactionDate)
	system.ShipStock("M1", "M", "D2", "D", 5, transactionDate)
	system.ShipStock("D1", "M", "D2", "D", 3, transactionDate)
	system.ShipStock("M1", "M", "D2", "D", 2, transactionDate)
}

