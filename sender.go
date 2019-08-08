package amqpsender

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

func NewSender() *Sender {
	return &Sender{}
}

func (sndr *Sender) Start(connectionString string, amqpConfig amqp.Config, noWait bool) error {
	var err error
	sndr.amqpConnection, err = amqp.DialConfig(connectionString, amqpConfig)
	if err != nil {
		return errors.Wrap(err, "failed to connect to RabbitMQ")
	}
	sndr.amqpChannel, err = sndr.amqpConnection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	err = sndr.amqpChannel.Confirm(noWait)
	if err != nil {
		return fmt.Errorf("failed to set channel to noWait mode: %v", err)
	}
	return nil
}

func (sndr *Sender) GetConnectionNotifyClose() chan *amqp.Error {
	if sndr.amqpConnectionNotifyClose == nil {
		sndr.amqpConnectionNotifyClose = sndr.amqpConnection.NotifyClose(make(chan *amqp.Error))
	}
	return sndr.amqpConnectionNotifyClose
}

func (sndr *Sender) GetConnectionNotifyCancel() chan amqp.Blocking {
	if sndr.amqpConnectionNotifyBlocked == nil {
		sndr.amqpConnectionNotifyBlocked = sndr.amqpConnection.NotifyBlocked(make(chan amqp.Blocking))
	}
	return sndr.amqpConnectionNotifyBlocked
}

func (sndr *Sender) GetChannelNotifyCancel() chan string {
	if sndr.amqpChannelNotifyCancel == nil {
		sndr.amqpChannelNotifyCancel = sndr.amqpChannel.NotifyCancel(make(chan string))
	}
	return sndr.amqpChannelNotifyCancel
}

func (sndr *Sender) GetChannelNotifyClose() chan *amqp.Error {
	if sndr.amqpChannelNotifyClose == nil {
		sndr.amqpChannelNotifyClose = sndr.amqpChannel.NotifyClose(make(chan *amqp.Error))
	}
	return sndr.amqpChannelNotifyClose
}

func (sndr *Sender) GetChannelNotifyConfirm() (chan uint64, chan uint64) {
	if sndr.amqpChannelNotifyConfirmAck == nil || sndr.amqpChannelNotifyConfirmNack == nil {
		sndr.amqpChannelNotifyConfirmAck, sndr.amqpChannelNotifyConfirmNack = sndr.amqpChannel.NotifyConfirm(make(chan uint64), make(chan uint64))
	}
	return sndr.amqpChannelNotifyConfirmAck, sndr.amqpChannelNotifyConfirmNack
}

func (sndr *Sender) GetChannelNotifyFlow() chan bool {
	if sndr.amqpChannelNotifyFlow == nil {
		sndr.amqpChannelNotifyFlow = sndr.amqpChannel.NotifyFlow(make(chan bool))
	}
	return sndr.amqpChannelNotifyFlow
}

func (sndr *Sender) GetChannelNotifyPublish() chan amqp.Confirmation {
	if sndr.amqpChannelNotifyPublish == nil {
		sndr.amqpChannelNotifyPublish = sndr.amqpChannel.NotifyPublish(make(chan amqp.Confirmation))
	}
	return sndr.amqpChannelNotifyPublish
}

func (sndr *Sender) GetChannelNotifyReturn() chan amqp.Return {
	if sndr.amqpChannelNotifyReturn == nil {
		sndr.amqpChannelNotifyReturn = sndr.amqpChannel.NotifyReturn(make(chan amqp.Return))
	}
	return sndr.amqpChannelNotifyReturn
}

func (sndr *Sender) SendTask(exchangeName string, routingKey string, mandatory bool, immediate bool, amqpPub amqp.Publishing) error {
	amqpPub.Timestamp = time.Now()
	return sndr.amqpChannel.Publish(exchangeName, routingKey, mandatory, immediate, amqpPub)
}

func (sndr *Sender) Close() error {
	if err := sndr.amqpConnection.Close(); err != nil {
		return errors.Wrap(err, "error with graceful close amqp connection")
	}
	_ = sndr.amqpChannel.Close()

	// clean channels
	close(sndr.amqpConnectionNotifyClose)
	close(sndr.amqpConnectionNotifyBlocked)
	close(sndr.amqpChannelNotifyCancel)
	close(sndr.amqpChannelNotifyClose)
	close(sndr.amqpChannelNotifyConfirmAck)
	close(sndr.amqpChannelNotifyConfirmNack)
	close(sndr.amqpChannelNotifyFlow)
	close(sndr.amqpChannelNotifyPublish)
	close(sndr.amqpChannelNotifyReturn)

	return nil
}

/*
todo Channel.NotifyReturn
Since publishings are asynchronous, any undeliverable message will get returned
by the server.  Add a listener with Channel.NotifyReturn to handle any
undeliverable message when calling publish with either the mandatory or
immediate parameters as true.
*/
