# Event sourcing

At this point you I hope you can see the implicit nature of Polygate's event sourcing pattern. We can just send the payload through streams and log the payload to `stdout` so you can store it on any event store using any tool to grab the logs from Polygate, for example:

* As grabber you can use Filebeat, Fluentd
* As event store you can use Elasticsearch, InfluxDB, MongoDB

> Logs are a powerful source of data from Polygate, don't underestimate it, don't treat them as generic system logs

The Seeder is on Polygate's project roadmap, feel free to contribute with ideas, it is on brainstorming stage.

## How events are logged

Events are logged using the `job-event.proto` (that can be found under the `polygate-data/` folder) and encoded as base64 strings to easily be transported across any channel. By default, Polygate will write every log line as JSON strings, for development time you can enable the `PRETTY_LOG=true` environment variable to understand the output.

## Prefer Polygate for commands

Polygate is designed for not-so-big payloads, I can't tell you "just for small payloads" because small is relative. To understand why, remember that this project is suitable for event sourcing, CQRS architectures, so usually you will have views for queries and don't need to pass the request at query time through services hops. In the future we can implement this using gRPC server streaming. Polygate complies with Protobuf's primitives and one is to not have big messages.