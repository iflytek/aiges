package main

import (
	"bufio"
	"flag"
	"fmt"
	"git.xfyun.cn/AIaaS/json_script"
	"io/ioutil"
	"os"
	"time"
)
var file = flag.String("f","","")
func main(){
	flag.Parse()
	vm:=jsonscpt.NewVm()
	if *file==""{
		sc:=bufio.NewScanner(os.Stdin)
		for sc.Scan(){
			func (){
				defer func() {
					if err := recover();err !=nil{
						fmt.Println(err)
					}
				}()
				vm.ExecJsonObject(vm.ExecJsonObject(sc.Text()))
			}()
		}
		return
	}
	b,err:=ioutil.ReadFile(*file)
	if err !=nil{
		fmt.Println(err.Error())
		return
	}
	cmd,err:=jsonscpt.CompileExpFromJson(b)
	if err !=nil{
		fmt.Println(err)
		return
	}
	start:=time.Now()
	for i:=0;i<1;i++{
		if err:=vm.SafeExecute(cmd, func(err interface{}) {
			fmt.Println(err)
		});err !=nil{
			fmt.Println(err)
		}
	}

	fmt.Println("total cost=>",time.Since(start))

}

