### Deploy Kafka to K8S
`kubectl apply -f ./deploy-kafka.yaml`


### Start Kafka & Zookeper
  
```
// Set HOSTIP env param. on evry terminal instance 
$ HOSTIP=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')

$ docker-compose up
$ docker run --rm ches/kafka kafka-topics.sh --list --zookeeper [host ip]:2181 // or -> 
$ docker run --rm ches/kafka kafka-topics.sh --list --zookeeper $(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+'):2181
```
  
```
$ docker run -d --name zookeeper -p 2181:2181 jplock/zookeeper
$ docker run -d --name kafka -p 7203:7203 -p 9092:9092 -e ZOOKEEPER_IP=$HOSTIP ches/kafka
// -e KAFKA_ADVERTISED_HOST_NAME=$HOSTIP -e KAFKA_MESSAGE_MAX_BYTES=3000000 -e KAFKA_REPLICA_FETCH_MAX_BYTES=3100000 
```

```
$ docker run --rm ches/kafka kafka-topics.sh --create --topic tags --replication-factor 1 --partitions 1 --zookeeper $HOSTIP:2181
$ docker run --rm ches/kafka kafka-topics.sh --list --zookeeper $HOSTIP:2181
$ docker run --rm --interactive ches/kafka kafka-console-producer.sh --topic commands --broker-list $HOSTIP:9092

// (From another terminal)  
$ docker run --rm ches/kafka kafka-console-consumer.sh --topic commands --from-beginning --bootstrap-server $HOSTIP:9092 --consumer-property group.id=testApp1

$ docker run --rm --interactive ches/kafka kafka-consumer-groups.sh --new-consumer --describe --group testApp1 --bootstrap-server $HOSTIP:9092
```
Reference: https://github.com/ches/docker-kafka  

### Check kafka instance details

```
docker exec -it zookeeper bash
bin/zkCli.sh -server 127.0.0.1:2181  
ls /brokers  
ls /brokers/topics  
ls /consumers
```

### Consumer with .Net Core

```
$ cd ./KafkaComparer.Consumer  
$ HOSTIP=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')
$ export TOPIC_NAME=commands  KAFKA_URL=$HOSTIP CONSUMER_GROUP=commands-consumers

dotnet run .

<<<<<<< HEAD
docker build -t kafkaconsumer:latest .  
$ docker run -e CONSUMER_GROUP=$CONSUMER_GROUP -e TOPIC_NAME=$TOPIC_NAME -e KAFKA_URL=$HOSTIP:9092 --rm -it kafkaconsumer
=======
docker build -t kafkacomparerconsumer:latest .  
$ docker run -e CONSUMER_GROUP=$CONSUMER_GROUP -e TOPIC_NAME=$TOPIC_NAME -e KAFKA_URL=$HOSTIP:9092 --rm -it kafkacomparerconsumer
>>>>>>> article
```

### Producer with .Net Core

```
$ cd ./KafkaComparer.Producer  
$ HOSTIP=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')
$ export TOPIC_NAME=commands  KAFKA_URL=$HOSTIP

dotnet run .

docker build -t kafkacomparerproducer:latest .
$ docker run -e TOPIC_NAME=$TOPIC_NAME -e KAFKA_URL=$HOSTIP:9092 --rm -it kafkacomparerproducer
```

### Consumer with Golang

```
$ HOSTIP=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')
$ cd ./KafkaComparer.Consumer.Golang
$ go run . -kafka_brokers=$HOSTIP:9092 -group_id="test"

docker build -t kafkacomparerconsumer:go .
$ docker run -it kafkacomparerconsumer:go -kafka_brokers=$HOSTIP:9092 -group_id="test"
```  

### Producer with Golang

```
$ HOSTIP=$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')
$ cd ./KafkaComparer.Producer.Golang
$ go run . -kafka_brokers=$HOSTIP:9092
  
docker build -t kafkacomparerproducer:go .  
$ docker run -it kafkacomparerproducer:go -kafka_brokers=$HOSTIP:9092
```

### Producer & Consumer with 'kafkacat'

```
docker run --interactive --rm confluentinc/cp-kafkacat kafkacat -b 192.168.20.180:9092 -t tags -K: -P  
docker run --tty --interactive --rm confluentinc/cp-kafkacat kafkacat -b 192.168.20.180:9092 -L https://github.com/edenhill/kafkacat
```

### Topic with 3 partition

`docker run --rm ches/kafka kafka-topics.sh --create --topic tagsPart3 --replication-factor 1 --partitions 3 --zookeeper 172.31.162.65:2181`
