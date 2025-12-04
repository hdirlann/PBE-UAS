package middleware

import (
	"sync"
	"time"
)

type permEntry struct {
	Perms    []string
	ExpireAt time.Time
}

var (
	permCache sync.Map // map[string]*permEntry
	cacheTTL  = 5 * time.Minute
)

func GetCachedPerms(roleID string) ([]string, bool) {
	if v, ok := permCache.Load(roleID); ok {
		if e, ok2 := v.(*permEntry); ok2 {
			if time.Now().Before(e.ExpireAt) {
				return e.Perms, true
			}
			permCache.Delete(roleID)
		}
	}
	return nil, false
}

func SetCachedPerms(roleID string, perms []string) {
	permCache.Store(roleID, &permEntry{Perms: perms, ExpireAt: time.Now().Add(cacheTTL)})
}
