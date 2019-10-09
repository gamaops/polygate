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
	Pattern     string
	Name        string
	Capped      uint64
	Stream      string
	Deadline    uint32
	ExpiresIn   uint32
	DropRequest bool
}

type ConfigurationServiceExpose struct {
	Service  string
	Consumer struct {
		Concurrency uint32
		Block       string
		Retry       struct {
			Limit    int64
			PageSize int64
			Deadline string
		}
	}
	Client struct {
		Address     string
		Port        uint16
		Healthcheck struct {
			MaxRetryTime uint32
		}
	}
	Methods []ConfigurationMethodExpose
}

type Configuration struct {
	Redis struct {
		Prefix      string
		JobPoolSize int `yaml:"jobPoolSize"`
		Nodes       []struct {
			Sequence uint16
			Address  string
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

	// TODO: Validate configuration

	return configuration

}
