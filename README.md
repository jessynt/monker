# Monker

monker is a [NSQ](https://github.com/youzan/nsq) to [Kafka](https://github.com/apache/kafka) bridge

# Configuring

### Config file

Config should have the following structure:

```yaml
logLevel: info
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