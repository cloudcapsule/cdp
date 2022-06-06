FROM golang:1.17.3-alpine3.15 as builder
ARG BUILD_SHA
ARG BUILD_VERSION
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY cmd/ cmd/
COPY gen/ gen/
COPY pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
    go build \
    -ldflags="-X 'main.Build=${BUILD_SHA}' -X 'main.Version=${BUILD_VERSION}'" \
    -v -o bin/cdp \
    cmd/cdp/*.go

FROM alpine:3.15.0
WORKDIR /opt/app-root
COPY --from=builder /workspace/bin/cdp /opt/app-root/cdp

