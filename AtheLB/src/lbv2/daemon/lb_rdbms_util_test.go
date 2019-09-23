package daemon

import (
	"fmt"
	"testing"
	"time"

	"git.xfyun.cn/AIaaS/xsf-external/utils"
)

func TestMysqlManager_RetrieveEx(t *testing.T) {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(true),
		utils.SetLevel("info"),
		utils.SetFileName("test.log"),
		utils.SetMaxSize(3),
		utils.SetMaxBackups(3),
		utils.SetMaxAge(3),
		utils.SetAsync(false),
		utils.SetCacheMaxCount(30000),
		utils.SetBatchSize(1024))
	if nil != err {
		t.Fatal(err)
	}
	_, _ = MysqlManagerInst.Init(
		1,
		10000,
		logger,
		"http://172.16.154.235:8081/ws",
		"xfyun",
		"12345678",
		time.Second*10,
		"100IME",
		"db-service-v3-3.0.0.1001",
		"bj",
		"ifly_cp_msp_balance",
		"seg_list_lbv2")
	//rst, err := MysqlManagerInst.GetSubSvcSegIdData("sms")
	//rows, err := MysqlManagerInst.Retrieve(nil)
	data, dataErr := MysqlManagerInst.GetSubSvcSegIdSrvipEx()
	std.Printf("dataLen:%v,data:%#v,dataErr:%v\n", len(data), data, dataErr)
}

func TestMysqlManager_Retrieve(t *testing.T) {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(true),
		utils.SetLevel("info"),
		utils.SetFileName("test.log"),
		utils.SetMaxSize(3),
		utils.SetMaxBackups(3),
		utils.SetMaxAge(3),
		utils.SetAsync(true),
		utils.SetCacheMaxCount(30000),
		utils.SetBatchSize(1024))
	if nil != err {
		t.Fatal(err)
	}
	_, _ = MysqlManagerInst.Init(
		1,
		10000,
		logger,
		"http://172.16.154.235:8081/ws",
		"xfyun",
		"12345678",
		time.Second*10,
		"100IME",
		"db-service-v3-3.0.0.1001",
		"bj",
		"ifly_cp_msp_balance",
		"seg_list_lbv2")
	//rst, err := MysqlManagerInst.GetSubSvcSegIdData("sms")
	//rows, err := MysqlManagerInst.Retrieve(nil)
	data, dataErr := MysqlManagerInst.GetSubSvcSegIdSrvip()
	std.Printf("data:%#v,dataErr:%v\n", data, dataErr)
}

func TestMysqlManager_AddNewSegIdDataAsync(t *testing.T) {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(true),
		utils.SetLevel("info"),
		utils.SetFileName("test.log"),
		utils.SetMaxSize(3),
		utils.SetMaxBackups(3),
		utils.SetMaxAge(3),
		utils.SetAsync(false),
		utils.SetCacheMaxCount(30000),
		utils.SetBatchSize(1024))
	if nil != err {
		t.Fatal(err)
	}

	_, _ = MysqlManagerInst.Init(
		1,
		10000,
		logger,
		"http://172.16.154.235:8081/ws",
		"xfyun",
		"12345678",
		time.Second*10,
		"100IME",
		"db-service-v3-3.0.0.1001",
		"bj",
		"ifly_cp_msp_balance",
		"seg_list_lbv2")

	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "0", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "1", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "2", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "3", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "4", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "5", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "6", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "7", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "8", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	MysqlManagerInst.AddNewSegIdDataAsync(RowData{segIdDb: "9", typeDb: "xxx", serverIpDb: "127.0.0.x"})
	time.Sleep(time.Second * 3)
}

func TestMysqlManager_AddNewSegIdData(t *testing.T) {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(true),
		utils.SetLevel("info"),
		utils.SetFileName("test.log"),
		utils.SetMaxSize(3),
		utils.SetMaxBackups(3),
		utils.SetMaxAge(3),
		utils.SetAsync(true),
		utils.SetCacheMaxCount(30000),
		utils.SetBatchSize(1024))
	if nil != err {
		t.Fatal(err)
	}

	_, _ = MysqlManagerInst.Init(
		1,
		10000,
		logger,
		"http://172.16.154.235:8081/ws",
		"xfyun",
		"12345678",
		time.Second*10,
		"100IME",
		"db-service-v3-3.0.0.1001",
		"bj",
		"ifly_cp_msp_balance",
		"seg_list_lbv2")

	std.Println(MysqlManagerInst.AddNewSegIdData(RowData{segIdDb: "10", typeDb: "xxx", serverIpDb: "127.0.0.x"}))
}

func TestMysqlManager_DelServer(t *testing.T) {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(true),
		utils.SetLevel("info"),
		utils.SetFileName("test.log"),
		utils.SetMaxSize(3),
		utils.SetMaxBackups(3),
		utils.SetMaxAge(3),
		utils.SetAsync(true),
		utils.SetCacheMaxCount(30000),
		utils.SetBatchSize(1024))
	if nil != err {
		t.Fatal(err)
	}

	_, _ = MysqlManagerInst.Init(
		1,
		10000,
		logger,
		"http://172.16.154.235:8081/ws",
		"xfyun",
		"12345678",
		time.Second*10,
		"100IME",
		"db-service-v3-3.0.0.1001",
		"bj",
		"ifly_cp_msp_balance",
		"seg_list_lbv2")

	std.Println(MysqlManagerInst.DelServer("10.1.87.21:5090"))
}

func TestMysqlManager_GetTableSize(t *testing.T) {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(true),
		utils.SetLevel("debug"),
		utils.SetFileName("test.log"),
		utils.SetMaxSize(3),
		utils.SetMaxBackups(3),
		utils.SetMaxAge(3),
		utils.SetAsync(false),
		utils.SetCacheMaxCount(30000),
		utils.SetBatchSize(1024))
	if nil != err {
		t.Fatal(err)
	}

	_, _ = MysqlManagerInst.Init(
		1,
		10000,
		logger,
		"http://172.16.154.235:8081/ws",
		"xfyun",
		"12345678",
		time.Second*10,
		"100IME",
		"db-service-v3-3.0.0.1001",
		"bj",
		"ifly_cp_msp_balance",
		"seg_list_lbv2")

	std.Println(MysqlManagerInst.GetTableSize())
}

func TestMysqlManager_RetrieveAll(t *testing.T) {
	logger, err := utils.NewLocalLog(
		utils.SetCaller(true),
		utils.SetLevel("debug"),
		utils.SetFileName("test.log"),
		utils.SetMaxSize(3),
		utils.SetMaxBackups(3),
		utils.SetMaxAge(3),
		utils.SetAsync(false),
		utils.SetCacheMaxCount(30000),
		utils.SetBatchSize(1024))
	if nil != err {
		t.Fatal(err)
	}

	_, _ = MysqlManagerInst.Init(
		1,
		10000,
		logger,
		"http://172.16.154.235:8081/ws",
		"xfyun",
		"12345678",
		time.Second*10,
		"100IME",
		"db-service-v3-3.0.0.1001",
		"bj",
		"ifly_cp_msp_balance",
		"seg_list_lbv2")

	rows, rowsErr := MysqlManagerInst.RetrieveAll(1)
	fmt.Printf("rowsLen:%d,rowsErr:%v,rows:%v\n", len(rows), rowsErr, rows)
}
