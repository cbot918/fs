package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"

	"fs/internal/database"
	"fs/internal/model"
)

const (
	// cache prefix key, must end with a colon
	roomCachePrefixKey = "room:"
	// RoomExpireTime expire time
	RoomExpireTime = 5 * time.Minute
)

var _ RoomCache = (*roomCache)(nil)

// RoomCache cache interface
type RoomCache interface {
	Set(ctx context.Context, id string, data *model.Room, duration time.Duration) error
	Get(ctx context.Context, id string) (*model.Room, error)
	MultiGet(ctx context.Context, ids []string) (map[string]*model.Room, error)
	MultiSet(ctx context.Context, data []*model.Room, duration time.Duration) error
	Del(ctx context.Context, id string) error
	SetPlaceholder(ctx context.Context, id string) error
	IsPlaceholderErr(err error) bool
}

// roomCache define a cache struct
type roomCache struct {
	cache cache.Cache
}

// NewRoomCache new a cache
func NewRoomCache(cacheType *database.CacheType) RoomCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Room{}
		})
		return &roomCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Room{}
		})
		return &roomCache{cache: c}
	}

	return nil // no cache
}

// GetRoomCacheKey cache key
func (c *roomCache) GetRoomCacheKey(id string) string {
	return roomCachePrefixKey + id
}

// Set write to cache
func (c *roomCache) Set(ctx context.Context, id string, data *model.Room, duration time.Duration) error {
	if data == nil {
		return nil
	}
	cacheKey := c.GetRoomCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *roomCache) Get(ctx context.Context, id string) (*model.Room, error) {
	var data *model.Room
	cacheKey := c.GetRoomCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *roomCache) MultiSet(ctx context.Context, data []*model.Room, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRoomCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *roomCache) MultiGet(ctx context.Context, ids []string) (map[string]*model.Room, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRoomCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Room)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.Room)
	for _, id := range ids {
		val, ok := itemMap[c.GetRoomCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *roomCache) Del(ctx context.Context, id string) error {
	cacheKey := c.GetRoomCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *roomCache) SetPlaceholder(ctx context.Context, id string) error {
	cacheKey := c.GetRoomCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *roomCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
