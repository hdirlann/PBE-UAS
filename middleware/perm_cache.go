package middleware

import (
	"sync"
	"time"
)

type cacheItem struct {
	perms  []string
	expire time.Time
}

var (
	permCache   = map[string]cacheItem{}
	permCacheMu sync.RWMutex
	cacheTTL    = 5 * time.Minute
)

// GetCachedPerms returns perms and true if present + not expired
func GetCachedPerms(roleID string) ([]string, bool) {
	permCacheMu.RLock()
	defer permCacheMu.RUnlock()
	it, ok := permCache[roleID]
	if !ok {
		return nil, false
	}
	if time.Now().After(it.expire) {
		return nil, false
	}
	return it.perms, true
}

func SetCachedPerms(roleID string, perms []string) {
	permCacheMu.Lock()
	defer permCacheMu.Unlock()
	permCache[roleID] = cacheItem{
		perms:  perms,
		expire: time.Now().Add(cacheTTL),
	}
}

func InvalidateCachedPerms(roleID string) {
	permCacheMu.Lock()
	delete(permCache, roleID)
	permCacheMu.Unlock()
}
