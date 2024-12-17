package service

import (
	"GoWeatherAPI/internal/models"
	"GoWeatherAPI/internal/pubsub"
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type Notifier struct{
	subscriber pubsub.Subscriber[string]
	logger *slog.Logger
	model models.Model
}


func NewNotifier(subscriber pubsub.Subscriber[string], logger *slog.Logger, model models.Model) *Notifier{
	return &Notifier{
		subscriber : subscriber,
		logger: logger.With("context", "NotifierService"),
		model: model,	
	}
	
}

func (ns *Notifier) Start(ctx context.Context){
	
	dataChan := ns.subscriber.Subscribe(ctx, "update")
	go func() {
		for {
			select{
				case <- ctx.Done():
					ns.logger.Info("Received Cancell signal")
					return
				case data := <-dataChan:	
					location, err := ns.model.Get(ctx, "95134")
					if err != nil {
						ns.logger.Error("Error getting location model", "error", err)
					 	return 
					};
					stringReader := strings.NewReader( fmt.Sprintf("location:%s temperature: %.2f city %s", location.Zipcode, location.Temperature, location.Name))
					ns.logger.Info("Making http call to NOTIFY.SH", "data", data.Data())
					req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://ntfy.sh/rahul-weather", stringReader)
					if err != nil {
						ns.logger.Error("Error notifying", "error" , err)
					 return   
					};
					rsp, err := http.DefaultClient.Do(req)
					ns.logger.Info("Success http call to NOTIFY.SH", "status", rsp.Status, "statusCode", rsp.StatusCode)
			}
		}
	}()	
}

func insecureHttpClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:  &tls.Config{
				InsecureSkipVerify: true, // Skip certificate validation
			},
		},
	}
}
