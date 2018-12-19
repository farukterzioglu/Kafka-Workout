# Start Kafka & Zookeper # 
https://github.com/ches/docker-kafka

docker run -d --name zookeeper -p 2181:2181 jplock/zookeeper  
docker run -d --name kafka -p 7203:7203 -p 9092:9092 -e KAFKA_ADVERTISED_HOST_NAME=172.31.162.65 -e KAFKA_MESSAGE_MAX_BYTES=3000000 -e KAFKA_REPLICA_FETCH_MAX_BYTES=3100000 -e ZOOKEEPER_IP=172.31.162.65 ches/kafka  
docker run --rm ches/kafka kafka-topics.sh --create --topic tags --replication-factor 1 --partitions 1 --zookeeper 172.31.162.65:2181  
docker run --rm ches/kafka kafka-topics.sh --list --zookeeper 172.31.162.65:2181  
docker run --rm --interactive ches/kafka kafka-console-producer.sh --topic tags --broker-list 172.31.162.65:9092  
docker run --rm ches/kafka kafka-console-consumer.sh --topic tags --from-beginning --zookeeper 172.31.162.65:2181  
docker run --rm --interactive ches/kafka kafka-consumer-groups.sh --new-consumer --describe --group group1 --bootstrap-server 172.31.162.65:9092  

# Check kafka instance details #
docker exec -it zookeeper bash  
bin/zkCli.sh -server 127.0.0.1:2181  
ls /brokers  
ls /brokers/topics  
ls /consumers   

# Consumer with .Net Core #
cd .\KafkaComparer.Consumer  
docker build -t kafkacomparerconsumer:latest .  
docker run -e CONSUMER_GROUP='tags-consumers' -e TOPIC_NAME='tags' -e KAFKA_URL='172.31.162.65:9092' --rm -it kafkacomparerconsumer  

# Consumer with Golang #
cd .\KafkaComparer.Consumer.Golang  
docker build -t kafkacomparerconsumer:go .  
docker run -it kafkacomparerconsumer:go  

# Producer with .Net Core #
cd .\KafkaComparer.Producer  
docker build -t kafkacomparerproducer:latest .  
docker run -e TOPIC_NAME='tags' -e KAFKA_URL='172.31.162.65:9092' --rm -it kafkacomparerproducer  

# Producer & Consumer with 'kafkacat' # (https://github.com/edenhill/kafkacat)
docker run --interactive --rm confluentinc/cp-kafkacat kafkacat -b 192.168.20.180:9092 -t tags -K: -P 
docker run --tty --interactive --rm confluentinc/cp-kafkacat kafkacat -b 192.168.20.180:9092 -L

# Topic with 3 partition #
docker run --rm ches/kafka kafka-topics.sh --create --topic tagsPart3 --replication-factor 1 --partitions 3 --zookeeper 172.31.162.65:2181  
