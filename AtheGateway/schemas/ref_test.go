package schemas

import (
	"testing"
	"fmt"
	"strconv"
	"regexp"
)

func Test_checkRule(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{args:args{""}, want:false,},  //00
		{args:args{"$"}, want:true,},  //01
		{args:args{"$.abc"}, want:true,}, //02
		{args:args{"$.name"}, want:true,},  //03
		{args:args{"$[0].name"}, want:true,},//04
		{args:args{"$[0].age"}, want:true,}, //05
		{args:args{"$[d].age"}, want:false,}, //06
		{args:args{"$[.age"}, want:false,}, // 07
		{args:args{"$].age"}, want:false,}, // 08
		{args:args{"$0].age"}, want:false,},
		{args:args{"$[1].age[0]"}, want:true,},
		{args:args{"$[1].age[0"}, want:false,},
		{args:args{"$[1].age0]"}, want:false,},
		{args:args{"$[1].age2[2]"}, want:true,},
		{args:args{"$[1].age2[2].[3]"}, want:false,},
		{args:args{"$[1].age2[2].v[3]"}, want:true,},
		{args:args{"$[1].age2.v[3]"}, want:true,},
		{args:args{"$.age2.v[3]"}, want:true,},
		{args:args{"$.age2.v[3]dsf"}, want:false,},
		{args:args{"$.age2.v[3.dsf"}, want:false,},
		{args:args{"$.age2.v[e].dsf"}, want:false,},
		{args:args{"$.age2.v[3].dsf"}, want:true,},
		{args:args{"$.data[0]"}, want:true,},
		{args:args{"$[0].data"}, want:true,},
		{args:args{"$[1].data"}, want:true,},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkRule(tt.args.token); got != tt.want {
				t.Errorf("checkRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
var insert = regexp.MustCompile(`("\w+")+$`)

func TestLoadMapping(t *testing.T) {
	fmt.Println(checkRule("$"))
}
func TestMarshals(t *testing.T) {
	str:=`"sdf""dfsd"`
	fmt.Println(insert.MatchString(str))
	r:=insert.FindAllSubmatch([]byte(str),-1)

	for _,v:=range r{
		//fmt.Println(string(v[0]))
		fmt.Println("----")
		for _,vv:=range v{
			fmt.Println(string(vv))
		}
	}

}

type Pt struct {

	Name string
	Sign string

}

func TestGetMappingKey(t *testing.T) {
	a:=[]Pt{
		{"都让人","发的搜集"},
		{"水电费水电费","防守打法"},
	}

	ia:=ToInterface(len(a), func(i int) interface{} {
		return a[i]
	})
	PrintArray(ia, func(i interface{}) map[string]string {
		v:=i.(Pt)
		return map[string]string{
			"name":v.Name,
			"sign":v.Sign,
		}
	})
}

func ToInterface(len int,f func(int)interface{})[]interface{}  {
	var is []interface{}
	for i:=0;i<len;i++{
		is = append(is,f(i))
	}
	return is
}

func PrintArray(is interface{},f func(i interface{})map[string]string)  {
	isa,ok:=is.([]interface{})
	if !ok||len(isa)==0{
		return
	}
	samlpe:=f(isa[0])
	var keys[]string
	for k,_:=range samlpe{
		keys = append(keys,k)
	}
	maxks:=make(map[string]int)
	var kv []interface{}
	for _,key:=range keys{
		maxks[key] = len(key)
		kv = append(kv,key)
		for _,v:=range isa{
			line:=f(v)
			len:=len(line[key])
			if maxks[key]<len{
				maxks[key] = len
			}
		}
	}
	fmt.Println(keys)
	fmt.Println(maxks)
	var format =""
	var splitLine= ""
	for _,v:=range maxks{
		format+="|%-"+strconv.Itoa(v)+"s"
		splitLine+="|"
		for i:=0;i<v;i++{
			splitLine+="-"
		}
	}
	splitLine+=""
	format+="\n"
	fmt.Printf(format,kv...)
	fmt.Println(splitLine)
	//fmt.Println(splitLine)
	for _,v:=range isa{
		var p []interface{}
		m:=f(v)
		for _,vv:=range keys{
			p = append(p,m[vv])
		}
	//	fmt.Println(splitLine)
	}
	//fmt.Println(format)

}