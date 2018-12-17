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
            if (string.IsNullOrEmpty(consumergroup)) 
                consumergroup = Guid.NewGuid().ToString();

            var topicName = Environment.GetEnvironmentVariable("TOPIC_NAME");
            var brokerList = Environment.GetEnvironmentVariable("KAFKA_URL");

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
                GroupId = consumergroup,
                BootstrapServers = brokerList,
                // Note: The AutoOffsetReset property determines the start offset in the event
                // there are not yet any committed offsets for the consumer group for the
                // topic/partitions of interest. By default, offsets are committed
                // automatically, so in this example, consumption will only start from the
                // eariest message in the topic 'my-topic' the first time you run the program.
                //AutoOffsetReset = AutoOffsetResetType.Earliest,
                //MessageMaxBytes = 3000000
            };

            using (var consumer = new Consumer<Null, string>(config))
            {
                consumer.OnError += (_, e) => Console.WriteLine($"Error: {e.Reason}");

                bool consuming = true;
                // The client will automatically recover from non-fatal errors. You typically
                // don't need to take any action unless an error is marked as fatal.
                consumer.OnError += (_, e) => consuming = !e.IsFatal;

                consumer.Assign(new List<TopicPartitionOffset>() { new TopicPartitionOffset(topicName, partition: 0, offset: Offset.Beginning) });

                // You can use this event to perform actions such as retrieving offsets
                // from an external source / manually setting start offsets using
                // the Assign method.
                //consumer.OnPartitionsAssigned += (obj, partitions) => {
                //    Console.WriteLine($"Assigned partitions: [{string.Join(", ", partitions)}], member id: {consumer.MemberId}");
                //    var beginning = partitions.Select(p => new TopicPartitionOffset(p.Topic, p.Partition, offset)).ToList();
                //    Console.WriteLine($"Updated assignment: [{string.Join(", ", beginning)}]");
                //    consumer.Assign(beginning);
                //};

                consumer.Subscribe(new string[] {topicName});
                while (consuming)
                {
                    try
                    {
                        var msg = consumer.Consume();
                        Console.WriteLine($"{msg.Value}\nTopic: {msg.Topic} Partition: {msg.Partition} Offset: {msg.Offset} TopicPartitionOffset: { msg.TopicPartitionOffset} \n");
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
