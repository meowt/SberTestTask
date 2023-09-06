package cMap

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestCMap_Put(t *testing.T) {
	type args struct {
		key string
		v   interface{}
		ttl time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "concurrent writing&reading",
			args: args{
				key: "key",
				v:   "value",
				ttl: time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := New(100 * time.Millisecond)
			for i := 0; i < 1000; i++ {
				go func() {
					cm.Put(tt.args.key, tt.args.v, tt.args.ttl)
					t.Logf("put %v\n", tt.args.v)
				}()
				go func() {
					t.Logf("got: %v\n", cm.Get(tt.args.key))
				}()
			}
		})
	}
}

func TestCMap_Get(t *testing.T) {
	type args struct {
		key   string
		value interface{}
		ttl   time.Duration
	}
	tests := []struct {
		name        string
		args        args
		wantValue   interface{}
		timeToSleep time.Duration
	}{
		// TODO: Add test cases.
		{
			args:      args{key: "1", value: "one"},
			wantValue: "one",
		},
		{
			name:        "expired value",
			args:        args{key: "1", value: "one", ttl: 200 * time.Millisecond},
			wantValue:   nil,
			timeToSleep: time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := New(100 * time.Millisecond)
			cm.Put(tt.args.key, tt.args.value, tt.args.ttl)
			time.Sleep(tt.timeToSleep)
			if gotValue := cm.Get(tt.args.key); !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("Get() = %v, want %v", gotValue, tt.wantValue)
			}
		})
	}
}

func TestCMap_StoreToFile(t *testing.T) {
	type args struct {
		interval time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{args: args{interval: time.Second}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Creating new map
			cm := New(100 * time.Millisecond)

			//Inserting some data
			cm.Put("1", "one", 0)
			cm.Put("2", "two", 2*time.Second)
			cm.Put("3", "three", 0)
			cm.Put("4", "four", 0)
			cm.Put("5", "five", 10*time.Second)

			//Opening a storage file
			f, err := os.OpenFile("cMapStorage.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer f.Close()

			//Executing function
			err = cm.StoreToFile(f, 100*time.Millisecond)

			//Checking the result
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			time.Sleep(2 * time.Second)
			//Random exiting case
			//os.Exit(10)
		})
	}
}

func TestCMap_LoadFromFile(t *testing.T) {
	type fields struct {
		m  map[string]Value
		mu sync.Mutex
	}
	type args struct {
		f *os.File
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CMap
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Opening a storage file
			f, err := os.OpenFile("cMapStorage.json", os.O_RDWR, 0755)
			defer f.Close()

			//Executing function
			cm, err := New(100 * time.Millisecond).LoadFromFile(f)

			//Checking results
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%v", cm)
		})
	}
}
