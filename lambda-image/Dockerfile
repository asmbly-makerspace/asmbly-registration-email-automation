FROM golang:1.22.4 AS builder
WORKDIR /

# Disable cgo to create a static binary.
ENV CGO_ENABLED="0"

# Compile for arm64 Linux
ENV GOOS="linux"
ENV GOARCH="arm64"

# Cache dependencies
ADD app.go go.mod go.sum lambda_events.go ./
ADD ./mailservice ./mailservice
RUN go mod download

# Build
RUN go build -a -o app

# Copy artifacts to a clean image
FROM scratch
WORKDIR /
EXPOSE 5000
COPY --from=builder /app /

CMD ["/app"]