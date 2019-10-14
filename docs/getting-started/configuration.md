# Configuration

There're two sepparate configuration sources to setup Polygate: environment variables and the configuration file.

## Environment Variables

- `CONFIGURATION_FILE` - The path to the configuration file, this setting is required.
- `PRETTY_LOG` - Prints pretty log instead of JSON lines. The default value is "false" you can enable this feature setting it to "true".
- `ENABLE_HOT_RELOAD` - Enables configuration file change detection, when the file is changed, Polygate will stop gracefully so your orchestrator can restart it. The default value is "false" you can enable this feature setting it to "true".
- `LOG_LEVEL` - Can be: debug, info, warn or error.

## Configuration File

```yaml
redis:
  prefix: polygate # Every item on Redis will have this prefix, this avoids collisions between two applications using the same Redis
  jobPoolSize: 5 # The size of Redis connection pool to send jobs related commands like XADD, PUBLISH
  nodes: # Redis nodes to distribute the workload
    - sequence: 1 # Standalone Redis server
      host: "127.0.0.1"
      port: 6379
      db: 1 # This is not required, by default it'll pick the db 0
      password: supersecret # This is not required
    - sequence: 2 # The sequence parameter ensures the ordering of this list (it's important to the partioning algorithm)
      sentinel: true # Enables sentinel mode
      master: redismaster # Master name of sentinel
      sentinelNodes:
        - host: "172.22.1.34"
          port: 26379
        - host: "172.22.1.33"
          port: 26379
        - host: "172.22.1.32"
          port: 26379
server:
  address: "0.0.0.0" # Address to bind server's listener
  port: 4774 # Port to bind server's listener
  enable: true # Enables Polygate's producer mode
  maxHeaderListSize: 1048576 # Maximum size (bytes) of headers from a gRPC call (this parameter interferes on metadata size)
client:
  enable: true # Enables Polygate's consumer mode
metrics:
  address: "0.0.0.0" # Address to bind the metrics server
  port: 9090 # Port to bind the metrics server
  shutdownTimeout: "15s" # When Polygate is gracefully stopping this is the maximum allowed time to wait before forcing a close
  routes: # Routes must start with /
    metrics: /metrics
    readiness: /ready
    liveness: /live
protos:
  services: # Defines the accepted gRPC calls
    - service: mock.v1.MockedUnaryService # package name+"."+service name
      consumer:
        concurrency: 200 # Maximum number of jobs to be picked from Redis Streams at the same time
        block: "5000ms" # How much we must block a XREADGROUP call while waiting for jobs
        retry: # Retry policy
          limit: 3 # Maximum number of times to retry a stalled out job
          pageSize: 1000 # This is the pagination setting to iterato over the PEL list, usually you can keep this setting as 1000
          deadline: "10000ms" # How many time must be passed since a job was not acknowledged to consider stalled out
      client:
        address: "127.0.0.1" # Upstream gRPC server address
        port: 4770 # Upstream gRPC server port
      methods:
        - name: SendUnaryMock # The gRPC method name
          pattern: queue # The queue pattern will behave like request-response flow
          capped: 500 # The maximum number of messages on stream, old messages will be removed from Redis, increase this number if you have high throughput
          stream: SendUnaryMock # The Redis Stream name
    - service: mock.v1.MockedClientStreamService
      consumer:
        concurrency: 100
        block: "5000ms"
        retry:
          limit: 3
          pageSize: 1000
          deadline: "10000ms"
      client:
        address: "127.0.0.1"
        port: 4770
      methods:
        - name: SendClientStreamMock
          pattern: fireAndForget # Fire and forget pattern will behave like a pub/sub with delivery guarantees
          capped: 500
          stream: SendClientStreamMock
          timeoutWaitForNext: "300ms" # Time to consider a client stream connected to the upstream as idle to remove from consumer's pool, this timer is reset when a new job arrives
```

Create a file named `configuration.yaml` and get the IP of your Redis container:

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' redis
```

Now just remove the Sentinel configuration as we're just using a standalone Redis and don't forget to replace the host IP. The Redis configuration will look like:

```yaml
redis:
  prefix: polygate
  jobPoolSize: 5
  nodes:
    - sequence: 1
      host: "IP_OF_YOUR_REDIS_CONTAINER"
      port: 6379
```

And do the same thing with the upstream IPs, first get the IP:

```bash
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' grpcserver
```

Replace the IP in clients specifications:

```yaml
protos:
  services:
    - service: mock.v1.MockedUnaryService
      consumer:
        concurrency: 200
        block: "5000ms"
        retry:
          limit: 3
          pageSize: 1000
          deadline: "10000ms"
      client:
        address: "YOUR_UPSTREAM_IP" # Replace the IP here
        port: 4770
      methods:
        - name: SendUnaryMock
          pattern: queue
          capped: 500
          stream: SendUnaryMock
    - service: mock.v1.MockedClientStreamService
      consumer:
        concurrency: 100
        block: "5000ms"
        retry:
          limit: 3
          pageSize: 1000
          deadline: "10000ms"
      client:
        address: "YOUR_UPSTREAM_IP" # Replace the IP here
        port: 4770
      methods:
        - name: SendClientStreamMock
          pattern: fireAndForget
          capped: 500
          stream: SendClientStreamMock
          timeoutWaitForNext: "300ms"
```
