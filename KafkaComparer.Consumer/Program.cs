using System;
using System.Collections.Generic;
using System.Text;
using System.Threading.Tasks;
using Confluent.Kafka;
using Confluent.Kafka.Serialization;

namespace KafkaComparer.Consumer
{
  class Program
  {
    static void Main(string[] args)
    {
      var consumergroup = Environment.GetEnvironmentVariable("CONSUMER_GROUP");
      var topicName = Environment.GetEnvironmentVariable("TOPIC_NAME");
      var kafkaUrl = Environment.GetEnvironmentVariable("KAFKA_URL");

      var config = new Dictionary<string, object>
            {
                { "group.id", consumergroup },
                { "bootstrap.servers", kafkaUrl },
                { "enable.auto.commit", "false"}
            };

      using (var consumer = new Consumer<Null, string>(config, null, new StringDeserializer(Encoding.UTF8)))
      {
        consumer.Subscribe(new string[] { topicName });

        consumer.OnMessage += (_, msg) =>
        {
          Console.WriteLine($"Topic: {msg.Topic} Partition: {msg.Partition} Offset: {msg.Offset} {msg.Value}");
          consumer.CommitAsync(msg);
        };

        while (true)
        {
          consumer.Poll(100);
        }
      }
    }
  }
}
