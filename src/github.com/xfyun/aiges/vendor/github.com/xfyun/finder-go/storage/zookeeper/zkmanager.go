package zookeeper

import (
	"strconv"
	"strings"
	"sync"
	"time"

	errors "github.com/xfyun/finder-go/errors"
	"github.com/xfyun/finder-go/log"
	"github.com/xfyun/finder-go/storage/common"
	"github.com/cooleric/go-zookeeper/zk"
)

//zk超时时间设置
const zk_connection_timeout = 5

type ZkManager struct {
	conn   *zk.Conn
	params map[string]string
	exit   chan bool
	//记录了临时路径
	tempPaths sync.Map
	//记录了path对应的Watcher
	dataWatcher sync.Map
}

func NewZkManager(params map[string]string) (*ZkManager, error) {
	zm := &ZkManager{
		tempPaths:   sync.Map{},
		params:      params,
		dataWatcher: sync.Map{},
	}
	return zm, nil
}

func (zm *ZkManager) GetServerAddr() string {
	return zm.params["servers"]
}

func (zm *ZkManager) GetTempPaths() sync.Map {
	return zm.tempPaths
}
func (zm *ZkManager) SetTempPaths(tempPath sync.Map) {
	zm.tempPaths = tempPath
}
func (zm *ZkManager) GetZkNodePath() (string, error) {
	if path, ok := zm.params["zk_node_path"]; ok {
		return path, nil
	} else {
		return "", errors.NewFinderError(errors.ZkInfoMissZkNodePath)
	}
}

func (zm *ZkManager) Init() error {

	//必要参数判断
	serverStr, exist := zm.params["servers"]
	if !exist || len(serverStr) == 0 {
		return errors.NewFinderError(errors.ZkParamsMissServers)
	}

	servers := strings.Split(serverStr, ",")
	timeout, exist := zm.params["session_timeout"]
	if !exist || len(timeout) == 0 {
		return errors.NewFinderError(errors.ZkParamsMissSessionTimeout)
	}
	sessionTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		return err
	}
	//新建zookeeper连接
	conn, _, err := zk.Connect(servers, time.Duration(sessionTimeout)*time.Millisecond, zk.WithEventCallback(zm.eventCallback), zk.WithLogger(log.Log))
	if err != nil {
		return err
	}
	_, _, err = conn.Exists(zm.params["zk_node_path"])
	if err != nil {
		return err
	}
	zm.conn = conn
	return nil
}

func (zm *ZkManager) eventCallback(e zk.Event) {
	switch e.Type {
	case zk.EventSession:
		switch e.State {
		case zk.StateDisconnected:
			return
		case zk.StateConnecting:
			return
		case zk.StateConnected:
			return
		case zk.StateHasSession:
			go zm.RecoverTempPaths()
			return
		case zk.StateExpired:
			return
		case zk.StateAuthFailed:
			return
		case zk.StateConnectedReadOnly:
			return
		case zk.StateSaslAuthenticated:
			return
		case zk.StateUnknown:
			return
		}
		return
	}
}

/**
 * 在恢复会话的时候进行调用，用来恢复临时路径
 */
func (zm *ZkManager) RecoverTempPaths() {
	defer recoverFunc()
	var err error
	zm.tempPaths.Range(func(key, value interface{}) bool {
		log.Log.Debugf("recover temp path %v ", key.(string))
		if value == nil {
			//如果path上的数据为空。则直接设置path
			for {
				err = zm.SetTempPath(key.(string))
				//TODO 如果在恢复临时路径的时候，挂了
				if err != nil {
					log.Log.Errorf("caught an error:zm.SetTempPath in recoverTempPaths: %v ", err)
					continue
				}
				break
			}
		} else {
			//如果path上的数据不为空。则直接设置path和对应的数据
			for {
				err = zm.SetTempPathWithData(key.(string), value.([]byte))
				if err != nil {
					log.Log.Errorf("caught an error:zm.SetTempPathWithData in recoverTempPaths: %s", err)
					continue
				}
				break
			}
		}

		return true
	})
}

func (zm *ZkManager) Destroy() error {
	log.Log.Debugf("exit send.")
	zm.params = nil

	zm.conn.Close()
	log.Log.Debugf("close end.")
	go func() {
		log.Log.Debugf("send exit sigterm.")
		zm.exit <- true
	}()

	log.Log.Debugf("destroied")
	return nil
}

func (zm *ZkManager) GetData(path string) ([]byte, error) {
	//获取节点数据
	data, _, err := zm.conn.Get(path)
	return data, err
}

func (zm *ZkManager) GetDataWithWatchV2(path string, callback common.ChangedCallback) ([]byte, error) {

	//获取数据，并注册Watch
	data, _, event, err := zm.conn.GetW(path)
	if err != nil {
		return nil, err
	}

	//监听是否有watch到达
	go watchEvent(zm, event, callback)

	return data, err
}

func (zm *ZkManager) GetDataWithWatchV3(path string, callback common.ChangedCallback) ([]byte, error) {
	data, _, event, err := zm.conn.GetW(path)
	if err != nil {
		log.Log.Infof("get data with watchV3 err  %v", err)
		return data, err
	}
	go func() {
		for {
			select {
			case e := <-event:
				switch e.Type {
				case zk.EventNodeDataChanged, zk.EventNodeCreated:
					for {
						data, _, event, err = zm.conn.GetW(path)
						if err != nil {
							log.Log.Errorf("get data error:%v", err)
							time.Sleep(1 * time.Second)
							continue
						}
						break
					}
					callback.DataChangedCallback(e.Path, getNodeFromPath(e.Path), data)
				case zk.EventNodeDeleted:
					log.Log.Errorf("node is deleted stop watch: path=%v", e.Path)
					return
				}
			case exit, ok := <-zm.exit:
				if exit || !ok {
					log.Log.Infof("zk exited")
					return
				}

			}

		}

	}()
	return nil, err
}

func watchEvent(zm *ZkManager, event <-chan zk.Event, callback common.ChangedCallback) {
	select {
	case e, ok := <-event:
		if !ok {
			log.Log.Infof("<-event; !ok")
			return
		}
		defer recoverFunc()
		log.Log.Debugf("recv event : %v", e)
		callback.Process(e.Path, getNodeFromPath(e.Path))
		break
	case exit, ok := <-zm.exit:
		if !ok {
			log.Log.Infof("<-exit; !ok")
			return
		}
		if exit {
			log.Log.Infof("received exit sigterm.")
			return
		}
	}
}
func (zm *ZkManager) GetDataWithWatch(path string, callback common.ChangedCallback) ([]byte, error) {

	data, _, event, err := zm.conn.GetW(path)
	if err != nil {
		log.Log.Infof("get data with watch err  %v", err)
	}
	//返回的event
	go func(zm *ZkManager, p string, event <-chan zk.Event) {
		for {
			select {
			case e, ok := <-event:
				if !ok {
					log.Log.Errorf("handler event --> path: %v , err:  %v", path, e)
					return
				}
				log.Log.Infof("recv event %v", e)
				if e.Type == zk.EventNodeDeleted {
					return
				}
				if e.State != zk.StateConnected {
					return
				}
				log.Log.Infof("handler event %v", e)
				var retryCount int32
				for {
					// 这个地方有问题，如果节点被删除的话，会成为死循环，修改为尝试三次
					data, _, event, err = zm.conn.GetW(path)
					if err == zk.ErrNoNode {
						log.Log.Errorf("node deleted , stop watch: paths:%s,err:%v", path, err)
						return
					}
					if err != nil {
						log.Log.Debugf("get data with watch , path: %v ,err: %v", path, err)
						retryCount++
						if retryCount > 30 {
							time.Sleep(time.Duration(retryCount) * time.Second)
							break
						}
						continue
					} else {
						callback.DataChangedCallback(e.Path, getNodeFromPath(e.Path), data)
					}

					break
				}
			case exit, ok := <-zm.exit:
				if !ok {
					log.Log.Infof("<-exit; !ok")
					return
				}
				if exit {
					log.Log.Infof("received exit sigterm.")
					return
				}
			}
		}
	}(zm, path, event)

	return data, err
}

func (zm *ZkManager) GetChildren(path string) ([]string, error) {
	nodes, _, err := zm.conn.Children(path)
	return nodes, err
}

func recoverFunc() {
	if err := recover(); err != nil {
		log.Log.Errorf("recover ： %v", err)
	}
}
func (zm *ZkManager) GetChildrenWithWatch(path string, callback common.ChangedCallback) ([]string, error) {
	data, _, event, err := zm.conn.ChildrenW(path)
	if err != nil {
		return nil, err
	}

	go func(zm *ZkManager, p string, event <-chan zk.Event) {
		defer recoverFunc()
		for {
			select {
			case e, ok := <-event:
				if !ok {
					log.Log.Infof("[ GetChildrenWithWatch ]  <-event; !ok")
					return
				}
				log.Log.Debugf("recv event ：[ GetChildrenWithWatch ]  %v", event)
				if e.State != zk.StateConnected {
					log.Log.Debugf("[ GetChildrenWithWatch ]  e.State != zk.StateConnected %v", "")
				}
				var retryCount int
				for {
					data, _, event, err = zm.conn.ChildrenW(path)
					retryCount++
					//if retryCount > 5 {
					//
					//	break
					//}
					// 一直watch 重试不要退出，否则会出现订阅不上的情况
					if err == zk.ErrNoNode {
						log.Log.Errorf("node deleted , stop watch dir: paths:%s,err:%v", path, err)
						return // 没有node，直接return
					}
					if err != nil {
						time.Sleep(time.Duration(retryCount%30) * time.Second)
						log.Log.Errorf("[ GetChildrenWithWatch ] retry get children err: %v, path: %v", err, path)
						continue
					} else {
						callback.ChildrenChangedCallback(e.Path, getNodeFromPath(e.Path), data)
					}
					break
				}
			case exit, ok := <-zm.exit:
				if !ok {
					log.Log.Infof("<-exit; !ok")
					return
				}
				if exit {
					log.Log.Infof("received exit sigterm.")
					return
				}
			}
		}
	}(zm, path, event)

	return data, err
}

func (zm *ZkManager) SetPath(path string) error {
	return zm.SetPathWithData(path, []byte{})
}

func (zm *ZkManager) CheckExists(path string) (bool, error) {
	exists, _, err := zm.conn.Exists(path)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (zm *ZkManager) SetPathWithData(path string, data []byte) error {
	if data == nil {
		return errors.NewFinderError(errors.ZkDataCanotNil)
	}
	_, err := zm.conn.Create(path, data, PERSISTENT, zk.WorldACL(zk.PermAll))
	if err == zk.ErrNoNode {
		err = makeDirs(zm.conn, path, false)
		if err != nil {
			return err
		}
		_, err = zm.conn.Create(path, data, PERSISTENT, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	}

	return nil
}

func (zm *ZkManager) SetTempPath(path string) error {
	err := zm.SetTempPathWithData(path, []byte{})
	if err == nil {
		zm.tempPaths.Store(path, nil)
	}
	return err
}

func (zm *ZkManager) SetTempPathWithData(path string, data []byte) error {
	if data == nil {
		return errors.NewFinderError(errors.ZkDataCanotNil)
	}
	_, err := zm.conn.Create(path, data, EPHEMERAL, zk.WorldACL(zk.PermAll))
	if err == zk.ErrNoNode {
		err = makeDirs(zm.conn, path, false)
		if err != nil {
			return err
		}
		_, err = zm.conn.Create(path, data, EPHEMERAL, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	} else if err == zk.ErrNodeExists {
		err = zm.RemoveInRecursive(path)
		if err != nil {
			return err
		}
		_, err = zm.conn.Create(path, data, EPHEMERAL, zk.WorldACL(zk.PermAll))
		if err != nil {
			return err
		}
	} else {
		return err
	}

	zm.tempPaths.Store(path, data)

	return nil
}

func (zm *ZkManager) SetData(path string, value []byte) error {
	if value == nil {
		return errors.NewFinderError(errors.ZkDataCanotNil)
	}
	_, err := zm.conn.Set(path, value, DEFAULT_VERSION)
	return err
}

func (zm *ZkManager) Remove(path string) error {
	return zm.conn.Delete(path, DEFAULT_VERSION)
}

func (zm *ZkManager) RemoveInRecursive(path string) error {
	return recursiveDelete(zm.conn, path, true)
}

func (zm *ZkManager) UnWatch(path string) error {
	return nil
}
