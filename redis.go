package main

import (
	"sort"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/howeyc/crc16"
	log "github.com/sirupsen/logrus"
)

// PubSubClient is a struct to hold pubsub mappings
type PubSubClient struct {
	pubsub *redis.PubSub
}

var redisClients = make(map[string]map[int]*redis.Client)
var pubsubClients = make(map[string]map[int]*PubSubClient)

func startRedisClient(key string, poolSize int) {

	_, exists := redisClients[key]

	if exists {
		log.Fatalf("Redis clients already exists for this key: %v", key)
	}

	redisClients[key] = make(map[int]*redis.Client)

	sort.Slice(configuration.Redis.Nodes, func(i int, j int) bool {
		return configuration.Redis.Nodes[i].Sequence > configuration.Redis.Nodes[j].Sequence
	})

	for index, node := range configuration.Redis.Nodes {

		var address strings.Builder

		address.WriteString(node.Host)
		address.WriteRune(':')
		address.WriteString(strconv.Itoa(int(node.Port)))

		redisClients[key][index] = redis.NewClient(&redis.Options{
			Addr:         address.String(),
			Password:     node.Password,
			DB:           int(node.Db),
			Network:      "tcp",
			PoolSize:     poolSize, // TODO: Test pool size
			MinIdleConns: poolSize,
		})
	}

}

func routeRedisClient(key string, route []byte) *redis.Client {

	checksum := crc16.Checksum([]byte(route), crc16.IBMTable)
	index := int(checksum) % len(redisClients[key])

	return redisClients[key][index]

}

func closeRedisClients(key string) {

	log.Warnf("Closing Redis client: %v", key)

	for _, client := range redisClients[key] {
		err := client.Close()
		if err != nil {
			log.Errorf("Error while closing Redis client: %v", err)
		}
	}

	log.Warnf("Redis client closed: %v", key)

	delete(redisClients, key)

}
