package curator

import "testing"

// TestClientStartRace tests race in Start(), please use -race with go test
func TestClientStartRace(t *testing.T) {
	NewClient("127.0.0.1:2181", nil).Start()
}
