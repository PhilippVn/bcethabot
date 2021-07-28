package cache

import "sync"

// This is a cache for owners of temporary channels
type CacheOwner struct {
	Cache *sync.Map
}

func (c *CacheOwner) NewCacheOwner() *CacheOwner {
	return &CacheOwner{
		Cache: &sync.Map{},
	}
}
