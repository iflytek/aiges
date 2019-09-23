package curator

import (
	"github.com/cooleric/go-zookeeper/zk"
)

var (
	OPEN_ACL_UNSAFE = zk.WorldACL(zk.PermAll)
	CREATOR_ALL_ACL = zk.AuthACL(zk.PermAll)
	READ_ACL_UNSAFE = zk.WorldACL(zk.PermRead)
)

type ACLProvider interface {
	// Return the ACL list to use by default
	GetDefaultAcl() []zk.ACL

	// Return the ACL list to use for the given path
	GetAclForPath(path string) []zk.ACL
}

type defaultACLProvider struct {
	defaultAcls []zk.ACL
}

func (p *defaultACLProvider) GetDefaultAcl() []zk.ACL {
	return p.defaultAcls
}

func (p *defaultACLProvider) GetAclForPath(path string) []zk.ACL {
	return p.defaultAcls
}

func NewDefaultACLProvider() ACLProvider {
	return &defaultACLProvider{OPEN_ACL_UNSAFE}
}

type acling struct {
	aclList     []zk.ACL
	aclProvider ACLProvider
}

func (a *acling) getAclList(path string) []zk.ACL {
	if a.aclList != nil {
		return a.aclList
	}

	if len(path) > 0 {
		if acls := a.aclProvider.GetAclForPath(path); acls != nil {
			return acls
		}
	}

	return a.aclProvider.GetDefaultAcl()
}

type getACLBuilder struct {
	client        *curatorFramework
	backgrounding backgrounding
	stat          *zk.Stat
}

func (b *getACLBuilder) ForPath(givenPath string) ([]zk.ACL, error) {
	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath, givenPath)

		return nil, nil
	} else {
		return b.pathInForeground(adjustedPath)
	}
}

func (b *getACLBuilder) pathInBackground(path string, givenPath string) {
	tracer := b.client.ZookeeperClient().StartTracer("getACLBuilder.pathInBackground")

	defer tracer.Commit()

	acls, err := b.pathInForeground(path)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: GET_ACL,
			err:       err,
			path:      b.client.unfixForNamespace(path),
			acls:      acls,
			stat:      b.stat,
			context:   b.backgrounding.context,
		}

		if err != nil {
			event.path = givenPath
		}

		event.name = GetNodeFromPath(event.path)

		b.backgrounding.callback(b.client, event)
	}
}

func (b *getACLBuilder) pathInForeground(path string) ([]zk.ACL, error) {
	zkClient := b.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			acls, stat, err := conn.GetACL(path)

			if stat != nil {
				if b.stat != nil {
					*b.stat = *stat
				} else {
					b.stat = stat
				}
			}

			return acls, err
		}
	})

	acls, _ := result.([]zk.ACL)

	return acls, err
}

func (b *getACLBuilder) StoringStatIn(stat *zk.Stat) GetACLBuilder {
	b.stat = stat

	return b
}

func (b *getACLBuilder) InBackground() GetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *getACLBuilder) InBackgroundWithContext(context interface{}) GetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *getACLBuilder) InBackgroundWithCallback(callback BackgroundCallback) GetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *getACLBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) GetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}

type setACLBuilder struct {
	client        *curatorFramework
	backgrounding backgrounding
	acling        acling
	version       int32
}

func (b *setACLBuilder) ForPath(givenPath string) (*zk.Stat, error) {
	adjustedPath := b.client.fixForNamespace(givenPath, false)

	if b.backgrounding.inBackground {
		go b.pathInBackground(adjustedPath, givenPath)

		return nil, nil
	} else {
		return b.pathInForeground(adjustedPath)
	}
}

func (b *setACLBuilder) pathInBackground(path string, givenPath string) {
	tracer := b.client.ZookeeperClient().StartTracer("setACLBuilder.pathInBackground")

	defer tracer.Commit()

	stat, err := b.pathInForeground(path)

	if b.backgrounding.callback != nil {
		event := &curatorEvent{
			eventType: SET_ACL,
			err:       err,
			path:      b.client.unfixForNamespace(path),
			acls:      b.acling.getAclList(path),
			stat:      stat,
			context:   b.backgrounding.context,
		}

		if err != nil {
			event.path = givenPath
		}

		event.name = GetNodeFromPath(event.path)

		b.backgrounding.callback(b.client, event)
	}
}

func (b *setACLBuilder) pathInForeground(path string) (*zk.Stat, error) {
	zkClient := b.client.ZookeeperClient()

	result, err := zkClient.NewRetryLoop().CallWithRetry(func() (interface{}, error) {
		if conn, err := zkClient.Conn(); err != nil {
			return nil, err
		} else {
			return conn.SetACL(path, b.acling.getAclList(path), b.version)
		}
	})

	stat, _ := result.(*zk.Stat)

	return stat, err
}

func (b *setACLBuilder) WithACL(acls ...zk.ACL) SetACLBuilder {
	b.acling.aclList = acls

	return b
}

func (b *setACLBuilder) WithVersion(version int32) SetACLBuilder {
	b.version = version

	return b
}

func (b *setACLBuilder) InBackground() SetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true}

	return b
}

func (b *setACLBuilder) InBackgroundWithContext(context interface{}) SetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context}

	return b
}

func (b *setACLBuilder) InBackgroundWithCallback(callback BackgroundCallback) SetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true, callback: callback}

	return b
}

func (b *setACLBuilder) InBackgroundWithCallbackAndContext(callback BackgroundCallback, context interface{}) SetACLBuilder {
	b.backgrounding = backgrounding{inBackground: true, context: context, callback: callback}

	return b
}
