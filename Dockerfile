# ---- Stage 1: Build ----
# Build context must be the workspace root (planx-lab/) so that local replace
# directives in go.mod can resolve sibling modules.
FROM --platform=$BUILDPLATFORM golang:1.27.1-alpine3.24@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS builder
ENV GOTOOLCHAIN=local
RUN apk add --no-cache git
ARG TARGETOS
ARG TARGETARCH

WORKDIR /build

# Copy dependency modules first for layer caching.
COPY planx-proto/ /build/planx-proto/
COPY planx-sdk-go/ /build/planx-sdk-go/

# Copy plugin module and download deps.
COPY planx-plugin-sink-stdout/go.mod planx-plugin-sink-stdout/go.sum ./planx-plugin-sink-stdout/
RUN cd planx-plugin-sink-stdout && go mod download

# Copy plugin source and build.
COPY planx-plugin-sink-stdout/ ./planx-plugin-sink-stdout/
RUN cd planx-plugin-sink-stdout && CGO_ENABLED=0 GOOS="$TARGETOS" GOARCH="$TARGETARCH" go build -mod=readonly -trimpath -buildvcs=false -ldflags="-s -w" -o /sink-stdout ./cmd/plugin

# ---- Stage 2: Runtime ----
FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

RUN apk add --no-cache ca-certificates

RUN addgroup -S planx && adduser -S -G planx planx

WORKDIR /plugin

COPY --from=builder /sink-stdout /plugin/sink-stdout

RUN chown -R planx:planx /plugin

USER planx

ENTRYPOINT ["/plugin/sink-stdout"]
