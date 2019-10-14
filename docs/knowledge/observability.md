# Observability

And finally Polygate provides an entire set of observability features. Again, Polygate was thought to fit well on Kubernetes environment.

> I need to say it again, pay attention on Polygate's log, distributed tracing can be achieved with information provided by logs

## Metrics server

The metrics server exposes Prometheus metrics to track things like:

* Histogram of jobs execution time
* Histogram of jobs/events payload size
* Jobs counter
* Failures counter
* Gauge to count current active client streams
* Gauge to count current listeners waiting for response

## Liveness probe

Polygate provides a dedicated route for liveness probe that responds:

* 200 if everything is up
* 503 if everything is down

## Readiness probe

Polygate provides a dedicated route for readiness probe that responds:

* 200 if everything is ready
* 206 if it is degraded (for example: if you're running consumer and producer mode on the same instance and consumer is down)
* 503 if everything is not ready

## Profiling

We also expose a [pprof](https://github.com/google/pprof) route through `/debug/pprof` so you can easily profile Polygate's runtime.