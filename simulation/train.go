package simulation

import (
	"fmt"
	"time"
)

const (
	SignalHold    Signal = 0
	SignalCaution Signal = 1
	SignalGo      Signal = 2
)

type Signal int

func (s Signal) String() string {
	switch s {
	case SignalGo:
		{
			return "O"
		}
	case SignalCaution:
		{
			return "o"
		}
	case SignalHold:
		{
			return "_"
		}
	}
	return "?"
}

func NewTrain(id int, line string, maximumSpeed, acceleration, braking float64, averageTimeInStation time.Duration) *Train {
	return &Train{
		ID:                   id,
		Line:                 line,
		Capacity:             128,
		MaximumSpeed:         maximumSpeed,
		Acceleration:         acceleration,
		AverageTimeInStation: averageTimeInStation,
		MinumumSafeDistance:  5.0,
		Braking:              braking,
		Signal:               SignalGo,
	}
}

type Train struct {
	ID         int
	Line       string
	IsOutbound bool
	Passengers []*Passenger

	Capacity int

	LeftYard         time.Duration
	ArrivedAtStation time.Duration
	HeldAtStation    time.Duration

	AverageTimeInStation time.Duration

	Signal Signal

	CautionSpeed float64
	MaximumSpeed float64
	Acceleration float64
	Braking      float64

	MinumumSafeDistance float64

	Position float64
	Speed    float64

	DistanceTraveled float64
	Speeds           []float64
	RoundTripTimes   []time.Duration
}

func (t *Train) String() string {
	direction := "↓"
	if !t.IsOutbound {
		direction = "↑"
	}
	if t.ArrivedAtStation != 0 {
		return fmt.Sprintf("%s [%d] %v", direction, t.ID, t.Signal)
	} else {
		return fmt.Sprintf("%s [%d] %0.2f @ %0.2fm/s", direction, t.ID, t.Position, t.Speed)
	}
}

func (t *Train) SendSignal(signal Signal) {
	t.Signal = signal
}

func (t *Train) Hold(wallClock time.Duration) {
	t.Signal = SignalHold
	t.HeldAtStation = wallClock
}

func (t *Train) Accelerate(stepLength time.Duration, targetSpeed float64) {
	stepLengthSeconds := float64(stepLength) / float64(time.Second)
	if t.Speed < targetSpeed {
		t.Speed += (t.Acceleration / stepLengthSeconds)
	}
}

func (t *Train) Decellerate(stepLength time.Duration) {
	stepLengthSeconds := float64(stepLength) / float64(time.Second)
	if t.Speed > 0 {
		t.Speed -= (t.Braking / stepLengthSeconds)
		if t.Speed < 0 {
			t.Speed = 0
		}
	}
}

func (t *Train) Move(stepLength time.Duration, track *Track) {
	t.Speeds = append(t.Speeds, t.Speed)
	stepLengthSeconds := float64(stepLength) / float64(time.Second)
	t.DistanceTraveled += t.Speed * stepLengthSeconds
	t.Position += t.Speed * stepLengthSeconds
}

func (t *Train) BrakingDistance() float64 {
	timeToDecellerate := t.Speed / t.Braking
	return (timeToDecellerate * t.Speed) / 2.0
}

func (t *Train) ShouldStartBrakingForStation(track *Track) bool {
	brakingDistance := t.BrakingDistance()
	return float64(track.DistanceMeters)-t.Position <= brakingDistance
}

func (t *Train) EvaluateSituation(stepLength time.Duration, track *Track) {
	trainAhead := track.GetNextTrain(t.Position)
	if trainAhead == nil {
		if t.Signal != SignalGo {
			t.SendSignal(SignalGo)
		}
		return
	}

	var distanceToNextTrain float64
	if trainAhead.Position != 0 {
		distanceToNextTrain = trainAhead.Position - t.Position
	} else {
		timeRemainingInStation := t.AverageTimeInStation - trainAhead.ArrivedAtStation
		distanceWhileWaiting := t.Speed * (float64(timeRemainingInStation) / float64(time.Second))
		distanceToNextTrain = track.DistanceMeters - t.Position - distanceWhileWaiting
	}

	brakingDistance := t.BrakingDistance() + t.MinumumSafeDistance

	switch trainAhead.Signal {
	case SignalHold:
		{
			if distanceToNextTrain < brakingDistance {
				t.SendSignal(SignalHold)
			} else {
				t.SendSignal(SignalCaution)
			}
		}
	case SignalGo, SignalCaution:
		{
			if distanceToNextTrain < brakingDistance {
				switch t.Signal {
				case SignalGo:
					{
						t.SendSignal(SignalCaution)
					}
				case SignalHold:
					{
						t.SendSignal(SignalHold)
					}
				}
			}

		}
	}

}

func (t *Train) Motion(stepLength time.Duration, track *Track) {
	switch t.Signal {
	case SignalGo:
		{
			if t.ShouldStartBrakingForStation(track) {
				t.Decellerate(stepLength)
			} else {
				t.Accelerate(stepLength, t.MaximumSpeed)
			}
		}
	case SignalCaution:
		{
			if t.ShouldStartBrakingForStation(track) {
				t.Decellerate(stepLength)
			} else if t.Speed > t.CautionSpeed {
				t.Decellerate(stepLength)
			} else {
				t.Accelerate(stepLength, t.CautionSpeed)
			}
		}
	case SignalHold:
		{
			if t.Speed > 0 {
				t.Decellerate(stepLength)
			}
		}
	}
	t.Move(stepLength, track)
}

func (t *Train) HasLeftYard(wallClock time.Duration) {
	t.LeftYard = wallClock
	t.IsOutbound = true
}

func (t *Train) HasReachedStation(wallClock time.Duration, track *Track) bool {
	return t.Position >= track.DistanceMeters
}

func (t *Train) ReturnsToYard(wallClock time.Duration, station *Station) {
	t.RoundTripTimes = append(t.RoundTripTimes, wallClock-t.LeftYard)
	t.Reset()
}

func (t *Train) ArrivesAtStation(wallClock time.Duration, station *Station) {
	t.Speed = 0
	t.Position = 0
	station.TrainEnters(t)
	t.ArrivedAtStation = wallClock
	t.DisembarkPassengers(wallClock, station)
	t.EmbarkPassengers(wallClock, station)
}

func (t *Train) Reset() {
	t.IsOutbound = true
	t.LeftYard = 0
	t.ArrivedAtStation = 0
}

func (t *Train) Depart(wallClock time.Duration, station *Station) {
	if wallClock-t.ArrivedAtStation >= t.AverageTimeInStation {
		switch t.Signal {
		case SignalGo, SignalCaution:
			{
				station.TrainDeparts(t)
				t.ArrivedAtStation = 0
			}
		}
	}
}

func (t *Train) EmbarkPassengers(wallClock time.Duration, station *Station) {
	currentCapacity := t.Capacity - len(t.Passengers)
	if currentCapacity > 0 {
		for x := 0; x < station.WaitingPassengers.Len(); x++ {
			p := station.WaitingPassengers.Dequeue()
			if p.IsOutBound == t.IsOutbound { //if going the same direction
				p.Boarding(wallClock, t)
				t.Passengers = append(t.Passengers, p)
			} else {
				station.WaitingPassengers.Enqueue(p)
			}
		}
	}
}

func (t *Train) DisembarkPassengers(wallClock time.Duration, station *Station) {
	var newPassengers []*Passenger
	for _, rider := range t.Passengers {
		if rider.Destination == station.Name {
			rider.Disembarking(wallClock, t)
			station.GeneralPopulation.Enqueue(rider)
		} else {
			newPassengers = append(newPassengers, rider)
		}
	}
	t.Passengers = newPassengers
}
