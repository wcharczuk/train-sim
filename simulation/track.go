package simulation

import "time"

type Track struct {
	DistanceMeters float64

	IsOutBound bool
	Trains     []*Train

	Begin *Station
	End   *Station
}

func (t *Track) AddTrain(train *Train) {
	t.Trains = append(t.Trains, train)
}

func (t *Track) RemoveTrain(trainID int) {
	t.Trains = filterTrains(t.Trains, func(t2 *Train) bool {
		return t2.ID != trainID
	})
}

func (t *Track) HasTrain(trainID int) bool {
	return anyTrains(t.Trains, func(t2 *Train) bool {
		return t2.ID == trainID
	})
}

func (t *Track) MoveTrains(stepLength time.Duration, wallClock time.Duration) {
	var trainsToRemove []*Train
	for x := 0; x < len(t.Trains); x++ {
		train := t.Trains[x]
		train.EvaluateSituation(stepLength, t)
		train.Motion(stepLength, t)

		if train.HasReachedStation(wallClock, t) {
			trainsToRemove = append(trainsToRemove, train)
			train.ArrivesAtStation(wallClock, t.End)
		}
	}

	for x := 0; x < len(trainsToRemove); x++ {
		t.RemoveTrain(trainsToRemove[x].ID)
	}
}

func (t *Track) GetNextTrain(position float64) *Train {
	onTracks := lastTrain(t.Trains, func(train *Train) bool {
		return train.Position > position
	})

	if onTracks != nil {
		return onTracks
	}

	if t.IsOutBound {
		return t.End.OutBoundTrain
	} else {
		return t.End.InBoundTrain
	}
}

type trainPredicate func(t *Train) bool

func anyTrains(trains []*Train, predicate trainPredicate) bool {
	for _, t := range trains {
		if predicate(t) {
			return true
		}
	}
	return false
}

func firstTrain(trains []*Train, predicate trainPredicate) *Train {
	for _, t := range trains {
		if predicate(t) {
			return t
		}
	}
	return nil
}

func lastTrain(trains []*Train, predicate trainPredicate) *Train {
	for x := len(trains) - 1; x >= 0; x-- {
		t := trains[x]
		if predicate(t) {
			return t
		}
	}
	return nil
}

func filterTrains(trains []*Train, predicate trainPredicate) []*Train {
	var newTrains []*Train
	for _, t := range trains {
		if predicate(t) {
			newTrains = append(newTrains, t)
		}
	}
	return newTrains
}
