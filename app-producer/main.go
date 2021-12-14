package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/segmentio/kafka-go"
)

type Entry struct {
	Id      int
	Created time.Time
	Message string
}

const (
	topic  = "entry-added"
	broker = "0.0.0.0:9092"
)

func send(ctx context.Context, dataChannel chan Entry) {
	kafkaWriter := &kafka.Writer{
		Addr:  kafka.TCP(broker),
		Topic: topic,
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done")
			return
		case entry := <-dataChannel:
			fmt.Printf("Sending data %+v", entry)

			err := kafkaWriter.WriteMessages(ctx, kafka.Message{
				Key:   []byte(strconv.Itoa(entry.Id)),
				Value: []byte(fmt.Sprintf("%+v", entry)),
			})

			if err != nil {
				fmt.Println("Message couldn't have been written" + err.Error())
			} else {
				fmt.Println("Message sent")
			}
		}
	}
}

func generateData(dataChannel chan Entry, quantity int) {
	for i := 0; i < quantity; i++ {
		entryStruct := Entry{}

		err := faker.FakeData(&entryStruct)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Generated data %+v", entryStruct)
		dataChannel <- entryStruct
	}
}

func main() {
	dataChannel := make(chan Entry, 100)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go generateData(dataChannel, 100)
	send(ctx, dataChannel)
}
