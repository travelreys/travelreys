# TiinyPlanet

TiinyPlanet is a travel platform that provides the following features:
* Real-time multi-player collaboration for trip itinerary planning and retrospective.
* Travel logistics integrator: Flights, Transists, Lodgings, Restaurants, Insurance, Activities
* Content management for user-generated travel guides, with comments, reviews and thumbs-up.
* Social platform with a media feed for sharing and viewing guides, pictures, videos and images.
* Merchant plaform for tourism merchants to manage bookings, services and visitors.
* Location-based gamification like Pokemon Go that has quests like the escape room thingy that will provide experience points and rewards.
* Traveller experience and rewards from games, comments, reviews etc.

## Dependencies

|: Dependency :|: Version :|
|--------------|-----------|
| Go           | v1.19.4   |
| NodeJS       | v18.12.1  |
| MongoDB      | v6.0.3    |
| Redis        | v7.0      |
| nats.io      | v2.9.10   |

## Getting Started

```bash
# NATS.io
docker run -d -p 4222:4222 -p 8222:8222 -p 6222:6222 --name nats nats:2.9.10

# Redis
docker run -d --name redis redis:7.0

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

> Remember to configure `.envrc` with the correct environment variables!

## Testing
```bash
make test
```
