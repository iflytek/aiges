package catch

import (
	"os/exec"
	"time"
)

func GetPstack(pid string) (cStack string) {
	var cmd *exec.Cmd
	endFlag := make(chan bool)
	go func() {
		defer func() { endFlag <- true }()
		cmd = exec.Command("pstack", pid)
		if output, err := cmd.Output(); err != nil {
			cStack = "failed to get cstack " + err.Error()
		} else {
			cStack = string(output)
		}
	}()
	select {
	case <-endFlag:
		return
	case <-time.After(time.Duration(3) * time.Second):
		cStack = "failed to get cstack timeout"
		_ = cmd.Process.Kill()
		return
	}
}
