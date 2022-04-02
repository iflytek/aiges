package catch

type EventType int

const (
	EVENT_STARTFAILED EventType = 1
	EVENT_DEADLOCK    EventType = 2
	EVENT_CRASH       EventType = 3
)

type EventHandle interface {
	Occur()
	collectStack(pid string) (cStack string, goStack []byte)
	reportEvent(cStack string, goStack []byte)
}
