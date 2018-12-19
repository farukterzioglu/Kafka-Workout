using System;
using System.Collections.Generic;
using System.Text;
using Confluent.Kafka;

namespace KafkaComparer.Producer
{
    public class Program
    {
        static void Main(string[] args)
        {
            var topicName = Environment.GetEnvironmentVariable("TOPIC_NAME");
            var kafkaUrl = Environment.GetEnvironmentVariable("KAFKA_URL");

            // https://kafka.apache.org/0110/documentation.html#producerconfigs
            var config = new ProducerConfig(){
                MessageMaxBytes = 3000000,
                BootstrapServers = kafkaUrl
            };

            using (var producer = new Producer<Null, string>(config))
            {
                string text = null;

                while (text != "exit")
                {
                    text = Console.ReadLine();
                    try
                    {
                        var dr = producer.ProduceAsync(topicName, new Message<Null, string> { Value = text }).GetAwaiter().GetResult();
                        Console.WriteLine($"{DateTime.Now} Delivered to '{dr.TopicPartitionOffset}'");
                    }
                    catch (KafkaException e)
                    {
                        Console.WriteLine($"Delivery failed: {e.Error.Reason}");
                    }
                }
                producer.Flush(TimeSpan.FromSeconds(10));
            }
        }
    }
}
