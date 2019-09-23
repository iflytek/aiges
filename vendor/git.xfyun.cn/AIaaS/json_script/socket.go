package jsonscpt

import (
	"net"
)

var SocketListen Func = func(i ...interface{}) interface{} {
	if len(i)>0{

		ls,err:=net.Listen("tcp4",ConvertToString(i[0]))
		if err !=nil{
			return -1
		}
		return ls
	}
	return -1
}

var SocketAccept Func = func(i ...interface{}) interface{} {
	if len(i)>0{
		if ls,ok:=i[0].(net.Listener);ok{
			app,err:=ls.Accept()
			if err !=nil{
				return -1
			}
			return app
		}
	}
	return -1
}

var SocketRead Func = func(i ...interface{}) interface{} {
	if len(i)>1{
		conn,ok:=i[0].(net.Conn)
		if ok{
			b:=make([]byte,int(number(i[1])))
			n,err:=conn.Read(b)
			if err !=nil{
				return -1
			}
			return string(b[:n])
		}
	}
	return -1
}