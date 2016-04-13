package simulation

import (
	"fmt"
	"math/rand"
	"time"
)

var firstNames = []string{
	"John",
	"Paul",
	"James",
	"William",
	"Matthew",
	"David",
	"George",
	"Scott",
	"Gregory",
	"Carl",
	"Gustav",
	"Chester",
	"Abraham",
	"Sterling",
	"Robert",
	"Sigmund",
	"Eli",
}

var lastNames = []string{
	"Smith",
	"Jones",
	"Doe",
	"Howard",
	"Rockafeller",
	"Cotterpot",
	"Grant",
	"Lincoln",
	"Franklin",
	"Vanderbilt",
	"Roebling",
	"Moses",
	"Malone",
	"Murphy",
	"Adams",
	"Merriweather",
	"Yeager",
	"McCue",
	"Hamilton",
	"Archer",
}

func randomFirstName(randomProvider *rand.Rand) string {
	index := randomProvider.Intn(len(firstNames))
	return firstNames[index]
}

func randomLastName(randomProvider *rand.Rand) string {
	index := randomProvider.Intn(len(lastNames))
	return lastNames[index]
}

func NewPassenger(provider *rand.Rand, id int) *Passenger {
	return &Passenger{
		ID:        id,
		FirstName: randomFirstName(provider),
		LastName:  randomLastName(provider),
	}
}

type Passenger struct {
	ID        int
	FirstName string
	LastName  string

	StartedWaiting time.Duration
	StartedRiding  time.Duration

	IsOutBound  bool
	Destination string

	Waiting  []time.Duration
	InMotion []time.Duration
}

func (p *Passenger) Boarding(wallClock time.Duration, train *Train) {
	p.Waiting = append(p.Waiting, wallClock-p.StartedWaiting)
	p.StartedWaiting = 0
	p.StartedRiding = wallClock
}

func (p *Passenger) Disembarking(wallClock time.Duration, train *Train) {
	p.InMotion = append(p.InMotion, wallClock-p.StartedRiding)
	p.StartedRiding = 0
}

func (p *Passenger) String() string {
	return fmt.Sprintf("%s %s", p.FirstName, p.LastName)
}
