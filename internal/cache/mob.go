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
	mobCachePrefixKey = "mob:"
	// MobExpireTime expire time
	MobExpireTime = 5 * time.Minute
)

var _ MobCache = (*mobCache)(nil)

// MobCache cache interface
type MobCache interface {
	Set(ctx context.Context, id uint64, data *model.Mob, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Mob, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Mob, error)
	MultiSet(ctx context.Context, data []*model.Mob, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// mobCache define a cache struct
type mobCache struct {
	cache cache.Cache
}

// NewMobCache new a cache
func NewMobCache(cacheType *database.CacheType) MobCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Mob{}
		})
		return &mobCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Mob{}
		})
		return &mobCache{cache: c}
	}

	return nil // no cache
}

// GetMobCacheKey cache key
func (c *mobCache) GetMobCacheKey(id uint64) string {
	return mobCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *mobCache) Set(ctx context.Context, id uint64, data *model.Mob, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetMobCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *mobCache) Get(ctx context.Context, id uint64) (*model.Mob, error) {
	var data *model.Mob
	cacheKey := c.GetMobCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *mobCache) MultiSet(ctx context.Context, data []*model.Mob, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetMobCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *mobCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Mob, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetMobCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Mob)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Mob)
	for _, id := range ids {
		val, ok := itemMap[c.GetMobCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *mobCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetMobCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *mobCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetMobCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *mobCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
