package utils

import (
	"sync"
	"time"
)

var tokenBlacklist sync.Map

func init() {
	go startBlacklistCleaner()
}

func AddToBlacklist(token string, expiresAt time.Time) {
	if token == "" {
		return
	}
	if expiresAt.IsZero() || expiresAt.Before(time.Now()) {
		expiresAt = time.Now().Add(1 * time.Hour)
	}
	tokenBlacklist.Store(token, expiresAt)
}

func IsBlacklisted(token string) bool {
	if token == "" {
		return false
	}
	val, exists := tokenBlacklist.Load(token)
	if !exists {
		return false
	}

	expiresAt, ok := val.(time.Time)
	if !ok {
		return false
	}
	if time.Now().After(expiresAt) {
		tokenBlacklist.Delete(token)
		return false
	}

	return true
}

func startBlacklistCleaner() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		now := time.Now()
		tokenBlacklist.Range(func(key, value any) bool {
			if exp, ok := value.(time.Time); ok {
				if now.After(exp) {
					tokenBlacklist.Delete(key)
				}
			}
			return true
		})
	}
}
