package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"fs/internal/database"
	"fs/internal/model"
)

func newMobCache() *gotest.Cache {
	record1 := &model.Mob{}
	record1.ID = 1
	record2 := &model.Mob{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewMobCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_mobCache_Set(t *testing.T) {
	c := newMobCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Mob)
	err := c.ICache.(MobCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(MobCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_mobCache_Get(t *testing.T) {
	c := newMobCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Mob)
	err := c.ICache.(MobCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MobCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(MobCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_mobCache_MultiGet(t *testing.T) {
	c := newMobCache()
	defer c.Close()

	var testData []*model.Mob
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Mob))
	}

	err := c.ICache.(MobCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(MobCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.Mob))
	}
}

func Test_mobCache_MultiSet(t *testing.T) {
	c := newMobCache()
	defer c.Close()

	var testData []*model.Mob
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Mob))
	}

	err := c.ICache.(MobCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_mobCache_Del(t *testing.T) {
	c := newMobCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Mob)
	err := c.ICache.(MobCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_mobCache_SetCacheWithNotFound(t *testing.T) {
	c := newMobCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Mob)
	err := c.ICache.(MobCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(MobCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewMobCache(t *testing.T) {
	c := NewMobCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewMobCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewMobCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
