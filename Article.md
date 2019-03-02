  
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
docker build -t kafkaconsumer:latest .
$ docker run -e CONSUMER_GROUP=$CONSUMER_GROUP -e TOPIC_NAME=$TOPIC_NAME -e KAFKA_URL=$HOSTIP:9092 --rm -it kafkacomparerconsumer
```

```
// Deploy to K8S
kubectl create -f zookeeper.yaml
kubectl create -f kafka-service.yaml
minikube tunnel // required to get external ip for Kafka-service
kubectl describe svc kafka-service // Note down 'LoadBalancer Ingress' & NodePort 
kubectl create -f kafka-broker.yaml

$ docker run -e TOPIC_NAME=commands -e KAFKA_URL=127.0.0.1:9092 --rm -it kafkacomparerproducer
$ docker run -e CONSUMER_GROUP=k8stest -e TOPIC_NAME=commands -e KAFKA_URL=10.110.251.182:9092 --rm -it kafkacomparerconsumer

$ TOPIC_NAME=commands KAFKA_URL=10.110.251.182 dotnet run .

kubectl apply -f ./KafkaComparer.Consumer/deployment.yaml
kubectl get pods // get pod name of kafkaconsumer 
kubectl logs -f kafkaconsumer-[***]

kubectl apply -f ./KafkaComparer.Producer/deployment.yaml
```

