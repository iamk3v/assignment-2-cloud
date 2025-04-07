package database

import (
	"time"
)

type CacheEntry struct {
	Key       string    `firestore:"key" json:"key"`
	Data      string    `firestore:"data" json:"data"`
	Timestamp time.Time `firestore:"timestamp" json:"timestamp"`
}

const ()

func GetCacheEntry(key string) (*CacheEntry, error) {}

func SetCacheEntry(key string, data interface{}) error {}

func IsCacheVaild(entry *CacheEntry) bool {}

func getCachedData(key string, dest interface{}) error {}
