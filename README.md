docker build -t kafkacomparerconsumer:latest .
docker run -e CONSUMER_GROUP='tags-consumers' -e TOPIC_NAME='tags' -e KAFKA_URL='172.31.162.65:9092' --rm -it kafkacomparerconsumer

