# Running Polygate

As simple as it can be, now you just need to start Polygate:

```bash
docker run --rm --name=polygate -p 4774:4774 -p 9090:9090 \
	-e "CONFIGURATION_FILE=/tmp/configuration.yaml" \
	-e PRETTY_LOG="true" \
	-v $(pwd)/configuration.yaml:/tmp/configuration.yaml \
	gamaops/polygate
```

An to consume Polygate you can use [BloomRPC](https://github.com/uw-labs/bloomrpc/releases), connecting to the `127.0.0.1:4774` address and using the following payload as request:

```json
{
  "success": true,
  "currentTime": "2019-10-14T20:52:35.177Z",
  "source": "bloomrpc",
  "fail": false,
  "delay": 0,
  "currentStatus": 1
}
```

Is quite simple to run Polygate, and your downstreams/upstreams will not see it, we designed to behave like a real transparent proxy. And you c*an use the metrics routes specified on `metrics` configuration to see some monitoring data:

* [http://127.0.0.1:9090/metrics](http://127.0.0.1:9090/metrics) - Prometheus metrics
* [http://127.0.0.1:9090/ready](http://127.0.0.1:9090/ready) - Readiness probe
* [http://127.0.0.1:9090/live](http://127.0.0.1:9090/live) - Liveness probe

--------------------

At this point you have knowledge about Polygate's design, setup and in the next sections you'll see more about architectures and engineering perspectives.