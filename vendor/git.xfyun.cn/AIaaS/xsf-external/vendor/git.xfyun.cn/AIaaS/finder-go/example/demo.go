package main

import (
	"errors"
	"fmt"
	"git.xfyun.cn/AIaaS/finder-go/log"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"git.xfyun.cn/AIaaS/finder-go"

	"encoding/json"
	common "git.xfyun.cn/AIaaS/finder-go/common"
	"git.xfyun.cn/AIaaS/finder-go/utils/httputil"
)

type ServiceItemTest struct {
	ServiceName string
	ApiVersion  string
}
type RegisterItemTest struct {
	ServiceAddr string
	ApiVersion  string
}
type TestConfig struct {
	Type               int               `json:"type"`
	CompanionUrl       string            `json:"companionUrl"`
	Address            string            `json:"address"`
	Project            string            `json:"project"`
	Group              string            `json:"group"`
	Service            string            `json:"service"`
	Version            string            `json:"version"`
	ProviderApiVersion string            `json:"providerApiVersion"`
	SubscribeFile      []string          `json:"subscribeFile"`
	UnSubscribeFile    []string          `json:"unSubscribeFile"`
	SubribeServiceItem []ServiceItemTest `json:"subribeServiceItem"`
	UnSubscribeTime    time.Duration     `json:"unSubscribeTime"`
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("参数错误")
		return
	}
	loggerConfig := zap.NewProductionConfig()
	//loggerConfig.EncoderConfig.EncodeTime = normalTimeEncoder
	//TODO 日志目录
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.EncoderConfig.TimeKey = "time"

	logger, _ := loggerConfig.Build()
	Logger := logger.Sugar()
	zkLog := ZkLogger{
		zap: Logger,
	}
	zkLog.Printf("dddd")
	file, _ := os.Open(args[1])
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := TestConfig{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(conf)
	if conf.Type == 1 {
		//订阅配置文件。和之前的区别不大，主要是回调函数增加了OnError方法，如：当配置文件不存在的话，可以进行通知。可为空实现
		newConfigFinder(conf)
	} else if conf.Type == 2 {
		//订阅服务。和之前的区别主要是增加了版本号的概念，用于指定服务的特定版本。。且回调函数的参数也增加了一个版本号的参数。。用于明确服务的版本信息
		newServiceFinder(conf, nil)
	} else if conf.Type == 3 {
		//注册服务.和之前的区别是注册服务的时候，必须制定版本号
		newProviderFinder(conf)
	} else if conf.Type == 4 {
		newConfigFinder(conf)
		newServiceFinder(conf, nil)
	} else if conf.Type == 5 {
		newConfigFinder(conf)

	} else if conf.Type == 6 {
		newQueryServiceFinder(conf)
		return
	}else if conf.Type==7 {
		newQueryServiceNoWatchFinder(conf)
		return
	}
	//newConfigFinder("127.0.0.1:10010", []string{"xsfs.toml"})
	//newProviderFinder("299.99.99.99:99")
	//newProviderFinder("299.99.99.99:100")
	//TODO  1. companion连不上怎么办，zk连不上怎么办
	for {
		time.Sleep(time.Minute * 20)
		fmt.Println("I'm running.")
	}

}

func pressureTest(conf TestConfig) {

}
func newServiceFinder(conf TestConfig, lg log.Logger) {

	cachePath, err := os.Getwd()
	if err != nil {
		return
	}
	//缓存信息的存放路径
	cachePath += "/findercache"
	config := common.BootConfig{
		//companion地址
		CompanionUrl: conf.CompanionUrl,
		//缓存路径
		CachePath: cachePath,
		//是否缓存服务信息
		CacheService: true,
		//是否缓存配置信息
		CacheConfig:   true,
		ExpireTimeout: 10 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: conf.Project,
			Group:   conf.Group,
			Service: conf.Service,
			Version: conf.Version,
			Address: conf.Address,
		},
	}

	//创建finder。
	f, err := finder.NewFinderWithLogger(config, lg)

	//init
	log.Println("----------------------------------------------")
	if err != nil {
		fmt.Println(err)
	} else {
		//订阅服务。和之前的区别是订阅服务的时候，除了指定服务名外，必须指定版本号
		testUseServiceAsync(f, conf.SubribeServiceItem)
		//testUseServiceAsync(f, conf.SubribeServiceItem)

	}
}
func newProviderFinder(conf TestConfig) {
	cachePath, err := os.Getwd()
	if err != nil {
		return
	}
	cachePath += "/findercache"
	config := common.BootConfig{
		CompanionUrl: conf.CompanionUrl,
		//缓存路径
		CachePath: cachePath,
		//是否缓存服务信息
		CacheService: true,
		//是否缓存配置信息
		CacheConfig: true,
		//和zk之间的会话超时时间
		ExpireTimeout: 5 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: conf.Project,
			Group:   conf.Group,
			Service: conf.Service,
			Version: conf.Version,
			Address: conf.Address,
		},
	}

	loggerConfig := zap.NewProductionConfig()
	//loggerConfig.EncoderConfig.EncodeTime = normalTimeEncoder
	//TODO 日志目录
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.EncoderConfig.TimeKey = "time"

	logger, _ := loggerConfig.Build()
	Logger := logger.Sugar()
	zkLog := ZkLogger{
		zap: Logger,
	}
	zkLog.Infof("ddd")
	//创建finder。
	f, err := finder.NewFinderWithLogger(config, nil)

	if err != nil {
		fmt.Println(err)
	} else {
		//和之前的区别是，必须指定对应的版本号
		testRegisterService(f, conf.Address, conf.ProviderApiVersion)
	}
}

func testRegisterService(f *finder.FinderManager, addr string, apiVersion string) {

	//必须指定版本号
	f.ServiceFinder.RegisterServiceWithAddr(addr, apiVersion)

}

type ZkLogger struct {
	zap *zap.SugaredLogger
}

func (l *ZkLogger) Infof(fmt string, v ...interface{}) {
	l.zap.Infof(fmt, v)
}
func (l *ZkLogger) Debugf(fmt string, v ...interface{}) {
	l.zap.Infof(fmt, v)
}
func (l *ZkLogger) Errorf(fmt string, v ...interface{}) {
	l.zap.Infof(fmt, v)
}
func (l *ZkLogger) Printf(fmt string, v ...interface{}) {
	l.zap.Infof(fmt, v)
}
func newQueryServiceNoWatchFinder(conf TestConfig) {
	cachePath, err := os.Getwd()
	if err != nil {
		return
	}
	cachePath += "/findercache"

	//元数据信息

	config := common.BootConfig{
		//CompanionUrl:     "http://companion.xfyun.iflytek:6868",
		//compaion地址
		CompanionUrl: conf.CompanionUrl,
		//缓存文件的地址
		CachePath: cachePath,
		//缓存服务提供者的信息，当为true的时候，和zk连接不上的话，使用缓存信息
		CacheService: true,
		//缓存配置文件信息。
		CacheConfig: true,
		//和zk之间的会话时间
		ExpireTimeout: 5 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: conf.Project,
			Group:   conf.Group,
			Service: conf.Service,
			Version: conf.Version,
			Address: conf.Address,
		},
	}
	loggerConfig := zap.NewProductionConfig()
	//loggerConfig.EncoderConfig.EncodeTime = normalTimeEncoder
	//TODO 日志目录
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.EncoderConfig.TimeKey = "time"

	logger, _ := loggerConfig.Build()
	Logger := logger.Sugar()
	zkLog := ZkLogger{
		zap: Logger,
	}
	zkLog.Infof("ddddd")
	f, err := finder.NewFinderWithLogger(config, nil)
	if err != nil {
		panic(err)

	}
	for {
		dat,_:=f.ServiceFinder.QueryService(conf.Project, conf.Group)

		fmt.Println(dat["db-proxy"])

		time.Sleep(1*time.Second)
	}
	//handler := new(ServiceChangedHandle)

}
func newQueryServiceFinder(conf TestConfig) {
	cachePath, err := os.Getwd()
	if err != nil {
		return
	}
	cachePath += "/findercache"

	//元数据信息

	config := common.BootConfig{
		//CompanionUrl:     "http://companion.xfyun.iflytek:6868",
		//compaion地址
		CompanionUrl: conf.CompanionUrl,
		//缓存文件的地址
		CachePath: cachePath,
		//缓存服务提供者的信息，当为true的时候，和zk连接不上的话，使用缓存信息
		CacheService: true,
		//缓存配置文件信息。
		CacheConfig: true,
		//和zk之间的会话时间
		ExpireTimeout: 5 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: conf.Project,
			Group:   conf.Group,
			Service: conf.Service,
			Version: conf.Version,
			Address: conf.Address,
		},
	}
	loggerConfig := zap.NewProductionConfig()
	//loggerConfig.EncoderConfig.EncodeTime = normalTimeEncoder
	//TODO 日志目录
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.EncoderConfig.TimeKey = "time"

	logger, _ := loggerConfig.Build()
	Logger := logger.Sugar()
	zkLog := ZkLogger{
		zap: Logger,
	}
	zkLog.Infof("ddddd")
	f, err := finder.NewFinderWithLogger(config, nil)
	if err != nil {
		panic(err)

	}
	handler := new(ServiceChangedHandle)
	dat,_:=f.ServiceFinder.QueryServiceWatch("AIaaS", "dx",handler)
	fmt.Println("---------------------------------------------------------------")
	fmt.Println(dat)
	fmt.Println("---------------------------------------------------------------")
	dd,_:=json.Marshal(dat)
	fmt.Println(string(dd))
	time.Sleep(1*time.Hour)
}
func newConfigFinder(conf TestConfig) {
	cachePath, err := os.Getwd()
	if err != nil {
		return
	}
	cachePath += "/findercache"

	//元数据信息

	config := common.BootConfig{
		//CompanionUrl:     "http://companion.xfyun.iflytek:6868",
		//compaion地址
		CompanionUrl: conf.CompanionUrl,
		//缓存文件的地址
		CachePath: cachePath,
		//缓存服务提供者的信息，当为true的时候，和zk连接不上的话，使用缓存信息
		CacheService: true,
		//缓存配置文件信息。
		CacheConfig: true,
		//和zk之间的会话时间
		ExpireTimeout: 5 * time.Second,
		MeteData: &common.ServiceMeteData{
			Project: conf.Project,
			Group:   conf.Group,
			Service: conf.Service,
			Version: conf.Version,
			Address: conf.Address,
		},
	}
	loggerConfig := zap.NewProductionConfig()
	//loggerConfig.EncoderConfig.EncodeTime = normalTimeEncoder
	//TODO 日志目录
	loggerConfig.OutputPaths = []string{"stdout"}
	loggerConfig.EncoderConfig.TimeKey = "time"

	logger, _ := loggerConfig.Build()
	Logger := logger.Sugar()
	zkLog := ZkLogger{
		zap: Logger,
	}
	zkLog.Infof("ddddd")
	f, err := finder.NewFinderWithLogger(config, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		//使用并订阅文件的变更。。
		testUseConfigAsyncByName(f, conf.SubscribeFile)
		if conf.Type == 5 {
			ss := conf.UnSubscribeTime
			tick := time.NewTicker(ss * time.Minute)
			select {
			case <-tick.C:
				fmt.Println("开始取消文件")
				go testUnscribeConfigfile(f, conf.UnSubscribeFile)

			}
		}

	}
}
func testUnscribeConfigfile(f *finder.FinderManager, names []string) {
	f.ConfigFinder.BatchUnSubscribeConfig(names)
}
func getLocalIP(url string) (string, error) {
	var host string
	var port string
	var localIP string
	items := strings.Split(url, ":")
	if len(items) == 3 {
		host = strings.Replace(items[1], "/", "", -1)
		port = items[2]
	} else if len(items) == 2 {
		host = strings.Replace(items[0], "/", "", -1)
		port = items[1]
	} else {
		host = url
		port = "80"
	}

	if len(host) == 0 {
		return "", errors.New("testRemote:invalid remote url")
	}
	if len(port) == 0 {
		port = "80"
	}
	ips, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}
	for _, ip := range ips {
		conn, err := net.Dial("tcp", ip+":"+port)
		if err != nil {
			fmt.Println("testRemote:", err)
			continue
		}
		localIP = conn.LocalAddr().String()
		fmt.Println("testRemote:ok")
		err = conn.Close()
		if err != nil {
			fmt.Println("testRemote:", err)
			break
		}
		break
	}

	if len(localIP) == 0 {
		return "", errors.New("testRemote:failed")
	}

	fmt.Println("local ip:", localIP)

	return localIP, nil
}

func testCache(cachepath string) {
	configFile := `[test]\r\n\titem = "value"`
	config := &common.Config{
		Name: "default.cfg",
		File: []byte(configFile),
	}
	err := finder.CacheConfig(cachepath, config)
	if err != nil {
		fmt.Println(err)
	}
	c, err := finder.GetConfigFromCache(cachepath, "default.cfg")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("default.cfg:", string(c.File))
	}

	zkInfo := &common.StorageInfo{
		Addr:            []string{"10.1.86.73:2181", "10.1.86.74:2181"},
		ConfigRootPath:  "/polaris/config/",
		ServiceRootPath: "/polaris/service/",
	}
	err = finder.CacheStorageInfo(cachepath, zkInfo)
	if err != nil {
		fmt.Println(err)
	}
	newZkInfo, err := finder.GetStorageInfoFromCache(cachepath)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ZkAddr:", newZkInfo.Addr)
		fmt.Println("ConfigRootPath:", newZkInfo.ConfigRootPath)
		fmt.Println("ServiceRootPath:", newZkInfo.ServiceRootPath)
	}

	service := &common.Service{
		ServiceName:  "xrpc",
		ProviderList: []*common.ServiceInstance{},
		Config:       &common.ServiceConfig{},
	}
	instance := &common.ServiceInstance{
		Addr: "127.0.0.0:9091",
		Config: &common.ServiceInstanceConfig{
			IsValid: true,
		},
	}
	service.ProviderList = append(service.ProviderList, instance)

	err = finder.CacheService(cachepath, service)
	if err != nil {
		fmt.Println(err)
	}

}

func testConfigFeedback() {
	url := "http://10.1.200.75:9080/finder/push_config_feedback"
	contentType := "application/x-www-form-urlencoded"
	hc := &http.Client{
		Transport: &http.Transport{
			Dial: func(nw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(1 * time.Second)
				c, err := net.DialTimeout(nw, addr, time.Second*1)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	params := []byte("pushId=123456&project=test&group=default&service=xrpc&version=1.0.0&config=default.cfg&addr=10.1.86.221:9091&update_status=1&update_time=1513044755&load_status=1&load_time=1513044757")
	result, err := httputil.DoPost(hc, contentType, url, params)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(result))
}

func testUseService(f *finder.FinderManager) {
	//	handler := new(ServiceChangedHandle)
	item := []common.ServiceSubscribeItem{}
	item = append(item, common.ServiceSubscribeItem{ServiceName: "test0803", ApiVersion: "1.0"})
	serviceList, err := f.ServiceFinder.UseService(item)
	//serviceList, err := f.ServiceFinder.UseAndSubscribeService([]string{"iatExecutor"}, handler)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, s := range serviceList {
			fmt.Println(s.ServiceName, ":")
			for _, item := range s.ProviderList {
				fmt.Println("addr:", item.Addr)
				fmt.Println("is_valid:", item.Config.IsValid)
			}
		}
	}

	count := 0
	for {
		count++
		if count > 200 {
			//f.ConfigFinder.UnSubscribeConfig("default.toml")
		}
		if count > 600 {
			//err = f.ServiceFinder.UnRegisterService()
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("UnRegisterService is ok.")
			}
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func testGrayData(f *finder.FinderManager) {
	f.ConfigFinder.UseConfig([]string{"ddd"})
}
func testServiceAsync(f *finder.FinderManager) {

	var err error
	err = f.ServiceFinder.RegisterService("1.0")
	//err = f.ServiceFinder.RegisterServiceWithAddr("10.1.203.36:50052")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("RegisterService is ok.")
	}
	time.Sleep(time.Second * 2)

}

func testUseServiceAsync(f *finder.FinderManager, items []ServiceItemTest) {
	handler := new(ServiceChangedHandle)
	subscri := make([]common.ServiceSubscribeItem, 0)
	for _, item := range items {
		subscri = append(subscri, common.ServiceSubscribeItem{ServiceName: item.ServiceName, ApiVersion: item.ApiVersion})
	}

	//订阅服务，订阅服务的还是，增加具体的版本号。。且必须制定订阅的是服务的那个版本。所有的回调函数增加服务版本号参数，用于说明具体的服务信息
	serviceList, err := f.ServiceFinder.UseAndSubscribeService(subscri, handler)
	//serviceList, err := f.ServiceFinder.UseService(subscri)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, s := range serviceList {
			fmt.Println("订阅的服务：", s.ServiceName, ":", s.ApiVersion, " --->")
			for _, item := range s.ProviderList {
				fmt.Println("----提供者地址 :")
				fmt.Println("--------:", item.Addr)
			}
		}
	}

}

func testUseConfigAsync(f *finder.FinderManager) {

	handler := ConfigChangedHandle{}
	count := 0

	//f.ConfigFinder.UseAndSubscribeConfig([]string{"test2.toml", "xsfc.toml.cfg"}, handler)
	configFiles, err := f.ConfigFinder.UseAndSubscribeConfig([]string{"2.yml"}, &handler)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range configFiles {
		fmt.Println(c.Name, ":\r\n", string(c.File))
	}

	for {
		//fmt.Println("The ", count, "th show:")
		//configFiles, err := f.ConfigFinder.UseAndSubscribeConfig([]string{"test2.toml", "xsfc.tmol"}, handler)

		//f.ConfigFinder.UseAndSubscribeConfig([]string{"test2.toml", "xsfc.tmol"}, handler)
		//configFiles, err := f.ConfigFinder.UseConfig([]string{"xsfc.tmol"})

		if count > 200 {
			f.ConfigFinder.UnSubscribeConfig("default.toml")
		}
		if count > 600 {
			break
		}
		count++
		time.Sleep(time.Second * 1)
	}

}
func testUserConfig(f *finder.FinderManager, name []string) {
	configFiles, err := f.ConfigFinder.UseConfig(name)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range configFiles {
		fmt.Println(c.Name, ":\r\n", string(c.File))
	}

}

func testUseConfigAsyncByName(f *finder.FinderManager, name []string) {

	handler := ConfigChangedHandle{}

	//使用并订阅文件的变更。回调函数相比之前多了一个OnError .用于在运行过程中出现解析文件错误的时候进行通知，可以为空的实现
	configFiles, err := f.ConfigFinder.UseAndSubscribeConfig(name, &handler)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range configFiles {
		fmt.Println("首次获取配置文件名称：", c.Name, "  、\r\n内容为:\r\n", string(c.File))
	}

	configFiles, err = f.ConfigFinder.UseConfig(name)
	if err != nil {
		fmt.Println(err)
	}
	for _, c := range configFiles {
		fmt.Println("首次获取配置文件名称：", c.Name, "  、\r\n内容为:\r\n", string(c.File))
	}
}
