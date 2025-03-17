package consumer

import (
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func SetupRabbitMQ() *amqp091.Connection {

	conn, err := amqp091.Dial(fmt.Sprintf("amqps://%s/%s", os.Getenv("RMQ_HOST"), os.Getenv("RMQ_VHOST")))

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"code_execution_job_queue",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	return conn

}
