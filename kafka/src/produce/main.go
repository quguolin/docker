package main

//https://medium.com/rahasak/kafka-producer-with-golang-fab7348a5f9a

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"
	"github.com/segmentio/kafka-go"
)

const (
	kafkaConn1 = "127.0.0.1:9092"
	kafkaConn2 = "127.0.0.1:9093"
	kafkaConn3 = "127.0.0.1:9094"
	topic      = "test2"
)

var brokerAddrs = []string{kafkaConn1, kafkaConn2, kafkaConn3}

func createTopic(topic string) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	admin, err := sarama.NewClusterAdmin(brokerAddrs, config)
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
	}
	defer func() { _ = admin.Close() }()
	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 3,
	}, false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("create topic(%s) success", topic)
}

func isTopicExists(topic string) bool {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	//kafka end point
	brokers := []string{kafkaConn1, kafkaConn2, kafkaConn3}

	//get broker
	cluster, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := cluster.Close(); err != nil {
			panic(err)
		}
	}()
	//get all topic from cluster
	topics, err := cluster.Topics()
	if err != nil {
		panic(err)
	}
	fmt.Printf("all topic %v\n", topics)
	for _, v := range topics {
		if v == topic {
			return true
		}
	}
	return false
}

func main() {
	if !isTopicExists(topic) {
		createTopic(topic)
	} else {
		fmt.Printf("top(%s) existed\n", topic)
	}
	// read command line input
	reader := bufio.NewReader(os.Stdin)
	writer := newKafkaWriter(brokerAddrs, topic)
	defer writer.Close()
	for {
		fmt.Print("Enter msg: ")
		msgStr, _ := reader.ReadString('\n')

		msg := kafka.Message{
			Value: []byte(msgStr),
		}
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

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer(brokerAddrs, config)

	return prd, err
}

func publish(message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
}
