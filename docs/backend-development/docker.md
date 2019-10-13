# Docker

Every configuration here was tested with **Docker version 18.09.6, build 481bc77**.

## Compose stack

The compose bellow is the current required full backend stack to run any service/API, you can customize this file removing any technology that you don't need.

```yml
version: "3.7"

networks:
  private_stack_net:
    driver: bridge

services:
  redis:
    image: redis:5
    networks:
      - private_stack_net
    ports:
      - "6379:6379"
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          cpus: '0.3'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M
  mongo:
    image: mongo:4-xenial
    networks:
      - private_stack_net
    ports:
      - "27017:27017"
    environment:
      MONGO_INITDB_ROOT_USERNAME: "root"
      MONGO_INITDB_ROOT_PASSWORD: "123456"
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          cpus: '0.3'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M
  mongo-express:
    image: mongo-express
    ports:
      - 8081:8081
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: "root"
      ME_CONFIG_MONGODB_ADMINPASSWORD: "123456"
    networks:
      - private_stack_net
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          cpus: '0.3'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:7.3.0
    environment:
      - node.name=es01
      - discovery.type=single-node
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
    ports:
      - "9200:9200"
    networks:
      - private_stack_net
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          cpus: '0.5'
          memory: 600M
        reservations:
          cpus: '0.25'
          memory: 512M
  kibana:
    image: docker.elastic.co/kibana/kibana:7.3.0
    environment:
      ELASTICSEARCH_HOSTS: http://elasticsearch:9200
    ports:
      - "5601:5601"
    networks:
      - private_stack_net
    deploy:
      mode: replicated
      replicas: 1
      resources:
        limits:
          cpus: '0.3'
          memory: 256M
        reservations:
          cpus: '0.25'
          memory: 128M
```

And you can start this stack executing the command: `docker-compose -f ./stack.yaml up`

Set the following domains on `/etc/hosts`:

```
127.0.0.1	kibana.docker.local
127.0.0.1	elasticsearch.docker.local
127.0.0.1	mongo.docker.local
127.0.0.1	redis.docker.local
```

You'll now have the following endpoints:

* Mongo Express - http://mongo.docker.local:8081/
* Kibana - http://kibana.docker.local:5601
* Elasticsearch - http://elasticsearch.docker.local:9200
* Redis - redis.docker.local:6379
* MongoDB - mongo.docker.local:27017

## Envoy

To run Envoy to act as proxy to gRPC-Web you need to follow these steps:

1. Get the Docker network gateway using the command: `ip route show | grep docker0`
  1. Example: in the line `172.17.0.0/16 dev docker0 proto kernel scope link src 172.17.0.1` the **172.17.0.1** is the gateway
2. Now you need to create an **envoy.yaml** file (don't forget to replace the gateway IP in upstreams):
  ```yaml
  admin:
    access_log_path: /tmp/admin_access.log
    address:
      socket_address: { address: 0.0.0.0, port_value: 9901 }

  static_resources:
    listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 5050 }
      filter_chains:
      - filters:
        - name: envoy.http_connection_manager
          config:
            codec_type: auto
            stat_prefix: ingress_http
            route_config:
              name: local_route
              virtual_hosts:
              - name: local_service
                domains: ["*"]
                routes:
                - match: { prefix: "/identity." }
                  route:
                    cluster: identity_api
                    max_grpc_timeout: 0s
                cors:
                  allow_origin:
                  - "*"
                  allow_methods: GET, PUT, DELETE, POST, OPTIONS
                  allow_headers: keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,errid,errno,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout
                  max_age: "1728000"
                  expose_headers: errid,errno,grpc-status,grpc-message
            http_filters:
            - name: envoy.grpc_web
            - name: envoy.cors
            - name: envoy.router
    clusters:
    - name: identity_api
      connect_timeout: 0.25s
      type: logical_dns
      http2_protocol_options: {}
      lb_policy: round_robin
      # win/mac hosts: Use address: host.docker.internal instead of address: localhost in the line below
      hosts: [{ socket_address: { address: '172.17.0.1', port_value: 35102 }}]
  ```
3. Now you can run Envoy:
  ```bash
  docker run -d -p 5050:5050 --rm --name=envoy \
  -v $(pwd)/envoy.yaml:/etc/envoy/envoy.yaml \
  envoyproxy/envoy:latest /usr/local/bin/envoy -c /etc/envoy/envoy.yaml
  ```
4. And if you want to change the configuration file (add more APIs for example) you can edit the configuration file and execute: `docker restart envoy` to reload Envoy