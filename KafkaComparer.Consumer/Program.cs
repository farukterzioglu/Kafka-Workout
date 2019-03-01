using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Confluent.Kafka;

namespace KafkaComparer.Consumer
{
    class Program
    {
        static void Main(string[] args)
        {
            var consumergroup = Environment.GetEnvironmentVariable("CONSUMER_GROUP");
            if (string.IsNullOrEmpty(consumergroup)) consumergroup = Guid.NewGuid().ToString();
            Console.WriteLine($"Consumer group : {consumergroup}");

            var topicName = Environment.GetEnvironmentVariable("TOPIC_NAME");
            if (string.IsNullOrEmpty(topicName)) throw new ArgumentNullException(topicName);
            Console.WriteLine($"Topic name : {topicName}");

            var brokerList = Environment.GetEnvironmentVariable("KAFKA_URL");
            if (string.IsNullOrEmpty(brokerList)) throw new ArgumentNullException(brokerList);
            Console.WriteLine($"Broker list : {brokerList}");

            Offset offset;
            var offSetEnvVar = Environment.GetEnvironmentVariable("OFFSET");
            if(string.IsNullOrEmpty(offSetEnvVar) && long.TryParse(offSetEnvVar, out long offsetValue)){
                offset = new Offset(offsetValue);
            } else {
                offset = Offset.End;
            }

            // https://kafka.apache.org/0110/documentation.html#newconsumerconfigs
            var config = new ConsumerConfig
            {
                ClientId = "KafkaComparer.Consumer",
                GroupId = consumergroup,
                BootstrapServers = brokerList,
                HeartbeatIntervalMs = 3000,
                SessionTimeoutMs = 6000, // default 3 second 
                EnableAutoCommit = false, // false -> consumer.Commit();

                // Note: The AutoOffsetReset property determines the start offset in the event
                // there are not yet any committed offsets for the consumer group for the
                // topic/partitions of interest. By default, offsets are committed
                // automatically, so in this example, consumption will only start from the
                // eariest message in the topic 'my-topic' the first time you run the program.
                // The default is “latest,” which means that lacking a valid offset, the consumer will start reading from the newest records
                AutoOffsetReset = AutoOffsetResetType.Latest,

                MessageMaxBytes = 4000000,
            };

            using (var consumer = new Consumer<string, string>(config))
            {
                bool consuming = true;
                consumer.OnError += (_, e) => {
                    Console.WriteLine($"Error: {e.Reason}");
                    consuming = !e.IsFatal;
                };


                // You can use this event to perform actions such as retrieving offsets
                // from an external source / manually setting start offsets using
                // the Assign method.
                consumer.OnPartitionsAssigned += (obj, partitions) => {
                   Console.WriteLine($"Assigned partitions: [{string.Join(", ", partitions)}], member id: {consumer.MemberId}");
                   //var beginning = partitions.Select(p => new TopicPartitionOffset(p.Topic, p.Partition, offset)).ToList();
                   //Console.WriteLine($"Updated assignment: [{string.Join(", ", beginning)}]");
                   //consumer.Assign(beginning);
                   
                   //consumer.Assign(new List<TopicPartitionOffset>() { new TopicPartitionOffset(topicName, partition: 0, offset: Offset.Beginning) });
                };

                consumer.Subscribe(new string[] {topicName});
                while (consuming)
                {
                    try
                    {
                        // Get data 
                        ConsumeResult<string, string> msg = consumer.Consume();
                        // consumer.Commit(); // Commit before processing?

                        if (string.IsNullOrEmpty(msg.Value)) continue;

                        // Process the data
                        Console.WriteLine($"{msg.Value}{(!string.IsNullOrEmpty(msg.Key) ?  " - " + msg.Key : "") } | Topic: {msg.Topic} Partition: {msg.Partition} Offset: {msg.Offset}"); // TopicPartitionOffset: { msg.TopicPartitionOffset}

                        if (msg.Value.Contains("crash")) throw new Exception("crash");
                        if (msg.Value == "dont")  continue;
                            
                        // Commit 
                        consumer.Commit();
                    }
                    catch (ConsumeException e)
                    {
                        Console.WriteLine($"Error occured: {e.Error.Reason}");
                    }
                }
                consumer.Close();
            }
        }
    }
}
