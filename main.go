package main

import (
	"fmt"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var (
	listenAddress    = flag.String("listen-address", ":8080", "The address on which to expose the web interface and generated Prometheus metrics.")
	metricsEndpoint  = flag.String("telemetry-path", "/metrics", "Path under which to expose metrics.")
	zookeeperConnect = flag.String("zookeeper-connect", "localhost:2181", "Zookeeper connection string")
	clusterName      = flag.String("cluster-name", "kafka-cluster", "Name of the Kafka cluster used in static label")
	topicsFilter     = flag.String("topics-filter", "", "Regex expression to export only topics that match expression")
	refreshInterval  = flag.Int("refresh-interval", 5, "Seconds to sleep in between refreshes")
)

var (

	consumerGroupGougeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "kafka",
		Subsystem:   "consumergroup",
		Name:        "current_offset",
		Help:        "Current Offset of a ConsumerGroup at Topic/Partition",
		ConstLabels: map[string]string{"cluster": *clusterName},
	},
		[]string{"consumergroup", "topic", "partition"},
	)

	consumerGroupLogSizeGougeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "kafka",
		Subsystem:   "consumergroup",
		Name:        "logsize",
		Help:        "Current logsize of a ConsumerGroup at Topic/Partition",
		ConstLabels: map[string]string{"cluster": *clusterName},
	},
		[]string{"consumergroup", "topic", "partition"},
	)

	consumerGroupLagGougeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "kafka",
		Subsystem:   "consumergroup",
		Name:        "lag",
		Help:        "Current approximate Lag of a ConsumerGroup at Topic/Partition",
		ConstLabels: map[string]string{"cluster": *clusterName},
	},
		[]string{"consumergroup", "topic", "partition"},
	)
)

// 注册prometheus
func init() {
	prometheus.MustRegister(consumerGroupGougeVec)
	prometheus.MustRegister(consumerGroupLogSizeGougeVec)
	prometheus.MustRegister(consumerGroupLagGougeVec)
}



func main() {
	fmt.Println("Running offset exporter")
	flag.Parse()

	fmt.Println("Settings: ")
	fmt.Println("listen-address: ", *listenAddress)
	fmt.Println("telemetry-path: ", *metricsEndpoint)
	fmt.Println("zookeeper-connect: ", *zookeeperConnect)
	fmt.Println("cluster-name: ", *clusterName)
	fmt.Println("refresh-interval: ", *refreshInterval)

	//Init Clients
	initClients()

	// Periodically record stats from Kafka
	go func() {
		for {
			updateOffsets()
			// 5s轮询一次
			time.Sleep(time.Duration(time.Duration(*refreshInterval) * time.Second))
		}
	}()

	// Expose the registered metrics via HTTP.
	http.Handle(*metricsEndpoint, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
