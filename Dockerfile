FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod verify
RUN GOOS=linux go build -o ./bin/log-alerting ./cmd/log-alerting

ENTRYPOINT /app/bin/log-alerting