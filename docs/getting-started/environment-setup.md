# Environment Setup

You'll need a Redis 5+ container running:

```bash
docker run --rm --name=redis -p 6379:6379 redis:5
```

And we're going to create a mocked gRPC server, to do this you need to download the [mock](https://github.com/gamaops/polygate/tree/master/examples/mock) path from Polygate's repository. Once you downloaded the folder, navigate into it and start the gRPC mock server:

```bash
docker run --rm --name=grpcserver -p 4770:4770 -p 4771:4771 -v $(pwd):/proto tkpd/gripmock /proto/mock.proto
```

And to add the stubs execute the `setup.sh` script that is in the **mock** folder.

Now you need to create the configuration file to tell Polygate how to behave, click [here](getting-started/configuration) to learn more about Polygate's configuration file and environment variables.