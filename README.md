# Travelreys

Travelreys is a travel platform.

## Dependencies

| : Dependency : | : Version : |
| -------------- | ----------- |
| Go             | v1.19.4     |
| MongoDB        | v6.0.3      |
| Redis          | v7.0        |
| nats.io        | v2.9.10     |

## Getting Started

```bash
# NATS.io
docker run -d -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats nats:2.9.10

# Redis
docker run -d --name redis -p6379:6379 redis:7.0

# etcd
docker run -d --name etcd \
  -p 2380:2380 -p 2379:2379 \
  quay.io/coreos/etcd:v3.3.26 etcd \
    --name etcd0 \
    --advertise-client-urls "http://${HostIP}:2379" \
    --listen-client-urls "http://0.0.0.0:2379" \
    --initial-advertise-peer-urls "http://${HostIP}:2380" \
    --listen-peer-urls "http://0.0.0.0:2380" \
    --initial-cluster "etcd0=http://${HostIP}:2380" \
    --initial-cluster-token cluster-1 \
    --initial-cluster-state new


# MongoDB
docker run -d --name mongo1 -p 27017:27017 mongo:6.0.2 mongod --replSet=rs0
docker run -d --name mongo2 -p 27018:27017 mongo:6.0.2 mongod --replSet=rs0
docker run -d --name mongo3 -p 27019:27017 mongo:6.0.2 mongod --replSet=rs0

docker exec -it mongo1 mongosh

rs.initiate({
  _id: "rs0",
  members: [
    { _id: 0, host: "docker.for.mac.host.internal:27017", priority: 1 },
    { _id: 1, host: "docker.for.mac.host.internal:27018", priority: 0.5 },
    { _id: 2, host: "docker.for.mac.host.internal:27019", priority: 0.5 },
  ]
});

rs.status()

# Minio
docker run --name minio -d -p 9000:9000 -p 9001:9001 \
  quay.io/minio/minio server /data --console-address ":9001"

# Mailhog
docker run --rm -d \
  --name mailhog \
  -p 1025:1025 \
  -p 8025:8025 \
  mailhog/mailhog


# Build Server
make
```

## Run Locally
```bash
$ ./build/server --help
Usage of ./build/server:
      --host string        host address to bind server
      --log-level string   log level
      --port string        http server port

$ ./build/coordinator --help
Usage of ./build/coordinator:
      --host string        host address to bind server
      --log-level string   log level
```

> Remember to configure `.envrc` with the correct environment variables!

## Testing
```bash
make test
```

## Environment variables

| : Env Vars :                        | : Description : |
| ----------------------------------- | --------------- |
| TRAVELREYS_CORS_ORIGIN              |                 |
| TRAVELREYS_MONGO_URL                |                 |
| TRAVELREYS_MONGO_DBNAME             |                 |
| TRAVELREYS_NATS_URL                 |                 |
| TRAVELREYS_REDIS_URL                |                 |
| TRAVELREYS_JWT_SECRET               |                 |
| TRAVELREYS_UNSPLASH_ACCESSKEY       |                 |
| TRAVELREYS_SKYSCANNER_APIKEY        |                 |
| TRAVELREYS_SKYSCANNER_APIHOST       |                 |
| TRAVELREYS_GOOGLE_MAPS_APIKEY       |                 |
| TRAVELREYS_OAUTH_GOOGLE_SECRET_FILE |                 |

