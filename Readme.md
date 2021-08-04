# Beurtbalkje

> Connection queue proxy for restarting services

`usage: beurtbalkje [-port=8888] [-timeout=30] target`

Beurtbalkje is a tcp proxy that proxies all traffic to another service, but keeps accepting connections while the target server is down.

Therfor instead of a connection refused error while the (development) service is restarting, the connection is accepted and the client waits until it gets to send and receive data.
As soon as the target server is up again the connection is proxied.

This makes it look to the client as if the service took some extra time to respond.

## Usage (npm)

```shell
npx beurtbalkje --port=8888 localhost:8080
```

## Usage (go)

```shell
go install github.com/bfanger/beurtbalkje

beurtbalkje --port=8888 localhost:8080
```

## Misc

A ["beurtbalkje"](https://nl.wikipedia.org/wiki/Beurtbalkje) is the dutch word a [checkout divider](https://en.wikipedia.org/wiki/Checkout_divider).
