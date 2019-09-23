package common

import (
	"os"
	"strings"
	"errors"
	"fmt"
)

func CreateFile(file string)   (f *os.File,err error){
	if file ==""{
		return nil,errors.New("cannot create file without name")
	}
	idx:=strings.LastIndex(file,"/")
	if idx == len(file)-1{
		err =  errors.New(fmt.Sprintf("%s: is directory ,not file",file))
		return
	}

	if idx>=0{
		dir:=file[:idx]
		if len(dir)>0{
			if err =os.MkdirAll(dir,0666);err!=nil{
				return
			}
		}
	}

	f,err = os.Create(file)
	return

}