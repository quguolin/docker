package main

//https://medium.com/rahasak/kafka-producer-with-golang-fab7348a5f9a

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaConn1 = "127.0.0.1:9092"
	kafkaConn2 = "127.0.0.1:9093"
	kafkaConn3 = "127.0.0.1:9094"
	topic      = "test_kafka"
)

var brokerAddrs = []string{kafkaConn1, kafkaConn2, kafkaConn3}

func main() {
	// read command line input
	//reader := bufio.NewReader(os.Stdin)
	writer := newKafkaWriter(brokerAddrs, topic)
	defer writer.Close()
	for {
		//fmt.Print("Enter msg: ")
		time.Sleep(time.Second)
		msg := kafka.Message{
			Value: []byte(strconv.FormatInt(time.Now().Unix(), 10)),
		}
		fmt.Println("msg is ", string(msg.Value))
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func newKafkaWriter(kafkaURL []string, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kafkaURL,
		Topic:    topic,
		Balancer: &kafka.Hash{},
	})
}
