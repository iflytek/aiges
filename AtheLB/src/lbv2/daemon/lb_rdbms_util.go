package daemon

//
////RDBMS(Relational Database Management System
//type RowData struct {
//	segIdDb    string
//	typeDb     string
//	serverIpDb string
//}
//
//func (R *RowData) String() string {
//	return fmt.Sprintf("segIdDb:%s,typeDb:%s,serverIpDb:%s", R.segIdDb, R.typeDb, R.serverIpDb)
//}
//
//type getCountOralRst struct {
//	Ret      int `json:"ret"`
//	Count    int `json:"count"`
//	CountMAX int `json:"count_MAX"`
//}
//
//var columnJson = map[string]string{
//	"seg_id":    "seg_id",
//	"type":      "type",
//	"server_ip": "server_ip",
//}
//
//type oriStruct struct {
//	Ret    string `json:"ret"`
//	Result []struct {
//		SegID    string `json:"seg_id"`
//		ServerIP string `json:"server_ip"`
//		Type     string `json:"type"`
//	} `json:"result"`
//}
//
//var columnJsonRwMu sync.RWMutex
//
//var MysqlManagerInst MysqlManager
//
//type MysqlManager struct {
//	able   int64
//	rowsch chan RowData
//	log    *utils.Logger
//	batch  int
//
//	WebService
//}
//
//func (m *MysqlManager) Init(able int64, batch int, log *utils.Logger, baseUrl string, caller, callerKey string, timeout time.Duration,
//	token string, version string, idc string, schema string, table string) (bool, error) {
//	m.able = able
//	if m.able == 0 {
//		return true, nil
//	}
//	m.log = log
//	m.batch = batch
//	m.rowsch = make(chan RowData, 1000)
//	go m.mysqlWriter()
//	return m.WebService.Init(baseUrl, caller, callerKey, timeout, token, version, idc, schema, table, log)
//}
//
///*
//1、将原始的行结构形式的数据处理为map形式
//2、k为seg_id，v为server_ip
//*/
//func (m *MysqlManager) GetSubSvcSegIdData(subSvc string) (rst map[string]string, err error) {
//	if m.able == 0 {
//		return nil, nil
//	}
//	if nil == rst {
//		rst = make(map[string]string)
//	}
//	rows, rowsErr := m.Retrieve(map[string]string{"type": subSvc})
//	if nil != rowsErr {
//		err = rowsErr
//		return
//	}
//	for _, v := range rows {
//		rst[v.segIdDb] = v.serverIpDb
//	}
//	return
//}
//func (m *MysqlManager) GetTableSize() (count, countMax int, err error) {
//	if m.able == 0 {
//		return 0, 0, nil
//	}
//	countJson := map[string]string{
//		"*": "count",
//	}
//	getCountRst, getCountRstErr := MysqlManagerInst.GetCount(countJson, nil)
//	m.log.Debugw("fn:GetTableSize", "getCountRst", getCountRst, "getCountRstErr", getCountRstErr)
//	if nil != getCountRstErr {
//		return 0, 0, getCountRstErr
//	}
//
//	getCountOralRstTmp := &getCountOralRst{}
//	if err := json.Unmarshal([]byte(getCountRst), getCountOralRstTmp); nil != err {
//		return 0, 0, err
//	}
//
//	return getCountOralRstTmp.Count, getCountOralRstTmp.CountMAX, nil
//}
//func (m *MysqlManager) GetSubSvcSegIdSrvipEx() (rst map[string]map[string]string, err error) {
//	if m.able == 0 {
//		return nil, nil
//	}
//	var rows []RowData
//	var rowsErr error
//	for retry := 0; retry < 3; retry++ {
//		rows, rowsErr = MysqlManagerInst.RetrieveAll(m.batch)
//		if nil != rowsErr {
//			err = fmt.Errorf("pos:MysqlManagerInst.Retrieve(nil) rowsErr:%v,retry:%v", rowsErr, retry)
//			continue
//		} else {
//			break
//		}
//	}
//
//	rst = make(map[string]map[string]string)
//
//	if nil != rowsErr {
//		return nil, rowsErr
//	}
//
//	for _, row := range rows {
//		if v, vOk := rst[row.typeDb]; vOk {
//			if nil == v {
//				v = make(map[string]string)
//			}
//			v[row.segIdDb] = row.serverIpDb
//		} else {
//			rst[row.typeDb] = map[string]string{row.segIdDb: row.serverIpDb}
//		}
//	}
//	return
//}
//
//func (m *MysqlManager) GetSubSvcSegIdSrvip() (rst map[string]map[string]string, err error) {
//	if m.able == 0 {
//		return nil, nil
//	}
//	var rows []RowData
//	var rowsErr error
//	for retry := 0; retry < 3; retry++ {
//		rows, rowsErr = MysqlManagerInst.Retrieve(nil)
//		if nil != rowsErr {
//			err = fmt.Errorf("pos:MysqlManagerInst.Retrieve(nil) rowsErr:%v,retry:%v", rowsErr, retry)
//			continue
//		} else {
//			break
//		}
//	}
//
//		rst = make(map[string]map[string]string)
//	for _, row := range rows {
//		if v, vOk := rst[row.typeDb]; vOk {
//			if nil == v {
//				v = make(map[string]string)
//			}
//			v[row.segIdDb] = row.serverIpDb
//		} else {
//			rst[row.typeDb] = map[string]string{row.segIdDb: row.serverIpDb}
//		}
//	}
//	return
//}
//
//func (m *MysqlManager) AddNewSegIdData(row RowData) (bool, error) {
//	if m.able == 0 {
//		return true, nil
//	}
//	return m.Insert(map[string]string{"seg_id": row.segIdDb, "type": row.typeDb, "server_ip": row.serverIpDb})
//}
//
//func (m *MysqlManager) mysqlWriter() {
//	for {
//		select {
//		case row := <-m.rowsch:
//			{
//				cnt := 3
//				for {
//					writeOk, writeReply := MysqlManagerInst.AddNewSegIdData(row)
//					m.log.Debugw(
//						"setServer",
//						"writeOk", writeOk, "writeReply", writeReply, "row", row, "cnt", cnt)
//					cnt++
//					if writeOk {
//						m.log.Errorw(
//							"setServer MysqlManagerInst.AddNewSegIdData success",
//							"row", row, "cnt", cnt)
//						break
//					}
//					if cnt >= 3 {
//						//todo 添加详细错误描述
//						m.log.Errorw(
//							"setServer MysqlManagerInst.AddNewSegIdData failed",
//							"row", row, "cnt", cnt, "writeReply", writeReply)
//						break
//					}
//				}
//			}
//		}
//	}
//}
//
///*
//	异步交互
//*/
//func (m *MysqlManager) AddNewSegIdDataAsync(row RowData) {
//	if m.able == 0 {
//		return
//	}
//	m.rowsch <- row
//}
//
//func (m *MysqlManager) DelServer(addr string) (bool, error) {
//	if m.able == 0 {
//		return true, nil
//	}
//	return m.Delete(map[string]string{"server_ip": addr})
//
//}
//
//func (m *MysqlManager) db2Rows(in string) (rst []RowData, err error) {
//	tmp := oriStruct{}
//	if err := json.Unmarshal([]byte(in), &tmp); nil != err {
//		return nil, fmt.Errorf("can't parse the json -> %v,err -> %v", in, err)
//	}
//	if tmp.Ret != "0" {
//		return nil, fmt.Errorf("the ret -> %v not equal 0", tmp.Ret)
//	}
//	for _, v := range tmp.Result {
//		rst = append(rst, RowData{segIdDb: v.SegID, typeDb: v.Type, serverIpDb: v.ServerIP})
//	}
//	return rst, nil
//}
//
///*
//	拉取所有数据
//*/
//func (m *MysqlManager) RetrieveAll(batch int) ([]RowData, error) {
//	if m.able == 0 {
//		return nil, nil
//	}
//	if batch == 0 {
//		return nil, fmt.Errorf("fn:%v,batch:%v", "RetrieveAll", batch)
//	}
//
//	tableCount, tableCountMax, tableErr := m.GetTableSize()
//	m.log.Infow(
//		"GetTableSize",
//		"tableCount", tableCount, "tableCountMax", tableCountMax, "tableErr", tableErr)
//
//	if nil != tableErr {
//		return nil, tableErr
//	}
//	if tableCountMax == 0 {
//		m.log.Warnw("db table is empty", "tableCount", tableCount, "tableCountMax", tableCountMax, "tableErr", tableErr)
//		return nil, nil
//	}
//	dataIx := 0
//	var rowDatas []RowData
//	var GetListStr string
//	var GetListStrErr error
//
//	pbCnt := func() int {
//		cntTmp := tableCountMax / batch
//		if tableCountMax%batch != 0 {
//			cntTmp++
//		}
//		return cntTmp
//	}()
//	m.log.Infow("fn:RetrieveAll", "pbCnt", pbCnt)
//	bar := pb.StartNew(pbCnt).Prefix("Mysql:")
//	defer bar.Finish()
//
//	for {
//		if dataIx >= tableCountMax {
//			m.log.Infow("fn:Query", "dataIx", dataIx, "tableCountMax", tableCountMax)
//			break
//		}
//		for retry := 0; retry < 3; retry++ {
//			GetListStr, GetListStrErr = m.WebService.Query(dataIx, batch)
//			if nil != GetListStrErr {
//				m.log.Infow("fn:Query failed", "dataIx", dataIx, "batch", batch, "retry", retry)
//				continue
//			} else {
//				m.log.Infow("fn:Query success", "dataIx", dataIx, "batch", batch, "retry", retry, "GetListStr", GetListStr)
//				break
//			}
//		}
//
//		if nil != GetListStrErr {
//			m.log.Errorw("fn:Query", "GetListStrErr", GetListStrErr)
//			return nil, GetListStrErr
//		}
//
//		rows, rowsErr := m.db2Rows(GetListStr)
//		if nil != rowsErr {
//			m.log.Errorw(
//				"db2Rows failed",
//				"GetListStr", GetListStr, "rows", rows, "rowsErr", rowsErr)
//			return nil, rowsErr
//		}
//		rowDatas = append(rowDatas, rows...)
//		bar.Increment()
//		dataIx += batch
//	}
//
//	return rowDatas, nil
//}
//func (m *MysqlManager) Retrieve(where map[string]string) ([]RowData, error) {
//	if m.able == 0 {
//		return nil, nil
//	}
//	columnJsonRwMu.RLock()
//	defer columnJsonRwMu.RUnlock()
//	GetListStr, GetListStrErr := m.WebService.GetList(columnJson, where, nil)
//	if nil != GetListStrErr {
//		return nil, GetListStrErr
//	}
//	return m.db2Rows(GetListStr)
//}
//
//func (m *MysqlManager) Create(increase []RowData) bool {
//	if m.able == 0 {
//		return true
//	}
//	panic("implement me")
//}
//
//func (m *MysqlManager) Update(set RowData, where map[string]string) bool {
//	if m.able == 0 {
//		return true
//	}
//	panic("implement me")
//}
//
//func (m *MysqlManager) Delete(where map[string]string) (bool, error) {
//	if m.able == 0 {
//		return true, nil
//	}
//	return m.WebService.Delete(where)
//}
