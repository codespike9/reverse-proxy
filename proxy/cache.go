package proxy

import (
	"sync"
	"time"
)

type cachedItem struct {
	data      []byte
	expiresAt time.Time
}

var (
	cache = make(map[string]cachedItem)
	mu    sync.RWMutex
)

func SaveToCache(key string, data []byte) {
	mu.Lock()
	cache[key] = cachedItem{
		data:      data,
		expiresAt: time.Now().Add(10 * time.Second),
	}
	mu.Unlock()
}

func GetFromCache(key string) ([]byte, bool) {
	mu.RLock()
	item, found := cache[key]
	mu.RUnlock()

	if !found || time.Now().After(item.expiresAt) {
		return nil, false
	}

	return item.data, true
}
