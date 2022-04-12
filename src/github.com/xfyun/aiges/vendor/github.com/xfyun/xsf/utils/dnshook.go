package utils

import (
	"net"
	"context"
)

type dnsDail struct{
	addr string
}


/*
Dail:
desc:
    重写DNS中的dail行为
*/
func (d *dnsDail) Dail(ctx context.Context, network, address string) (net.Conn, error) {
    var n net.Dialer
	if len(d.addr) <= 0  {
      return  n.DialContext(ctx, network, address)
	}
	return n.DialContext(ctx, network, d.addr)
	}

	/*
	LookupHost:
	desc:
	根据host到指定的DNS上拉取对应的ip信息。windows 下默认使用系统DNS配置
	*/
func LookupHost(host string, dns string) ( []string,  error) {
  d :=new(dnsDail)
  d.addr = dns
  net.DefaultResolver.Dial = d.Dail
  return net.LookupHost(host)
}
