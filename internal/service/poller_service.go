package service

import (
	"GoWeatherAPI/internal/poller"
	"GoWeatherAPI/internal/pubsub"
	"GoWeatherAPI/internal/queue"

	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Caller interface {
	Call(context.Context)
}

type Poller struct {
	pollPeriod   time.Duration
	tickerCount  prometheus.Counter
	queue        queue.Queue // is this a redis queue?
	pollerLocker poller.Locker
	queueLocker  poller.Locker
	logger       *slog.Logger
	myPubSub     pubsub.PubSub[string]
}

func NewPoller(pollPeriod time.Duration, logger *slog.Logger, myPubSub pubsub.PubSub[string], queue queue.Queue, queueLocker, pollerLocker poller.Locker) *Poller {

	return &Poller{
		queue:        queue,
		logger:       logger.With("Service", "Poller"),
		queueLocker:  queueLocker,
		pollerLocker: pollerLocker,
		myPubSub:     myPubSub,
		pollPeriod:   pollPeriod,
		tickerCount: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "Weather",
			Name:      "poller_ticker_count",
			Help:      "Counter for poller ticker",
		}),
	}
}

func (p *Poller) Start(ctx context.Context) {
	p.startPoller(ctx)
	p.startSubscriber(ctx)

}

func (p *Poller) startSubscriber(ctx context.Context) {

	dataChan := p.myPubSub.Subscribe(ctx, "create", "delete")
	p.logger.Info("Suscribed to create and delete chan")

	go func() {
		for {
			select {
			case <-ctx.Done():
				p.logger.Info("Received cancellation signal")
				return
			case data := <-dataChan:
				if p.queueLocker.Lock(ctx) {
					switch data.Type() {
					case "create":
						p.logger.Info("About to Enqueue", "data", data.Data())
						if err := p.queue.Enqueue(ctx, data.Data()); err != nil {
							p.logger.Error("Error enqueing", "error", err)
						}
					case "delete":
						p.logger.Info("About to Delete", "data", data.Data())
						if err := p.queue.Delete(ctx, data.Data()); err != nil {
							p.logger.Error("Error deleting", "error", err)
						}
					}
				}
			}
		}
	}()

}

func (p *Poller) startPoller(ctx context.Context) {
	ticker := time.NewTicker(p.pollPeriod)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				p.logger.Info("Reveived cancel signal ")
				return

			case <-ticker.C:

				// Implementing locks and queues.
				p.logger.Debug(" Trying to lock")
				if p.pollerLocker.Lock(ctx) {
					timeContext, cancel := context.WithTimeout(ctx, p.pollPeriod)

					defer cancel()
					zipcode, err := p.queue.Next(timeContext)
					if err != nil {
						p.logger.Error("Error getting the next item in queue", "error", err)
						continue
					}
					p.logger.Info("PollerService: Trying to publish call event")
					if err := p.myPubSub.Publish(timeContext, "call", zipcode); err != nil {
						p.logger.Error("Error publishing", "error", err)
					}

				}
			}
		}
	}()
}
