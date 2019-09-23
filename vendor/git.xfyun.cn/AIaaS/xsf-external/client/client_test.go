package xsf

import (
	"fmt"
	"reflect"
)

func t(){
//ssh
    /*cli := NewClient()
    conn,e := cli.GetConn("ats")
	if e != nil {
		c := NewXsfCallClient(conn)//连接坏掉怎么移除
		var rd ReqData
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
        resd,e := c.Call(ctx,&rd)
        fmt.Println(resd,e)
	}*/
	var x map[string]string
	_,ok :=x["hello"]
	fmt.Println(reflect.TypeOf(ok))

//基于busin，
/*
1、只需要传一个data，收一个data
2、

*/
}
