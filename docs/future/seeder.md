# Seeder

Seeder is the Polygate addon to operate event storing mechanisms. The goal here is to provide a simple and reliable way to store events and replay them when the necessary keeping Polygate as a stateless application.

## Backlog

This isn't a prioritized backlog and any item in this list is subject to change or cancellation:

- [ ] Add event store file format specification (data structure related issue)
- [ ] Add metrics about stored events
- [ ] Allow a point-in-time restore
- [ ] Add support to dry run
- [ ] Add support to file rotation
- [ ] Add CLI commands/options to run Polygate to restore events (or maybe a gRPC endpoint?)