package main

import (
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

type Parameters struct {
	enableHotReload            bool
	logLevel                   log.Level
	terminateTimeout           uint32
	configurationFile          string
	terminateGRPCServerTimeout uint32
	serializePayloadBase64     bool
}

func loadParameters() Parameters {

	parameters := Parameters{
		enableHotReload:            true,
		logLevel:                   log.DebugLevel,
		terminateTimeout:           15000,
		terminateGRPCServerTimeout: 7000,
		serializePayloadBase64:     true,
		configurationFile:          "",
	}

	configurationFile, exists := os.LookupEnv("CONFIGURATION_FILE")

	if !exists {
		log.Fatal("You must specify the CONFIGURATION_FILE environment variable")
	}

	parameters.configurationFile = configurationFile

	enableHotReload, exists := os.LookupEnv("ENABLE_HOT_RELOAD")

	if exists && enableHotReload == "false" {
		parameters.enableHotReload = false
	}

	logLevel, exists := os.LookupEnv("LOG_LEVEL")

	if exists {
		switch logLevel {
		case "info":
			parameters.logLevel = log.InfoLevel
		case "warn":
			parameters.logLevel = log.WarnLevel
		case "error":
			parameters.logLevel = log.ErrorLevel
		}
	}

	terminateTimeout, exists := os.LookupEnv("TERMINATE_TIMEOUT")

	if exists {
		value, err := strconv.ParseUint(terminateTimeout, 10, 32)
		if err != nil {
			log.WithField("TERMINATE_TIMEOUT", terminateTimeout).Fatal("Invalid timeout value, must be a number")
		}
		parameters.terminateTimeout = uint32(value)
	}

	terminateGRPCServerTimeout, exists := os.LookupEnv("TERMINATE_GRPC_SERVER_TIMEOUT")

	if exists {
		value, err := strconv.ParseUint(terminateGRPCServerTimeout, 10, 32)
		if err != nil {
			log.WithField("TERMINATE_GRPC_SERVER_TIMEOUT", terminateGRPCServerTimeout).Fatal("Invalid timeout value, must be a number")
		}
		parameters.terminateGRPCServerTimeout = uint32(value)
	}

	serializePayloadBase64, exists := os.LookupEnv("SERIALIZE_PAYLOAD_BASE64")

	if exists && serializePayloadBase64 == "false" {
		parameters.serializePayloadBase64 = false
	}

	return parameters

}

type ConfigurationMethodExpose struct {
	// TODO: Add pattern validation: queue, fireAndForget
	Pattern            string
	Name               string
	Capped             uint64
	Stream             string
	TimeoutWaitForNext string `yaml:"timeoutWaitForNext"`
}

type ConfigurationServiceExpose struct {
	Service  string
	Consumer struct {
		Concurrency uint32
		Block       string
		Retry       struct {
			Limit    int64
			PageSize int64 `yaml:"pageSize"`
			Deadline string
		}
	}
	Client struct {
		Address string
		Port    uint16
	}
	Methods []ConfigurationMethodExpose
}

type Configuration struct {
	Redis struct {
		Prefix      string
		JobPoolSize int `yaml:"jobPoolSize"`
		Nodes       []struct {
			Sequence uint16
			Host     string
			Port     uint16
			Db       uint8
			Password string
		}
	}
	Server struct {
		Address           string
		Port              uint16
		Enable            bool
		MaxHeaderListSize uint32 `yaml:"maxHeaderListSize"`
	}
	Client struct {
		Enable bool
	}
	Metrics struct {
		Address string
		Port    uint16
		Routes  struct {
			Metrics   string
			Readiness string
			Liveness  string
		}
	}
	Protos struct {
		Services []ConfigurationServiceExpose
	}
}

func loadConfiguration() Configuration {

	content, err := ioutil.ReadFile(parameters.configurationFile)

	if err != nil {
		log.Fatalf("Error while reading configuration file: %v", err)
	}

	configuration := Configuration{}

	err = yaml.Unmarshal(content, &configuration)

	if err != nil {
		log.Fatalf("Error while parsing YAML from configuration file: %v", err)
	}

	defaultRedisValues(&configuration)
	defaultServerValues(&configuration)
	defaultMetricsValues(&configuration)
	defaultProtosValues(&configuration)

	return configuration

}

func defaultRedisValues(conf *Configuration) {

	if len(conf.Redis.Prefix) == 0 {
		conf.Redis.Prefix = "polygate"
	}

	if conf.Redis.JobPoolSize <= 0 {
		conf.Redis.JobPoolSize = 5
	}

	if len(conf.Redis.Nodes) == 0 {
		log.Fatalf("You must specify at least one Redis node")
	}

	for i := range conf.Redis.Nodes {

		node := &conf.Redis.Nodes[i]

		if len(node.Host) == 0 {
			log.Fatalf("Redis node %v must have a valid address", i)
		} else if node.Port <= 0 {
			log.Fatalf("Redis node %v must have a valid port", i)
		}

	}

}

func defaultServerValues(conf *Configuration) {

	if len(conf.Server.Address) == 0 {
		conf.Server.Address = "0.0.0.0"
	}

	if conf.Server.Port <= 0 {
		conf.Server.Port = 4774
	}

	if conf.Server.MaxHeaderListSize <= 0 {
		log.Warnf("A maxHeaderListSize of %v bytes will disable metadata support, do you really want this?", conf.Server.MaxHeaderListSize)
	}

}

func defaultMetricsValues(conf *Configuration) {

	if len(conf.Metrics.Address) == 0 {
		conf.Metrics.Address = "0.0.0.0"
	}

	if conf.Metrics.Port <= 0 {
		conf.Metrics.Port = 2112
	}

	if len(conf.Metrics.Routes.Metrics) == 0 {
		conf.Metrics.Routes.Metrics = "/metrics"
	}

	if len(conf.Metrics.Routes.Readiness) == 0 {
		conf.Metrics.Routes.Readiness = "/ready"
	}

	if len(conf.Metrics.Routes.Liveness) == 0 {
		conf.Metrics.Routes.Liveness = "/live"
	}

}

func defaultProtosValues(conf *Configuration) {

	if len(conf.Protos.Services) == 0 {
		log.Fatal("You must specify at least one service")
	}

	for i := range conf.Protos.Services {

		service := &conf.Protos.Services[i]

		if len(service.Service) == 0 {
			log.Fatalf("Empty service name for service %v", i)
		}

		if service.Consumer.Concurrency <= 0 {
			service.Consumer.Concurrency = 50
		}

		if len(service.Consumer.Block) == 0 {
			service.Consumer.Block = "5000ms"
		}

		if service.Consumer.Retry.Limit <= 0 {
			service.Consumer.Retry.Limit = 3
		}

		if service.Consumer.Retry.PageSize <= 0 {
			service.Consumer.Retry.PageSize = 1000
		}

		if len(service.Consumer.Retry.Deadline) == 0 {
			service.Consumer.Retry.Deadline = "10000ms"
		}

		if len(service.Client.Address) == 0 {
			log.Fatalf("Empty client address for service %v", i)
		}

		if service.Client.Port <= 0 {
			log.Fatalf("Invalid client port for service %v", i)
		}

		if len(service.Methods) == 0 {
			log.Fatalf("You must specify at least one method for service %v", i)
		}

		for t := range service.Methods {
			method := &service.Methods[t]

			if len(method.Name) == 0 {
				log.Fatalf("Invalid method name for service %v, method %v", i, t)
			}

			if method.Pattern != "queue" && method.Pattern != "fireAndForget" {
				log.Fatalf("Invalid method pattern for service %v, method %v", i, t)
			}

			if method.Capped <= 0 {
				method.Capped = 10000
			}

			if len(method.Stream) == 0 {
				method.Stream = method.Name
			}

			if len(method.TimeoutWaitForNext) == 0 {
				method.TimeoutWaitForNext = "1s"
			}

		}

	}

}
