package heartbeatmanagement

import (
	"log"
	"strconv"
	"time"

	"github.com/Contargo/chamqp"
	"github.com/coreos/go-systemd/daemon"
)

type HeartbeatSender[T any] struct {
	chamqpChannel            *chamqp.Channel
	heartbeatContentGatherer HeartbeatContent[T]
	exchangeName             string
	routingKey               string
}

type HeartbeatContent[T any] interface {
	GetHeartbeatContent() T
}

func NewHeartbeatSender[T any](chamqpChannel *chamqp.Channel, heartbeatContentGatherer HeartbeatContent[T], exchangeName, routingKey string) *HeartbeatSender[T] {
	return &HeartbeatSender[T]{
		chamqpChannel,
		heartbeatContentGatherer,
		exchangeName,
		routingKey,
	}
}

func (h *HeartbeatSender[T]) StartSending() {
	h.StartSendingWithParams(false, false, false, false)
}

func (h *HeartbeatSender[T]) StartSendingWithParams(durable, autodelete, internal, nowait bool) {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		interval = 50 * time.Second
	}

	h.chamqpChannel.ExchangeDeclare(h.exchangeName, "topic", durable, autodelete, internal, nowait, nil, nil)
	for {
		heartbeatObj := h.heartbeatContentGatherer.GetHeartbeatContent()
		err := h.chamqpChannel.PublishJSONWithProperties(h.exchangeName, h.routingKey, false, false, heartbeatObj, chamqp.Properties{
			Expiration: strconv.Itoa(2 * 60 * 1000),
		})
		if err == nil {
			notificationSupported, err := daemon.SdNotify(false, daemon.SdNotifyWatchdog)
			if !notificationSupported {
				log.Println("Systemd watchdog notification not supported")
			} else {
				if err != nil {
					log.Println("err during sending out systemd notificaiton", err)
				}
			}
		}

		time.Sleep(interval)
	}
}
