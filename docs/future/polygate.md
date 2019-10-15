# Polygate

When I first thought about Polygate I thought about it being forever community-first, if in the future some form of monetization appears it will not be through feature limitation or something like a "shareware"/"freemium" style. Talking more about the direction of the project now:

* I plan to maintain it forever, as it brings me motivation to enhance my knowledge about distributed computing
* Anyone who wants to contribute and follows the contributing directives can be a contributor of project
* I don't want Polygate to be just an application that you just download, setup and run, I expect that we can talk about use cases, troubleshooting, enhancements, new features and ideas about distributed computing

## Backlog

This isn't a prioritized backlog and any item in this list is subject to change or cancellation:

- [ ] Add option to retry jobs on queue pattern even the upstream is failing (currently we also fail the job and send the rejection back to producer)
- [ ] Add a "shadow consumer" option to provide a way to two consumers consume the same stream but just one resolve (on queue pattern, fire and forget already supports this feature)
- [ ] Add a way to send back the response to producer from fire and forget pattern
- [ ] Add support to duplex and server stream calls
- [ ] Add support to secure context gRPC (using TLS)
- [ ] Add Helm chart to repository and an all-in-one yaml to deploy on Kubernetes
- [ ] Add support to gRPC compression