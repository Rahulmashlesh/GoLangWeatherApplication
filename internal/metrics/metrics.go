package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	TempGage *prometheus.GaugeVec
}

func NewMetrics(*prometheus.Registry) Metrics {
	tempGage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "Weather",
			Name:      "temperature",
			Help:      "GaugeVec for temperature",
		},
		[]string{
			"City",
			"Zip",
		},
	)

	prometheus.MustRegister(tempGage)

	return Metrics{TempGage: tempGage}
}
