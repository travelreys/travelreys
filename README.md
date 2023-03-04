# TiinyPlanet

TiinyPlanet is a travel platform.

## Dependencies

|: Dependency :|: Version :|
|--------------|-----------|
| Go           | v1.19.4   |
| MongoDB      | v6.0.3    |
| Redis        | v7.0      |
| nats.io      | v2.9.10   |

## Getting Started

```bash
# NATS.io
docker run -d -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats nats:2.9.10

# Redis
docker run -d --name redis -p6379:6379 redis:7.0

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

|: Env Vars :|: Description :|
|--------------|-----------|
|TIINYPLANET_CORS_ORIGIN||
|TIINYPLANET_MONGO_URL||
|TIINYPLANET_MONGO_DBNAME||
|TIINYPLANET_NATS_URL||
|TIINYPLANET_REDIS_URL||
|TIINYPLANET_JWT_SECRET||
|TIINYPLANET_UNSPLASH_ACCESSKEY||
|TIINYPLANET_SKYSCANNER_APIKEY||
|TIINYPLANET_SKYSCANNER_APIHOST||
|TIINYPLANET_GOOGLE_MAPS_APIKEY||
|TIINYPLANET_OAUTH_GOOGLE_SECRET_FILE||


