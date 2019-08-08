package amqpsender

import "github.com/streadway/amqp"

type Sender struct {
	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	// notification channels
	amqpConnectionNotifyClose    chan *amqp.Error
	amqpConnectionNotifyBlocked  chan amqp.Blocking
	amqpChannelNotifyCancel      chan string
	amqpChannelNotifyClose       chan *amqp.Error
	amqpChannelNotifyConfirmAck  chan uint64
	amqpChannelNotifyConfirmNack chan uint64
	amqpChannelNotifyFlow        chan bool
	amqpChannelNotifyPublish     chan amqp.Confirmation
	amqpChannelNotifyReturn      chan amqp.Return
}
