package scheduler

import (
	"encoding/json"
	"fmt"
	"lc-code-execution-service/types"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func ConsumerJobRequests(conn *amqp091.Connection, jobQueue chan types.Job) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		"code_execution_job_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer started")
	for msg := range msgs {
		var job types.Job
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Printf("Error parsing job JSON: %v, body: %s", err, msg.Body)
			continue
		}

		// log.Printf("Received job: %+v", job)
		jobQueue <- job

		// job -> go svc -> directory(id) -> creaste the files of user code ->  mount container -> container run code -> output.json -> go svc reads output -> db meain store -> ws mein message
	}
}
