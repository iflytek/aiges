package schemas

import (
	"fmt"
	"github.com/oliveagle/jsonpath"
	"sort"
	"strings"
)

type mapRule struct {
	jpath         string
	fieldJsonPath *jsonpath.Compiled
	value         map[string]bool
}

type routeItem struct {
	rules        []*mapRule
	subServiceId string
	priority     int
}

func (r *routeItem) exec(req interface{}) bool {
	for _, rule := range r.rules {
		val, err := rule.fieldJsonPath.Lookup(req)
		if err != nil {
			return false
		}
		if !rule.value[String(val)] {
			return false
		}
	}
	return true
}

type routeMap struct {
	ruleItems  []*routeItem
	routerKeys []string
}

func (r *routeMap) getSubServiceId(req interface{}) (serviceId string, routerInfo string) {
	if r == nil {
		return "", ""
	}
	for _, item := range r.ruleItems {
		if item.exec(req) {
			routerValues := make([]string, len(r.routerKeys))
			for i, key := range r.routerKeys {
				for _, rule := range item.rules {
					if strings.HasSuffix(rule.jpath, key) {
						//routerValues[i] = rule.value
						val, _ := rule.fieldJsonPath.Lookup(req)

						routerValues[i] = String(val)
					}
				}
			}
			return item.subServiceId, strings.Join(routerValues, "_")
		}
	}
	return "", ""
}

/*
// 规则例子，满足其中之一即可
{
	"subRoute":[
		{
			"subServiceId":"123456",
			"rules":{
				"$.parameter.domain":"iat",
				"$.parameter.accent":"mandarin"
			}
		}
	]
}
*/
/*
 */
type routeKey struct {
	key string
	jp  *jsonpath.Compiled
}

func newRouteMap(meta Meta) (*routeMap, error) {
	subRoute, ok := meta["subRoute"].([]interface{})
	if !ok {
		// sub route 为空，不返回错误，兼容老的schema
		return nil, nil
	}
	subRoutes := make([]*routeItem, 0, len(subRoute))
	for _, item := range subRoute {
		itm, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("sub Route item must be map")
		}
		serviceId := String(itm["subServiceId"])
		if serviceId == "" {
			return nil, fmt.Errorf("sub serviceId cannot be empty")
		}
		rules, ok := itm["rules"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("rules must be map")
		}
		mrls := make([]*mapRule, 0, len(rules))
		for key, val := range rules {
			jp, err := jsonpath.Compile(key)
			if err != nil {
				return nil, fmt.Errorf("jsonpath compile error! path:%s,err:%w", key, err)
			}
			mrls = append(mrls, &mapRule{
				fieldJsonPath: jp,
				value:         mapOfValue(val),
				jpath:         key,
			})
		}
		subRoutes = append(subRoutes, &routeItem{
			rules:        mrls,
			subServiceId: serviceId,
			priority:     int(Number(itm["priority"])),
		})
	}

	routerKeys, ok := meta["routeKey"].([]interface{})
	rks := make([]string, 0, len(routerKeys))
	if ok {
		for _, key := range routerKeys {
			keyString, _ := key.(string)
			rks = append(rks, keyString)
		}
	}
	// 按照优先级降序排列，优先级高的在前面，会被优先执行
	sort.Slice(subRoutes, func(i, j int) bool {
		return subRoutes[i].priority > subRoutes[j].priority
	})
	return &routeMap{
		ruleItems:  subRoutes,
		routerKeys: rks,
	}, nil
}

func mapOfValue(v interface{}) map[string]bool {
	res := make(map[string]bool)
	for _, s := range Strings(v) {
		res[s] = true
	}
	return res
}
