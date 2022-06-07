package bvt

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	finder "github.com/xfyun/finder-go/common"
	"github.com/xfyun/xsf/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"
)

var bvtLogger = log.New(os.Stderr, "", log.LstdFlags)
var bvtCfgNotFound = errors.New("can't find cfg")

var headerKvs = map[string]string{"Content-Type": "application/json;charset=utf-8"}
var (
	cmdGProject       = flag.String("bvtProject", "", "bvtProject")
	cmdGGroup         = flag.String("bvtGroup", "", "bvtGroup")
	cmdGService       = flag.String("bvtService", "", "bvtService")
	cmdGVersion       = flag.String("bvtVersion", "", "bvtVersion")
	cmdGCfgFile       = flag.String("bvtCfgFile", "", "bvtCfgFile")
	cmdGCompanionUrl  = flag.String("bvtCompanionUrl", "", "bvtCompanionUrl")
	cmdTimeout        = flag.Int("bvtTimeout", 0, "bvtTimeout,unit:ms")
	cmdPlatformAddr   = flag.String("bvtPlatform", "", "bvtPlatform")
	cmdId             = flag.String("bvtId", "", "bvtId")
	cmdEngIp          = flag.String("bvtEngIp", "", "bvtEngIp")
	cmdCallback       = flag.String("bvtCb", "", "bvt callback")
	cmdAsync          = flag.Int("bvtAsync", -1, "bvt async,1:true,2:false")
	cmdServiceAddress = flag.String("bvtServiceAddress", "", "bvtServiceAddress")
	cmdBvt            = flag.Int("bvtAble", -1, "bvt able; 0:disable,1:enable")
)

const (
	defaultGEnable = true

	ExceptionDef        = 100 //默认异常退出码
	ExceptionNewRequest = 101 //构建请求对象失败
	ExceptionDeadline   = 102 //超时
	ExceptionBody       = 103 //读取body失败
	ExceptionBvt        = 104 //bvt结果不符合预期
)

type bvtVerifier struct {
	ready bool
	/*
		g开头表示一些共有的变量，独立存储
	*/
	gAble         bool
	gProject      string
	gGroup        string
	gService      string
	gVersion      string
	gCfgFile      string
	gCompanionUrl string
	gConfigure    *utils.Configure

	timeout        time.Duration
	ctx            context.Context
	platformAddr   string
	id             string
	engIp          string
	callback       string
	async          bool
	serviceAddress string

	licMax    int
	service   string
	namespace string

	bvtAsync bool //回滚场景下，先响应成功，然后异步执行bvt，失败仍然退出
}

func (bvt *bvtVerifier) String() string {
	return fmt.Sprintf(
		"ready:%v,"+
			"gAble:%v,"+
			"gProject:%v,"+
			"gGroup:%v,"+
			"gService:%v,"+
			"gVersion:%v,"+
			"gCfgFile:%v,"+
			"gCompanionUrl:%v,"+
			"timeout:%v,"+
			"platformAddr:%v,"+
			"id:%v,"+
			"engIp:%v,"+
			"serviceAddress:%v",
		bvt.ready,
		bvt.gAble,
		bvt.gProject,
		bvt.gGroup,
		bvt.gService,
		bvt.gVersion,
		bvt.gCfgFile,
		bvt.gCompanionUrl,
		bvt.timeout,
		bvt.platformAddr,
		bvt.id,
		bvt.engIp,
		bvt.serviceAddress,
	)
}
func (bvt *bvtVerifier) init(
	gProject string,
	gGroup string,
	gService string,
	gVersion string,
	gCfgFile string,
	gCompanionUrl string,

	timeout time.Duration,
	platformAddr string,
	id string,
	engIp string,
	callback string,
	async bool,
	serviceAddress string,

	licMax int,
	service string,
	namespace string,
) {
	{
		bvt.gProject = gProject
		if len(*cmdGProject) != 0 {
			bvt.gProject = *cmdGProject
		}
		bvt.gGroup = gGroup
		if len(*cmdGGroup) != 0 {
			bvt.gGroup = *cmdGGroup
		}
		bvt.gService = gService
		if len(*cmdGService) != 0 {
			bvt.gService = *cmdGService
		}
		bvt.gVersion = gVersion
		if len(*cmdGVersion) != 0 {
			bvt.gVersion = *cmdGVersion
		}
		bvt.gCfgFile = gCfgFile
		if len(*cmdGCfgFile) != 0 {
			bvt.gCfgFile = *cmdGCfgFile
		}
		bvt.gCompanionUrl = gCompanionUrl
		if len(*cmdGCompanionUrl) != 0 {
			bvt.gCompanionUrl = *cmdGCompanionUrl
		}
		bvt.timeout = timeout
		if *cmdTimeout != 0 {
			bvt.timeout = time.Millisecond * time.Duration(*cmdTimeout)
		}
		bvt.ctx, _ = context.WithTimeout(context.Background(), bvt.timeout)
		bvt.platformAddr = platformAddr
		if len(*cmdPlatformAddr) != 0 {
			bvt.platformAddr = *cmdPlatformAddr
		}
		bvt.id = id
		if len(*cmdId) != 0 {
			bvt.id = *cmdId
		}
		bvt.engIp = engIp
		if len(*cmdEngIp) != 0 {
			bvt.engIp = *cmdEngIp
		}
		bvt.callback = callback
		if len(*cmdCallback) != 0 {
			bvt.callback = *cmdCallback
		}

		bvt.async = async
		if *cmdAsync != -1 {
			bvt.async = func() bool {
				if *cmdAsync == 0 {
					return false
				}
				return true
			}()
		}

		bvt.serviceAddress = serviceAddress
		if len(*cmdServiceAddress) != 0 {
			bvt.serviceAddress = *cmdServiceAddress
		}

		bvt.licMax = licMax
		bvt.service = service
		bvt.namespace = namespace
	}
	{
		bvtLogger.Printf("bvtVerifier dump:%s\n", bvt.String())
	}
	{
		bvt.getConfigure()
		bvt.gAble = defaultGEnable
		enable, enableErr := bvt.gConfigure.GetInt("bvt", "able")
		if enableErr == nil && enable == 0 {
			bvt.gAble = false
		}
		if !bvt.gAble {
			bvtLogger.Printf("global bvt is disable.\n")
			return
		}

		switch *cmdBvt {
		case 0:
			{
				bvtLogger.Printf("bvtAble is disable.\n")
				return
			}
		}
	}
	{
		//检查部署平台
		if !bvt.checkDeploy() {
			bvtLogger.Printf("deploy platform reject bvt.\n")
			return
		}
	}
	{
		bvt.ready = true
		bvtLogger.Printf("bvtVerifier ready:%v\n", bvt.ready)
	}

}
func (bvt *bvtVerifier) checkDeploy() bool {

	//查验部署平台的状态
	bvtLogger.Printf("ready to get DEPLOY_STATUS_API...")
	getDeployStatusApi := func() string {
		deployStatusApi := os.Getenv("DEPLOY_STATUS_API")
		if !strings.HasPrefix(deployStatusApi, "http") {
			deployStatusApi = `http://` + deployStatusApi
		}
		return deployStatusApi + `/dcos/aipaas/app/status`
	}
	getClient := func() *http.Client {
		return &http.Client{}
	}
	getUrl := func() string {
		Url, err := url.Parse(getDeployStatusApi())
		if err != nil {
			bvt.checkErr(bvt.errWrapper(fmt.Sprintf("parase %v failed", getDeployStatusApi()), err))
		}
		params := Url.Query()
		params.Set("idc", bvt.gGroup)
		params.Set("namespace", bvt.namespace)
		params.Set("name", bvt.service)
		Url.RawQuery = params.Encode()
		return Url.String()
	}
	analyzer := func(data []byte) bool {
		bvtLogger.Printf("deploy platform rst:%v\n", string(data))
		m := map[string]interface{}{}
		if json.Unmarshal(data, &m) != nil {
			bvtLogger.Println("can't unmarshal data")
			return false
		}
		bvtLogger.Printf("deploy unpacked:%v\n", m)
		status, statusOk := m["status"]
		if !statusOk {
			return false
		}
		statusFloat64, statusFloat64Ok := status.(float64)
		if !statusFloat64Ok {
			return false
		}
		if statusFloat64 == 2 {
			bvt.bvtAsync = true
			return true
		}
		if statusFloat64 == 0 {
			return true
		}
		return false
	}
	bvtLogger.Printf("deployStatusApi:%s\n", getUrl())
	resp, respErr := getClient().Get(getUrl())
	bvt.checkErr(bvt.errWrapper(fmt.Sprintf("get %v failed", getUrl()), respErr))
	defer resp.Body.Close()
	body, bodyErr := ioutil.ReadAll(resp.Body)
	bvt.checkErr(bvt.errWrapper("read body failed", bodyErr))

	return analyzer(body)
}
func (bvt *bvtVerifier) getConfigure() {
	getRawCfg := func(
		project string,
		group string,
		service string,
		version string,
		cfgFile string,
		cfgUrl string,
	) []byte {
		srvInst, srvInstErr := NewService(CreateCfgOpt(
			utils.WithCfgTick(time.Second),
			utils.WithCfgSessionTimeOut(time.Second),
			utils.WithCfgURL(cfgUrl),
			utils.WithCfgCachePath("finderCache"),
			utils.WithCfgCacheConfig(true),
			utils.WithCfgCacheService(true),
			utils.WithCfgPrj(project),
			utils.WithCfgGroup(group),
			utils.WithCfgService(service),
			utils.WithCfgVersion(version),
		))
		bvt.checkErr(srvInstErr)
		cfgContent, cfgContentErr := srvInst.GetRawCfg(cfgFile)
		bvt.checkErr(cfgContentErr)
		bvtLogger.Printf("cfgFile:%s\n", cfgContent)
		return cfgContent
	}

	rawCfg := getRawCfg(
		bvt.gProject,
		bvt.gGroup,
		bvt.gService,
		bvt.gVersion,
		bvt.gCfgFile,
		bvt.gCompanionUrl,
	)
	var configureErr error
	bvt.gConfigure, configureErr = utils.NewCfgWithBytes(string(rawCfg))
	bvt.checkErr(configureErr)
}

func (bvt *bvtVerifier) checkErr(err error) {
	if err == nil {
		return
	}
	fnName := func() string {
		pc, _, _, _ := runtime.Caller(2)
		return runtime.FuncForPC(pc).Name()
	}()

	bvtLogger.Printf("fn:%v,err:%v\n", fnName, err)
	bvtLogger.Println("failure")

	if strings.Contains(err.Error(), "http.NewRequest") {
		bvtLogger.Println("http.NewRequest")
		os.Exit(ExceptionNewRequest)
	}
	if strings.Contains(err.Error(), "context deadline exceeded") {
		bvtLogger.Println("context deadline exceeded")
		os.Exit(ExceptionDeadline)
	}
	if strings.Contains(err.Error(), "read body failed") {
		bvtLogger.Println("read body failed")
		os.Exit(ExceptionBody)
	}
	if strings.Contains(err.Error(), "analyze bvtRst failed") {
		bvtLogger.Println("analyze bvtRst failed")
		os.Exit(ExceptionBvt)
	}
	os.Exit(ExceptionDef)

}
func (bvt *bvtVerifier) errWrapper(desc string, prev error) error {
	if prev == nil {
		return nil
	}
	return fmt.Errorf("%v;prev err:%v", desc, prev)
}
func (bvt *bvtVerifier) bvtCheck() error {
	if !bvt.ready {
		return nil
	}
	bvtLogger.Println("start bvt...")

	getClient := func() *http.Client {
		return &http.Client{}
	}

	bvtLogger.Println("ready to generate getBody func...")
	getBody := func() (*strings.Reader, error) {
		rawData := map[string]interface{}{
			"id":       bvt.id,
			"async":    bvt.async,
			"callback": bvt.callback,
			"auth_num": bvt.licMax,
			"params": map[string]interface{}{
				"eng_ip":          bvt.engIp,
				"service_address": bvt.serviceAddress,
			},
		}
		bvtLogger.Printf("about to generate request raw data:%#v", rawData)
		marshalRst, marshalRstErr := json.Marshal(rawData)
		if marshalRstErr != nil {
			bvtLogger.Printf("marshal rawData failed:%v", marshalRstErr)
			return nil, fmt.Errorf("marshal failed,marshalRst:%v,marshalRstErr:%v,rawData:%#v", string(marshalRst), marshalRstErr, rawData)
		}
		bvtLogger.Printf("generate request body successfully,rawData:%v,marshalRst:%v", rawData, string(marshalRst))
		return strings.NewReader(string(marshalRst)), nil
	}

	bvtLogger.Println("ready to newRequest...")
	getBodyRst, getBodyRstErr := getBody()
	bvt.checkErr(bvt.errWrapper(fmt.Sprintf("get body failed,err:%v", getBodyRstErr), getBodyRstErr))

	bvtLogger.Printf("about new request,target:%v", bvt.platformAddr)
	req, err := http.NewRequest("POST", bvt.platformAddr, getBodyRst)
	bvt.checkErr(bvt.errWrapper(fmt.Sprintf("http.NewRequest(POST, %v, %v) failed", bvt.platformAddr, getBodyRst), err))

	bvtLogger.Println("ready to set headers...")
	for k, v := range headerKvs {
		req.Header.Set(k, v)
	}

	bvtLogger.Println("ready to send request...")
	resp, err := getClient().Do(req.WithContext(bvt.ctx))
	bvt.checkErr(bvt.errWrapper("http.Do failed", err))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	bvt.checkErr(bvt.errWrapper("read body failed", err))

	type scheme struct {
		Code int `json:"code"`
		Data struct {
			MsID  string `json:"ms_id"`
			MshID int    `json:"msh_id"`
		} `json:"data"`
		Message string `json:"message"`
	}

	bvtLogger.Printf("ready to parse replay body:%s...\n", body)
	bvtRst, bvtRstErr := func() (*scheme, error) {
		schemeInst := &scheme{}
		err := json.Unmarshal(body, schemeInst)
		if err != nil {
			bvtLogger.Printf("unmarshal failed,body:%v,schema:%#v", string(body), schemeInst)
			return nil, fmt.Errorf("unmarshal err:%w", err)
		}
		return schemeInst, nil
	}()
	bvt.checkErr(bvt.errWrapper("analyze bvtRst failed", bvtRstErr))

	return func(bvtRst *scheme) error {
		if bvtRst.Code != 0 {
			bvtLogger.Printf("bvtRst.code exception,code:%v", bvtRst.Code)
			return fmt.Errorf("bvt check fail,rst:%+v", bvtRst)
		}
		return nil
	}(bvtRst)
}

type cfgInstCallback struct{}

func (c *cfgInstCallback) OnConfigFileChanged(con *finder.Config) bool { return true }
func (c *cfgInstCallback) OnError(errInfo finder.ConfigErrInfo)        {}

type Service struct {
	*utils.FindManger

	cfgInstCallback
}

func (s *Service) OnConfigFilesAdded(configs map[string]*finder.Config) bool {
	return true
}

func (s *Service) OnConfigFilesRemoved(configNames []string) bool {
	return true
}

func CreateCfgOpt(opts ...utils.CfgOpt) *utils.CfgOption {
	optInst := &utils.CfgOption{}
	for _, opt := range opts {
		opt(optInst)
	}
	return optInst
}
func NewService(cfgOpt *utils.CfgOption) (*Service, error) {
	srvInst := Service{}
	if err := srvInst.Init(cfgOpt); err != nil {
		return nil, err
	}
	return &srvInst, nil
}

func (s *Service) Init(cfgOpt *utils.CfgOption) (err error) {
	s.FindManger, err = utils.NewFinder(cfgOpt)
	return
}

func (s *Service) GetRawCfg(fileName string) ([]byte, error) {
	cfgInst, cfgInstErr := s.UseCfgAndSub(fileName, s)
	if cfgInstErr != nil {
		return nil, cfgInstErr
	}

	for k, v := range cfgInst {
		if k == fileName {
			return v.File, nil
		}
	}
	return nil, bvtCfgNotFound
}

var bvtVerifierInst *bvtVerifier

func Init(
	gProject string,
	gGroup string,
	gService string,
	gVersion string,
	gCfgFile string,
	gCompanionUrl string,

	timeout time.Duration,
	platformAddr string,
	id string,
	engIp string,
	callback string,
	async bool,
	serviceAddress string,

	licMax int,
	service string,
	namespace string,
) {
	bvtVerifierInst = &bvtVerifier{}

	bvtVerifierInst.init(
		gProject,
		gGroup,
		gService,
		gVersion,
		gCfgFile,
		gCompanionUrl,

		timeout,
		platformAddr,
		id,
		engIp,
		callback,
		async,
		serviceAddress,

		licMax,
		service,
		namespace,
	)
}

func Check() error {

	if bvtVerifierInst == nil {
		bvtLogger.Printf("bvtVerifierInst is disable,ignore...")
		return nil
	}
	if bvtVerifierInst.bvtAsync {
		go func() {
			err := bvtVerifierInst.bvtCheck()
			if err != nil {
				bvtLogger.Println("bvtCheck:", err)
				os.Exit(ExceptionDef)
			}
		}()
		return nil
	}
	return bvtVerifierInst.bvtCheck()
}
