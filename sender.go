package amqpsender

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// NewSender returns new Sender instance
func NewSender() *Sender {
	return &Sender{}
}

// Start is for initialize connection to AMQP and create new channel
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

// GetConnectionNotifyClose returns NotifyClose channel for AMQP Connection
func (sndr *Sender) GetConnectionNotifyClose() chan *amqp.Error {
	if sndr.amqpConnectionNotifyClose == nil {
		sndr.amqpConnectionNotifyClose = sndr.amqpConnection.NotifyClose(make(chan *amqp.Error))
	}
	return sndr.amqpConnectionNotifyClose
}

// GetConnectionNotifyBlocked returns NotifyBlocked channel for AMQP Connection
func (sndr *Sender) GetConnectionNotifyBlocked() chan amqp.Blocking {
	if sndr.amqpConnectionNotifyBlocked == nil {
		sndr.amqpConnectionNotifyBlocked = sndr.amqpConnection.NotifyBlocked(make(chan amqp.Blocking))
	}
	return sndr.amqpConnectionNotifyBlocked
}

// GetChannelNotifyCancel returns NotifyCancel channel for AMQP Channel
func (sndr *Sender) GetChannelNotifyCancel() chan string {
	if sndr.amqpChannelNotifyCancel == nil {
		sndr.amqpChannelNotifyCancel = sndr.amqpChannel.NotifyCancel(make(chan string))
	}
	return sndr.amqpChannelNotifyCancel
}

// GetChannelNotifyClose returns NotifyClose channel for AMQP Channel
func (sndr *Sender) GetChannelNotifyClose() chan *amqp.Error {
	if sndr.amqpChannelNotifyClose == nil {
		sndr.amqpChannelNotifyClose = sndr.amqpChannel.NotifyClose(make(chan *amqp.Error))
	}
	return sndr.amqpChannelNotifyClose
}

// GetChannelNotifyConfirm returns NotifyConfirm channels for Ack and Nack for AMQP Channel
func (sndr *Sender) GetChannelNotifyConfirm() (chan uint64, chan uint64) {
	if sndr.amqpChannelNotifyConfirmAck == nil || sndr.amqpChannelNotifyConfirmNack == nil {
		sndr.amqpChannelNotifyConfirmAck, sndr.amqpChannelNotifyConfirmNack = sndr.amqpChannel.NotifyConfirm(make(chan uint64), make(chan uint64))
	}
	return sndr.amqpChannelNotifyConfirmAck, sndr.amqpChannelNotifyConfirmNack
}

// GetChannelNotifyFlow returns NotifyFlow channel for AMQP Channel
func (sndr *Sender) GetChannelNotifyFlow() chan bool {
	if sndr.amqpChannelNotifyFlow == nil {
		sndr.amqpChannelNotifyFlow = sndr.amqpChannel.NotifyFlow(make(chan bool))
	}
	return sndr.amqpChannelNotifyFlow
}

// GetChannelNotifyPublish returns NotifyPublish channel for AMQP Channel
func (sndr *Sender) GetChannelNotifyPublish() chan amqp.Confirmation {
	if sndr.amqpChannelNotifyPublish == nil {
		sndr.amqpChannelNotifyPublish = sndr.amqpChannel.NotifyPublish(make(chan amqp.Confirmation))
	}
	return sndr.amqpChannelNotifyPublish
}

// GetChannelNotifyReturn returns NotifyReturn channel for AMQP Channel
func (sndr *Sender) GetChannelNotifyReturn() chan amqp.Return {
	if sndr.amqpChannelNotifyReturn == nil {
		sndr.amqpChannelNotifyReturn = sndr.amqpChannel.NotifyReturn(make(chan amqp.Return))
	}
	return sndr.amqpChannelNotifyReturn
}

// SendTask sends job task to AMQP
func (sndr *Sender) SendTask(exchangeName string, routingKey string, mandatory bool, immediate bool, amqpPub amqp.Publishing) error {
	return sndr.amqpChannel.Publish(exchangeName, routingKey, mandatory, immediate, amqpPub)
}

// SendTaskWithCurrentTimestamp sends job task to AMQP
func (sndr *Sender) SendTaskWithCurrentTimestamp(exchangeName string, routingKey string, mandatory bool, immediate bool, amqpPub amqp.Publishing) error {
	amqpPub.Timestamp = time.Now()
	return sndr.amqpChannel.Publish(exchangeName, routingKey, mandatory, immediate, amqpPub)
}

// SendTaskSimple sends job task to AMQP
func (sndr *Sender) SendTaskSimple(exchangeName string, routingKey string, amqpPub amqp.Publishing) error {
	return sndr.amqpChannel.Publish(exchangeName, routingKey, false, false, amqpPub)
}

// SendTaskSimpleWithCurrentTimestamp sends job task to AMQP
func (sndr *Sender) SendTaskSimpleWithCurrentTimestamp(exchangeName string, routingKey string, amqpPub amqp.Publishing) error {
	amqpPub.Timestamp = time.Now()
	return sndr.amqpChannel.Publish(exchangeName, routingKey, false, false, amqpPub)
}

// Close is for gracefull closing amqp connection, amqp channel and all notify channels
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
