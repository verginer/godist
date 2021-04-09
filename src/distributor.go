package godist

import (
	"log"
	"time"
)

//
type Distributor struct {
	store        []*Package
	deaId        string
	deaAct       string
	manufactured int
}

func NewDistributor(deaID string, deaAct string) *Distributor {
	return &Distributor{deaId: deaID, deaAct: deaAct}
}

// AddPackages appends the packages to the end of the store
func (d *Distributor) AddPackages(packs []*Package) {
	d.store = append(d.store, packs...)
}

// RemoveDepletedPackages removes the packages from the store if depleted==true
func (d *Distributor) RemoveDepletedPackages() {
	var newStore []*Package

	for _, pack := range d.store {
		if !pack.depleted {
			newStore = append(newStore, pack)
		}
	}
	d.store = newStore
}

// PreparePackages returns a list of packages that satisfy the quantity from the earliest to the newest.
func (d *Distributor) PreparePackages(requestedQuantity int, customerId string, packagedDate time.Time) []*Package {
	var packages []*Package
	residualDemand := requestedQuantity
	for _, pack := range d.store {
		preparedPackage, _ := pack.Take(residualDemand, customerId, packagedDate)
		residualDemand -= preparedPackage.quantity
		packages = append(packages, preparedPackage)
		if residualDemand < 0 {
			log.Fatal("The residual demand is negative")
		}
		if residualDemand == 0 {
			break
		}
	}

	// if there is still residual demand manufacture it
	if residualDemand > 0 {
		finalPackage := d.Manufacture(residualDemand, customerId, packagedDate)
		packages = append(packages, finalPackage)
	}

	// Clean up depleted packages
	d.RemoveDepletedPackages()
	return packages
}

// Manufacture creates a pristine package with trace set to his deaId and the customerId
func (d *Distributor) Manufacture(quantity int, customerId string, date time.Time) *Package {
	d.manufactured += quantity
	return NewPackage(quantity, d.deaId, customerId, date)
}

// ExtractTraces will return the aggregated traces counting the quantity flowing on it.
func (d *Distributor) ExtractTraces(traceCount AggregateTrace) {
	for _, pack := range d.store {

		finalNode := pack.trace[len(pack.trace)-1]
		if finalNode != d.deaId {
			panic(1)
		}

		traceKey := pack.trace.String()
		traceCount[traceKey] += pack.quantity
	}
}

func (d *Distributor) TotalStock() int {
	var total int
	for _, pack := range d.store {
		total += pack.quantity
	}
	return total
}
