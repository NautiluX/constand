# Constand

Small helper to generate random standup order, deterministic by date.
That means the same result for everyone curling the team URL at the same (UTC) day.

Useful for standups in self-organized distributed teams.

## TL;DR

```
docker build . -t order:latest
docker run -t -i -p 8081:8081 order:latest
curl "localhost:8081?team=Alice&team=Bob"
curl "localhost:8081/pick/one/for?team=Alice&team=Bob&purpose=retro"
```

## The name

**Con**sistent + **stand**up order = **Constand**
