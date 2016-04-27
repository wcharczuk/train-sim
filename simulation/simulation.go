package simulation

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/blendlabs/go-util"
)

func New(stepLength time.Duration, totalTime time.Duration, pauseTime *time.Duration) *Simulation {
	return &Simulation{
		StepLength: stepLength,
		TotalTime:  totalTime,
		PauseTime:  pauseTime,

		Stasis:   false,
		Complete: false,

		Provider: rand.New(rand.NewSource(time.Now().Unix())),

		TotalPassengerCount: 1 << 20,
		TotalTrainCount:     32,

		TrainCapacity:            256,
		TrainAverageAcceleration: 1.5,
		TrainAverageBraking:      1.5,
		TrainMaximumSpeed:        24.5872,

		StationIncidentLikelihood: 1.0,
		AverageIncidentDelay:      1 * time.Minute,

		AverageTimeBetweenTrains: 150 * time.Second,
		AverageTimeInStation:     30 * time.Second,
	}
}

type Simulation struct {
	StepLength time.Duration
	TotalTime  time.Duration
	PauseTime  *time.Duration

	TotalPassengerCount int
	TotalTrainCount     int

	TrainCapacity            int
	TrainAverageAcceleration float64
	TrainAverageBraking      float64
	TrainMaximumSpeed        float64

	// StationIncidentLikelihood is the likelihood adjuster an incident happens in a station.
	// typicall it is given in incidents per hour. 0.001 is one incident every 1000 hours.
	StationIncidentLikelihood float64

	// AverageIncidentDelay is the time an incident typically delays a train.
	AverageIncidentDelay time.Duration

	// AverageTimeBetweenTrains is the average time between trains leaving the yard.
	AverageTimeBetweenTrains time.Duration

	// AverageTimeInStation is the average time the train waits in the station.
	AverageTimeInStation time.Duration

	WallClock time.Duration

	Stasis   bool
	Complete bool

	Stations []*Station
	People   *QueueOfPassenger
	Yard     *QueueOfTrain

	TotalAverageRidership int
	LastTrainReleased     time.Duration

	Provider *rand.Rand

	LogEntries []string
}

func (s *Simulation) logf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	s.LogEntries = append(s.LogEntries, fmt.Sprintf("%v - %s\n", s.WallClock, message))
}

func (s *Simulation) GeneratePassengers() {
	s.People = NewQueueOfPassenger()
	for x := 0; x < s.TotalPassengerCount; x++ {
		s.People.Enqueue(NewPassenger(s.Provider, x+1))
	}
}

func (s *Simulation) GenerateStations() {
	s.Stations = []*Station{}

	//NYC Subway "3" Line
	//-----Terminus------
	s.Stations = append(s.Stations, NewStation("Harlem-148 Street", 3952, s.People))               //0
	s.Stations = append(s.Stations, NewStation("145 Street", 3631, s.People))                      //1
	s.Stations = append(s.Stations, NewStation("135 Street", 15335, s.People))                     //2
	s.Stations = append(s.Stations, NewStation("125 Street", 15744, s.People))                     //3
	s.Stations = append(s.Stations, NewStation("116 Street", 11787, s.People))                     //4
	s.Stations = append(s.Stations, NewStation("Central Park North (110 Street)", 9559, s.People)) //5
	s.Stations = append(s.Stations, NewStation("96 Street", 39254, s.People))                      //6
	s.Stations = append(s.Stations, NewStation("72 Street", 40639, s.People))                      //7
	s.Stations = append(s.Stations, NewStation("Times Square-42 Street", 204908, s.People))        //8
	s.Stations = append(s.Stations, NewStation("34 Street-Penn Station", 92693, s.People))         //9
	s.Stations = append(s.Stations, NewStation("14 Street", 49990, s.People))                      //10
	s.Stations = append(s.Stations, NewStation("Chambers Street", 23862, s.People))                //11
	s.Stations = append(s.Stations, NewStation("Park Place", 55683, s.People))                     //12
	s.Stations = append(s.Stations, NewStation("Fulton Street", 69444, s.People))                  //13
	s.Stations = append(s.Stations, NewStation("Wall Street", 28075, s.People))                    //14
	//------MNH/BK-------
	s.Stations = append(s.Stations, NewStation("Clark Street", 6083, s.People))                    //15
	s.Stations = append(s.Stations, NewStation("Borough Hall", 38944, s.People))                   //16
	s.Stations = append(s.Stations, NewStation("Hoyt Street-Fulton Mall", 7138, s.People))         //17
	s.Stations = append(s.Stations, NewStation("Nevins Street", 11752, s.People))                  //18
	s.Stations = append(s.Stations, NewStation("Atlantic Avenue", 41645, s.People))                //19
	s.Stations = append(s.Stations, NewStation("Bergen Street", 3923, s.People))                   //20
	s.Stations = append(s.Stations, NewStation("Grand Army Plaza", 7971, s.People))                //21
	s.Stations = append(s.Stations, NewStation("Eastern Parkway-Brooklyn Museum", 4889, s.People)) //22
	s.Stations = append(s.Stations, NewStation("Franklin Avenue", 15787, s.People))                //23
	s.Stations = append(s.Stations, NewStation("Nostrand Avenue", 4268, s.People))                 //24
	s.Stations = append(s.Stations, NewStation("Kingston Avenue", 5017, s.People))                 //25
	s.Stations = append(s.Stations, NewStation("Crown Heights-Utica Avenue", 28287, s.People))     //26
	s.Stations = append(s.Stations, NewStation("Sutter Avenue-Rutland Road", 8084, s.People))      //27
	s.Stations = append(s.Stations, NewStation("Saratoga Avenue", 5933, s.People))                 //28
	s.Stations = append(s.Stations, NewStation("Rockaway Avenue", 5735, s.People))                 //29
	s.Stations = append(s.Stations, NewStation("Junius Street", 2361, s.People))                   //30
	s.Stations = append(s.Stations, NewStation("Pennslyvania Avenue", 5718, s.People))             //31
	s.Stations = append(s.Stations, NewStation("Van Siclen Avenue", 3438, s.People))               //32
	s.Stations = append(s.Stations, NewStation("New Lots Avenue", 6626, s.People))                 //33
	//-----Terminus------

	//-----Terminus------
	s.Stations[0].LinkWith(s.Stations[1], 868)
	s.Stations[1].LinkWith(s.Stations[2], 773)
	s.Stations[2].LinkWith(s.Stations[3], 825)
	s.Stations[3].LinkWith(s.Stations[4], 722)
	s.Stations[4].LinkWith(s.Stations[5], 474)
	s.Stations[5].LinkWith(s.Stations[6], 2061)
	s.Stations[6].LinkWith(s.Stations[7], 1947)
	s.Stations[7].LinkWith(s.Stations[8], 2579)
	s.Stations[8].LinkWith(s.Stations[9], 634)
	s.Stations[9].LinkWith(s.Stations[10], 1578)
	s.Stations[10].LinkWith(s.Stations[11], 2714)
	s.Stations[11].LinkWith(s.Stations[12], 327)
	s.Stations[12].LinkWith(s.Stations[13], 365)
	s.Stations[13].LinkWith(s.Stations[14], 380)
	s.Stations[14].LinkWith(s.Stations[15], 1706)
	//------MNH/BK-------
	s.Stations[15].LinkWith(s.Stations[16], 535)
	s.Stations[16].LinkWith(s.Stations[17], 511)
	s.Stations[17].LinkWith(s.Stations[18], 467)
	s.Stations[18].LinkWith(s.Stations[19], 471)
	s.Stations[19].LinkWith(s.Stations[20], 475)
	s.Stations[20].LinkWith(s.Stations[21], 714)
	s.Stations[21].LinkWith(s.Stations[22], 659)
	s.Stations[22].LinkWith(s.Stations[23], 549)
	s.Stations[23].LinkWith(s.Stations[24], 629)
	s.Stations[24].LinkWith(s.Stations[25], 712)
	s.Stations[25].LinkWith(s.Stations[26], 806)
	s.Stations[26].LinkWith(s.Stations[27], 1106)

	s.Stations[27].LinkWith(s.Stations[28], 696)
	s.Stations[28].LinkWith(s.Stations[29], 632)
	s.Stations[29].LinkWith(s.Stations[30], 554)
	s.Stations[30].LinkWith(s.Stations[31], 652)
	s.Stations[31].LinkWith(s.Stations[32], 552)
	s.Stations[32].LinkWith(s.Stations[33], 411)
	//-----Terminus------
}

func (s *Simulation) GenerateTrains() {
	s.Yard = NewQueueOfTrain()
	for x := 0; x < s.TotalTrainCount; x++ {
		t := NewTrain(x, "IRT 3", s.TrainMaximumSpeed, s.TrainAverageAcceleration, s.TrainAverageBraking, s.AverageTimeInStation)
		t.Capacity = s.TrainCapacity
		s.Yard.Enqueue(t)
	}
}

func (s *Simulation) CalculateTotalAverageRidership() {
	var totalAverageRidership int
	for _, station := range s.Stations {
		totalAverageRidership += station.RidersPerDayMean
	}
	s.TotalAverageRidership = totalAverageRidership
}

func (s *Simulation) ShouldReleaseTrainFromYard() bool {
	if s.Complete {
		return false
	}

	if s.Yard.Len() < 1 {
		return false
	}

	return s.WallClock-s.LastTrainReleased >= s.AverageTimeBetweenTrains
}

func (s *Simulation) StopsToDestination(headingOutBound bool, from *Station, to string) int {
	var count int
	var station *Station = from
	if headingOutBound {
		for station.OutBoundTrack != nil {
			station = station.OutBoundTrack.End
			count = count + 1

			if station.Name == to {
				return count
			}
		}
	} else {
		for station.InBoundTrack != nil {
			station = station.InBoundTrack.End
			count = count + 1

			if station.Name == to {
				return count
			}
		}
	}
	return len(s.Stations) + 1
}

func (s *Simulation) DestinationIsOutbound(from *Station, to string) bool {
	outboundStops := s.StopsToDestination(true, from, to)
	inboundStops := s.StopsToDestination(false, from, to)
	return outboundStops < inboundStops
}

func (s *Simulation) StationIncident(station *Station) {
	if station.OutBoundTrain == nil && station.InBoundTrain == nil {
		return
	}

	if s.Complete {
		return
	}

	if !s.Stasis {
		return
	}

	var trainAffected *Train
	if station.OutBoundTrain != nil && station.InBoundTrain != nil {
		if s.Provider.Float64() <= 0.5 {
			trainAffected = station.OutBoundTrain
		} else {
			trainAffected = station.InBoundTrain
		}
	} else if station.OutBoundTrain != nil {
		trainAffected = station.OutBoundTrain
	} else {
		trainAffected = station.InBoundTrain
	}

	ridershipRatio := float64(station.RidersPerDayMean) / float64(s.TotalAverageRidership)
	nominalLikelihood := (float64(s.StepLength) / float64(time.Hour)) * s.StationIncidentLikelihood

	didIncidentOccur := s.Provider.Float64() < ridershipRatio*nominalLikelihood

	if didIncidentOccur {
		s.logf("Station Incident at %s, holding train for %v", station.Name, s.AverageIncidentDelay)
		trainAffected.Hold(s.WallClock)
	}
}

func (s *Simulation) ReleaseTrainsOnHold(station *Station) {
	if station.OutBoundTrain == nil && station.InBoundTrain == nil {
		return
	}

	s.ReleaseTrainOnHold(station, station.OutBoundTrain)
	s.ReleaseTrainOnHold(station, station.InBoundTrain)
}

func (s *Simulation) ReleaseTrainOnHold(station *Station, train *Train) {
	if train != nil {
		if train.Signal == SignalHold {
			if s.WallClock-train.HeldAtStation >= s.AverageIncidentDelay {
				s.logf("Station Incident resolved at %s", station.Name)
				train.SendSignal(SignalGo)
			}
		}
	}
}

func (s *Simulation) PassengersArrive(station *Station) {
	if s.People.Len() == 0 {
		return
	}

	if !s.Stasis {
		return
	}

	if s.Complete {
		return
	}

	pdf := station.PassengerArrivalPDF(s.Provider, s.StepLength)
	if s.Provider.Float64() <= pdf {
		passenger := s.People.Dequeue()
		s.PassengerArrivesAtStation(station, passenger)
	}
}

func (s *Simulation) PassengerArrivesAtStation(station *Station, passenger *Passenger) {
	destinationIndex := s.Provider.Intn(len(s.Stations))
	passenger.Destination = s.Stations[destinationIndex].Name

	//figure out outbound | inbound
	passenger.IsOutBound = s.DestinationIsOutbound(station, passenger.Destination)

	if passenger.IsOutBound && station.OutBoundTrain != nil && len(station.OutBoundTrain.Passengers) < station.OutBoundTrain.Capacity {
		passenger.Waiting = append(passenger.Waiting, 0)
		station.OutBoundTrain.Passengers = append(station.OutBoundTrain.Passengers, passenger)
	} else if !passenger.IsOutBound && station.InBoundTrain != nil && len(station.InBoundTrain.Passengers) < station.InBoundTrain.Capacity {
		passenger.Waiting = append(passenger.Waiting, 0)
		station.InBoundTrain.Passengers = append(station.InBoundTrain.Passengers, passenger)
	} else {
		passenger.StartedWaiting = s.WallClock
		station.WaitingPassengers.Enqueue(passenger)
	}
}

func (s *Simulation) IsComplete() {
	s.logf("Simulation Complete, draining line.")
	s.Complete = true
}

func (s *Simulation) IsAtStasis() {
	s.logf("Simulation At Stasis, starting passenger arrivals.")
	s.Stasis = true
}

func (s *Simulation) AllTrainsReturned() bool {
	return s.Yard.Len() == s.TotalTrainCount
}

func (s *Simulation) OutBoundTerminus() *Station {
	return s.Stations[len(s.Stations)-1]
}

func (s *Simulation) InBoundTerminus() *Station {
	return s.Stations[0]
}

func (s *Simulation) Step() {
	if s.ShouldReleaseTrainFromYard() {
		s.LastTrainReleased = s.WallClock
		t := s.Yard.Dequeue()
		t.HasLeftYard(s.WallClock)
		t.ArrivesAtStation(s.WallClock, s.Stations[0])
		s.logf("Releasing [%d] from yard, %d left in yard", t.ID, s.Yard.Len())
	}

	for _, station := range s.Stations {
		s.PassengersArrive(station)
		station.CheckWaitingTrains(s.WallClock)
		s.StationIncident(station)
		s.ReleaseTrainsOnHold(station)
	}

	// do outbound trains
	station := s.InBoundTerminus()
	for station.OutBoundTrack != nil {
		station.OutBoundTrack.MoveTrains(s.StepLength, s.WallClock)
		station = station.OutBoundTrack.End
	}

	// do inbound trains
	station = s.OutBoundTerminus()
	for station.InBoundTrack != nil {
		station.InBoundTrack.MoveTrains(s.StepLength, s.WallClock)
		station = station.InBoundTrack.End
	}

	station = s.InBoundTerminus()
	if station.InBoundTrain != nil {
		train := station.InBoundTrain
		s.logf("Returning [%d] to the yard", train.ID)
		station.TrainDeparts(train)
		train.ReturnsToYard(s.WallClock, station)
		s.Yard.Enqueue(train)
		if !s.Stasis {
			s.IsAtStasis()
		}
	}

	s.WallClock += s.StepLength
}

func (s *Simulation) Run() {
	s.GeneratePassengers()
	s.GenerateTrains()
	s.GenerateStations()
	s.CalculateTotalAverageRidership()

	for s.WallClock < s.TotalTime {
		s.Step()
		if s.PauseTime != nil {
			s.Display()
			time.Sleep(*s.PauseTime)
		}
	}

	s.IsComplete()

	for !s.AllTrainsReturned() {
		s.Step()
		if s.PauseTime != nil {
			s.Display()
			time.Sleep(*s.PauseTime)
		}
	}

	stats := s.ComputeStats()
	fmt.Printf("Simulation Stats:\n%v", stats)
}

// --------------------------------------------------------------------------------
// Stats & Display Methods
// --------------------------------------------------------------------------------

func (s *Simulation) ComputeStats() *SimulationStats {
	return &SimulationStats{
		AveragePassengerWaitingTime: s.computeMeanPassengerWaitingTime(),
		AveragePassengerTripTime:    s.computeMeanPassengerTripTime(),
		AverageTrainRoundTripTime:   s.computeMeanRoundTripTime(),
	}
}

func (s *Simulation) computeMeanPassengerWaitingTime() time.Duration {
	var times []time.Duration
	for x := 0; x < s.People.Len(); x++ {
		p := s.People.Dequeue()
		if len(p.Waiting) != 0 {
			times = append(times, util.MeanOfDuration(p.Waiting))
		}
		s.People.Enqueue(p)
	}
	return util.MeanOfDuration(times)
}

func (s *Simulation) computeMeanPassengerTripTime() time.Duration {
	var times []time.Duration
	for x := 0; x < s.People.Len(); x++ {
		p := s.People.Dequeue()
		if len(p.InMotion) != 0 {
			times = append(times, util.MeanOfDuration(p.InMotion))
		}
		s.People.Enqueue(p)
	}
	return util.MeanOfDuration(times)
}

func (s *Simulation) computeMeanRoundTripTime() time.Duration {
	var times []time.Duration
	for x := 0; x < s.Yard.Len(); x++ {
		t := s.Yard.Dequeue()
		if len(t.RoundTripTimes) != 0 {
			times = append(times, util.MeanOfDuration(t.RoundTripTimes))
		}
		s.Yard.Enqueue(t)
	}
	return util.MeanOfDuration(times)
}

func (s *Simulation) Display() {
	clear()

	var status string
	if s.Complete {
		status = " Complete"
	}
	if !s.Stasis {
		status = " Warming Up"
	}
	fmt.Printf("Clock: %v%s\n\n", s.WallClock, status)

	for _, station := range s.Stations {
		if station.OutBoundTrain != nil && station.InBoundTrain != nil {
			outWaiting := s.WallClock - station.OutBoundTrain.ArrivedAtStation
			inWaiting := s.WallClock - station.InBoundTrain.ArrivedAtStation
			fmt.Printf("%s - Waiting: %d %v %v %v %v\n", station.Name, station.WaitingPassengers.Len(), station.OutBoundTrain, outWaiting, station.InBoundTrain, inWaiting)
		} else if station.OutBoundTrain != nil {
			outWaiting := s.WallClock - station.OutBoundTrain.ArrivedAtStation
			fmt.Printf("%s - Waiting: %d %v %v\n", station.Name, station.WaitingPassengers.Len(), station.OutBoundTrain, outWaiting)
		} else if station.InBoundTrain != nil {
			inWaiting := s.WallClock - station.InBoundTrain.ArrivedAtStation
			fmt.Printf("%s - Waiting: %d %v %v\n", station.Name, station.WaitingPassengers.Len(), station.InBoundTrain, inWaiting)
		} else {
			fmt.Printf("%s - Waiting: %d\n", station.Name, station.WaitingPassengers.Len())
		}
		if station.OutBoundTrack != nil {
			for _, train := range station.OutBoundTrack.Trains {
				fmt.Printf("%v ", train)
			}
			for _, train := range station.OutBoundTrack.End.InBoundTrack.Trains {
				fmt.Printf("%v ", train)
			}
			fmt.Println()
		}
	}

	fmt.Println()
	fmt.Println("Log Entries:")
	for _, entry := range last(s.LogEntries, 10) {
		fmt.Print(entry)
	}
}

func last(values []string, count int) []string {
	var lastValues []string
	var valuesLen = len(values)
	if valuesLen == 0 {
		return lastValues
	}

	if valuesLen <= count {
		for x := valuesLen - 1; x >= 0; x-- {
			lastValues = append(lastValues, values[x])
		}
	} else {
		for x := valuesLen - 1; x >= valuesLen-(count+1); x-- {
			lastValues = append(lastValues, values[x])
		}
	}

	return lastValues
}

func clear() {
	fmt.Print("\033[H\033[2J")
}
