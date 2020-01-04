FROM golang:alpine AS builder

ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/server ./server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/server ./
RUN chmod +x ./server
ENTRYPOINT ["./server"]
EXPOSE 8888
