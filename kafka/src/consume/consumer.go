package main

//https://leel0330.github.io/golang/%E5%9C%A8go%E4%B8%AD%E4%BD%BF%E7%94%A8kafka/

import (
	"context"
	"flag"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaConn1 = "127.0.0.1:9092"
	kafkaConn2 = "127.0.0.1:9093"
	kafkaConn3 = "127.0.0.1:9094"
)

var (
	topic = flag.String("t", "test_kafka", "kafka topic")
	group = flag.String("g", "test-group", "kafka consumer group")
)

func main() {
	flag.Parse()
	config := kafka.ReaderConfig{
		Brokers:  []string{kafkaConn1, kafkaConn2, kafkaConn3},
		Topic:    *topic,
		MinBytes: 1e3,
		MaxBytes: 1e6,
		GroupID:  *group,
	}
	reader := kafka.NewReader(config)
	ctx := context.Background()
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("fail to get msg:%v", err)
			continue
		}
		log.Printf("msg content:topic=%v,partition=%v,offset=%v,content=%v",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("fail to commit msg:%v", err)
		}
	}
}
