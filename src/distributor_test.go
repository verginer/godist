package godist

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	DateFormat = "01022006"
	//DateFormat = "2006-01-02"
)

func TestDistributor_AddPackage(t *testing.T) {
	startDate := NewTransactionDate("01012012", DateFormat)
	myDist := NewDistributor("M0", "M")
	myPack1 := NewPackage(10, "M1", "D1", startDate)
	myPack2 := NewPackage(5, "M1", "D2", startDate)
	myDist.AddPackages([]*Package{myPack1})
	myDist.AddPackages([]*Package{myPack2})
	assert.Equal(t, 2, len(myDist.store))
	assert.Equal(t, 15, myDist.TotalStock())
}

func TestDistributor_PreparePackages(t *testing.T) {
	startDate := NewTransactionDate("01012021", DateFormat)
	myMan := NewDistributor("M1", "M")
	myDist1 := NewDistributor("D1", "D")
	myDist2 := NewDistributor("D2", "D")

	preparedPackages1 := myMan.PreparePackages(10, myDist1.deaId, startDate)
	preparedPackages2 := myMan.PreparePackages(15, myDist2.deaId, startDate)

	assert.Equal(t, 1, len(preparedPackages1))
	assert.Equal(t, 10, preparedPackages1[0].quantity)
	assert.Equal(t, 25, myMan.manufactured)

	myDist1.AddPackages(preparedPackages1)
	myDist2.AddPackages(preparedPackages2)

	preparedPackages3 := myDist1.PreparePackages(15, myDist2.deaId, startDate)
	myDist2.AddPackages(preparedPackages3)

}
