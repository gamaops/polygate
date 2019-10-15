package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	logrus "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"

	dto "github.com/prometheus/client_model/go"
)

func TestInitialMetrics(t *testing.T) {
	if producerJobCount != nil {
		t.Error("producerJobCount must be initialized as nil")
	}
	if producerFailedJobCount != nil {
		t.Error("producerFailedJobCount must be initialized as nil")
	}
	if consumerJobCount != nil {
		t.Error("consumerJobCount must be initialized as nil")
	}
	if producerWaitingResponsesCount != nil {
		t.Error("producerWaitingResponsesCount must be initialized as nil")
	}
	if producerClientStreamsCount != nil {
		t.Error("producerClientStreamsCount must be initialized as nil")
	}
	if consumerClientStreamsCount != nil {
		t.Error("consumerClientStreamsCount must be initialized as nil")
	}
	if producerJobExecutionSeconds != nil {
		t.Error("producerJobExecutionSeconds must be initialized as nil")
	}
	if consumerJobExecutionSeconds != nil {
		t.Error("consumerJobExecutionSeconds must be initialized as nil")
	}
	if producerJobPayloadBytes != nil {
		t.Error("producerJobPayloadBytes must be initialized as nil")
	}
	if producerJobEventBytes != nil {
		t.Error("producerJobEventBytes must be initialized as nil")
	}
}

func TestLoadProducerMetrics(t *testing.T) {

	loadProducerMetrics()

	if producerJobCount == nil {
		t.Error("producerJobCount must not be nil after loading")
	}
	if producerFailedJobCount == nil {
		t.Error("producerFailedJobCount must not be nil after loading")
	}
	if producerWaitingResponsesCount == nil {
		t.Error("producerWaitingResponsesCount must not be nil after loading")
	}
	if producerClientStreamsCount == nil {
		t.Error("producerClientStreamsCount must not be nil after loading")
	}
	if producerJobExecutionSeconds == nil {
		t.Error("producerJobExecutionSeconds must not be nil after loading")
	}
	if producerJobPayloadBytes == nil {
		t.Error("producerJobPayloadBytes must not be nil after loading")
	}
	if producerJobEventBytes == nil {
		t.Error("producerJobEventBytes must not be nil after loading")
	}

	m := &dto.Metric{}
	producerReady.Write(m)

	if *m.Gauge.Value != 0 {
		t.Error("producerReady gauge must be 0 after initialization")
	}

}

func TestLoadConsumerMetrics(t *testing.T) {

	loadConsumerMetrics()

	if consumerJobCount == nil {
		t.Error("consumerJobCount must not be nil after loading")
	}
	if consumerClientStreamsCount == nil {
		t.Error("consumerClientStreamsCount must not be nil after loading")
	}
	if consumerJobExecutionSeconds == nil {
		t.Error("consumerJobExecutionSeconds must not be nil after loading")
	}

	m := &dto.Metric{}
	consumerReady.Write(m)

	if *m.Gauge.Value != 0 {
		t.Error("consumerReady gauge must be 0 after initialization")
	}

}

func TestLivenessHandler(t *testing.T) {

	req, err := http.NewRequest("GET", configuration.Metrics.Routes.Liveness, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(livenessHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	polygateUp.Set(1)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestReadinessHandler(t *testing.T) {

	req, err := http.NewRequest("GET", configuration.Metrics.Routes.Readiness, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(readinessHandler)

	configuration.Client.Enable = true
	configuration.Server.Enable = true

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	producerReady.Set(1)
	consumerReady.Set(0)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusPartialContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPartialContent)
	}

	producerReady.Set(0)
	consumerReady.Set(1)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusPartialContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusPartialContent)
	}

	producerReady.Set(1)
	consumerReady.Set(1)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	configuration.Server.Enable = false
	producerReady.Set(0)
	consumerReady.Set(1)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	configuration.Server.Enable = true
	configuration.Client.Enable = false
	producerReady.Set(1)
	consumerReady.Set(0)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

}

func TestLoadMetricsServer(t *testing.T) {

	configuration.Metrics.Routes.Metrics = "/metrics"
	configuration.Metrics.Routes.Readiness = "/ready"
	configuration.Metrics.Routes.Liveness = "/liveness"

	var hook *logtest.Hook
	log, hook = logtest.NewNullLogger()

	exitCode := 0

	log.ExitFunc = func(code int) {
		exitCode = code
	}

	configuration.Metrics.ShutdownTimeout = "invalid"

	loadMetricsServer()

	if hook.LastEntry().Level != logrus.FatalLevel {
		t.Error("Metrics server should validate timeout duration")
	}

	if exitCode != 1 {
		t.Error("Invalid duration must exit process")
	}

	m := &dto.Metric{}
	polygateUp.Write(m)

	if *m.Gauge.Value != 0 {
		t.Error("polygateUp gauge must be 0 after initialization")
	}

	configuration.Metrics.ShutdownTimeout = "15s"
	exitCode = 0

	configuration.Metrics.Routes.Metrics = "/metrics2"
	configuration.Metrics.Routes.Readiness = "/ready2"
	configuration.Metrics.Routes.Liveness = "/liveness2"

	loadMetricsServer()

	if exitCode != 0 {
		t.Error("Process exited without reason")
	}

}

func TestStartStopMetricsServer(t *testing.T) {

	if metricsServer != nil {
		t.Error("Metrics server must no be set before starting")
	}

	configuration.Metrics.Address = "0.0.0.0"
	configuration.Metrics.Port = 35500

	startMetricsServer()

	if metricsServer == nil {
		t.Error("Metrics server must be se after starting")
	}

	if metricsServer.Addr != "0.0.0.0:35500" {
		t.Error("Metrics server listening incorred address")
	}

	var hook *logtest.Hook
	log, hook = logtest.NewNullLogger()

	stopMetricsServer()

	if hook.LastEntry().Level != logrus.WarnLevel {
		t.Error("Error when finishing the metrics server")
	}
}
