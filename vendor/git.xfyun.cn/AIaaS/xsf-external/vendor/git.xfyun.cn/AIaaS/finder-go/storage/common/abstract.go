package common

// case zk.EventNodeCreated:
// 		return
// 	case zk.EventNodeDeleted:
// 		return
// 	case zk.EventNodeDataChanged:
// 		return
// 	case zk.EventNodeChildrenChanged:
// 		return
// 	case zk.EventNotWatching:
// 		return
type ChangedCallback interface {
	DataChangedCallback(path string, node string, data []byte)
	ChildrenChangedCallback(path string, node string, children []string)
	Process(path string, node string)
	ChildDeleteCallBack(path string)
}
