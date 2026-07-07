package cache

import (
	"sync"
)

// MemoryCache โครงสร้างเก็บข้อมูลในหน่วยความจำ
type MemoryCache struct {
	mu   sync.RWMutex
	data []byte
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{}
}

// Set เก็บข้อมูลลงแคช
func (c *MemoryCache) Set(jsonData []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = jsonData
}

// Get ดึงข้อมูลจากแคช
func (c *MemoryCache) Get() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}