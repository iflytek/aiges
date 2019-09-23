package common

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"
)

func TestMapToString(t *testing.T) {
	var f float64 = 3
	var a = map[string]interface{}{
		"hello":   "thank you",
		"dfssdf1": f,
	}
	fmt.Println(MapToString(a))
	for i:=0;i<1000000;i++{
		ConvertToString(true)
		//MapToString(a)
		fmt.Sprintf("%v",true)
	}

}

func TestEncodingTobase64String(t *testing.T) {
	type args struct {
		src []byte
	}

	var t1 = "hello thank you"
	var t2 = "how are you"

	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			args: args{[]byte(t1)},
			want: base64.StdEncoding.EncodeToString([]byte("hello thank you")),
		},
		{
			args: args{[]byte(t2)},
			want: base64.StdEncoding.EncodeToString([]byte("how are you")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodingTobase64String(tt.args.src); got != tt.want {
				t.Errorf("EncodingTobase64String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToString(t *testing.T) {
	type args struct {
		buf []byte
	}
	var hello = "hello thank you"
	var test = "hello thank you test"
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{buf: []byte(hello)},
			want: hello,
		},
		{
			args: args{buf: []byte(test)},
			want: test,
		},
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToString(tt.args.buf); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToString(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			args:args{"hello"},
			want:"hello",
		},
		{
			args:args{15},
			want:"15",
		},
		{
			args:args{true},
			want:"true",
		},
		{
			args:args{false},
			want:"false",
		},
		{
			args:args{16.14},
			want:"16.14",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToString(tt.args.v); got != tt.want {
				t.Errorf("ConvertToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringBuildler_Len(t *testing.T) {
	time.AfterFunc(time.Second*2, func() {
		fmt.Println("end")
	})
	fmt.Println("----------------------")
	time.Sleep(5*time.Second)
}
