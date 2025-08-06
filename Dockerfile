FROM golang:1.24 AS builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app ./cmd/service/main.go

FROM alpine:latest
WORKDIR /opt/
ENV APP_PATH="/opt/app"
COPY --from=builder /go/bin/app ${APP_PATH}
CMD ["/opt/app"]
EXPOSE 8080 50051