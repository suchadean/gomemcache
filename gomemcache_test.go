package gomemcache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	type fields struct {
		lock sync.RWMutex
		data map[string][]byte
	}

	tests := []struct {
		description      string
		key              string
		value            []byte
		ttl              time.Duration
		expected         []byte
		concurrentWrites int
		throwsErr        bool
	}{
		{
			description:      "set value and check if retrieved value is equal",
			key:              "testKey",
			value:            []byte("testValue"),
			ttl:              0,
			expected:         []byte("testValue"),
			concurrentWrites: 1,
			throwsErr:        false,
		},
		{
			description:      "set value concurrent times and check if lock works",
			key:              "testKey",
			value:            []byte("testValue"),
			ttl:              0,
			expected:         []byte("testValue"),
			concurrentWrites: 20,
			throwsErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cache := New()

			if tt.concurrentWrites >= 1 {
				cache.SetValue(tt.key, tt.value, tt.ttl)
			} else {
				for i := 1; i <= tt.concurrentWrites; i++ {
					go cache.SetValue(tt.key, tt.value, tt.ttl)
				}
			}

			actualExists := cache.KeyExists(tt.key)
			actualGet, err := cache.GetValue(tt.key)
			if err != nil {
				if tt.throwsErr {
					assert.Error(t, err, "expected error thrown")
					return
				}

				assert.Fail(t, "unexpected error thrown", err.Error())
				return
			}

			assert.Equal(t, true, actualExists)
			assert.Equal(t, tt.expected, actualGet)

			delErr := cache.DeleteKey(tt.key)
			if delErr != nil {
				assert.Fail(t, "unexpected error thrown", delErr.Error())
			}

			actualNotExists := cache.KeyExists(tt.key)
			assert.Equal(t, false, actualNotExists)
		})
	}
}

func TestSetValueDuration(t *testing.T) {
	t.Run("check ttl and expect key to be not present", func(t *testing.T) {
		cache := New()

		key := "testKey"
		value := []byte("testValue")
		ttl := time.Duration(time.Millisecond * 100)

		cache.SetValue(key, value, ttl)

		time.Sleep(ttl + time.Millisecond*20)

		assert.Equal(t, false, cache.KeyExists(key))
	})
}
