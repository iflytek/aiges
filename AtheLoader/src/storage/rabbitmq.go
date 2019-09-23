package storage

import (
	"git.xfyun.cn/AIaaS/xsf-external/utils"
	"github.com/amqp"
	"strings"
)

const rabbitSep = "://"
var (
	rabUrl string
	rabQueue string
	rabLog *utils.Logger
	con *amqp.Connection
	channel *amqp.Channel
	queue amqp.Queue
)

func RabInit(host string, usr string, pass string, name string, logger* utils.Logger) (err error){
	if len(host) > 0 {
		rabLog = logger
		rabQueue = name
		rabUrl = host
		if index := strings.Index(host, rabbitSep); index != -1 {
			rabUrl = host[:index+len(rabbitSep)] + usr +":"+ pass + "@" + host[index+len(rabbitSep):]
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

func RabPuslish(body []byte) (err error) {
	// TODO 确认是否channel, connection or socket被关闭导致报错,重连场景
	// TODO channel publish锁竞争,多实例channel
	if len(rabUrl) > 0 {
		if err = channel.Publish(
			"",
			rabQueue,
			false,
			false,
			amqp.Publishing{
				Body: body,
			}); err != nil {
			rabLog.Errorw("rabbitmq publish fail", "err", err.Error(), "url", rabUrl)
		}
	}
	return
}

func RabFini() {
	if len(rabUrl) > 0 {
		channel.Close()
		con.Close()
	}
}