# Monitoring

Our monitoring stack consists of:

* Elastic Stack
* Prometheus

## Elastic Stack

The elastic stack should have capabilities to handle:

* Event sourcing
* Tracing
* Caching views

Applications should only write directly to Elasticsearch when they're dealing with cache, otherwise they will log the entries to standard output (stdout) and Filebeat + Logstash will be responsible to store this data in Elasticsearch. This architecture is aware of Elasticsearch backpressure and enables us to choose different locations where we can put these kinds of data using the output plugins of Logstash.

## Prometheus

Evey application must expose a HTTP endpoint that can be scraped by Prometheus:

* `/metrics` - Must expose application and business metrics.
* `/metrics/heath` - Must expose application health (like a liveness probe) and respond with `200` status code only if it's alive.
* `/metrics/ready` - Must expose application readiness state and respond with `200` status code only if it's ready to work/receive traffic.