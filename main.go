package main

import (
	"lc-code-execution-service/config"
	"lc-code-execution-service/consumer"
	"lc-code-execution-service/scheduler"
	"lc-code-execution-service/types"
	"lc-code-execution-service/util"
	"os"
	"strconv"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection

func init() {
	if os.Getenv("LOAD_ENV") == "true" {
		config.LoadEnv()
	}
	config.LoadImages()
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

	var worker = os.Getenv("WORKERS")
	w, err := strconv.Atoi(worker)
	if err != nil {
		panic(err)
	}
	jobQueue := make(chan types.Job, w)

	go scheduler.ConsumerJobRequests(conn, jobQueue)

	scheduler.StartWorkerPool(jobQueue, w)

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
