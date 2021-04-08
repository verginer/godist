package godist

import (
	"fmt"
	"log"
	"time"
)

type SupplySystem struct {
	currentDate  time.Time
	distributors map[string]*Distributor
}

func NewSupplySystem() *SupplySystem {
	return &SupplySystem{distributors: make(map[string]*Distributor)}
}

func (s *SupplySystem) AddDistributor(deaId, deaAct string) {
	s.distributors[deaId] = NewDistributor(deaId, deaAct)
}

func (s *SupplySystem) GetDistributor(deaId, deaAct string) *Distributor {
	focalDistributor := s.distributors[deaId]
	if focalDistributor == nil {
		s.AddDistributor(deaId, deaAct)
		focalDistributor = s.distributors[deaId]
	}
	return focalDistributor
}

//func (s *SupplySystem) ShipStock(sourceId, sourceAct, targetId, targetAct string, quantity int, date time.Time) error {
func (s *SupplySystem) ShipStock(t Transaction) error {
	if t.date.Before(s.currentDate) {
		return fmt.Errorf(
			"transaction from %s to %s on %s happened before current date %s",
			t.sendingId, t.receivingId, t.date, s.currentDate,
		)
	}

	s.currentDate = t.date
	sendingDistributor := s.GetDistributor(t.sendingId, t.sendingAct)
	receivingDistributor := s.GetDistributor(t.receivingId, t.receivingAct)

	packages := sendingDistributor.preparePackages(t.quantity, receivingDistributor.deaID, t.date)
	receivingDistributor.addPackages(packages)
	return nil
}

func (s *SupplySystem) ExtractTraces() map[string]int {
	allTraces := make(map[string]int)
	for _, dist := range s.distributors {
		dist.extractTraces(allTraces)
	}
	return allTraces
}

func (s *SupplySystem) TotalInStock() int {
    var totSock int
    for _, dist := range s.distributors {
        totSock += dist.TotalStock()
    }
    return totSock
}

func (s *SupplySystem) TotalManufactured() int {
    var totManuf int
    for _, dist := range s.distributors {
        totManuf += dist.manufactured
    }
    return totManuf
}

func ReplayTransactionsFromFile(transactionsPath string) *SupplySystem {
	supSystem := NewSupplySystem()

	transactions := LoadTransactionsFromCSV(transactionsPath)

	log.Println("Replaying TransactionsCollection")
	for _, t := range transactions {
		err := supSystem.ShipStock(t)
		if err != nil {
			log.Fatal(err)
}
	}
	return supSystem
}
