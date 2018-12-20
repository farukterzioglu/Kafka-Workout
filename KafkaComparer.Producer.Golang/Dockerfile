FROM golang:1.10.1 AS builder
RUN go version

COPY . /go/src/github.com/farukterzioglu/kafka-go/
WORKDIR /go/src/github.com/farukterzioglu/kafka-go/
RUN set -x && \
    go get github.com/Shopify/sarama 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/farukterzioglu/kafka-go/app .

ENTRYPOINT ["./app"]