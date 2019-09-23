package xsf

import (
	"reflect"
	"testing"

	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

func Test_getVCpuCfg(t *testing.T) {
	type args struct {
		cfgVer     string
		cfgPrj     string
		cfgGroup   string
		cfgService string
		cfgName    string
		cfgUrl     string
		cfgMode    utils.CfgMode
	}
	tests := []struct {
		name        string
		args        args
		wantVCpuMap map[string]interface{}
		wantErr     bool
	}{
		{
			name: "normal",
			args:
			args{cfgVer: "1.0.0", cfgPrj: "guiderAllService", cfgGroup: "common", cfgService: "xsf", cfgName: "xsf.toml", cfgUrl: "http://10.1.87.69:6868", cfgMode: utils.Centre},
			wantVCpuMap: map[string]interface{}{
				"model1": int64(1),
				"model2": int64(2),
				"model3": int64(3),
			},
			wantErr: false,
		},
	}

	mapEqual := func(m1, m2 map[string]interface{}) bool {
		for k, v := range m1 {
			v1, v1Ok := m2[k]
			if !v1Ok {
				return false
			}
			if !reflect.DeepEqual(v, v1) {
				return false
			}
		}
		for k, v := range m2 {
			v1, v1Ok := m1[k]
			if !v1Ok {
				return false
			}
			if !reflect.DeepEqual(v, v1) {
				return false
			}
		}
		return true
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVCpuMap, err := getVCpuCfg(tt.args.cfgVer, tt.args.cfgPrj, tt.args.cfgGroup, tt.args.cfgService, tt.args.cfgName, tt.args.cfgUrl, tt.args.cfgMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("getVCpuCfg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !mapEqual(gotVCpuMap, tt.wantVCpuMap) {
				t.Errorf("getVCpuCfg() = %v, want %v", gotVCpuMap, tt.wantVCpuMap)
			}
		})
	}
}
