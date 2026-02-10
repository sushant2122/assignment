FROM golang:1.21-alpine3.19 AS tier-1

WORKDIR /app
ENV CGO_ENABLED=0
RUN go mod init backend
COPY main.go .
RUN go mod tidy
RUN go build -o backend main.go

FROM alpine:3.23.2
WORKDIR /app
COPY --from=tier-1 /app/backend .
ENTRYPOINT ["/app/backend"]
