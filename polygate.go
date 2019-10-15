package main

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	stdlog "log"
	_ "net/http/pprof"

	"github.com/fsnotify/fsnotify"
	"github.com/go-redis/redis"
	logrus "github.com/sirupsen/logrus"
)

var parameters Parameters
var configuration Configuration
var instanceProducerID string
var log = logrus.New()

func setup() {

	prettyLog, exists := os.LookupEnv("PRETTY_LOG")

	if exists && prettyLog == "true" {
		log.SetFormatter(&logrus.TextFormatter{})
	} else {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	log.SetOutput(os.Stdout)

	parameters = loadParameters()

	log.SetLevel(parameters.logLevel)

	redisLogger := stdlog.New(logrus.StandardLogger().Writer(), "", 0)

	redis.SetLogger(redisLogger)

	configuration = loadConfiguration()

	if configuration.Server.Enable {
		loadProducerMetrics()
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
		loadConsumerMetrics()
		if !configuration.Server.Enable {
			startRedisClient("job", configuration.Redis.JobPoolSize)
		}
		loadClientJobHandlers()
		loadConsumers()
	}

}

func main() {

	setup()

	sigs := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	loadMetricsServer()
	startMetricsServer()

	if parameters.enableHotReload {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatalf("Error enabling the hot reload configuration watcher: %v", err)
		}
		defer watcher.Close()
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					log.Warnf("Configuration change detected: %v", event)
					if event.Op&fsnotify.Write == fsnotify.Write {
						syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				}
			}
		}()
		err = watcher.Add(parameters.configurationFile)
		if err != nil {
			log.Fatalf("Error while adding the configuration file to watcher: %v", err)
		}
		log.Infof("Hot reload enabled, watching file: %v", parameters.configurationFile)
	}

	stopConsumer := func(sig os.Signal) {
		log.Warnf("Stopping client, signal received: %v", sig)
		consumerReady.Set(0)
		consumersStopped = true
		consumersStopWait.Wait()
		log.Warn("Stopped consumers")
		clientJobsWait.Wait()
		log.Warn("No more jobs to wait")
		if !configuration.Server.Enable {
			closeRedisClients("job")
		}
		stopConsumersRedisConnections()
		stopMetricsServer()
		wg.Done()
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	if configuration.Server.Enable {
		startProducerListener()
		server := createServer()
		wg.Add(1)
		go func() {
			sig := <-sigs
			log.Warnf("Stopping server, signal received: %v", sig)
			producerReady.Set(0)
			server.GracefulStop()
			log.Warn("gRPC server stopped")
			closeProducerListener()
			closeRedisClients("producer")
			if configuration.Client.Enable {
				stopConsumer(sig)
				wg.Done()
				return
			}
			closeRedisClients("job")
			stopMetricsServer()
			wg.Done()
		}()
	}

	if configuration.Client.Enable {
		startConsumers()
		wg.Add(1)
		if !configuration.Server.Enable {
			go func() {
				sig := <-sigs
				stopConsumer(sig)
			}()
		}
	}

	polygateUp.Set(1)

	wg.Wait()

}
