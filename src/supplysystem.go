package godist

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type SupplySystem struct {
	currentDate  time.Time
	distributors map[string]*Distributor
}

func NewSupplySystem() *SupplySystem {
	return &SupplySystem{distributors: make(map[string]*Distributor)}
}

// ReplayTransactions replays the transaction record to replicate the real shipment events
func (s *SupplySystem) ReplayTransactions(transactions TransactionsCollection) error {
	for _, t := range transactions {
		err := s.ShipStock(t)
		if err != nil {
			log.Fatal(err)
		}
	}

	manufactured := s.TotalManufactured()
	inStock := s.TotalInStock()
	if manufactured != inStock {
		return fmt.Errorf("%d Units manufcaturd, but only %d in stock, something is wrong", manufactured, inStock)
	}
	return nil
}

// ShipStock applies a given transaction to the system by moving stock from
// source to destination, if a date happening before the lastest transaction is
// encountered an error is raised.
func (s *SupplySystem) ShipStock(t Transaction) error {
	if t.date.Before(s.currentDate) {
		return fmt.Errorf(
			"transaction from %s to %s on %s happened before current Date %s",
			t.sendingId, t.receivingId, t.date, s.currentDate,
		)
	}

	s.currentDate = t.date
	sendingDistributor := s.GetDistributor(t.sendingId, t.sendingAct)
	receivingDistributor := s.GetDistributor(t.receivingId, t.receivingAct)

	packages := sendingDistributor.PreparePackages(t.quantity, receivingDistributor.deaId, t.date)
	receivingDistributor.AddPackages(packages)
	return nil
}

// GetDistributor returns the given distributor if it does not exist it will create it and add it to the store
func (s *SupplySystem) GetDistributor(deaId, deaAct string) *Distributor {
	focalDistributor := s.distributors[deaId]
	if focalDistributor == nil {
		s.AddDistributor(deaId, deaAct)
		focalDistributor = s.distributors[deaId]
	}
	return focalDistributor
}

func (s *SupplySystem) AddDistributor(deaId, deaAct string) {
	s.distributors[deaId] = NewDistributor(deaId, deaAct)
}

func (s *SupplySystem) TotalManufactured() int {
	var total int
	for _, dist := range s.distributors {
		total += dist.manufactured
	}
	return total
}

// Traces
type AggregateTrace map[string]int

func (a AggregateTrace) Sum() int {
	total := 0
	for _, q := range a {
		total += q
	}
	return total
}

func (s *SupplySystem) ExtractTraces() AggregateTrace {
	allTraces := make(AggregateTrace)
	for _, dist := range s.distributors {
		dist.ExtractTraces(allTraces)
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

func (a AggregateTrace) ToJson(path string) {
	jsonData, err := json.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}

	jsonFile, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}
