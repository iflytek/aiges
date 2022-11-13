package schemas

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

func String(v interface{}) string {
	switch v := v.(type) {
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case nil:
		return ""
	case []interface{}:
		arr := make([]string, len(v))
		for i, e := range v {
			arr[i] = String(e)
		}
		return strings.Join(arr, ";")
	default:
		return fmt.Sprintf("%v", v)
	} // [ 1 2 3 4 5 ]
}

func Strings(v interface{}) []string {
	switch vv := v.(type) {
	case []interface{}:
		arr := make([]string, len(vv))
		for i, i2 := range vv {
			arr[i] = String(i2)
		}
		return arr
	case []string:
		return vv
	default:
		return []string{String(v)}
	}
}

func Number(i interface{}) float64 {
	switch i.(type) {
	case float64:
		return i.(float64)
	case int:
		return float64(i.(int))
	case int32:
		return float64(i.(int32))
	case int64:
		return float64(i.(int64))
	case string:
		n, _ := strconv.ParseFloat(i.(string), 64)
		return float64(n)
	}
	return 0
}

func Bool(v interface{}) bool {
	switch v.(type) {
	case bool:
		return v.(bool)
	case string:
		return v.(string) == "true"
	case float64:
		return int(v.(float64)) > 0
	}
	if v != nil {
		return true
	}
	return false
}

type Queue struct {
}

type HandleCache struct {
	handles map[int]interface{}
	maxSize int32
	index   int32
	lock    sync.RWMutex
}

func (h *HandleCache) idleInst() int {
	h.lock.RLock()
	defer h.lock.RUnlock()
	for i := 0; i < int(h.maxSize); i++ {
		hs := (atomic.AddInt32(&h.index, 1) % h.maxSize)
		if h.handles[int(hs)] == nil {
			if hs == h.maxSize {
				h.index = 0
			}
			return int(hs)
		}
	}
	return -1
}

func (h *HandleCache) NewHandle(bind interface{}) int {
	h.lock.Lock()
	handle := h.idleInst()
	h.handles[handle] = bind
	h.lock.Unlock()
	return handle
}

func (h *HandleCache) GetHandle(handle int) interface{} {
	h.lock.RLock()
	ha := h.handles[handle]
	h.lock.RUnlock()
	return ha
}

func (h *HandleCache) ReleaseHandle(handle int) {
	h.lock.Lock()
	delete(h.handles, handle)
	h.lock.Unlock()
}

func NewInst() (handle int) {

	return 0
}

func In(handle int) {

}
