package main

import (
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var parameters Parameters
var configuration Configuration
var instanceProducerID string

func init() {

	prettyLog, exists := os.LookupEnv("PRETTY_LOG")

	if exists && prettyLog == "true" {
		log.SetFormatter(&log.TextFormatter{})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}

	log.SetOutput(os.Stdout)

	parameters = loadParameters()

	log.SetLevel(parameters.logLevel)

	configuration = loadConfiguration()

	if configuration.Server.Enable {
		startRedisClient("job", configuration.Redis.JobPoolSize)
		startRedisClient("producer", 1)
		var producerID strings.Builder
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("Unable to get hostname: %v", err)
		}
		producerID.WriteString(configuration.Redis.Prefix)
		producerID.WriteRune(':')
		producerID.WriteString("prod:")
		producerID.WriteString(hostname)
		instanceProducerID = producerID.String()
	}

	if configuration.Client.Enable {
		if !configuration.Server.Enable {
			startRedisClient("job", configuration.Redis.JobPoolSize)
		}
		loadClientJobHandlers()
		loadConsumers()
	}

}

func main() {

	sigs := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	startConsumers()
	if configuration.Server.Enable {
		startProducerListener()
		server := createServer()
		wg.Add(1)
		go func() {
			sig := <-sigs
			log.Warnf("Stopping server, signal received: %v", sig)
			server.GracefulStop()
			log.Warn("gRPC server stopped")
			closeProducerListener()
			closeRedisClients("producer")
			wg.Done()
		}()
	}

	//http.ListenAndServe("localhost:8088", nil)
	wg.Wait()

}
