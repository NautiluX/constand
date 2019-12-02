# standup

Small helper to generate random standup order

## TL;DR

```
docker build . -t order:latest
docker run -t -i -p 8081:8081 order:latest
curl localhost:8081?team=Alice&team=Bob
```
