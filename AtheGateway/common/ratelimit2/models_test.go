package ratelimit2

import (
	"reflect"
	"testing"
)

func TestLimitConfigCache_GetConfig(t *testing.T) {
	type fields struct {
		cache        map[string]*Config
		globalConfig *GlobalConfig
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LimitConfigCache{
				cache:        tt.fields.cache,
				globalConfig: tt.fields.globalConfig,
			}
			if got := c.GetConfig(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LimitConfigCache.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Hashstr(s string) int64 {
	var h int64 = 0
	if h==0 && len(s)>0{
		for i:=0;i< len(s);i++{
			h = (31 * h +int64(s[i]))&0x7fffffff
		}
	}


	return h
}

