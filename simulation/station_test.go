package simulation

import (
	"testing"
	"time"

	"github.com/blendlabs/go-assert"
)

func TestStationPassengerArrivalPDF(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.CalculateTotalAverageRidership()

	station := sim.Stations[8] //times square

	pdf := station.PassengerArrivalPDF(sim.Provider, 1*time.Second)
	assert.True(pdf > 0, pdf)
	assert.True(pdf <= 1.0, pdf)
}
