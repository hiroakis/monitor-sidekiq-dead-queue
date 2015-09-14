# monitor-sidekiq-dead-queue

If the tool finds new dead queue, post it to slack.

## Build

```
make build
```

## Usage

```
monitor-sidekiq-dead-queue \
    -host REDIS-HOSTNAME \
    -port 6379 \
    -posfile /tmp/monitor-sidekiq-dead-queue.pos \
    -url https://hooks.slack.com/services/xxxxxxxx \
    -channel ALERT-CHANNEL \
    -slackusername sidekiq dead queue \
    -iconemoji interrobang
```

## License

MIT