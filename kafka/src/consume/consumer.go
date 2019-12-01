package main

//https://leel0330.github.io/golang/%E5%9C%A8go%E4%B8%AD%E4%BD%BF%E7%94%A8kafka/

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

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

const (
	count = 10
)

type Service struct {
	wg          *sync.WaitGroup
	cs          []chan int
	close       bool
	kafkaReader *kafka.Reader
	mapAid      map[int]int
}

func New() *Service {
	kafkaConfig := kafka.ReaderConfig{
		Brokers:  []string{kafkaConn1, kafkaConn2, kafkaConn3},
		Topic:    *topic,
		MinBytes: 1e3,
		MaxBytes: 1e6,
		GroupID:  *group,
	}
	s := &Service{
		wg:          &sync.WaitGroup{},
		kafkaReader: kafka.NewReader(kafkaConfig),
		mapAid:      make(map[int]int),
	}
	for i := 0; i < count; i++ {
		s.cs = append(s.cs, make(chan int, 1024))
	}
	return s
}

func (s *Service) ConsumeKafka() {
	defer s.wg.Done()
	ctx := context.Background()
	for {
		msg, err := s.kafkaReader.FetchMessage(ctx)
		if err != nil {
			log.Printf("fail to get msg:%v", err)
			continue
		}
		//log.Printf("msg content:topic=%v,partition=%v,offset=%v,content=%v",
		//	msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		err = s.kafkaReader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("fail to commit msg:%v", err)
		}
		v, _ := strconv.Atoi(string(msg.Value))
		s.cs[v%count] <- v
		fmt.Println("ConsumeKafka ", int64(v))
	}
}

func (s *Service) ConsumeCs(i int) {
	defer s.wg.Done()
	c := s.cs[i]
	for {
		fmt.Println("ConsumeCs")
		msg, ok := <-c
		fmt.Println("msg is ", msg)

		if !ok || s.close {
			return
		}
		s.mapAid[msg%10]++
		fmt.Println("consume ", s.mapAid)
	}
}

func (s *Service) Close() {
	s.close = true
	for i := 0; i < count; i++ {
		close(s.cs[i])
	}
}

func main() {
	flag.Parse()
	s := New()
	s.wg.Add(1)
	go s.ConsumeKafka()
	for i := 0; i < count; i++ {
		s.wg.Add(1)
		go s.ConsumeCs(i)
	}
	s.wg.Wait()
	signalHandler()
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
