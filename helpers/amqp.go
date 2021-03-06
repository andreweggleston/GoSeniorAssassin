package helpers

import (
	"github.com/sirupsen/logrus"
	"github.com/andreweggleston/GoSeniorAssassin/config"
	"github.com/streadway/amqp"
)

var AMQPChannel *amqp.Channel
var AMQPConn *amqp.Connection

func ConnectAMQP() {
	var err error

	AMQPConn, err = amqp.Dial(config.Constants.RabbitMQURL)
	if err != nil {
		logrus.Fatal(err)
	}

	AMQPChannel, err = AMQPConn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Connected to RabbitMQ on ", config.Constants.RabbitMQURL)
}
