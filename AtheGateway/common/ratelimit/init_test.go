package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	f := 0x7fffffffffffffff
	fmt.Println(f)
	Init()
	//fmt.Println(CheckLimit("100IME1","sdf"))
	//ReleaseLimit("100IME1","sdf")
	//fmt.Println(CheckLimit("100IME1","sdf"))

	//for i:=0;i<11;i++{
	//	 func() {
	//		b:=CheckMaxConnLimit()
	//		if !b{
	//			fmt.Println("f-------f")
	//		}
	//		time.Sleep(2*time.Second)
	//		//ReleaseMaxConnLimit()
	//	}()
	//}
	fmt.Println(CheckMaxConnLimit())
	ReleaseMaxConnLimit()
	fmt.Println(CheckMaxConnLimit())
	b := CheckMaxConnLimit()
	fmt.Println(b)
	time.Sleep(10 * time.Second)
}

func TestInitByConf(t *testing.T) {
	InitByConf(1)

	fmt.Println(CheckMaxConnLimit())
	ReleaseMaxConnLimit()
	fmt.Println(CheckMaxConnLimit())
	ReleaseMaxConnLimit()
	fmt.Println(CheckMaxConnLimit())
	ReleaseMaxConnLimit()
	fmt.Println(CheckMaxConnLimit())
}
