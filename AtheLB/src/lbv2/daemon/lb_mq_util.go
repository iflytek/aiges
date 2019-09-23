package daemon

/*
Message Queue
*/
type MqManager interface {
	Consume(topic, group string) (string, *MTRMessage, error)
	Init([]string) error
}

type RmqManager struct {
	rmqAdapter
}

func (r *RmqManager) Init(addr []string) error {
	return r.rmqAdapter.Init(addr)
}
func (r *RmqManager) Consume(topic, group string) (RemoteRmqAddr string, ConsumeR *MTRMessage, ConsumeE error) {
	return r.rmqAdapter.Consume(topic, group)
}
