package simulation

import (
	"testing"
	"time"

	"github.com/blendlabs/go-assert"
)

func createTestSimulation() *Simulation {
	sim := New(1*time.Second, 1*time.Hour, nil)
	sim.TotalPassengerCount = 128
	sim.TotalTrainCount = 32
	sim.Stasis = true
	sim.GeneratePassengers()
	sim.GenerateStations()
	sim.GenerateTrains()
	return sim
}

func TestSimulationNew(t *testing.T) {
	assert := assert.New(t)
	sim := New(1*time.Second, 1*time.Hour, nil)
	assert.NotNil(sim)
	assert.Equal(1*time.Second, sim.StepLength)
	assert.Equal(1*time.Hour, sim.TotalTime)
	assert.Nil(sim.PauseTime)
}

func TestSimulationLogf(t *testing.T) {
	assert := assert.New(t)
	sim := New(1*time.Second, 1*time.Hour, nil)
	sim.logf("test")
	assert.NotEmpty(sim.LogEntries)
}

func TestSimulationGeneratePassengers(t *testing.T) {
	assert := assert.New(t)
	sim := New(1*time.Second, 1*time.Hour, nil)
	sim.TotalPassengerCount = 32
	sim.GeneratePassengers()
	assert.Equal(32, sim.People.Len())
}

func TestSimulationGenerateStations(t *testing.T) {
	assert := assert.New(t)

	sim := New(1*time.Second, 1*time.Hour, nil)
	sim.GenerateStations()

	//walk the track forward
	var index int
	var station *Station = sim.InBoundTerminus()
	for station.OutBoundTrack != nil {
		assert.NotNil(station.OutBoundTrack.End)

		if index != 0 {
			assert.False(station.IsInboundTerminus())
		}
		if index != len(sim.Stations)-1 {
			assert.False(station.IsOutboundTerminus())
		}

		station = station.OutBoundTrack.End
		index++
	}

	station = sim.OutBoundTerminus()
	for station.InBoundTrack != nil {
		assert.NotNil(station.InBoundTrack.End)

		if index != 0 {
			assert.False(station.IsInboundTerminus())
		}
		if index != len(sim.Stations)-1 {
			assert.False(station.IsOutboundTerminus())
		}

		station = station.InBoundTrack.End
		index--
	}
}

func TestSimulationGenerateTrains(t *testing.T) {
	assert := assert.New(t)
	sim := New(1*time.Second, 1*time.Hour, nil)
	sim.TotalTrainCount = 32
	sim.GenerateTrains()
	assert.Equal(32, sim.Yard.Len())
}

func TestSimulationCalculateTotalAverageRidership(t *testing.T) {
	assert := assert.New(t)

	sim := New(1*time.Second, 1*time.Hour, nil)
	sim.GenerateStations()
	sim.CalculateTotalAverageRidership()
	assert.NotZero(sim.TotalAverageRidership)
}

func TestSimulationShouldReleaseTrainFromYard(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.AverageTimeBetweenTrains = 1 * time.Minute
	sim.LastTrainReleased = 1 * time.Minute
	sim.WallClock = 1 * time.Hour

	assert.True(sim.ShouldReleaseTrainFromYard())
}

func TestSimulationShouldReleaseTrainFromYardBelowAverageTime(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.AverageTimeBetweenTrains = 5 * time.Minute
	sim.LastTrainReleased = 1 * time.Minute
	sim.WallClock = 2 * time.Minute

	assert.False(sim.ShouldReleaseTrainFromYard())
}

func TestSimulationShouldReleaseTrainFromYardComplete(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.AverageTimeBetweenTrains = 1 * time.Minute
	sim.LastTrainReleased = 1 * time.Minute
	sim.WallClock = 2 * time.Hour
	sim.Complete = true

	assert.False(sim.ShouldReleaseTrainFromYard())
}

func TestSimulationShouldReleaseTrainFromYardEmptyYard(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.AverageTimeBetweenTrains = 1 * time.Minute
	sim.LastTrainReleased = 1 * time.Minute
	sim.WallClock = 2 * time.Hour
	sim.Yard = NewQueueOfTrain()

	assert.False(sim.ShouldReleaseTrainFromYard())
}

func TestSimulationStopsToDestination(t *testing.T) {
	assert := assert.New(t)

	sim := createTestSimulation()

	wallStreet := sim.Stations[14]
	home := "14 Street"

	assert.Equal(4, sim.StopsToDestination(false, wallStreet, home))
	assert.Equal(len(sim.Stations)+1, sim.StopsToDestination(true, wallStreet, home))
}

func TestSimulationDestinationIsOutBound(t *testing.T) {
	assert := assert.New(t)

	sim := createTestSimulation()

	wallStreet := sim.Stations[14]
	home := "14 Street"
	clarkStreet := "Clark Street"

	assert.False(sim.DestinationIsOutbound(wallStreet, home))
	assert.True(sim.DestinationIsOutbound(wallStreet, clarkStreet))
}

func TestSimulationDestinationIsOutBoundInBoundTerminus(t *testing.T) {
	assert := assert.New(t)

	sim := createTestSimulation()

	terminus := sim.InBoundTerminus()
	home := "14 Street"
	clarkStreet := "Clark Street"
	outboundTerminus := "New Lots Avenue"

	assert.True(sim.DestinationIsOutbound(terminus, home))
	assert.True(sim.DestinationIsOutbound(terminus, clarkStreet))
	assert.True(sim.DestinationIsOutbound(terminus, outboundTerminus))
}

func TestSimulationDestinationIsOutBoundOutBoundTerminus(t *testing.T) {
	assert := assert.New(t)

	sim := createTestSimulation()

	terminus := sim.OutBoundTerminus()
	home := "14 Street"
	clarkStreet := "Clark Street"
	inboundTerminus := "Harlem-148 Street"

	assert.False(sim.DestinationIsOutbound(terminus, home))
	assert.False(sim.DestinationIsOutbound(terminus, clarkStreet))
	assert.False(sim.DestinationIsOutbound(terminus, inboundTerminus))
}

func TestSimulationStationIncident(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.StationIncidentLikelihood = float64(1 << 20) //this should cause a problem.
	timeSquare := sim.Stations[8]
	train := sim.Yard.Dequeue()
	train.HasLeftYard(sim.WallClock)
	timeSquare.OutBoundTrain = train
	sim.StationIncident(timeSquare)
	assert.Equal(SignalHold, train.Signal)
}

func TestSimulationStationIncidentComplete(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.Stasis = true
	sim.Complete = true
	sim.StationIncidentLikelihood = float64(1 << 20) //this should cause a problem.
	timeSquare := sim.Stations[8]
	train := sim.Yard.Dequeue()
	train.HasLeftYard(sim.WallClock)
	timeSquare.OutBoundTrain = train
	sim.StationIncident(timeSquare)
	assert.Equal(SignalGo, train.Signal)
}

func TestSimulationStationIncidentZeroLikelihood(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.StationIncidentLikelihood = 0.0 //this should cause a problem.
	timeSquare := sim.Stations[8]
	train := sim.Yard.Dequeue()
	train.HasLeftYard(sim.WallClock)
	timeSquare.OutBoundTrain = train
	sim.StationIncident(timeSquare)
	assert.Equal(SignalGo, train.Signal)
}

func TestSimulationReleaseTrainsOnHold(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.StationIncidentLikelihood = float64(1 << 20) //this should cause a problem.
	timeSquare := sim.Stations[8]
	train := sim.Yard.Dequeue()
	train.HasLeftYard(sim.WallClock)
	timeSquare.OutBoundTrain = train
	sim.StationIncident(timeSquare)
	assert.Equal(SignalHold, train.Signal)

	train.HeldAtStation = 5 * time.Minute
	sim.AverageIncidentDelay = 1 * time.Minute
	sim.WallClock = 10 * time.Minute

	sim.ReleaseTrainsOnHold(timeSquare)
	assert.Equal(SignalGo, train.Signal)
}

func TestSimulationReleaseTrainsOnHoldBeforeAverageDelay(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.StationIncidentLikelihood = float64(1 << 20) //this should cause a problem.
	timeSquare := sim.Stations[8]
	train := sim.Yard.Dequeue()
	train.HasLeftYard(sim.WallClock)
	timeSquare.OutBoundTrain = train
	sim.StationIncident(timeSquare)
	assert.Equal(SignalHold, train.Signal)

	train.HeldAtStation = 5 * time.Minute
	sim.AverageIncidentDelay = 3 * time.Minute
	sim.WallClock = 6 * time.Minute

	sim.ReleaseTrainsOnHold(timeSquare)
	assert.Equal(SignalHold, train.Signal)
}

func TestSimulationPassengerArrivesAtStation(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	p := sim.People.Dequeue()
	station := sim.Stations[8]
	sim.PassengerArrivesAtStation(station, p)
	assert.NotEmpty(p.Destination)
	assert.NotZero(station.WaitingPassengers.Len())
}

func TestSimulationPassengerArrivesAtStationWithWaitingTrain(t *testing.T) {
	assert := assert.New(t)

	sim := createTestSimulation()
	p := sim.People.Dequeue()
	station := sim.Stations[8]
	station.OutBoundTrain = sim.Yard.Dequeue()
	station.InBoundTrain = sim.Yard.Dequeue()
	sim.PassengerArrivesAtStation(station, p)
	assert.NotEmpty(p.Destination)
	assert.Zero(station.WaitingPassengers.Len())

	if p.IsOutBound {
		assert.NotEmpty(station.OutBoundTrain.Passengers)
	} else {
		assert.NotEmpty(station.InBoundTrain.Passengers)
	}
}

func TestSimulationPassengersArrive(t *testing.T) {
	assert := assert.New(t)
	sim := createTestSimulation()
	sim.StepLength = 1 * time.Hour
	station := sim.Stations[8]
	sim.PassengersArrive(station)
	assert.NotZero(station.WaitingPassengers.Len())
}
