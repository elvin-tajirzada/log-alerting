FROM golang:1.22-alpine

WORKDIR /go/src/app

COPY . .

RUN go mod tidy
RUN go mod verify
RUN GOOS=linux go build -o ./bin/log-alerting ./cmd/log-alerting

ENTRYPOINT /go/src/app/bin/log-alerting