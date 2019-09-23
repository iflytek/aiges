package common

import (
	"runtime"
	"os"
	"io"
	"fmt"
)
//把环境变量保存到文件中，用户动态设置环境变量
func Setenv(k,v string)  (err error){
	var envfileName string
	if runtime.GOOS == "windows"{
		envfileName  = "C:/users/admin/webgate-env"
	}else{
		envfileName = "/etc/webgate-env"
	}
	var file *os.File
	file,err=os.OpenFile(envfileName,os.O_WRONLY,0666)
	if err != nil{
		file,err = os.Create(envfileName)
		if err !=nil{
			return
		}
	}
	defer file.Close()
	n,err:=file.Seek(0,io.SeekEnd)
	if err !=nil{
		return
	}
	_,err = file.WriteAt([]byte(fmt.Sprintf("export %s=%s\n",k,v)),n)

	return
}
