# Development Guide

This section will cover the onboarding for developers on GamaOps technology stack. You'll need Docker to run the required dependencies (like Elasticsearch, MongoDB, Redis) and some tooling to work on services/APIs (a.k.a microservices). An important mindset to keep is: everything must run and be analyzed locally before publishing, the CI/CD tools are just to ensure reliability of software, but the first word is up to you so you need to ensure that everything is fine before pushing your code. Some small tasks and over time they become intrinsic are:

* Build Docker images in your machine before publishing them
* Document everything while developing
* Analyze the performance of your application locally (using cAdvisor)
* Think that your application or it's dependencies can fail without breaking the entire domain

## Recipes

1. [Docker](backend-development/docker.md)
2. [NodeJS Runtime](backend-development/nodejs-runtime.md)
2. [Tools](backend-development/tools.md)