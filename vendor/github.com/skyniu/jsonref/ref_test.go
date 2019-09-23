package jsonref

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func TestLoad2(t *testing.T) {
	var m = map[string]interface{}{}
	Marshal("$.zhansan",m,User{"zhangsan",20})
	log.Println(Marshal("$.class[5].name",m,"biaoge"))
	Marshal("$.class[0]",m,User{"lisi",11})
	Marshal("$.class[1]",m,User{"wangwu",18})
	Marshal("$.class[2]",m,User{"dajj",18})
	log.Println(Marshal("$.class[3].age",m,23))
	Marshal("$.group[5].age",m,12)
	Marshal("$.group[5].son.son.name",m,"bgnb")
	Marshal("$.group[5].son.son.age",m,33)
	Marshal("$.nii.sss.ggg.hhh.jjj.kkk.ll.sss.mmm.ggg",m,23)
    Marshal("$.nii.sss.ggg.hhh.jjj.kkk.ll.sss.mmm.ggg.ff",m,23)
    Marshal("$.nii.sss.ggs[1].hhh[0].jjj[0].kkk[1].ll[2].sss.mmm.ggg.ff[1].ss[0]",m,"12")
	s,_:=json.Marshal(m)
	fmt.Println(string(s))
}

func Test_yy(t *testing.T)  {
	var m = map[string]interface{}{}
	Marshal("$.result",m,map[string]interface{}{
		"status2":2,
	},)
	Marshal("$.result.status",m,1)

	s,_:=json.Marshal(m)
	fmt.Println(string(s))
}

func TestLoad(t *testing.T) {
	var a []string
	var b =a
	b=append(b,"sdf")
	fmt.Println(a)
}

func TestMarshals(t *testing.T) {
	//s:=time.Now()
	marshal1()
	//fmt.Println("time",time.Since(s).Nanoseconds())
}

func marshal1()  {
	tmp,err:=Marshals([]QueryProp{
		//{"$.biaoge.name.sss[0].sdg","biaoge"},
		{"$.result.status",1},
		{"$.result",map[string]interface{}{
			"status":2,
		}},
		//{"$.biaoge.say","bgnb"},
		//{"$.dajj.name","dajj"},
		//{"$.dajj.say","dajj niubi"},
		//{"$.group[0]","biaoge"},
		//{"$.group[1]","dajj"},
		//{"$.group[2]","hg"},
		////{"$.group[2].gg","hg"},
		//{"$.less[0].hgfs.had.pg[0].hhh[1].fs","hg"},
	})
	if err !=nil{
		fmt.Println(err)
		return
	}
	s,_:=json.Marshal(tmp)
	fmt.Sprintf(string(s))
}
func Test_bench(t *testing.T)  {
	for i:=0;i<1000000;i++{
		marshal1()
	}
}

func TestNI(t *testing.T) {
	var m = map[string]interface{}{}
	log.Println(Marshal("$.s[0]",m,1))
	log.Println(Marshal("$.s[1]",m,2))
	log.Println(Marshal("$.s[2]",m,3))
	s,_:=json.Marshal(m)
	fmt.Println(string(s))
}

func TestInterface(t *testing.T) {
	var i interface{}
	inter(&i)
	log.Println(i)
}

func inter(v interface{})  {
	*v.(*interface{}) = map[string]interface{}{
		"kk":"dsf",
	}
	var b = *v.((*interface{}))
	b.(map[string]interface{})["sdf"]="sdf"
}