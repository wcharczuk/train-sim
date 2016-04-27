package main

import (
	"time"

	"github.com/wcharczuk/train-sim/simulation"
)

func delay(d time.Duration) *time.Duration {
	return &d
}

func main() {
	sim := simulation.New(1*time.Second, 3*time.Hour, nil)

	sim.TrainCapacity = 512

	/*sim.TotalTrainCount = 64
	sim.AverageTimeBetweenTrains = 45 * time.Second
	sim.AverageTimeInStation = 5 * time.Second

	sim.TrainAverageAcceleration = 5.0
	sim.TrainAverageBraking = 3.0
	sim.TrainMaximumSpeed = 50.0 // ~111 mph*/

	sim.Run()
}
