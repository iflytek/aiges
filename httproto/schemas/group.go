package schemas

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type AISchemaGroup struct {
	lock    sync.RWMutex
	schemas map[string]*AISchema
}

func NewAISchemaGroup() *AISchemaGroup {
	return &AISchemaGroup{
		lock:    sync.RWMutex{},
		schemas: make(map[string]*AISchema, 10),
	}
}

func (g *AISchemaGroup) Update(name string, schema *AISchema) error {
	if err := schema.Init(); err != nil {
		return fmt.Errorf("init schema error:%s", err.Error())
	}
	g.lock.Lock()
	g.schemas[name] = schema
	g.lock.Unlock()
	return nil
}

func (g *AISchemaGroup) Get(name string) *AISchema {
	g.lock.RLock()
	sc := g.schemas[name]
	g.lock.RUnlock()
	return sc
}

func (g *AISchemaGroup) Range(f func(key string, schema *AISchema) bool) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	for k, v := range g.schemas {
		if !f(k, v) {
			return
		}
	}
}

var sgm = &schemaGroupManager{
	routeIndex:     NewAISchemaGroup(),
	stdIndex:       NewAISchemaGroup(),
	servieIdIndex:  NewAISchemaGroup(),
	companionRoute: NewAISchemaGroup(),
	categoryRoute:  NewAISchemaGroup(),
}

func LoadSchema(schema *AISchema) error {
	if schema.Meta == nil {
		return fmt.Errorf("load schema error, meta is nil")
	}
	routesMap, err := newRouteMap(schema.Meta)
	if err != nil {
		return fmt.Errorf("load schema error,sub route load error:%w", err)
	}
	schema.headerSchema, _ = schema.SchemaOutput.Get("properties").Get("header").Get("properties").Interface().(map[string]interface{})
	schema.subRouteMap = routesMap
	if err := sgm.Update(schema.Meta.GetHost(), schema.Meta.GetRoute(), schema); err != nil {
		return err
	}
	return nil
}

//
func LoadAISchema(data []byte) error {
	schema := &AISchema{}
	if err := json.Unmarshal(data, schema); err != nil {
		return err
	}
	return LoadSchema(schema)
}

//
func GetSchema(host, route string) *AISchema {
	return sgm.Get(host, route)
}

//
func GetSchemaByServiceId(serviceId string, cloudId string) *AISchema {
	return sgm.GetByServiceIdAndCloudId(serviceId, cloudId)
}

func GetCompanionSchema(path string, cloudId string) *AISchema {
	return sgm.getCompanionRoute(path, cloudId)
}

func GetCategorySchema(path string, cloudId string) *AISchema {
	return sgm.getCategorySchema(path, cloudId)
}

//
type schemaGroupManager struct {
	routeIndex     *AISchemaGroup // host 为空时
	stdIndex       *AISchemaGroup // host 不为空时
	servieIdIndex  *AISchemaGroup // 根据serviceId 建立索引
	companionRoute *AISchemaGroup
	categoryRoute  *AISchemaGroup
}

func (s *schemaGroupManager) GetCompanionSchemaGroup(path string) *AISchema {
	return s.companionRoute.Get(path)
}

func generateKey(host, path string) string {
	return host + path
}

func AssembleKey(cloudId string, serviceId string) string {
	b := strings.Builder{}
	b.Grow(len(cloudId) + len(serviceId) + 1)
	b.WriteString(cloudId)
	b.WriteString("_")
	b.WriteString(serviceId)
	return b.String()
}

func assembleRouterKey(path, cloudId string) string {
	return path + "." + cloudId
}

func (m *schemaGroupManager) addCompanionRoute(cr string, cloudId string, s *AISchema) error {
	if err := m.companionRoute.Update(assembleRouterKey(cr, cloudId), s); err != nil {
		return err
	}
	return nil
}

func (m *schemaGroupManager) getCompanionRoute(cr string, cloudId string) *AISchema {
	return m.companionRoute.Get(assembleRouterKey(cr, cloudId))
}

func (m *schemaGroupManager) addCategorySchema(route []string, cloudId string, sc *AISchema) error {
	for _, s := range route {
		if err := m.categoryRoute.Update(assembleRouterKey(s, cloudId), sc); err != nil {
			return err
		}
	}
	return nil
}

func (m *schemaGroupManager) getCategorySchema(path string, cloudId string) *AISchema {
	return m.categoryRoute.Get(assembleRouterKey(path, cloudId))
}

func (m *schemaGroupManager) Update(hosts []string, paths []string, s *AISchema) error {
	cloudId := s.Meta.GetCloudId()
	if cloudId == "0" {
		cloudId = ""
	}

	if cm := s.Meta.GetCompanion(); cm != "" {
		return m.addCompanionRoute(cm, cloudId, s)
	}

	//if s.subRouteMap != nil {
	//	return m.addCategorySchema(s.Meta.GetRoute(), cloudId, s)
	//}

	if err := m.servieIdIndex.Update(AssembleKey(cloudId, s.Meta.GetServiceId()), s); err != nil {
		return err
	}
	if len(hosts) == 0 {
		for _, path := range paths {
			err := m.routeIndex.Update(path, s)
			if err != nil {
				return err
			}
		}
		return nil
	}
	for _, host := range hosts {
		for _, path := range paths {
			err := m.stdIndex.Update(generateKey(host, path), s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *schemaGroupManager) Get(host, path string) *AISchema {
	if s := m.stdIndex.Get(generateKey(host, path)); s != nil {
		return s
	}
	return m.routeIndex.Get(path)
}

func (m *schemaGroupManager) GetByServiceIdAndCloudId(serviceId string, cloudId string) *AISchema {
	return m.servieIdIndex.Get(AssembleKey(cloudId, serviceId))
}
