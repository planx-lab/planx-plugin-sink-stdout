# Planx Plugin Template (Go)

This repository is a template for Planx 4.0 Go plugins.

## How it works

- SDK starts gRPC server
- Engine launches plugin as a process
- Plugin only implements business logic

## Run locally

```bash
go run ./cmd/plugin
````

The plugin will start and print a handshake JSON to STDOUT.