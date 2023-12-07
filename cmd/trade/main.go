package main

import (
	"fmt"
	"sync"
	"encoding/json"

	"github.com/caiohportella/imersao-fullstack-fullcycle/go/internal/market/dto"
	"github.com/caiohportella/imersao-fullstack-fullcycle/go/internal/market/entities"
	"github.com/caiohportella/imersao-fullstack-fullcycle/go/internal/market/infra/kafka"
	"github.com/caiohportella/imersao-fullstack-fullcycle/go/internal/market/transformer"
	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {
	ordersIn := make(chan *entities.Order)
	ordersOut := make(chan *entities.Order)
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "consumergroup",
		"auto.offset.reset": "earliest",
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewKafkaConsumer(configMap, []string{"new-orders"})

	go kafka.Consume(kafkaMsgChan)  //creates second thread

	//receives from kafka, sends to input channel, processes and sends to output channel, then publishes to kafka
	book := entities.NewBook(ordersIn, ordersOut, wg)

	go book.Trade()  //creates third thread

	go func() {
		for msg := range kafkaMsgChan {
			wg.Add(1)
			fmt.Println(string(msg.Value))
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)

			if err != nil {
				panic(err)
			}

			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJSON, err := json.Marshal(output)

		if err != nil {
			fmt.Println(err)
		}

		producer.Publish(outputJSON, []byte("processed-orders") , "output")
	}
}