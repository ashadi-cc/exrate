package memory

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_memoryStore_Set(t *testing.T) {
	type fields struct {
		mu   *sync.RWMutex
		item map[string]interface{}
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "set",
			fields:    fields{&sync.RWMutex{}, make(map[string]interface{})},
			args:      args{"1", 1},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &memoryStore{
				mu:   tt.fields.mu,
				item: tt.fields.item,
			}
			tt.assertion(t, store.Set(tt.args.key, tt.args.value))
		})
	}
}

func Test_memoryStore_Get(t *testing.T) {
	type fields struct {
		mu   *sync.RWMutex
		item map[string]interface{}
	}
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      interface{}
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "key exist",
			fields:    fields{&sync.RWMutex{}, map[string]interface{}{"2": true}},
			args:      args{"2"},
			want:      true,
			assertion: assert.NoError,
		},
		{
			name:      "not exist",
			fields:    fields{&sync.RWMutex{}, map[string]interface{}{"2": true}},
			args:      args{"4"},
			want:      nil,
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &memoryStore{
				mu:   tt.fields.mu,
				item: tt.fields.item,
			}
			got, err := store.Get(tt.args.key)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
