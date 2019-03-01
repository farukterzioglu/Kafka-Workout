  
```
// Start Kafka

// Set HOSTIP on every terminal instances
$ HOSTIP=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')

docker run -d --name zookeeper -p 2181:2181 jplock/zookeeper
$ docker run -d --name kafka -p 7203:7203 -p 9092:9092 -e ZOOKEEPER_IP=$HOSTIP ches/kafka
```
  
```
// Go inside of Kafka
docker exec -it kafka bash
ls ./bin
```
 
```
// Create new topic 
$ docker run --rm ches/kafka kafka-topics.sh --create --topic commands --replication-factor 1 --partitions 3 --zookeeper $HOSTIP:2181
$ docker run --rm ches/kafka kafka-topics.sh --list --zookeeper $HOSTIP:2181
```

```
// Start produder
$ docker run --rm --interactive ches/kafka kafka-console-producer.sh --topic commands --broker-list $HOSTIP:9092

// Start consumer
$ docker run --rm ches/kafka kafka-console-consumer.sh --topic commands --from-beginning --bootstrap-server $HOSTIP:9092
$ docker run --rm ches/kafka kafka-console-consumer.sh --topic commands --from-beginning --bootstrap-server $HOSTIP:9092 --consumer-property group.id=testApp1
$ docker run --rm ches/kafka kafka-console-consumer.sh --topic commands --from-beginning --bootstrap-server $HOSTIP:9092 --consumer-property group.id=testApp2

$ docker run --rm --interactive ches/kafka kafka-consumer-groups.sh --new-consumer --describe --group testApp1 --bootstrap-server $HOSTIP:9092
```

```
// Start .Net consumer
$ cd ./KafkaComparer.Consumer
$ export TOPIC_NAME=commands  KAFKA_URL=$HOSTIP CONSUMER_GROUP=commands-consumers
$ dotnet run .

// Start a new consumer with same group id (from another terminal) 
$ export TOPIC_NAME=commands  KAFKA_URL=$HOSTIP CONSUMER_GROUP=commands-consumers
$ dotnet run .

// Start a new consumer with different group id
$ TOPIC_NAME=commands  KAFKA_URL=$HOSTIP CONSUMER_GROUP=commands-consumers-1 dotnet run .

// Start .Net producer
$ cd ./KafkaComparer.Producer  
$ TOPIC_NAME=commands KAFKA_URL=$HOSTIP dotnet run .
```

```
// Build & run producer on Docker
$ export TOPIC_NAME=commands  KAFKA_URL=$HOSTIP

$ cd ./KafkaComparer.Producer
docker build -t kafkacomparerproducer:latest .
$ docker run -e TOPIC_NAME=$TOPIC_NAME -e KAFKA_URL=$HOSTIP:9092 --rm -it kafkacomparerproducer
```

```
// Build & run consumer on Docker
$ export TOPIC_NAME=commands  KAFKA_URL=$HOSTIP CONSUMER_GROUP=commands-consumers

$ cd ./KafkaComparer.Consumer
docker build -t kafkacomparerconsumer:latest .
$ docker run -e CONSUMER_GROUP=$CONSUMER_GROUP -e TOPIC_NAME=$TOPIC_NAME -e KAFKA_URL=$HOSTIP:9092 --rm -it kafkacomparerconsumer
```