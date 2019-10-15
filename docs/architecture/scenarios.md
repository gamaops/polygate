# Scenarios

This section describes many scenarios that Polygate fits well.

## Unary Call (a.k.a Asynchronous Queue)

When you want to just have a request-response flow, you just need to set the method's **pattern** parameter to **queue**. Each request will added as message in Redis Streams and acquired as jobs by consumers. The consumer then sends the same payload to upstream and sends the response through Pub/Sub to producer, so you downstream receives the response. This pattern supports metadata.

## Client Stream (a.k.a Fire and Forget)

Fire and forget pattern is kind of reliable Pub/Sub. Once your client opens a client stream to Polygate it can send as many messages as it wants and Polygate will send to upstream those messages. We recommend using Google's [empty.proto](https://github.com/protocolbuffers/protobuf/blob/master/src/google/protobuf/empty.proto) as response 'cause Polygate will just send an empty buffer as response. Messages that could not be sent to upstream will be retried by Polygate using the retry policy specified in configuration file.

<!-- ## Shadow Consumer (a.k.a Exactly Once Processing -->

## Sidecar

Usually you'll deploy Polygate as standalone instances that can be easily scalled by just adding more instances, but there're some scenarios where you want to deploy as a sidecar of your application. To do this, you just need to set the flag **enable** to false on client or server. When you're deploying the sidecar of downstream (clients) set the **server.enable** to `false`, and on upstreams (server) set **client.enable** to `false`.