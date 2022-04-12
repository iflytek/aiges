package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	xsfInitStatus  = "xsf_initializing"
	userInitStatus = "user_initializing"
	bvtCheckStatus = "bvt_checking"
	finishStatus   = "xsf_done"
)
const (
	workDir  = "xsf_status"
	workFile = "recording"
)

var (
	statusRequired = false
)

func init() {
	switch goos := runtime.GOOS; goos {
	case "linux":
		{

			statusRequired = true
			_ = exec.Command("mkdir", "-p", workDir).Run()
		}
	default:
		{
			statusRequired = false
		}
	}
}
func flushStatus(status string) {
	if !statusRequired {
		return
	}
	f, e := os.Create(filepath.Join(workDir, workFile))
	if e != nil {
		panic(e.Error())
	}
	_, e = f.WriteString(status)
	if e != nil {
		panic(e.Error())
	}
	f.Close()
}
func SyncXsfInitStatus() {
	flushStatus(xsfInitStatus)
}
func SyncUserInitStatus() {
	flushStatus(userInitStatus)
}
func SyncBvtInitStatus() {
	flushStatus(bvtCheckStatus)
}
func SyncFinishStatus() {
	flushStatus(finishStatus)
}
