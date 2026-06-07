# ratelx

![CI](https://github.com/istvzsig/ratelx/actions/workflows/ci.yml/badge.svg)

> Lightweight, dependency-free rate limiter for Go (token bucket)

ratelx provides simple, production-ready primitives to control request flow in backend systems.

It is designed for real distributed systems where traffic must be shaped and protected.

---

## Why ratelx?

In modern backend systems:

- APIs can be overloaded
- downstream services have rate limits
- traffic spikes must be controlled
- services need protection from cascading load

Without rate limiting, systems become unstable under load.

ratelx helps prevent overload and maintain system stability.

---

## Features

### Token Bucket Algorithm

- Smooth request throttling
- Burst support
- Predictable rate control

### Modes

- Non-blocking mode (`Allow`)
- Blocking mode (`Wait`)

### Production-ready design

- Thread-safe (mutex protected)
- Context-aware cancellation
- Lightweight and dependency-free
- Easy to embed in services

---

## Installation

```bash
go get github.com/istvzsig/ratelx
```

## Quick Example

### Non-blocking usage

```go
limiter := ratelx.New(10, time.Second)

if limiter.Allow() {
	// process request
}
```

### Blocking usage

```go
err := limiter.Wait(ctx)
if err != nil {
	// context cancelled or timeout
	return err
}
```

## Blocking Mode (Wait)

Wait blocks until a request is allowed by the limiter.

It also supports context cancellation and measures real waiting time.

Example:

```text
wait request 1 → <nil> (waited 0s)
wait request 2 → <nil> (waited 0s)
wait request 3 → <nil> (waited 0s)
wait request 4 → <nil> (waited 1.2s)
wait request 5 → <nil> (waited 0.8s)
```

### Example Output

```bash
request 1 → ALLOWED
request 2 → ALLOWED
request 3 → ALLOWED
request 4 → RATE LIMITED
request 5 → RATE LIMITED
request 6 → ALLOWED
```

## How it works

ratelx uses a token bucket model:

- Tokens are added over time at a fixed rate
- Each request consumes one token
- Requests are allowed only if tokens are available
- Excess requests are either blocked or delayed

## Architecture

```text
time → → → →

tokens
  ↑        refill over time
  |
  |  consume on request
  ↓

[ Bucket Capacity ]
```

## Use Cases

ratelx is ideal for:

- API gateways
- microservices protection
- external API consumption
- background job throttling
- distributed systems traffic shaping

## Design Goals

- Simple and predictable
- No external dependencies
- Easy to embed in Go services
- Safe for concurrent use
- Production-oriented design

## Comparison with retryx

```text
| retryx           | ratelx             |
| ---------------- | ------------------ |
| Handles failures | Prevents overload  |
| Retry logic      | Request throttling |
| Circuit breaker  | Traffic shaping    |
| Resilience       | Protection         |
```

Together they form a backend reliability toolkit.

## Example Use with retryx (concept)

```text
Incoming Request
        ↓
   ratelx (limit traffic)
        ↓
   retryx (handle failures)
        ↓
   External Service
```

## Testing

Run:

```bash
go test ./...
```

Run with race detector:

```bash
go test -race ./...
```

## Behavior

- Smooth token refill over time
- Burst capacity supported
- Context-aware blocking mode
- Thread-safe implementation

## Future improvements

Planned enhancements:

- per-key rate limiting (IP / API key)
- sliding window limiter
- distributed Redis-backed mode
- metrics / observability hooks

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Summary

ratelx is a minimal, production-oriented rate limiter for Go systems.

No frameworks. No dependencies. No complexity.
