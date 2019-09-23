package cache

import "fmt"

// TreeCacheEventType represents the type of change to a path
type TreeCacheEventType int

const (
	// TreeCacheEventNodeAdded indicates a node was added
	TreeCacheEventNodeAdded TreeCacheEventType = iota
	// TreeCacheEventNodeUpdated indicates a node's data was changed
	TreeCacheEventNodeUpdated
	// TreeCacheEventNodeRemoved indicates a node was removed from the tree
	TreeCacheEventNodeRemoved

	// TreeCacheEventConnSuspended is called when the connection has changed to SUSPENDED
	TreeCacheEventConnSuspended
	// TreeCacheEventConnReconnected is called when the connection has changed to RECONNECTED
	TreeCacheEventConnReconnected
	// TreeCacheEventConnLost is called when the connection has changed to LOST
	TreeCacheEventConnLost
	// TreeCacheEventInitialized is posted after the initial cache has been fully populated
	TreeCacheEventInitialized
)

// String returns the string representation of TreeCacheEventType
// "Unknown" is returned when event type is unknown
func (et TreeCacheEventType) String() string {
	switch et {
	case TreeCacheEventNodeAdded:
		return "NodeAdded"
	case TreeCacheEventNodeUpdated:
		return "NodeUpdated"
	case TreeCacheEventNodeRemoved:
		return "NodeRemoved"
	case TreeCacheEventConnSuspended:
		return "ConnSuspended"
	case TreeCacheEventConnReconnected:
		return "ConnReconnected"
	case TreeCacheEventConnLost:
		return "ConnLost"
	case TreeCacheEventInitialized:
		return "Initialized"
	default:
		return "Unknown"
	}
}

// TreeCacheEvent represents a change to a path
type TreeCacheEvent struct {
	Type TreeCacheEventType
	Data *ChildData
}

// String returns the string representation of TreeCacheEvent
func (e TreeCacheEvent) String() string {
	var path string
	var data []byte
	if e.Data != nil {
		path = e.Data.Path()
		data = e.Data.Data()
	}
	return fmt.Sprintf("TreeCacheEvent{%s %s '%s'}", e.Type, path, data)
}
