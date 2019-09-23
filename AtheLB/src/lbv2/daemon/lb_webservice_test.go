package daemon

//
//func TestMysqlManager_basic(t *testing.T) {
//	logger, err := utils.NewLocalLog(
//		utils.SetCaller(true),
//		utils.SetLevel("info"),
//		utils.SetFileName("test.log"),
//		utils.SetMaxSize(3),
//		utils.SetMaxBackups(3),
//		utils.SetMaxAge(3),
//		utils.SetAsync(false),
//		utils.SetCacheMaxCount(30000),
//		utils.SetBatchSize(1024))
//	if err != nil {
//		t.Fatal(err)
//	}
//	WebServiceInst := WebService{}
//	WebServiceInst.Init("http://172.16.154.235:8081/ws", "xfyun", "12345678", time.Second*3,
//		`100IME`, "db-service-v3-3.0.0.1001", "bj", "ifly_cp_msp_balance",
//		"seg_list", logger)
//
//	columnJson := map[string]string{
//		"seg_id":    "seg_id",
//		"type":      "type",
//		"server_ip": "server_ip",
//	}
//
//	whereJson := map[string]string{
//		"type": "xxx",
//	}
//
//	limitJson := map[string]string{
//		"2": "3",
//	}
//	std.Println(WebServiceInst.GetList(columnJson, whereJson, limitJson))
//	//fmt.Println(WebServiceInst.GetListNoWhereJson(columnJson))
//	//fmt.Println(WebServiceInst.Insert(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//	//fmt.Println(WebServiceInst.Insert(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//	//fmt.Println(WebServiceInst.Delete(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//}
//
//func TestMysqlManager_getCount(t *testing.T) {
//	logger, err := utils.NewLocalLog(
//		utils.SetCaller(true),
//		utils.SetLevel("info"),
//		utils.SetFileName("test.log"),
//		utils.SetMaxSize(3),
//		utils.SetMaxBackups(3),
//		utils.SetMaxAge(3),
//		utils.SetAsync(true),
//		utils.SetCacheMaxCount(30000),
//		utils.SetBatchSize(1024))
//	if err != nil {
//		t.Fatal(err)
//	}
//	WebServiceInst := WebService{}
//	WebServiceInst.Init("http://172.16.154.235:8081/ws", "xfyun", "12345678", time.Second*3,
//		`100IME`, "db-service-v3-3.0.0.1001", "bj", "ifly_cp_msp_balance",
//		"seg_list", logger)
//	countJson := map[string]string{
//		"*": "count",
//	}
//	whereJson := map[string]string{
//		//"type": "xxx",
//	}
//	std.Println(WebServiceInst.GetCount(countJson, whereJson))
//	//fmt.Println(WebServiceInst.GetListNoWhereJson(columnJson))
//	//fmt.Println(WebServiceInst.Insert(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//	//fmt.Println(WebServiceInst.Insert(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//	//fmt.Println(WebServiceInst.Delete(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//}
//
////output:
///*
//{"ret":"0","result":[{"seg_id":"21","server_ip":"127.0.0.9","type":"vmail"}]} <nil>
//{"ret":"0","result":[{"seg_id":"21","server_ip":"127.0.0.9","type":"vmail"},{"seg_id":"11","server_ip":"127.0.0.1","type":"sms"}]} <nil>
//*/
//
//func TestMysqlManager_insert(t *testing.T) {
//	logger, err := utils.NewLocalLog(
//		utils.SetCaller(true),
//		utils.SetLevel("info"),
//		utils.SetFileName("test.log"),
//		utils.SetMaxSize(3),
//		utils.SetMaxBackups(3),
//		utils.SetMaxAge(3),
//		utils.SetAsync(true),
//		utils.SetCacheMaxCount(30000),
//		utils.SetBatchSize(1024))
//	if err != nil {
//		t.Fatal(err)
//	}
//	WebServiceInst := WebService{}
//	WebServiceInst.Init("http://172.16.154.235:8081/ws", "xfyun", "12345678", time.Second*3,
//		`100IME`, "db-service-v3-3.0.0.1001", "bj", "ifly_cp_msp_balance",
//		"seg_list", logger)
//	//columnJson := map[string]string{
//	//	"seg_id":    "seg_id",
//	//	"type":      "type",
//	//	"server_ip": "server_ip",
//	//}
//	//whereJson := map[string]string{
//	//	"type": "sms",
//	//}
//	//limitJson := map[string]string{
//	//	//"1": "0",
//	//}
//	//fmt.Println(WebServiceInst.GetList(columnJson, whereJson, limitJson))
//	//fmt.Println(WebServiceInst.GetListNoWhereJson(columnJson))
//	for i := 1; i <= 100; i++ {
//		std.Println(WebServiceInst.Insert(map[string]string{"seg_id": strconv.Itoa(i), "type": "xxx", "server_ip": "x.x.x.x"}))
//	}
//	//fmt.Println(WebServiceInst.Insert(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//	//fmt.Println(WebServiceInst.Delete(map[string]string{"seg_id": "xxx", "type": "xxx", "server_ip": "x.x.x.x"}))
//}
//
//func TestMysqlManager_query(t *testing.T) {
//	logger, err := utils.NewLocalLog(
//		utils.SetCaller(true),
//		utils.SetLevel("info"),
//		utils.SetFileName("test.log"),
//		utils.SetMaxSize(3),
//		utils.SetMaxBackups(3),
//		utils.SetMaxAge(3),
//		utils.SetAsync(false),
//		utils.SetCacheMaxCount(30000),
//		utils.SetBatchSize(1024))
//	if err != nil {
//		t.Fatal(err)
//	}
//	WebServiceInst := WebService{}
//	WebServiceInst.Init("http://172.16.154.235:8081/ws", "xfyun", "12345678", time.Second*3,
//		`100IME`, "db-service-v3-3.0.0.1001", "bj", "ifly_cp_msp_balance",
//		"seg_list", logger)
//
//	//fmt.Println(WebServiceInst.query("select * from `seg_list` order by seg_id desc LIMIT 1,2;"))
//	std.Println(WebServiceInst.Query(0, 1))
//}
