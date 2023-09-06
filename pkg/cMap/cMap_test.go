package cMap

import (
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
			cm := New()
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
