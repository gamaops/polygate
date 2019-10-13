# The desired solution

So I started to think about a sidecar application to make my services asynchronous. But a sidecar application is just an additional hop to our requests and creates unnecessary complexity, obvious I'm not talking about all the cases (we can see the Istio performs great with this approach). So what about a proxy? Well, it seems that a proxy is the ideal solution, we can keep it stateless, it's totally abstract to application, and we can use as sidecar as well when we want. But a proxy is not asynchronous and if I want to have failover/HA features I'll have to deal with a temporary transient state (the requests/responses).

And finally I got my insight: a hub-like application, where it can operate in various modes! And now we need to define some standards here.

> We must talk about two parts of the problem: who is producing events and who is consuming events

## Producers

Producers will produce events from gRPC server requests and deal with the downstream (services connected to producers). They will be the "gate" of hub and must provide event sourcing patterns and monitoring about downstreams.

## Consumers

Consumers will get events from streams and interact with upstreams through gRPC clients, it will send the feedback to producers and retry stalled out jobs. The feedback is only needed when jobs are produced on queue pattern, fire and forget pattern doesn't need this. Queues will be translated as unary calls from gRPC specification and fire and forget will be client streaming. Currently I don't have ideas to implement duplex or server streaming calls (maybe a pub/sub pattern?) so feel free to propose an implementation.