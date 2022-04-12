package profiler

import (
	"container/heap"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/pyroscope-io/dotnetdiag/nettrace"
)

// SampleProfiler processes event stream from Microsoft-DotNETCore-SampleProfiler
// provider and calculates time for every call stack.
type SampleProfiler struct {
	trace *nettrace.Trace
	sym   *symbols

	md      map[int32]*nettrace.Metadata
	stacks  map[int32][]uint64
	threads map[int64]*thread

	events  events
	samples []sample

	managedOnly bool
}

type Option func(*SampleProfiler)

// WithManagedCodeOnly prescribes SampleProfiler to ignore the time that
// was spent in native (unmanaged) code.
func WithManagedCodeOnly() Option {
	return func(p *SampleProfiler) {
		p.managedOnly = true
	}
}

type sample struct {
	stack []uint64
	value int64
}

type event struct {
	typ          clrThreadSampleType
	threadID     int64
	stackID      int32
	timestamp    int64
	relativeTime int64
}

type events []*event

func (e events) Len() int { return len(e) }

func (e events) Less(i, j int) bool { return e[i].timestamp < e[j].timestamp }

func (e events) Swap(i, j int) { e[i], e[j] = e[j], e[i] }

func (e *events) Push(x interface{}) { *e = append(*e, x.(*event)) }

func (e *events) Pop() interface{} {
	old := *e
	n := len(old)
	x := old[n-1]
	*e = old[0 : n-1]
	return x
}

// clrThreadSampleTraceData describes ThreadSample event payload for
// Microsoft-DotNETCore-SampleProfiler Event ID 0.
type clrThreadSampleTraceData struct {
	Type clrThreadSampleType
}

type clrThreadSampleType int32

const (
	_ clrThreadSampleType = iota - 1

	sampleTypeError
	sampleTypeExternal
	sampleTypeManaged
)

func NewSampleProfiler(trace *nettrace.Trace, options ...Option) *SampleProfiler {
	p := &SampleProfiler{
		trace:   trace,
		sym:     newSymbols(),
		md:      make(map[int32]*nettrace.Metadata),
		threads: make(map[int64]*thread),
		stacks:  make(map[int32][]uint64),
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (s *SampleProfiler) Samples() map[string]time.Duration {
	samples := make(map[string]time.Duration)
	for _, x := range s.samples {
		name := make([]string, len(x.stack))
		for i := range x.stack {
			name[i] = s.sym.resolve(x.stack[i])
		}
		samples[strings.Join(name, ";")] += time.Duration(x.value)
	}
	return samples
}

func (s *SampleProfiler) EventHandler(e *nettrace.Blob) error {
	md, ok := s.md[e.Header.MetadataID]
	if !ok {
		return fmt.Errorf("metadata not found")
	}

	switch {
	case md.Header.ProviderName == "Microsoft-DotNETCore-SampleProfiler" && md.Header.EventID == 0:
		return s.addSample(e)

	case md.Header.ProviderName == "Microsoft-Windows-DotNETRuntimeRundown":
		switch {
		case md.Header.EventID == 144:
			return s.sym.addMethod(e)

		case md.Header.EventID == 152:
			return s.sym.addModule(e)
		}
	}

	return nil
}

func (s *SampleProfiler) MetadataHandler(md *nettrace.Metadata) error {
	s.md[md.Header.MetaDataID] = md
	return nil
}

func (s *SampleProfiler) StackBlockHandler(sb *nettrace.StackBlock) error {
	for _, stack := range sb.Stacks {
		s.stacks[stack.ID] = stack.InstructionPointers(s.trace.PointerSize)
	}
	return nil
}

func (s *SampleProfiler) SequencePointBlockHandler(*nettrace.SequencePointBlock) error {
	for s.events.Len() != 0 {
		x := heap.Pop(&s.events).(*event)
		s.thread(x.threadID).addSample(x.typ, x.relativeTime, x.stackID)
	}
	for _, t := range s.threads {
		for stackID, value := range t.samples {
			s.samples = append(s.samples, sample{
				stack: s.stacks[stackID],
				value: value,
			})
		}
		t.samples = make(map[int32]int64)
	}
	s.stacks = make(map[int32][]uint64)
	return nil
}

// https://github.com/microsoft/perfview/blob/8a34d2d64bc958902b2fa8ea5799437df57d8de2/src/TraceEvent/TraceEvent.cs#L440-L460
func (s *SampleProfiler) addSample(e *nettrace.Blob) error {
	var d clrThreadSampleTraceData
	if err := binary.Read(e.Payload, binary.LittleEndian, &d); err != nil {
		return err
	}
	rel := e.Header.TimeStamp - s.trace.SyncTimeQPC
	if rel < 0 {
		return nil
	}
	heap.Push(&s.events, &event{
		typ:          d.Type,
		threadID:     e.Header.ThreadID,
		stackID:      e.Header.StackID,
		timestamp:    e.Header.TimeStamp,
		relativeTime: rel * (int64(time.Second) / s.trace.QPCFrequency),
	})
	return nil
}

func (s *SampleProfiler) thread(tid int64) *thread {
	t, ok := s.threads[tid]
	if ok {
		return t
	}
	t = &thread{
		samples:     make(map[int32]int64),
		managedOnly: s.managedOnly,
	}
	s.threads[tid] = t
	return t
}
