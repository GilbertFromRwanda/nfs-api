package auth

import (
	"sync"
	"time"
)

// tokenBlacklist stores revoked JWT tokens until they expire naturally.
var tokenBlacklist = struct {
	sync.RWMutex
	tokens map[string]time.Time
}{tokens: make(map[string]time.Time)}

// BlacklistToken adds a token to the blacklist with its expiration time.
func BlacklistToken(token string, expAt time.Time) {
	tokenBlacklist.Lock()
	defer tokenBlacklist.Unlock()
	tokenBlacklist.tokens[token] = expAt
}

// IsBlacklisted checks whether a token has been revoked.
func IsBlacklisted(token string) bool {
	tokenBlacklist.RLock()
	defer tokenBlacklist.RUnlock()
	_, exists := tokenBlacklist.tokens[token]
	return exists
}

// CleanupBlacklist removes expired tokens from the blacklist periodically.
func CleanupBlacklist(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			now := time.Now()
			tokenBlacklist.Lock()
			for token, expAt := range tokenBlacklist.tokens {
				if now.After(expAt) {
					delete(tokenBlacklist.tokens, token)
				}
			}
			tokenBlacklist.Unlock()
		}
	}()
}
