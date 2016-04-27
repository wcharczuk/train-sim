package simulation

import (
	"fmt"
	"time"
)

type SimulationStats struct {
	AveragePassengerTripTime    time.Duration
	AveragePassengerWaitingTime time.Duration
	AverageTrainRoundTripTime   time.Duration
}

func (ss *SimulationStats) String() string {
	return fmt.Sprintf("Mean Passenger Wait Time: %v\nMean Passenger Trip Time: %v\nMean Train Round Trip Time: %v\n", ss.AveragePassengerWaitingTime, ss.AveragePassengerTripTime, ss.AverageTrainRoundTripTime)
}
