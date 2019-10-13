# Going asynchronous

Have you read about [Redis Streams](https://redis.io/topics/streams-intro)? I worked with Kafka and RabbitMQ and I don't want to appear like I'm against these two, but they create a big infrastructure complexity to provide streaming with single delivery guarantees. I know they're not made to solve this, but Redis is a much more friendly application to developers, has high performance and it's stream implementation is quite simple. But one point about Redis: is not the most painless and clear thing to horizontally scale it. Redis Cluster is missing features like peer auto discovery and an automated sharding mechanism, I think it's hard to implement these features on it and the decision to not put effort to achieve a clustering mode like Elasticsearch is by design.

The major issue scaling Redis horizontally is sharding data. We can easily solve it using a client-side partitioning algorithm. The partitioning algorithm we choose here is quite simple:

* Assume that we have a numeric indexed list of Redis servers: `REDIS_NODE_LIST`
* We also have a string ID for each received request (Polygate uses [xid](https://github.com/rs/xid) to generate these IDs)

> WHICH_NODE_TO_PICK = CRC16(JOB_ID) % LENGTH(REDIS_NODE_LIST)

And that's it! We have our partitioning algorithm, and we can easily reshard our events by just adding/removing nodes to our Redis list. Polygate supports standalone Redis servers or Redis Sentinel as nodes in this list.

## Queues: how to send back the response?

Each producer will listen to a Redis Pub/Sub channel identified by its hostname, so when a stream having the queue pattern is consumed by upstream, the consumer will publish an event with the resolved/rejected job. And keep in mind that Polygate prioritizes the request flow, that's why only requests are sourced through logs.

## Retrying stalled out jobs

It's not so clear how to retry stalled out messages on Redis Streams, the documentation says that you can use PEL messages listed from XPENDING command and XCLAIM to reclaim those messages, but there are a lot of details and I think that in the future Redis will be more mature about this topic. But to deal with this problem I created a little Lua script to do all the logic to retry jobs. There're two important parameters to "XRETRY": deadline and maximum number of retries. The deadline specify how much time must be passed since a job acquired by a consumer but not acknowledged to be candidate to be retried. The maximum number of retries defines a threshold to avoid retrying jobs that are failing every time that they're acquired by consumers (this prevents from anomalies to break the system).

## Concurrency

Asynchronous concurrency is hard, but the simple way that Polygate operates it makes easy to think about multiple Redis nodes being consumed. When you specify the concurrency number for a stream this number is divided by the number of available Redis nodes to be consumed and the resulting number will be maximum number of messages to be acquired by consumers from each node.

## Fire and forget, client streaming

The fire and forget pattern is useful when you want to just publish events and multiple consumers can consume that message to do their job, it's kind of reliable Pub/Sub pattern. But how to know how many clients streams the consumer must open to send their messages? At the consumer side we don't have a notion of the state of the call from producers, when it starts or when it ends. It's not a predictable problem, at least not in a reliable way. So what Polygate does is to keep a pool of client streams connected to upstreams and each client stream have an idle timeout. This pool grows dynamically depending on the throughput of jobs and it'll keep the best size to attend all consumed jobs. When the client stream goes idle (not sending any job) it's automatically closed and removed from the pool. So keep in mind: **Polygate doesn't keep one to one relation of client streams from downstreams to upstreams**. That's why we can't handle metadata on fire and forget pattern, 'cause we can send multiple metadata related to many streams in just one client stream.