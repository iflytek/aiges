package cache

import "sync/atomic"

// NodeState represents state of TreeNode.
// TODO: make this a private type?
type NodeState int32

// Available node states.
const (
	NodeStatePENDING NodeState = iota
	NodeStateLIVE
	NodeStateDEAD
)

// Load loads then returns the state atomically.
func (s *NodeState) Load() NodeState {
	return NodeState(atomic.LoadInt32((*int32)(s)))
}

// CompareAndSwap set the state to new if value is old atomatically.
func (s *NodeState) CompareAndSwap(old, new NodeState) bool {
	return atomic.CompareAndSwapInt32((*int32)(s), int32(old), int32(new))
}

// Store sets the state to given value atomically.
func (s *NodeState) Store(new NodeState) {
	atomic.StoreInt32((*int32)(s), int32(new))
}

// Swap sets the state to given value and returns the old state atomically.
func (s *NodeState) Swap(new NodeState) NodeState {
	return NodeState(atomic.SwapInt32((*int32)(s), int32(new)))
}
