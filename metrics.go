package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

var producerJobCount *prometheus.CounterVec
var producerFailedJobCount *prometheus.CounterVec
var consumerJobCount *prometheus.CounterVec
var producerClientStreamsCount *prometheus.GaugeVec
var consumerClientStreamsCount *prometheus.GaugeVec
var producerJobExecutionSeconds *prometheus.HistogramVec
var consumerJobExecutionSeconds *prometheus.HistogramVec
var producerJobPayloadBytes *prometheus.HistogramVec
var producerJobEventBytes *prometheus.HistogramVec
var producerReady prometheus.Gauge
var consumerReady prometheus.Gauge
var polygateUp = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "polygate_up",
	Help: "Indicates if Polygate is up (0 = down, 1 = up)",
})

var metricsServer *http.Server
var metricsServerShutdownTimeout time.Duration
var metricsStatusUp = []byte("{\"status\":\"up\"}")
var metricsStatusDegraded = []byte("{\"status\":\"degraded\"}")
var metricsStatusDown = []byte("{\"status\":\"down\"}")

func loadProducerMetrics() {
	producerJobCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "polygate_producer_job_count",
			Help: "The total number of processed events",
		},
		[]string{"service", "method", "stream"},
	)
	producerFailedJobCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "polygate_producer_failed_job_count",
			Help: "The total number of failed events",
		},
		[]string{"service", "method", "stream"},
	)
	producerClientStreamsCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "polygate_producer_client_streams_count",
			Help: "The total number of current open client streams",
		},
		[]string{"service", "method", "stream"},
	)
	producerJobExecutionSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "polygate_producer_job_execution_seconds",
			Help: "The execution time in seconds of each job",
		},
		[]string{"service", "method", "stream"},
	)
	producerJobPayloadBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "polygate_producer_job_payload_bytes",
			Help: "The size of job payload (only the protobuf message)",
		},
		[]string{"service", "method", "stream"},
	)
	producerJobEventBytes = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "polygate_producer_job_event_bytes",
			Help: "The total size of each event (before encoding base64)",
		},
		[]string{"service", "method", "stream"},
	)
	producerReady = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "polygate_producer_ready",
		Help: "Indicates if Polygate's producer is ready (0 = no, 1 = yes)",
	})
	producerReady.Set(0)
}

func loadConsumerMetrics() {
	consumerJobCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "polygate_consumer_job_count",
			Help: "The total number of processed events",
		},
		[]string{"service", "method", "stream"},
	)
	consumerClientStreamsCount = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "polygate_consumer_client_streams_count",
			Help: "The total number of current open client streams",
		},
		[]string{"service", "method", "stream"},
	)
	consumerJobExecutionSeconds = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "polygate_consumer_job_execution_seconds",
			Help: "The execution time in seconds of each job",
		},
		[]string{"service", "method", "stream"},
	)
	consumerReady = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "polygate_consumer_ready",
		Help: "Indicates if Polygate's consumer is ready (0 = no, 1 = yes)",
	})
	consumerReady.Set(0)
}

func loadMetricsServer() {
	duration, err := time.ParseDuration(configuration.Metrics.ShutdownTimeout)
	if err != nil {
		log.Fatalf("Invalid metrics server shutdown duration: %v", err)
	}
	metricsServerShutdownTimeout = duration

	polygateUp.Set(0)
	http.Handle(configuration.Metrics.Routes.Metrics, promhttp.Handler())
	livenessMetric := &dto.Metric{}
	http.HandleFunc(configuration.Metrics.Routes.Liveness, func(w http.ResponseWriter, r *http.Request) {
		polygateUp.Write(livenessMetric)
		w.Header().Add("Content-Type", "application/json")
		if *livenessMetric.Gauge.Value == 0 {
			w.WriteHeader(503)
			w.Write(metricsStatusDown)
			return
		}
		w.WriteHeader(200)
		w.Write(metricsStatusUp)
	})

	readinessConsumerMetric := &dto.Metric{}
	readinessProducerMetric := &dto.Metric{}
	http.HandleFunc(configuration.Metrics.Routes.Readiness, func(w http.ResponseWriter, r *http.Request) {
		consumerStatus := 1
		producerStatus := 1

		if configuration.Server.Enable {
			producerReady.Write(readinessProducerMetric)
			if *readinessProducerMetric.Gauge.Value == 0 {
				producerStatus = 0
			}
		}
		if configuration.Client.Enable {
			consumerReady.Write(readinessConsumerMetric)
			if *readinessConsumerMetric.Gauge.Value == 0 {
				consumerStatus = 0
			}
		}
		w.Header().Add("Content-Type", "application/json")
		if consumerStatus == 1 && producerStatus == 1 {
			w.WriteHeader(200)
			w.Write(metricsStatusUp)
			return
		}
		if (consumerStatus == 0 && producerStatus == 0) || (consumerStatus == 0 && !configuration.Server.Enable) || (producerStatus == 0 && !configuration.Client.Enable) {
			w.WriteHeader(503)
			w.Write(metricsStatusDown)
			return
		}
		w.WriteHeader(206)
		w.Write(metricsStatusDegraded)
	})
}

func startMetricsServer() {
	var address strings.Builder

	address.WriteString(configuration.Metrics.Address)
	address.WriteRune(':')
	address.WriteString(strconv.Itoa(int(configuration.Metrics.Port)))

	metricsServer = &http.Server{
		Addr: address.String(),
	}

	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the metrics server: %v", err)
		}
	}()

	log.Infof("Metrics server is listening on: %v", address.String())

}

func stopMetricsServer() {
	log.Warn("Stopping metrics server")
	ctx, cancel := context.WithTimeout(context.Background(), metricsServerShutdownTimeout)
	defer cancel()
	if err := metricsServer.Shutdown(ctx); err != nil {
		log.Errorf("Error while stopping metrics server: %v", err)
	}
}
