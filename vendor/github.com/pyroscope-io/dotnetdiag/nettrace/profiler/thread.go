package profiler

type thread struct {
	lastExternalTime int64
	lastManagedTime  int64
	samples          map[int32]int64
	managedOnly      bool
}

type threadState int

const (
	_ threadState = iota - 1

	uninitialized
	managed
	external
)

func (t *thread) state() threadState {
	switch {
	case t.lastExternalTime < 0:
		return managed
	case t.lastExternalTime > 0:
		return external
	default:
		return uninitialized
	}
}

func (t *thread) addSample(sampleType clrThreadSampleType, relativeTime int64, stackID int32) {
	switch sampleType {
	case sampleTypeError:
		return

	case sampleTypeManaged:
		switch t.state() {
		case uninitialized:
			t.managedSample(stackID, relativeTime)
			t.lastExternalTime = -1
		case managed:
			t.managedSample(stackID, relativeTime)
		case external:
			t.externalSample(stackID, relativeTime)
			t.lastExternalTime = -relativeTime
		}
		t.lastManagedTime = relativeTime

	case sampleTypeExternal:
		switch t.state() {
		case external, uninitialized:
			t.externalSample(stackID, relativeTime)
		case managed:
			t.managedSample(stackID, relativeTime)
		}
		t.lastExternalTime = relativeTime
	}
}

func (t *thread) managedSample(stackID int32, rt int64) {
	if t.lastManagedTime > 0 {
		t.samples[stackID] += rt - t.lastManagedTime
	}
}

func (t *thread) externalSample(stackID int32, rt int64) {
	if t.lastExternalTime > 0 && !t.managedOnly {
		t.samples[stackID] += rt - t.lastExternalTime
	}
}
