package lb_client

import (
	"encoding/json"
	"errors"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/xfyun/redisgo"
	"github.com/cihub/seelog"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	Err_Zklist_Is_Null        = errors.New("zklist is null")
	Err_RouterType_Is_Null    = errors.New("routerType is null")
	Err_SubRouterType_Is_Null = errors.New("subrouterType is null")
	Err_Redishost_Is_Null     = errors.New("redishost is null")
)

var (
	default_maxIdle       = 10                 //默认redis的最大空闲连接数
	default_idleTimeOut   = time.Second * 3600 //默认redis的空闲连接超时时间
	default_zkConnTimeout = time.Second * 10   //默认zk连接超时时间
)

type LbUtil struct {
	LbOpt          LbCfgOption
	ZkConn         *zk.Conn       //zk连接实例
	RedisInst      *redisgo.Redis //redis连接实例
	RouterTypeAbs  string         //路由类型绝对路径
	SubRouterTypes []string       //子路由类型列表(后续一个进程可能启动多个子路由类型)(如：["iat_gray","iat_hefei"])
}

func (lu *LbUtil) Init(o ...LbCfgOpt) (err error) {
	var lbOpt LbCfgOption
	for _, opt := range o {
		opt(&lbOpt)
	}
	lu.LbOpt = lbOpt

	var (
		maxIdle       int           = lbOpt.MaxIdle
		idleTimeOut   time.Duration = lbOpt.IdleTimeOut
		zkConnTimeOut time.Duration = lbOpt.ZkConnTimeOut
	)

	if lbOpt.MaxIdle <= 0 {
		maxIdle = default_maxIdle
	}

	if lbOpt.IdleTimeOut <= 0 {
		idleTimeOut = default_idleTimeOut
	}

	if lbOpt.ZkConnTimeOut <= 0 {
		zkConnTimeOut = default_zkConnTimeout
	}

	//初始化zk相关操作
	if err = lu.initZk(lbOpt.Root, lbOpt.RouterType, lbOpt.SubRouterTypes,
		lbOpt.ZkServerList, zkConnTimeOut); err != nil {
		return
	}

	//初始化redis相关操作
	if err = lu.initRedis(lbOpt.RedisHost, lbOpt.RedisPasswd, lbOpt.MaxActive,
		maxIdle, lbOpt.Db, idleTimeOut); err != nil {
		return
	}

	return
}

//redis的相关初始化
func (lu *LbUtil) initZk(root, routerType string, subRouterTypes []string, zkList []string, zkConnTimeout time.Duration) (err error) {
	if len(zkList) == 0 {
		err = Err_Zklist_Is_Null
		return
	}
	if len(routerType) == 0 {
		err = Err_RouterType_Is_Null
		return
	}

	if len(subRouterTypes) == 0 {
		err = Err_SubRouterType_Is_Null
		return
	}

	//初始化zk相关配置
	lu.LbOpt.RouterType = strings.Replace(routerType, "/", "", -1)
	lu.RouterTypeAbs = path.Join("/"+root, lu.LbOpt.RouterType)
	for _, subRouterType := range subRouterTypes {
		lu.SubRouterTypes = append(lu.SubRouterTypes, strings.Replace(subRouterType, "/", "", -1))
	}

	//连接zk
	zkConn, err := lu.connectZk()
	if err != nil {
		return err
	}
	lu.ZkConn = zkConn
	return
}

//redis的相关初始化
func (lu *LbUtil) initRedis(redisHost, redisPasswd string, maxActive, maxIdle, db int, idleTimeOut time.Duration) (err error) {
	if len(redisHost) == 0 {
		err = Err_Redishost_Is_Null
		return
	}

	//初始化redis相关配置
	lu.RedisInst, err = redisgo.NewRedisInst(
		redisgo.WithRedisHost(redisHost),
		redisgo.WithRedisPwd(redisPasswd),
		redisgo.WithMaxIdle(maxIdle),
		redisgo.WithDb(db),
		redisgo.WithIdleTimeout(idleTimeOut),
	)
	if err != nil {
		return
	}
	return
}

//创建引擎存活节点
func (lu *LbUtil) createAliveNode(svc string, subRouterType string, data []byte) (err error) {
	_, err = lu.CreateR(lu.RouterTypeAbs, nil, 0, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	//子路由节点存储负载策略的数据
	svcAlivePathAbs := path.Join(lu.RouterTypeAbs, subRouterType)
	subRouterDirExists, _, err := lu.ZkConn.Exists(svcAlivePathAbs)
	if err != nil {
		return
	}
	if !subRouterDirExists {
		var subRouterParam = make(map[string]string)
		subRouterParam["lb_strategy"] = strconv.Itoa(lu.LbOpt.LbStrategy)

		paramJson, err := json.Marshal(subRouterParam)
		if err != nil {
			return err
		}
		_, err = lu.ZkConn.Create(svcAlivePathAbs, paramJson, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}

	//在子路由节点目录下创建引擎节点
	svcAliveAddr := path.Join(svcAlivePathAbs, svc)
	aliveNodeExists, _, err := lu.ZkConn.Exists(svcAliveAddr)
	if err != nil {
		return err
	}
	if aliveNodeExists {
		err = lu.ZkConn.Delete(svcAliveAddr, -1)
		if err != nil {
			return err
		}
	}

	//创建zk临时节点
	_, err = lu.ZkConn.Create(svcAliveAddr, data, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}

	return err
}

//删除存活节点
func (lu *LbUtil) deleteAliveNode(svc string, subRouterType string) (err error) {
	svcAliveAddr := path.Join(lu.RouterTypeAbs, subRouterType, svc)
	exists, _, err := lu.ZkConn.Exists(svcAliveAddr)
	if err != nil {
		return err
	}
	if exists {
		err = lu.ZkConn.Delete(svcAliveAddr, -1)
		if err != nil {
			return err
		}
	}
	return err
}

//循环创建目录
func (lu *LbUtil) CreateR(dir string, data []byte, flags int32, acl []zk.ACL) (string, error) {
	pathTmp := strings.TrimPrefix(dir, "/")
	pathSlice := strings.Split(pathTmp, "/")
	for i := 0; i < len(pathSlice); i++ {
		pathCur := strings.Join(pathSlice[0:i+1], "/")
		pathCur = "/" + pathCur
		exist, _, err := lu.ZkConn.Exists(pathCur)
		if err != nil {
			return "", err
		}
		if !exist {
			_, err = lu.ZkConn.Create(pathCur, data, flags, acl)
			if err != nil {
				return "", err
			}
		}
	}
	return "", nil
}

//连接zk
func (lu *LbUtil) connectZk() (zkConn *zk.Conn, err error) {
	var zkConnTimeOut time.Duration = lu.LbOpt.ZkConnTimeOut

	if lu.LbOpt.ZkConnTimeOut <= 0 {
		zkConnTimeOut = default_zkConnTimeout
	}

	zkConn, _, err = zk.Connect(lu.LbOpt.ZkServerList, zkConnTimeOut)
	if err != nil {
		seelog.Error("zk Connect error: %s", err.Error())
		return nil, err
	}
	return zkConn, err
}

//zk健康检查
func (lu *LbUtil) zkHealthCheck() (statusOk bool) {
	defer func() {
		if errErr := recover(); errErr != nil {
			seelog.Error("occur panic,err is:", errErr)
		}
	}()
	statusOk = false
	state := lu.ZkConn.State()
	//连接正常,并且有心跳
	if state == zk.StateHasSession {
		statusOk = true
		//断开连接
	} else if state == zk.StateDisconnected {
		seelog.Critical("zk Disconnected")
		//重新连接zk
		zkConn, err := lu.connectZk()
		if err != nil {
			return
		}
		lu.ZkConn = zkConn
		statusOk = true
		//会话过期
	} else if state == zk.StateExpired {
		seelog.Critical("zk session Expired")
		zkConn, err := lu.connectZk()
		if err != nil {
			return
		}
		lu.ZkConn = zkConn
		statusOk = true
		//zk自身在尝试不停重连
	} else if state == zk.StateConnecting {
		seelog.Critical("zk Connecting")
		//todo
		//权限校验失败
	} else if state == zk.StateAuthFailed {
		seelog.Critical("zk AuthFailed")
		//todo
	}
	return statusOk
}
