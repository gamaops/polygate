# Environment Setup

You'll need a Redis 5+ container running:

```bash
docker run --rm --name=redis -d -p 6379:6379 redis:5
```

Now you need to create the configuration file to tell Polygate how to behave, click [here](getting-started/configuration) to learn more about Polygate's configuration file and environment variables.