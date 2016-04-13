package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func NewStation(name string, averageRidership int, generalPopulation *QueueOfPassenger) *Station {
	return &Station{
		Name:               name,
		WaitingPassengers:  NewQueueOfPassenger(),
		GeneralPopulation:  generalPopulation,
		RidersPerDayMean:   averageRidership,
		RidersPerDayStdDev: math.Sqrt(float64(averageRidership) * 0.5), //totally bogus
	}
}

type Station struct {
	Name               string
	RidersPerDayMean   int
	RidersPerDayStdDev float64

	WaitingPassengers *QueueOfPassenger
	GeneralPopulation *QueueOfPassenger

	OutBoundTrain *Train
	InBoundTrain  *Train

	OutBoundTrack *Track
	InBoundTrack  *Track
}

func (s *Station) LinkWith(next *Station, distanceMeters float64) {
	s.OutBoundTrack = &Track{
		IsOutBound:     true,
		Begin:          s,
		End:            next,
		DistanceMeters: distanceMeters,
	}
	next.InBoundTrack = &Track{
		IsOutBound:     false,
		End:            s,
		Begin:          next,
		DistanceMeters: distanceMeters,
	}
}

func (s *Station) IsOutboundTerminus() bool {
	return s.OutBoundTrack == nil
}

func (s *Station) IsInboundTerminus() bool {
	return s.InBoundTrack == nil
}

func (s *Station) CheckForCollision(train *Train) {
	if s.OutBoundTrain != nil && train.IsOutbound {
		panic(fmt.Sprintf("Out Bound Train [%d] collides with Train [%d] at %s", train.ID, s.OutBoundTrain.ID, s.Name))
	}

	if s.InBoundTrain != nil && !train.IsOutbound {
		panic(fmt.Sprintf("In Bound Train [%d] collides with Train [%d] at %s", train.ID, s.InBoundTrain.ID, s.Name))
	}
}

func (s *Station) TrainEnters(train *Train) {
	s.CheckForCollision(train)

	if s.IsOutboundTerminus() { // flip the train around
		train.IsOutbound = false
	}

	if train.IsOutbound {
		s.OutBoundTrain = train
	} else {
		s.InBoundTrain = train
	}
}

func (s *Station) CheckWaitingTrains(wallClock time.Duration) {
	if s.OutBoundTrain != nil {
		s.OutBoundTrain.Depart(wallClock, s)
	}

	if s.InBoundTrain != nil {
		s.InBoundTrain.Depart(wallClock, s)
	}
}

func (s *Station) TrainDeparts(train *Train) {
	if train.IsOutbound {
		if s.OutBoundTrack != nil {
			s.OutBoundTrack.AddTrain(train)
		}

		s.OutBoundTrain = nil

	} else {
		if s.InBoundTrack != nil {
			s.InBoundTrack.AddTrain(train)
		}
		s.InBoundTrain = nil
	}
}

func (s *Station) ExpectedPassengersInQuantum(provider *rand.Rand, stepLength time.Duration) float64 {
	q := float64(time.Hour/stepLength) * 24.0
	u := float64(s.RidersPerDayMean) / q
	o := float64(s.RidersPerDayStdDev) / q

	return (provider.NormFloat64() * o) + u
}
