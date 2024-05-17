package cache

import (
	"database/sql"
	"log"
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

func LoadUserCache(db *sql.DB) {
	sqlQuery := "SELECT email FROM user_details"
	rows, err := db.Query(sqlQuery)
	defer rows.Close()
	if err != nil {
		log.Panic("Unable to get user Cache")
	}
	for rows.Next() {
		var email string
		rows.Scan(&email)
		AddUserToCache(email)
	}

}
