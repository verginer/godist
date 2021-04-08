package godist

import (
    "log"
    "strings"
	"time"
)

type Distributor struct {
	store        []*Package
	deaID        string
	deaAct       string
	manufactured int
}

func NewDistributor(deaID string, deaAct string) *Distributor {
	return &Distributor{deaID: deaID, deaAct: deaAct}
}

func (d *Distributor) addPackages(packs []*Package) {
	d.store = append(d.store, packs...)
}

func (d *Distributor) removeDepletedPackages() {
    var newStore []*Package

	for _, pack := range d.store {
		if !pack.depleted {
			newStore = append(newStore, pack)
		}
	}
	d.store = newStore
}

func (d *Distributor) TotalStock() int {
	var total int
	for _, pack := range d.store {
		total += pack.quantity
	}
	return total
}

func (d *Distributor) manufacture(quantity int, customerId string, date time.Time) *Package {
	d.manufactured += quantity
	return NewPackage(quantity, d.deaID, customerId, date)
}

func (d *Distributor) preparePackages(requestedQuantity int, customerId string, packagedDate time.Time) []*Package {
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
		finalPackage := d.manufacture(residualDemand, customerId, packagedDate)
		packages = append(packages, finalPackage)
	}

	// Clean up depleted packages
	d.removeDepletedPackages()

	return packages
}

func (d *Distributor) extractTraces(traceCount map[string]int) {
	for _, pack := range d.store {
		traceKey := strings.Join(pack.trace.sequence, "|")
		traceCount[traceKey] += pack.quantity
	}
}
