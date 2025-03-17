package main

import (
	"lc-code-execution-service/config"
	"lc-code-execution-service/consumer"
	"lc-code-execution-service/scheduler"
	"lc-code-execution-service/types"
	"lc-code-execution-service/util"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection

func init() {
	config.LoadEnv()
	util.NewHttpClient()
}

func main() {

	conn = consumer.SetupRabbitMQ()
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	scheduler.NewStatusProducer(channel, "code_execution_status_queue")

	jobQueue := make(chan types.Job, 1)

	go scheduler.ConsumerJobRequests(conn, jobQueue)

	scheduler.StartWorkerPool(jobQueue, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
