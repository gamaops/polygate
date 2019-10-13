# NodeJS Application

The following content explains how the engineering for NodeJS applications must be designed.

* Every NodeJS project will be written in Typescript

## Project structure

* **assets** - Any resource (file) needed to run the application
* **build** - Typescript compiling result
* **src** - The source folder with all code of application
   * **models** - Business models
   * **processors** - When the application is a queue consumer the processors must be here
   * **workers** - If the application deals with worker threads, the thread code will be placed here
   * **views** - Any code related to data presentation (normally only on edge APIs)
   * **validators** - Code related to data structure/business rule validation (normally only on edge APIs)
   * **services** - Representation of RPC (gRPC) services and methods
   * **admin** - Administrative tasks related code (like on-the-fly configurations management)

## Shared Libraries

We need to be careful to not couple services through shared libraries, however they still useful to avoid code duplication and centralize some interfaces description. You must treat any shared library as any other external dependecy.

### Backend Framework

The [backend framework](https://github.com/gamaops/backend-framework) repository contains common stacks, wrappers and utilities for backend applications, **this repository can't have any business rule**.

### Definitions

Proto definitions, database schemas, mappings and JSON Schemas go here, the main point to create a [definitions](https://github.com/gamaops/definitions) repository is to have a central place where everyone agrees about data structure.

#### Setup Elasticsearch Indexes

You can use the definitions repository to setup Elasticsearch indexes settings:

1. Create a JSON file (`elastic.json`) file with the Elasticsearch client options:
   ```json
   {
      "node": "http://elasticsearch.docker.local:9200/",
      "requestTimeout": 60000
   }
   ```
2. Export a variable with which environment you want to setup: `export ENVIRONMENT=development`
3. Execute the setup script: `npm run setup:elastic -- ./elastic.json`

## CI/CD

Our stack to build pipelines is:

* TravisCI
* DockerHub
* GitHub

### Setting up CI

1. Create a `.travis` folder in the root path of repository
2. Create a `.travis/semver.sh` file an put the content bellow:
   ```bash
   #!/bin/bash

   set -o errexit

   if [[ "$TRAVIS_BRANCH" == "master" && "$TRAVIS_PULL_REQUEST" == "false" ]]
   then

      sudo apt update -y
      sudo apt-get install jq -y

      git config --local user.name "TravisCI"
      git config --local user.email "travis@travis.org"

      git checkout master
      npm run build

      export COMMIT_LOG=`git log -1`
      export TRAVIS_BUILD=`echo $COMMIT_LOG | jq -r -s -R 'split("\n") | .[] | capture("Travis build: (?<buildno>.*)") | .buildno'`

      if [ -z "$TRAVIS_BUILD" ]
      then

         git add -A
         git commit --allow-empty -m "Travis build: $TRAVIS_BUILD_NUMBER"

         export SEMVER_LABEL=`echo $COMMIT_LOG | jq -r -s -R 'split("\n") | .[] | capture("#version:(?<semver>.*)") | .semver'`
         export PRERELEASE_LABEL=`echo $COMMIT_LOG | jq -r -s -R 'split("\n") | .[] | capture("#preid:(?<semver>.*)") | .semver'`

         if [ -z "$PRERELEASE_LABEL" ]
         then
            npm version $SEMVER_LABEL -m "Release v%s"
         else
            npm version prerelease --preid=$PRERELEASE_LABEL -m "Release v%s"
         fi

         git remote add origin-remote https://${GH_TOKEN}@github.com/$GH_REPOSITORY.git > /dev/null 2>&1
         git push --quiet --set-upstream origin-remote HEAD:$TRAVIS_BRANCH
         git push --quiet --set-upstream origin-remote --tags HEAD:$TRAVIS_BRANCH
      fi

   fi
   ```
3. Now turn this file into an executable executing the following command: `chmod +x .travis/semver.sh`
4. An your `.travis.yml`:
   ```yml
   language: node_js
   node_js:
   - '12'
   if: "(branch != master AND type = pull_request) OR (branch = master AND commit_message !~ /^Release v/) OR tag =~ ^v"
   stages:
   - test
   - Build and Tag
   script:
   - npm run lint
   - npm run test:unit
   env:
   - GH_REPOSITORY=gamaops/identity-service
   jobs:
   include:
   - stage: Build and Tag
      node_js: '12'
      name: Build and Tag
      script:
      - ".travis/semver.sh"
      if: tag !~ ^v AND branch = master AND type = push
   notifications:
   slack:
      rooms:
      - secure: "SECURE_TOKEN"
      on_success: always
      on_failure: always
      on_pull_requests: true

   ```

Don't forget to generate the Slack token and customize the tests if needed.