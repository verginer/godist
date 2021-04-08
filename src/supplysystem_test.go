package godist

import (
    "github.com/stretchr/testify/assert"
    "log"
    "testing"
)

func TestSupplySystem_ShipStock(t *testing.T) {
    system := NewSupplySystem()

    ts := TransactionsCollection{
        Transaction{
            sendingId:   "M1",
            receivingId: "D1",
            quantity:    10,
        },
        Transaction{
            sendingId:   "D1",
            receivingId: "D2",
            quantity:    8,
        },
        Transaction{
            sendingId:   "D1",
            receivingId: "D3",
            quantity:    6,
        },
        Transaction{
            sendingId:   "D3",
            receivingId: "D4",
            quantity:    10,
        },
    }

    for _, trans := range ts {
        err := system.ShipStock(trans)
        if err != nil {
            log.Fatal(err)
        }
    }

    traces := system.ExtractTraces()
    assert.Equal(t, system.TotalInStock(), system.TotalManufactured())
    assert.Equal(t, len(traces), 4)

}
