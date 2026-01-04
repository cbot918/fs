package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"fs/internal/database"
	"fs/internal/model"
)

const (
	// cache prefix key, must end with a colon
	itemCachePrefixKey = "item:"
	// ItemExpireTime expire time
	ItemExpireTime = 5 * time.Minute
)

var _ ItemCache = (*itemCache)(nil)

// ItemCache cache interface
type ItemCache interface {
	Set(ctx context.Context, id uint64, data *model.Item, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Item, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Item, error)
	MultiSet(ctx context.Context, data []*model.Item, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// itemCache define a cache struct
type itemCache struct {
	cache cache.Cache
}

// NewItemCache new a cache
func NewItemCache(cacheType *database.CacheType) ItemCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Item{}
		})
		return &itemCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Item{}
		})
		return &itemCache{cache: c}
	}

	return nil // no cache
}

// GetItemCacheKey cache key
func (c *itemCache) GetItemCacheKey(id uint64) string {
	return itemCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *itemCache) Set(ctx context.Context, id uint64, data *model.Item, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetItemCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *itemCache) Get(ctx context.Context, id uint64) (*model.Item, error) {
	var data *model.Item
	cacheKey := c.GetItemCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *itemCache) MultiSet(ctx context.Context, data []*model.Item, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetItemCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *itemCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Item, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetItemCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Item)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Item)
	for _, id := range ids {
		val, ok := itemMap[c.GetItemCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *itemCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetItemCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *itemCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetItemCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *itemCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
