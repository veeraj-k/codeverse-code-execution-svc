package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type SubmissionStatus struct {
	SubmissionId uuid.UUID `json:"job_id"`
	Status       string    `json:"status"`
}

type StatusProducer struct {
	channel *amqp091.Channel
	queue   string
}

var statusProducer *StatusProducer

func NewStatusProducer(channel *amqp091.Channel, queue string) *StatusProducer {
	statusProducer = &StatusProducer{
		channel: channel,
		queue:   queue,
	}
	return statusProducer
}

func (j *StatusProducer) ProduceJob(job *SubmissionStatus) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(job)
	if err != nil {
		panic(err)

	}
	err = j.channel.PublishWithContext(
		ctx,
		"",
		j.queue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/json",
			Body:        body,
		})

	if err != nil {
		panic(err)
	}
}
