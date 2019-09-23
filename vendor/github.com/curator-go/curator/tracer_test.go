package curator

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTracerDriver(t *testing.T) {
	var logs []string

	d := newDefaultTracerDriver()
	d.logger = func(format string, args ...interface{}) {
		logs = append(logs, fmt.Sprintf(format, args...))
	}

	d.AddTime("time", time.Second*15)
	d.AddCount("count", 100)
	d.AddCount("total", 10)
	d.AddCount("count", 20)

	assert.Equal(t, []string{
		"Trace time: 15s",
		"Counter count: 0 + 100",
		"Counter total: 0 + 10",
		"Counter count: 100 + 20",
	}, logs)
}

func TestTimeTrace(t *testing.T) {
	d := &mockTracerDriver{}

	d.On("AddTime", "test", time.Second*5).Return().Once()

	tracer := newTimeTracer("test", d)
	tracer.CommitAt(tracer.startTime.Add(time.Second * 5))

	d.AssertExpectations(t)
}
