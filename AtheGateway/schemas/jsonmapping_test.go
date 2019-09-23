package schemas

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

type Mapping struct {
	Version     string              `json:"version"`
	Service     string              `json:"service"`
	Route       string              `json:"route"`
	ReqMapping  []map[string]string `json:"request.data.mapping"`
	RespMapping []map[string]string `json:"response.data.mapping"`
}

//解析请求
func (mp *Mapping) ResolveUpCallReq(data []byte) (*[]map[string]interface{}, error) {
	var reqMap []map[string]interface{}
	err := json.Unmarshal(data, &reqMap)
	if err != nil {
		return nil, err
	}
	resolveResult := make([]map[string]interface{}, len(mp.ReqMapping))
	for index, elem := range mp.ReqMapping {
		var reqDate = make(map[string]interface{})
		for key, value := range elem {
			if reqMap[index] != nil {
				reqDate[key] = reqMap[index][value]
			}
			fmt.Printf("%s:%v\n", key, reqDate[key])
		}
		resolveResult[index] = reqDate
		fmt.Println("==========================")
	}
	return &resolveResult, nil
}

//解析Resp
func (mp *Mapping) ResolveResp(data []byte) (*[]map[string]interface{}, error) {
	var respMap []map[string]interface{}
	err := json.Unmarshal(data, &respMap)
	if err != nil {
		return nil, err
	}
	resolveResult := make([]map[string]interface{}, len(mp.RespMapping))
	for index, elem := range mp.RespMapping {
		var respDate = make(map[string]interface{})
		for key, value := range elem {
			if respMap[index] != nil {
				respDate[value] = respMap[index][key]
				fmt.Printf("%s:%v\n", value, respDate[value])
			}
		}
		fmt.Println("==========================")
		resolveResult[index] = respDate
	}

	return &resolveResult, nil
}
func main() {
	file, err := os.Open("json_mapping")
	if err != nil {
		fmt.Printf("打开文件失败,失败得原因是:%s\n", err.Error())
		return
	}
	defer file.Close()
	fileStat, err := file.Stat()
	if err != nil {
		fmt.Printf("获取文件状态失败,失败得原因是:%s\n", err.Error())
		return
	}
	var bufer []byte = make([] byte, fileStat.Size())
	reader := bufio.NewReader(file)
	reader.Read(bufer)
	mappings := jsonMapping(bufer)

	//================================
	req := `[{
		"status": 1,
		"format":"audio/L16;rate=16000",
		"encoding": "raw",
		"audio": "123"
	},{
		"id": 1,
		"status":2,
		"format":"audio/L16;rate=16000",
		"encoding": "raw",
		"audio": "1asdfadfasdfadf"
	}]`

	resp := `
		[{
		"status":1,
		"data":"resultasdfadfadfad"
		}]
	`
	for _, elem := range *mappings {
		elem.ResolveUpCallReq([]byte(req))
		fmt.Println("--------------------------------------")
		elem.ResolveResp([]byte(resp))
	}
	//fmt.Printf("mapping:%v\n", mapping)
}

func jsonMapping(buffer []byte) *[]Mapping {
	//fmt.Println("buffer size:", len(buffer))
	mappings := &[]Mapping{}
	json.Unmarshal(buffer, mappings)
	//for _, elem := range *mappings {
	//	fmt.Printf("Version:%s\nService:%s\nRoute:%s\nReqMapping:%s\nRespMapping:%s\n", elem.Version, elem.Service, elem.Route, elem.ReqMapping, elem.RespMapping)
	//	for _, data := range elem.ReqMapping {
	//		fmt.Println("req===>", data)
	//		for key, value := range data {
	//			fmt.Printf("%s:%s\n", key, value)
	//		}
	//	}
	//
	//	for _, data := range elem.RespMapping {
	//		fmt.Println("resp===>", data)
	//		for key, value := range data {
	//			fmt.Printf("%s:%s\n", key, value)
	//		}
	//	}
	//}
	return mappings
}
func TestCachedJsonpathLookUp(t *testing.T) {
	st:=time.Now()
	a:=0
	a++
	for a <1000000{
		a++
	}

	fmt.Println("hsdfsdf")
	fmt.Println("hsdfsdf")

	//time.Sleep(1*time.Nanosecond)
	fmt.Println(time.Now().Sub(st).Nanoseconds())
}