package poller

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

type Caller interface {
	Call()
}

type Poller struct {
	items       []Caller
	pollPeriod  time.Duration
	tickerCount prometheus.Counter
}

func NewPoller(pollPeriod time.Duration) *Poller {
	return &Poller{
		pollPeriod: pollPeriod,
		tickerCount: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "Weather",
			Name:      "poller_ticker_count",
			Help:      "Counter for poller ticker",
		}),
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

			//Increment the counter
			p.tickerCount.Inc()
		}
	}
}
