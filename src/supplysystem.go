package godist

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
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

func (s *SupplySystem) ShipStock(sourceId, sourceAct, targetId, targetAct string, quantity int, date time.Time) error {
	if date.Before(s.currentDate) {
		return fmt.Errorf(
			"transaction from %s to %s on %s happened before current date %s",
			sourceId, targetId, date, s.currentDate,
		)
	}

	s.currentDate = date
	sendingDistributor := s.GetDistributor(sourceId, sourceAct)
	receivingDistributor := s.GetDistributor(targetId, targetAct)

	packages := sendingDistributor.preparePackages(quantity, receivingDistributor.deaID, date)
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

func ReplayTransactionsFromFile(transactionsPath string) *SupplySystem {
	supSystem := NewSupplySystem()
	nTransactions, err := lineCounter(transactionsPath)

	if err != nil {
		log.Fatal(err)
	}

	transactions := LoadTransactionsFromCSV(transactionsPath, true)

	pBar := progressbar.Default(int64(nTransactions))
	log.Println("Replaying TransactionsCollection")
	for _, t := range transactions {
		err := supSystem.ShipStock(t.sendingID, t.sendingAct, t.receivingID, t.receivingAct, t.quantity, t.date)
		if err != nil {
			log.Fatal(err)
		}
		_ = pBar.Add(1)
	}
	return supSystem
}
