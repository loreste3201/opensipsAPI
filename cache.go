package main

import (
	"github.com/bradfitz/gomemcache/memcache"
)

//CacheAvailable show the status of caching server
var CacheAvailable = false
var cacheClient *memcache.Client

//StopCache ()
func StopCache() {
	cacheClient = nil
	CacheAvailable = false
}

//InitializeCache connect memcache server
func InitializeCache() {
	value, ok := Config["CACHE_URIS"]
	if ok {
		cacheClient = memcache.New(value)

		err := cacheClient.Set(&memcache.Item{Key: "test", Value: []byte("test cache")})
		if err != nil {
			Error("Cache not connected. Error: " + err.Error())
			return
		}

		CacheAvailable = true
		Debug("Cache connected")
		ClearCache()
	}
}

//CacheSet (Context string, ID string, Value string) add/update data to cache
func CacheSet(ItemID string, Data []byte) bool {
	if CacheAvailable {
		cacheClient.Set(&memcache.Item{Key: ItemID, Value: Data})
		return true
	}
	Error("CachePut: Cache not available")
	return false
}

//CacheGet (Context string, ID string) return (data as []bytes, status as bool)
func CacheGet(ItemID string) ([]byte, bool) {
	if !CacheAvailable {
		Error("CacheGet: Cache not available")
		return nil, false
	}
	Data, err := cacheClient.Get(ItemID)
	if err == nil {
		return Data.Value, true
	}
	Error("CacheGet: Error getting cache for " + ItemID)
	return nil, false
}

//CacheDelete (ItemID string) delete item from cache
func CacheDelete(ItemID string) bool {
	if CacheAvailable {
		_ = cacheClient.Delete(ItemID)
		return true
	}
	Error("CacheGet: Cache not available")
	return false
}

//ClearCache API
func ClearCache() {
	cacheClient.FlushAll()
}
