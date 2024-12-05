# Graceful Shutdown Example in Go

This repository demonstrates how to implement graceful shutdown patterns in a Go HTTP service. The example shows best practices for handling shutdown signals and ensuring in-flight requests are completed before the service terminates.

## Key Features

- Signal handling for graceful shutdown (SIGTERM/SIGINT)
- Configurable shutdown timeout (30 seconds default)
- Proper context cancellation chain
- Monitoring for forced shutdown signals during graceful shutdown
- Clean error handling and logging using `slog`

## Project Structure

```
.
├── LICENSE
├── README.md
├── cmd
│   └── cmd.go
│   └── main.go
├── go.mod
└── service
    └── service.go
```

## How It Works

### Service Initialization

The service is initialized with a base context and address in `service/service.go`:

```go
func New(ctx context.Context, addr string) *Service {
    return &Service{
        httpserver: &http.Server{
            Addr:        addr,
            BaseContext: func(net.Listener) context.Context { return ctx },
        },
        addr: addr,
    }
}
```

### Graceful Shutdown Process

The main shutdown logic in `cmd/main.go` follows these steps:

1. **Signal Handling**: The service listens for OS interrupt signals (CTRL+C) and SIGTERM
2. **Shutdown Initiation**: Upon receiving a shutdown signal, a timeout context is created
3. **Graceful Shutdown**: The service attempts to shut down gracefully within the timeout period
4. **Forced Shutdown**: A secondary signal during shutdown will trigger immediate termination

Key shutdown code:

```go
shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

go func() {
    slog.Info("initiating graceful shutdown")
    shutdownErrCh <- s.Close(shutdownCtx)
}()
```

## Usage

To run the service:

```bash
go run cmd/
```

The service will:

1. Start an HTTP server on port 8080
2. Serve a simple "Hello, World!" response at the root endpoint
3. Wait for shutdown signals
4. Attempt graceful shutdown when a signal is received

## Shutdown Behavior

- The service will attempt to complete all in-flight requests during shutdown
- Default shutdown timeout is 30 seconds
- Sending a second interrupt signal during shutdown will force immediate termination
- All shutdown events are logged using structured logging

## Error Handling

The service handles several types of errors:

- Service startup errors
- Shutdown timeout errors
- Errors during the shutdown process itself

All errors are properly wrapped and logged with context.
