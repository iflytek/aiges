package xsf

import "sync"

var VersionCheckList map[string]VersionChecker
var VersionRwMu sync.RWMutex

type VersionChecker interface {
	Version() string
}

func AddVersionCheck(name string, ver VersionChecker) {
	VersionRwMu.Lock()
	VersionCheckList[name] = ver
	VersionRwMu.Unlock()
}

func init() {
	VersionRwMu.Lock()
	VersionCheckList = make(map[string]VersionChecker)
	VersionRwMu.Unlock()
}
