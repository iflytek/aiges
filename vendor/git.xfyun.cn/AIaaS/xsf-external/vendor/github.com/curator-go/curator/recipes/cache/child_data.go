package cache

import "github.com/samuel/go-zookeeper/zk"

// ChildData contains data of a node including: stat, data, path
type ChildData struct {
	path string
	stat *zk.Stat
	data []byte
}

// NewChildData creates ChildData
func NewChildData(path string, stat *zk.Stat, data []byte) *ChildData {
	return &ChildData{
		path: path,
		stat: stat,
		data: data,
	}
}

// Path returns the full path of this child
func (cd ChildData) Path() string {
	return cd.path
}

// SetStat sets the stat data for this child
func (cd ChildData) SetStat(s *zk.Stat) {
	cd.stat = s
}

// Stat returns the stat data for this child
func (cd ChildData) Stat() *zk.Stat {
	return cd.stat
}

// Data returns the node data for this child when the cache mode is set to cache data
func (cd ChildData) Data() []byte {
	return cd.data
}
