# Constand

Small helper to generate random standup order, deterministic by date.
That means the same result for everyone curling the team URL at the same (UTC) day, or running the CLI locally.

Useful for standups or finding volunteers in self-organized distributed teams.

## How to use

Constand can be run on a server or as cli. To build locally, run

```
go build .
```

### Running a web server with docker

```
docker build . -t order:latest
docker run -t -i -p 8081:8081 order:latest
curl "localhost:8081?team=Alice&team=Bob"
curl "localhost:8081/pick/one/for?team=Alice&team=Bob&purpose=retro"
```

### Running a web server using cli

```
go build .
./constand -l
curl "localhost:8081?team=Alice&team=Bob"
curl "localhost:8081/pick/one/for?team=Alice&team=Bob&purpose=retro"
```

### Run cli

To run as CLI, the team needs to be defined in the config file `~/.constand.yaml`:

```yaml
---
team:
- Alice
- Bob
```

Afterwards you can quickly generate a standup order with

```
./constand 
```

or pick a volunteer with

```
./constand -p teammeeting -1 #this is a number
```

## The name

**Con**sistent + **stand**up order = **Constand**
