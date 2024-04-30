package main

import (
	_ "GoWeatherAPI/app"
	app2 "GoWeatherAPI/app"
	_ "GoWeatherAPI/config"
	_ "github.com/spf13/viper"
	"strings"
)

// Node represents a node in the topology.
type Node struct {
	Name string
}

// Topology represents the overall structure of the nodes.
type Topology struct {
	Nodes []Node
}

// groupNodesBySuffix groups nodes within the Topology based on their suffix.
func (t *Topology) groupNodesBySuffix() map[string][]Node {
	groups := make(map[string][]Node)
	for _, node := range t.Nodes {
		parts := strings.Split(node.Name, "-")
		if len(parts) > 1 {
			suffix := parts[len(parts)-1]
			groups[suffix] = append(groups[suffix], node)
		}
	}
	return groups
}
func main() {

	app := app2.App{}
	app.Run()
}

/*client := client.NewOpenWeatherMapClient(config.AppConfig.Apikey)
p := poller.NewPoller(config.AppConfig.PollInterval * time.Second)
metric := metrics.NewMetrics()
unit := config.AppConfig.Unit*/

// use a new hook, on before serve.
/*fmt.Println("Starting UI ...")
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	zipCodeHandler(w, r, p, client, logger, unit, metric)
})

addToPoller(client, logger, "78759", unit, metric, p)

// use a new hook, on before bootstrap.
go p.StartPollingWeatherAPI()*/

// TODO need to move prometheus stuff to hooks.
/*reg := prometheus.NewRegistry()
reg.MustRegister(
	collectors.NewGoCollector(),
	collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
)*/

//http.Handle("/metrics", promhttp.Handler())
//log.Fatal(http.ListenAndServe(":8081", nil))
