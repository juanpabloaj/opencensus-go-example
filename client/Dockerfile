FROM golang:alpine AS builder

ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/client ./client

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/client ./
RUN chmod +x ./client
ENTRYPOINT ["./client"]
