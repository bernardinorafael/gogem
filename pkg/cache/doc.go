// Package cache provides a Redis-backed caching layer with a generic
// cache-aside (GetOrSet) pattern for type-safe cache operations.
//
// Creating a cache client:
//
//	client := cache.New(redisClient, logger)
//
// The GetOrSet pattern fetches from cache first, falling back to the callback
// on a miss and storing the result for subsequent requests:
//
//	user, err := cache.GetOrSet(ctx, cache.SetParams{
//	    Client: client,
//	    Key:    "user:123",
//	    TTL:    5 * time.Minute,
//	}, func() (User, error) {
//	    return db.GetUser(ctx, "123")
//	})
//
// Deleting cache entries:
//
//	err := cache.Delete(ctx, client, "user:123", "user:456")
//
// All values are serialized as JSON. Cache write failures are logged as warnings
// but do not propagate errors, ensuring the primary data source remains the
// source of truth.
package cache
