package godist

import (
	"fmt"
	"testing"
)
import "github.com/stretchr/testify/assert"



func CreatePackageWithTrace(depth int) *Package {
	startDate := NewTransactionDate("01012012", DateFormat)
	var myPack = NewPackage(10, "M1", "D0", startDate)

	takenPackage, _ := myPack.Take(5, "D1", startDate)
	for i := 3; i <= depth; i++ {
		destination := "D" + fmt.Sprint(i)
		takenPackage, _ = takenPackage.Take(1, destination, startDate)
	}

	return takenPackage
}

func TestPackage_Take(t *testing.T) {
	startDate := NewTransactionDate("01012012", DateFormat)
	var myPack = NewPackage(30, "M1", "D0", startDate)
	takenPackage, complete := myPack.Take(15, "D1", startDate)
	assert.Equal(t, complete, true)
	assert.Equal(t, takenPackage.Quantity(), myPack.Quantity())
	takenPackage, complete = myPack.Take(30, "D2", startDate)
	assert.Equal(t, complete, false)
	assert.Equal(t, 15, takenPackage.quantity)
}

func TestPackage_Ancestors(t *testing.T) {
	startDate := NewTransactionDate("01012021", DateFormat)
	var myPack = NewPackage(10, "M1", "D0", startDate)

	takenPackage, _ := myPack.Take(5, "D1", startDate)
	for i := 3; i <= 20; i++ {
		destination := "D" + fmt.Sprint(i)
		takenPackage, _ = takenPackage.Take(1, destination, startDate)
	}
	var ancestors = takenPackage.Ancestors()
	assert.Equal(t, 20, len(ancestors))

	var quantityInAncestors int
	for _, pack := range ancestors {
		quantityInAncestors += pack.quantity
	}
	assert.Equal(t, 10, quantityInAncestors)
}

func TestPackage_Depleted(t *testing.T) {
	startDate := NewTransactionDate("01012021", DateFormat)
	var myPack = NewPackage(10, "M1", "D0", startDate)

	takenPackage, satisfied := myPack.Take(20, "D0", startDate)
	assert.False(t, satisfied)
	assert.Equal(t, 10, takenPackage.quantity)

	assert.Panics(t, func() {
		myPack.Take(20, "D0", startDate)
	})
}

func BenchmarkName(b *testing.B) {
	finalPackage := CreatePackageWithTrace(100)
	for i := 0; i < b.N; i++ {
		finalPackage.Ancestors()
	}
}
