# NodeJS Runtime Environment

This page covers all you need to know about the runtime environment for NodeJS backend applications. We prefer NodeJS LTS versions and everything here was developed/tested with NodeJS v12+.

## Setup

1. Install [pm2](http://pm2.keymetrics.io/) globally
2. To run any service/API you first need to build it:
   1. Pull any service (for example: https://github.com/gamaops/identity-service)
   2. Execute: `npm run build`
3. Now you can start the app using pm2: `pm2 start build/index.j --name identity-service`
   1. To view the dashboard: `pm2 dash`
   2. To view logs: `pm2 logs`
   3. To view pretty logs: `pm2 logs --raw | node_modules/.bin/bunyan`

## Docker

1. Read the best practices guide to dockerize NodeJS applications: https://github.com/nodejs/docker-node/blob/master/docs/BestPractices.md