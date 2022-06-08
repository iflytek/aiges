package utils

type configurator struct {
	cfg   string
	addrs []string
}

func (c *configurator) init(host string, dns string, srv string, tag string) error {
	var err error
	c.addrs, err = LookupHost(host, dns)
	if err != nil {
		return err
	}
	//todo 访问这些地址，然后向这些地址请求服务的配置
	return err
}
func (c *configurator) getcfg(srv string, tag string) string {
	return c.cfg
}
