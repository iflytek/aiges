package conf

import (
	"git.xfyun.cn/AIaaS/xsf/utils"
	"testing"
)

/*
	测试场景
	1. 使用本地配置(+-用户配置)
	2. 使用配置中心配置(+-用户配置)
	3. 无用户配置
	4. 配置value类型(仅支持string,非string将做内部类型转换至string)
	5. 多个用户配置,合并(TODO 不支持,优先级低)
	6. 异常配置测试,无配置段,空配置项,wrapper段;
	5. etc
*/

var config utils.Configure

const (
	// TODO 补充测试环境配置中心地址;
	cfgUrl   = ""
	cfgPrj   = "testPrj"
	cfgGroup = "testGrp"
	cfgSrv   = "testSrv"
	cfgName  = "aiges.toml"
)

// 本地配置&本地wrapper配置;
func TestLocalCfgWithWrapper(t *testing.T) {
	cfgOpt := &utils.CfgOption{}
	utils.WithCfgName(cfgName)
	cfg, err := utils.NewCfg(utils.Native, cfgOpt)
	if err != nil {
		t.Errorf("read local cfg %s fail with %s", cfgName, err.Error())
		return
	}

	err = Construct(cfg)
	if err != nil {
		t.Errorf("construct aiconfig fail with %s", err.Error())
		return
	}

	// TODO 校验本地文件(框架+用户)有效性;
	// 与本地期望文件及数据校验;

	return
}

// 远端配置文件&wrapper配置
func TestRemoteCfgWithWrapper(t *testing.T) {
	cfgOpt := &utils.CfgOption{}
	utils.WithCfgURL(cfgUrl)
	utils.WithCfgPrj(cfgPrj)
	utils.WithCfgGroup(cfgGroup)
	utils.WithCfgService(cfgSrv)
	utils.WithCfgName(cfgName)
	cfg, err := utils.NewCfg(utils.Centre, cfgOpt)
	if err != nil {
		t.Errorf("down remote cfg %s fail with %s", cfgName, err.Error())
		return
	}

	err = Construct(cfg)
	if err != nil {
		t.Errorf("construct aiconfig fail with %s", err.Error())
		return
	}

	// TODO 远端本地文件(框架+用户)有效性;
	// 与本地期望文件及数据校验;

	return
}

// 无用户配置场景
func TestNullWrapperCfg(t *testing.T) {

	return
}

// 用户自定义配置项key类型(int,uint,float,book,string), 内部转换ToString;
func TestWrapperCfgType(t *testing.T) {

	return
}

// 用户自定义配置section扁平化降级;
func TestWrapperCfgSec(t *testing.T) {

	return
}
