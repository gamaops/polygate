# Architecture Overview

Let's split the Polygate's implementation into some parts:

* Event store
* Payload streaming
* Networking
* Client, server, hub

## Event store

One advantages to have asynchronous communication is to have the chance to recover from upstream failures without sending failing responses to downstream. We can achieve a better control over the growth of system and using streams we can easily upgrade the downstream or the upstream without impact to the conterparty. But in the stateless era we can start to think about services as just states machines. Obviously you'll use databases to store some state, but this direct communication of commands to services and services to database still vulnerable to transient failures. When you have all the commands sent to your services stored **before** that command is delivered increases resilience to data loss and inconsistency on our state.

Currently Polygate logs each request as base64 encoded events so you can pick this piece of data and store in your preferred event store, something like Elasticsearch, InfluxDB, MongoDB and so on.

Seeder will be the Polygate feature to create a zero-dependency event store. What Seeder will do is much like what Redis does with AOF data format. It'll create append only files (that will be automatically rotated) and add all received events in these files. We're still researching the best way to do this to keep the main goal of Polygate: make distributed computing simple. If you have ideas about this operation contact us through GitHub issues or our chat at Gitter.

Once Seeder is released you'll be able to simply backup these files to any secure storage (like S3) and if you need, use Polygate to restore the state of your upstreams by replaying the events.

## Payload streaming

With Redis Streams we can achieve many, many patterns, like queues, fire and forget, scheduling and etc. Streaming on Redis is very fast (and you can use KeyDB if you're not satisfied with Redis performance). The **capped** parameter of Polygate's configuration allows us to define a maximum size for streams, discarding old messages when this number is exceeded, you can have this number really high if you have enough memory for your payloads and throughput. But we recommend you to not think about Redis Streams as persistent event store, the whole point about making the implementation of multiple Redis nodes to distribute the workload from client side is to make horizontally scaling Redis as simple as possible and this can be really hard if you have strict requirements for persistence on Redis.

Streaming allow us to better know our transaction rate, troughput growth and have an easy way to throttle the request/response flow. Using the metrics provided by Polygate you can see near real time what is happening in your mesh.

## Networking

Most time of Polygate's added latency comes from networking IO. You can think about three paths that Polygate uses to work:

* Redis connections
* Downstream (producer) connections
* Upstream (consumer) connections

You must have exactly this order in mind when prioritizing the reliability and performance of your networking topology. Here's why:

* Without Redis, Polygate will fatally fail
* Without downstreams events can not be produced
* Without upstreams events can not be consumed/recovered

Prefer to put Polygate next to your mesh endpoints, you can use it as sidecar (although it increases the complexity to maintain it) or simply put in the same network of your services.

## Client, server, hub

Polygate was designed to be handled by orchestrators so when any failure occurs in the client->hub->server flow it'll emit a fatal log and exit. Clients and server don't need to be aware of Polygate existence, there're just few things that you need to know when creating your services:

* Polygate don't support metadata on fire and forget pattern (the reason why is explained at the **Knowledge**) section
* We currently don't support duplex or server streaming calls, but is just a matter of thinking how they will be handled by Polygate, so if you have any ideas feel free to open an issue on GitHub or chat with us through Gitter