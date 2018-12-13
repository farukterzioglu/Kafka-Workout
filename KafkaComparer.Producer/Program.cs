using System;
using System.Collections.Generic;
using System.Text;
using Confluent.Kafka;
using Confluent.Kafka.Serialization;

namespace KafkaComparer.Producer
{
  public class Program
  {
    static void Main(string[] args)
    {
      var topicName = Environment.GetEnvironmentVariable("TOPIC_NAME");
      var kafkaUrl = Environment.GetEnvironmentVariable("KAFKA_URL");

      var config = new Dictionary<string, object>
      {
				{ "bootstrap.servers", kafkaUrl }
			};

      using (var producer = new Producer<Null, string>(config, null, new StringSerializer(Encoding.UTF8)))
      {
        string text = null;

        while (text != "exit")
        {
          text = Console.ReadLine();
          producer.ProduceAsync(topicName, null, text);
        }

        producer.Flush(100);
      }
    }
  }
}
