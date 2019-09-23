package curator

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type namespaceImpl struct {
	client     *curatorFramework
	namespace  string
	ensurePath EnsurePath
}

func newNamespace(client *curatorFramework, namespace string) *namespaceImpl {
	n := &namespaceImpl{
		client:    client,
		namespace: namespace,
	}

	if len(namespace) > 0 {
		if err := ValidatePath("/" + namespace); err != nil {
			client.logError(fmt.Errorf("Invalid namespace: %s, %s", namespace, err))

			return newNamespace(client, "")
		}

		n.ensurePath = NewEnsurePath(JoinPath("/", namespace))
	}

	return n
}

// Apply the namespace to the given path
func FixForNamespace(namespace, path string, isSequential bool) (string, error) {
	if len(namespace) > 0 {
		return JoinPath(namespace, path), nil
	}

	return path, nil
}

func (n *namespaceImpl) fixForNamespace(path string, isSequential bool) string {
	if n.ensurePath != nil {
		n.ensurePath.Ensure(n.client.ZookeeperClient())
	}

	s, _ := FixForNamespace(n.namespace, path, isSequential)

	return s
}

func (n *namespaceImpl) unfixForNamespace(path string) string {
	if len(n.namespace) > 0 && len(path) > 0 {
		prefix := JoinPath(n.namespace)

		if strings.HasPrefix(path, prefix) {
			if len(prefix) < len(path) {
				return path[len(prefix):]
			} else {
				return PATH_SEPARATOR
			}
		}
	}

	return path
}

type namespaceFacade struct {
	curatorFramework
}

func newNamespaceFacade(client *curatorFramework, namespace string) *namespaceFacade {
	facade := &namespaceFacade{
		curatorFramework: *client,
	}

	facade.namespace = newNamespace(client, namespace)
	facade.fixForNamespace = facade.namespace.fixForNamespace
	facade.unfixForNamespace = facade.namespace.unfixForNamespace

	return facade
}

func (f *namespaceFacade) Start() error {
	return errors.New("the requested operation is not supported")
}

func (f *namespaceFacade) Close() error {
	return errors.New("the requested operation is not supported")
}

func (f *namespaceFacade) CuratorListenable() CuratorListenable {
	f.logError(errors.New("CuratorListenable() is only available from a non-namespaced CuratorFramework instance"))

	return f.curatorFramework.listeners
}

func (f *namespaceFacade) Namespace() string {
	return f.namespace.namespace
}

type namespaceFacadeCache struct {
	client *curatorFramework
	cache  map[string]*namespaceFacade
	lock   sync.Mutex
}

func newNamespaceFacadeCache(client *curatorFramework) *namespaceFacadeCache {
	return &namespaceFacadeCache{
		client: client,
		cache:  make(map[string]*namespaceFacade),
	}
}

func (c *namespaceFacadeCache) Get(namespace string) *namespaceFacade {
	c.lock.Lock()
	defer c.lock.Unlock()

	if facade, exists := c.cache[namespace]; exists {
		return facade
	}

	facade := newNamespaceFacade(c.client, namespace)

	c.cache[namespace] = facade

	return facade
}
