package common

import "os/exec"

func SetSystemEnv(k,v string) error {
	cmd:=exec.Command("export" ,k+"="+v)
	return cmd.Run()
}
