package daemon

//
//type WebService struct {
//	caller    string
//	callerKey string
//	baseUrl   string
//	timeout   time.Duration
//	token     string
//	version   string
//	idc       string
//	schema    string
//	table     string
//	log       *utils.Logger
//}
//
//func (w *WebService) Init(baseUrl string, caller, callerKey string, timeout time.Duration, token string,
//	version string, idc string, schema string, table string, log *utils.Logger) (bool, error) {
//	w.baseUrl = baseUrl
//	w.caller = caller
//	w.callerKey = callerKey
//	w.timeout = timeout
//	w.token = token
//	w.version = version
//	w.idc = idc
//	w.schema = schema
//	w.table = table
//	w.log = log
//	return true, nil
//}
//func (w *WebService) getMd5(in string) string {
//	h := md5.New()
//	h.Write([]byte(in))
//	return hex.EncodeToString(h.Sum(nil))
//}
//
//func (w *WebService) Insert(insertJson map[string]string) (bool, error) {
//	queryStruct := w.generateQueryStruct(w.version, w.idc, w.schema, w.table, nil, nil,
//		insertJson, nil, nil, "")
//	requestUrl := w.generateUrl("insert", map[string]string{
//		"caller":      w.caller,
//		"checksum":    w.getMd5(queryStruct + w.callerKey),
//		`token`:       w.token,
//		`queryStruct`: queryStruct,
//	})
//	client := http.Client{Timeout: w.timeout}
//	req, err := http.NewRequest("GET", requestUrl, nil)
//	if err != nil {
//		return false, err
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return false, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return false, err
//	}
//	bodyStr := string(body)
//	if !strings.Contains(bodyStr, `"ret":"0"`) {
//		return false, errors.New(fmt.Sprintf("body:%s,err:%v", body, err))
//	}
//	return true, nil
//}
//func (w *WebService) Delete(whereJson map[string]string) (bool, error) {
//	queryStruct := w.generateQueryStruct(w.version, w.idc, w.schema, w.table, nil, whereJson,
//		nil, nil, nil, "")
//	requestUrl := w.generateUrl("delete", map[string]string{
//		"caller":      w.caller,
//		"checksum":    w.getMd5(queryStruct + w.callerKey),
//		`token`:       w.token,
//		`queryStruct`: queryStruct,
//	})
//	client := http.Client{Timeout: w.timeout}
//	req, err := http.NewRequest("GET", requestUrl, nil)
//	if err != nil {
//		return false, err
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return false, err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return false, err
//	}
//	if !strings.Contains(string(body), `"ret":"0"`) {
//		return false, errors.New(fmt.Sprintf("body:%s,err:%v", body, err))
//	}
//	return true, nil
//}
//func (w *WebService) GetList(columnJson map[string]string, whereJson map[string]string, limitJson map[string]string) (string, error) {
//	queryStruct := w.generateQueryStruct(w.version, w.idc, w.schema, w.table, columnJson, whereJson,
//		nil, limitJson, nil, "")
//	requestUrl := w.generateUrl("getList", map[string]string{
//		"caller":      w.caller,
//		"checksum":    w.getMd5(queryStruct + w.callerKey),
//		`token`:       w.token,
//		`queryStruct`: queryStruct,
//	})
//	client := http.Client{Timeout: w.timeout}
//	req, err := http.NewRequest("GET", requestUrl, nil)
//	w.log.Infow("GetList", "fn", "GetList", "requestUrl", requestUrl)
//	if err != nil {
//		return "", err
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	return fmt.Sprintf("%s", body), nil
//}
//
//func (w *WebService) GetCount(countJson map[string]string, whereJson map[string]string) (string, error) {
//	queryStruct := w.generateQueryStruct(w.version, w.idc, w.schema, w.table,
//		columnJson, whereJson, nil, nil, countJson, "")
//	requestUrl := w.generateUrl("getCount", map[string]string{
//		"caller":      w.caller,
//		"checksum":    w.getMd5(queryStruct + w.callerKey),
//		`token`:       w.token,
//		`queryStruct`: queryStruct,
//	})
//	w.log.Infow("fn:GetCount", "requestUrl", requestUrl)
//	client := http.Client{Timeout: w.timeout}
//	req, err := http.NewRequest("GET", requestUrl, nil)
//	//fmt.Printf("getList => requestUrl:%v\n", requestUrl)
//	if err != nil {
//		w.log.Infow("fn:GetCount http.NewRequest fail", "err", err)
//		return "", err
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		w.log.Infow("fn:GetCount client.Do fail", "err", err)
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		w.log.Infow("fn:GetCount ioutil.ReadAll fail", "err", err)
//		return "", err
//	}
//	return fmt.Sprintf("%s", body), nil
//}
//
//func (w *WebService) Query(offset, limit int) (string, error) {
//	return w.query(fmt.Sprintf("select * from `%v` order by seg_id desc LIMIT %d,%d;", w.table, offset, limit))
//}
//
//func (w *WebService) query(oralSql string) (string, error) {
//	queryStruct := w.generateQueryStruct(w.version, w.idc, w.schema, w.table,
//		columnJson, nil, nil, nil, nil, oralSql)
//	requestUrl := w.generateUrl("query", map[string]string{
//		"caller":      w.caller,
//		"checksum":    w.getMd5(queryStruct + w.callerKey),
//		`token`:       w.token,
//		`queryStruct`: queryStruct,
//	})
//	w.log.Infow("fn:query", "requestUrl", requestUrl)
//	client := http.Client{Timeout: w.timeout}
//	req, err := http.NewRequest("GET", requestUrl, nil)
//	//fmt.Printf("getList => requestUrl:%v\n", requestUrl)
//	if err != nil {
//		w.log.Infow("fn:query http.NewRequest fail", "err", err)
//		return "", err
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		w.log.Infow("fn:query client.Do fail", "err", err)
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		w.log.Infow("fn:query ioutil.ReadAll fail", "err", err)
//		return "", err
//	}
//	return fmt.Sprintf("%s", body), nil
//}
//
//func (w *WebService) GetListNoWhereJson(columnJson map[string]string) (string, error) {
//	queryStruct := w.generateGetlistQuerystructNowherejson(columnJson)
//	requrestUrl := w.generateUrl("getList", map[string]string{
//		`token`:       w.token,
//		`queryStruct`: queryStruct,
//	})
//	client := http.Client{Timeout: w.timeout}
//	req, err := http.NewRequest("GET", requrestUrl, nil)
//	if err != nil {
//		return "", err
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return "", err
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", err
//	}
//	return fmt.Sprintf("%s", body), nil
//}
//func (w *WebService) generateUrl(op string, params map[string]string) string {
//	u, _ := url.Parse(w.baseUrl + "/" + op)
//	q := u.Query()
//	for k, v := range params {
//		q.Set(k, v)
//	}
//	u.RawQuery = q.Encode()
//	urlReturn := u.String()
//	return urlReturn
//}
//func (w *WebService) generateQueryStruct(version string, idc string, schema string, table string,
//	columnJson map[string]string, whereJson map[string]string, insertJson map[string]string,
//	limitJson map[string]string, countJson map[string]string, oralSql string) string {
//	querystructMap := make(map[string]interface{})
//	querystructMap[`version`] = version
//	querystructMap[`meta`] = func() interface{} {
//		metaMap := make(map[string]interface{})
//		metaMap[`idc`] = idc
//		metaMap[`schema`] = schema
//		metaMap[`table`] = table
//		return metaMap
//	}()
//	querystructMap[`param`] = func() interface{} {
//		paramMap := make(map[string]interface{})
//
//		if len(columnJson) > 0 {
//			paramMap[`columnJson`] = func() interface{} {
//				columnjsonMap := make(map[string]interface{})
//				for k, v := range columnJson {
//					columnjsonMap[k] = v
//				}
//				return columnjsonMap
//			}()
//		}
//
//		if len(whereJson) > 0 {
//			paramMap[`whereJson`] = func() interface{} {
//				wherejsonMap := make(map[string]interface{})
//				for k, v := range whereJson {
//					wherejsonMap[k] = v
//				}
//				return wherejsonMap
//			}()
//		}
//
//		if len(insertJson) > 0 {
//			paramMap[`insertJson`] = func() interface{} {
//				wherejsonMap := make(map[string]interface{})
//				for k, v := range insertJson {
//					wherejsonMap[k] = v
//				}
//				return wherejsonMap
//			}()
//		}
//
//		if len(limitJson) > 0 {
//			paramMap[`limitJson`] = func() interface{} {
//				limitjsonMap := make(map[string]interface{})
//				for k, v := range limitJson {
//					limitjsonMap[k] = v
//				}
//				return limitjsonMap
//			}()
//		}
//
//		if len(countJson) > 0 {
//			paramMap[`countJson`] = func() interface{} {
//				countjsonMap := make(map[string]interface{})
//				for k, v := range countJson {
//					countjsonMap[k] = v
//				}
//				return countjsonMap
//			}()
//		}
//
//		if len(oralSql) > 0 {
//			paramMap[`sql`] = oralSql
//		}
//
//		return paramMap
//	}()
//	queryStruct, err := json.Marshal(querystructMap)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return string(queryStruct)
//}
//func (w *WebService) generateGetlistQuerystructNowherejson(columnJson map[string]string) string {
//	querystructMap := make(map[string]interface{})
//	querystructMap[`version`] = w.version
//	querystructMap[`meta`] = func() interface{} {
//		metaMap := make(map[string]interface{})
//		metaMap[`idc`] = w.idc
//		metaMap[`schema`] = w.schema
//		metaMap[`table`] = w.table
//		return metaMap
//	}()
//	querystructMap[`param`] = func() interface{} {
//		paramMap := make(map[string]interface{})
//		paramMap[`columnJson`] = func() interface{} {
//			columnjsonMap := make(map[string]interface{})
//			for k, v := range columnJson {
//				columnjsonMap[k] = v
//			}
//			return columnjsonMap
//		}()
//		return paramMap
//	}()
//	queryStruct, err := json.Marshal(querystructMap)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return string(queryStruct)
//}
