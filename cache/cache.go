package cache

import (
	"sync"
)

// Define the UserCache type as a map with string keys and boolean values
type UserCache map[string]bool

// Declare the global variable of type UserCache and a mutex
var (
    Cache      UserCache
    cacheMutex sync.Mutex
)

func init() {
    // Initialize the cache
    Cache = make(UserCache)
}

// AddUser adds a user to the cache
func AddUserToCache(email string) {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    Cache[email] = true
}

// IsUserInCache checks if a user is in the cache
func IsUserInCache(email string) bool {
    cacheMutex.Lock()
    defer cacheMutex.Unlock()
    return Cache[email]
}
