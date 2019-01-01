FROM golang:1.10.1 AS builder
RUN go version

COPY ./CommandEngine /go/src/github.com/farukterzioglu/KafkaComparer/CommandEngine/
WORKDIR /go/src/github.com/farukterzioglu/KafkaComparer/CommandEngine/
RUN set -x && \
    go get github.com/bsm/sarama-cluster 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/farukterzioglu/KafkaComparer/CommandEngine/app .

ENTRYPOINT ["./app"]