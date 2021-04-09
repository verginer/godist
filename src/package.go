package godist

import (
	"fmt"
	"strings"
	"time"
)

// Package is the container which is used to transport the pills
type Package struct {
	quantity         int
	parent           *Package
	originalQuantity int
	trace            PackageTrace
	depleted         bool
	packagedDate     time.Time
}

func NewPackage(quantity int, deaId string, customerId string, date time.Time) *Package {
	trace := PackageTrace{deaId, customerId}
	return &Package{
		quantity:         quantity,
		parent:           nil,
		originalQuantity: quantity,
		trace:            trace,
		depleted:         false,
		packagedDate:     date,
	}
}

// DerivePackage creates and returns a new package taking setting the parent
// equal to the package it has take package the content has been taken from. Note
// if parent is nil it will panic. In this case you should use NewPackage
func DerivePackage(quantity int, parent *Package, destinationDeaID string, date time.Time) *Package {
	if parent == nil {
		panic(1)
	}
	trace := parent.trace.Extend(destinationDeaID)
	return &Package{
		quantity:         quantity,
		parent:           parent,
		originalQuantity: quantity,
		trace:            trace,
		depleted:         false,
		packagedDate:     date,
	}
}

// Take creates a new package with the given quantity if available
// the residual quantity otherwise.
// Trying to taking something from a depleted package will result in a panic.
// If the demanded quantity cannot be satisfied take will return false and true otherwise.
func (p *Package) Take(quantity int, destinationDea string, date time.Time) (*Package, bool) {
	var takenQuantity int
	var satisfied bool

	if p.depleted || (p.quantity < 0) {
		panic(1)
	}

	if p.quantity >= quantity {
		p.quantity -= quantity
		takenQuantity = quantity
		if p.quantity == 0 {
			p.depleted = true
		}
		satisfied = true
	} else {
		takenQuantity = p.quantity
		p.quantity = 0
		p.depleted = true
		satisfied = false
	}
	takenPackage := DerivePackage(takenQuantity, p, destinationDea, date)
	return takenPackage, satisfied
}

// Ancestors returns an array of Package to the parent of the given node.
func (p *Package) Ancestors() []*Package {
	return p.iterativeAncestors(nil)
}

func (p *Package) iterativeAncestors(children []*Package) []*Package {
	ancestralLine := append(children, p)
	if p.parent != nil {
		return p.parent.iterativeAncestors(ancestralLine)
	} else {
		return append(children, p)
	}
}

func (p *Package) String() string {
	repr := fmt.Sprintf("P(%s, q=%d)", p.trace, p.quantity)
	return repr
}

// Traces

// PackageTrace represents the sequence of locations the package has passed through
type PackageTrace []string

func (t PackageTrace) String() string {
	return strings.Join(t, "|")
}

func (t PackageTrace) Extend(nextID string) PackageTrace {
	newSequence := make(PackageTrace, len(t), len(t)+1)
	copy(newSequence, t)
	newSequence = append(newSequence, nextID)
	return newSequence
}
