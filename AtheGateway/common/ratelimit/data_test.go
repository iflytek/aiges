package ratelimit

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestLocalMetricData_Add(t *testing.T) {
	var a = LocalMetricDataService{data: map[string]int{}}
	for i := 0; i < 100000; i++ {
		go a.Add("haha")
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println(a.Add("haha"))

}

func TestLocalMetricDataService_Add(t *testing.T) {
	s:=RedisMetricDataService{redis:redis.NewClient(&redis.Options{
		Addr:"10.1.87.70:9054",
	})}

	fmt.Println(s.AddDelta("ha2",4))
	fmt.Println(s.Add("ha2"))
	fmt.Println(s.Release("ha2"))
}

func TestLocalMetricDataService_AddDelta(t *testing.T) {
	type fields struct {
		data map[string]int
		lock sync.Mutex
	}
	type args struct {
		key  string
		deta int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LocalMetricDataService{
				data: tt.fields.data,
				lock: tt.fields.lock,
			}
			if err := l.AddDelta(tt.args.key, tt.args.deta); (err != nil) != tt.wantErr {
				t.Errorf("LocalMetricDataService.AddDelta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedisMetricDataService_Add(t *testing.T) {
	type fields struct {
		redis redis.Cmdable
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			this := &RedisMetricDataService{
				redis: tt.fields.redis,
			}
			got, err := this.Add(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("RedisMetricDataService.Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RedisMetricDataService.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}
