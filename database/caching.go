package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/iterator"
)

/*
CacheEntry Defines structs for cached items
*/
type CacheEntry struct {
	Key       string    `firestore:"key" json:"key"`
	Data      string    `firestore:"data" json:"data"`
	Timestamp time.Time `firestore:"timestamp" json:"timestamp"`
}

const (
	// Firestore collection for caching
	cacheCollection = "cache"
	// CacheExpiration How long a cache entry is valid. Set to 10 hours.
	CacheExpiration = 10 * time.Hour
)

/*
GetCacheEntry Retrieves a cache entry using a key
*/
func GetCacheEntry(key string) (*CacheEntry, error) {
	// Retrieve the document using the key from cache collection
	doc, err := Client.Collection(cacheCollection).Doc(key).Get(Ctx)
	if err != nil {
		return nil, err
	}
	// Make the document into a cacheEntry struct
	var entry CacheEntry
	if err := doc.DataTo(&entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

/*
SetCacheEntry Caches data under a key
*/
func SetCacheEntry(key string, data interface{}) error {
	// Marshal the provided data into JSON
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Create the Cache Entry
	entry := CacheEntry{
		Key:       key,
		Data:      string(bytes),
		Timestamp: time.Now(),
	}
	// Saving the cache entry to Firestore (can overwrite if it exists)
	_, err = Client.Collection(cacheCollection).Doc(key).Set(Ctx, entry)
	return err
}

/*
IsCacheValid Checks if the cache is valid
*/
func IsCacheValid(entry *CacheEntry) bool {
	return time.Since(entry.Timestamp) < CacheExpiration
}

/*
GetCachedData Retrieves the cached data with a key and unmarshals it into dest
*/
func GetCachedData(key string, dest interface{}) error {
	entry, err := GetCacheEntry(key)
	if err != nil {
		return err
	}
	// If still valid
	if !IsCacheValid(entry) {
		return fmt.Errorf("cache is expired")
	}
	// Unmarshal the JSON stored in the cache to the destination
	return json.Unmarshal([]byte(entry.Data), dest)
}

/*
PurgeExpiredCacheEntries Queries firestore for cache entries that have expired and deletes them
*/
func PurgeExpiredCacheEntries(ctx context.Context) error {
	// Calculate expiration time
	expirationThreshold := time.Now().Add(-CacheExpiration)

	// Query the firestore collection for documents with expired timestamps
	iter := Client.Collection(cacheCollection).Where("timestamp", "<", expirationThreshold).Documents(ctx)
	defer iter.Stop()

	var purgeCounter int

	// Iterate over the result
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("there was an error iterating cache documents: %w", err)
		}
		// Delete expired cache documents
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete cache entry %s: %w", doc.Ref.ID, err)
		}
		purgeCounter++
	}
	fmt.Printf("Purged %d cache entries\n", purgeCounter)
	return nil
}
