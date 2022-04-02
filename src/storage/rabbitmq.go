package storage

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/xfyun/xsf/utils"
	"strings"
	"sync"
)

const rabbitSep = "://"

var (
	rabUrl   string
	rabQueue string
	rabLog   *utils.Logger
	con      *amqp.Connection
	channel  *amqp.Channel
	queue    amqp.Queue
	lock     sync.Mutex
)

// TODO pre-init multi rmq‘s channel
// TODO publish retry if fail
func RabInit(host string, usr string, pass string, name string, logger *utils.Logger) (err error) {
	if len(host) > 0 {
		rabLog = logger
		rabQueue = name
		rabUrl = host
		if index := strings.Index(host, rabbitSep); index != -1 {
			rabUrl = host[:index+len(rabbitSep)] + usr + ":" + pass + "@" + host[index+len(rabbitSep):]
		}
		con, err = amqp.Dial(rabUrl)
		if err != nil {
			rabLog.Errorw("rabbitmq init Dial fail", "err", err.Error(), "url", rabUrl)
			return
		}
		channel, err = con.Channel()
		if err != nil {
			rabLog.Errorw("rabbitmq init Channel fail", "err", err.Error(), "url", rabUrl)
			return
		}

		if queue, err = channel.QueueDeclare(
			rabQueue,
			true,
			false,
			false,
			false,
			nil); err != nil {
			rabLog.Errorw("rabbitmq init QueueDeclare fail", "err", err.Error(), "url", rabUrl)
		}
	}
	return
}

func RabPublish(body []byte, retry int) (err error) {
	for len(rabUrl) > 0 && retry >= 0 {
		lock.Lock()
		if err = channel.Publish(
			"",
			rabQueue,
			false,
			false,
			amqp.Publishing{
				Body: body,
			}); err != nil {
			rabRecon() // reconnect
			retry--
			rabLog.Errorw("rabbitmq publish fail", "err", err.Error(), "url", rabUrl)
		} else {
			retry = -1 // 成功则无需重试
		}
		lock.Unlock()
	}
	return
}

func RabFini() {
	if len(rabUrl) > 0 {
		channel.Close()
		con.Close()
		fmt.Println("aiService.Finit: fini rabbit success!")
	}
}

func rabRecon() {
	var err error
	if con, err = amqp.Dial(rabUrl); err != nil {
		rabLog.Errorw("rabbitmq recon Dial fail", "err", err.Error(), "url", rabUrl)
		return
	}

	if channel, err = con.Channel(); err != nil {
		rabLog.Errorw("rabbitmq recon Channel fail", "err", err.Error(), "url", rabUrl)
		return
	}
	return
}
