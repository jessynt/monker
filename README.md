# Monker

monker is a [YouZan NSQ](https://github.com/youzan/nsq) to [Kafka](https://github.com/apache/kafka) bridge

## Architecture

> The Monkey is designed and implemented with a distributed architecture.

![](./doc/architecture/monker.png)

## Configuring

Config should have the following structure:

```yaml
logLevel: info
# You can also use redis as storage (e.g. `redis://redis:6379?key=yourkey`) 
storageDSN: "inmem://unknown"
worker:
  cycleTimeout: "2s"
  cacheSize: 10
  cacheFlushTimeout: "5s"
  storageReadTimeout: "10s"
kafka:
  brokers:
  - 127.0.0.1:1234
  topic: "exchanges"
nsq:
  lookupds:
  - 127.0.0.1:12345
  topic: "example-topic"
  channel: "default"
  concurrency: 1
```