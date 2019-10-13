# Architecture Guide

First you need to have a background knowledge over the following keywords:

* [CQRS](https://martinfowler.com/bliki/CQRS.html)
* [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)
* [Hexagonal Architecture](https://fideloper.com/hexagonal-architecture)
* [SCS](https://scs-architecture.org/)
* Pub/Sub, queuing, messaging, RPC and distributed computing in general

Every application on GamaOps is built using the architectures explained above and the services topology reflects the hexagonal architecture where there are two distinct kinds of backends: **internal services** and **edge APIs**.

![cqrs](../assets/images/cqrs-architecture.png)

## Internal Services

Internal services are the building blocks of everything that needs to be stateful or persistent or have complex business rules involving asynchronous flows. Services don't expose APIs, instead they communicate through queues ([Redis Streams](https://redis.io/topics/streams-intro)), using pub/sub to announce their jobs completions. 

## Edge APIs

Edge APIs are responsible for expose a domain to be consumed by clients, they have data structure validations, any cacheable business rule. Caching and composing views of data are responsibilities of edge APIs. When they need to deal with persistence they will dispatch jobs to internal services. Another important role of edge APIs is to publish events that are commands to replicate the machine state of persisted data.

## Events

Not everything needs to be an event, the choice of what needs to be an event is based on a simple criteria: does this payload change any state? In other words: is this a command?

Commands normally are requests that generate jobs to services, but can be handled inside the edge APIs boundaries if is just a cache change for example. So there are two things that certanly aren't events:

* Queries (search requests)
* Invalid commands (here we're talking just about edge APIs validations)