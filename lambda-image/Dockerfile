FROM golang:1.22.4-alpine3.20 AS builder
WORKDIR /

# Disable cgo to create a static binary.
ENV CGO_ENABLED="0"

# Compile for arm64 Linux
ENV GOOS="linux"
ENV GOARCH="arm64"

# Update CA Certs
RUN apk --update add ca-certificates

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
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app /

CMD ["/app"]