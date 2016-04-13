package simulation

import (
	"testing"
	"time"

	"github.com/blendlabs/go-assert"
)

func TestStationExpectedPassengersInQuantum(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.CalculateTotalAverageRidership()

	station := sim.Stations[8] //times square

	perHour := station.ExpectedPassengersInQuantum(sim.Provider, 1*time.Hour)
	perSecond := station.ExpectedPassengersInQuantum(sim.Provider, 1*time.Second)
	assert.True(perHour > 0)
	assert.True(perSecond > 0)
	assert.True(perHour > perSecond)
}
