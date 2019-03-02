using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Confluent.Kafka;

namespace KafkaComparer.Producer
{
    public class Program
    {
        static void Main(string[] args)
        {
            var topicName = Environment.GetEnvironmentVariable("TOPIC_NAME");
            var kafkaUrl = Environment.GetEnvironmentVariable("KAFKA_URL");

            Console.WriteLine($"Broker address : {kafkaUrl}");
            Console.WriteLine($"Topic name: {topicName}");
            
            // https://kafka.apache.org/0110/documentation.html#producerconfigs
            var config = new ProducerConfig(){
                MessageMaxBytes = 3000000,
                BootstrapServers = kafkaUrl,
                MessageTimeoutMs = 1000
            };

            using (var producer = new Producer<string, string>(config))
            {
                string text = null;

                while (text != "exit")
                {
                    Console.Write("Enter message: ");
                    text = Console.ReadLine();
                    if(string.IsNullOrEmpty(text)) continue;
                    
                    try
                    {
                        var values = text.Split('-');
                        var value = values[0];
                        string key = default;
                        if(values.Count() > 1){
                            key = values[1];
                        }

                        Message<string, string> message;
                        Task<DeliveryReport<string, string>> dr;

                        if(!string.IsNullOrEmpty(key)) {
                            message = new Message<string, string> { Key = key, Value = value };
                            dr = producer.ProduceAsync(topicName, message);
                        } else {
                            message = new Message<string, string> { Value = value };
                            dr = producer.ProduceAsync(topicName, message);
                        }
                        var result = dr.GetAwaiter().GetResult();
                        Console.WriteLine($"{DateTime.Now} Delivered to '{result.TopicPartitionOffset}'");
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
