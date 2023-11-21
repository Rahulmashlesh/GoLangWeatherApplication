package poller

import (
	"time"
)

type Caller interface {
	Call()
}

type Poller struct {
	items      []Caller
	pollPeriod time.Duration
}

func NewPoller(pollPeriod time.Duration) *Poller {
	return &Poller{
		pollPeriod: pollPeriod,
	}
}

func (p *Poller) Add(c Caller) {
	p.items = append(p.items, c)
}

func (p *Poller) StartPollingWeatherAPI() {
	ticker := time.NewTicker(p.pollPeriod)
	currentIndex := 0
	for {
		select {
		case <-ticker.C:
			go p.items[currentIndex].Call()
			// Update the current index using the modulus method to stay with in the limit
			currentIndex = (currentIndex + 1) % len(p.items)
		}
	}
}

// TODO: for next time
//set log levels
//add debug logger ,

// testing file to test api logs with out making actual calls to api

//read about io readers and writers
//std lib doc, .

// mock the client . / dependency injection.
