FROM golang:1.23-alpine3.20 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/schedule cmd/schedule/main.go

FROM alpine:3.20

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/schedule /app/schedule
COPY --from=builder /app/config/config-schedule.yaml /app/config-schedule.yaml

ENTRYPOINT ["./schedule", "-config", "config-schedule.yaml"]