package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"fs/internal/cache"
	"fs/internal/database"
	"fs/internal/model"
)

var _ MobDao = (*mobDao)(nil)

// MobDao defining the dao interface
type MobDao interface {
	Create(ctx context.Context, table *model.Mob) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.Mob) error
	GetByID(ctx context.Context, id uint64) (*model.Mob, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Mob, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Mob) (uint64, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Mob) error
}

type mobDao struct {
	db    *gorm.DB
	cache cache.MobCache      // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewMobDao creating the dao interface
func NewMobDao(db *gorm.DB, xCache cache.MobCache) MobDao {
	if xCache == nil {
		return &mobDao{db: db}
	}
	return &mobDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *mobDao) deleteCache(ctx context.Context, id uint64) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new mob, insert the record and the id value is written back to the table
func (d *mobDao) Create(ctx context.Context, table *model.Mob) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a mob by id
func (d *mobDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Mob{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a mob by id, support partial update
func (d *mobDao) UpdateByID(ctx context.Context, table *model.Mob) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *mobDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Mob) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.MobID != "" {
		update["mob_id"] = table.MobID
	}
	if table.MobName != "" {
		update["mob_name"] = table.MobName
	}
	if table.MobCname != "" {
		update["mob_cname"] = table.MobCname
	}
	if table.MobDesc != "" {
		update["mob_desc"] = table.MobDesc
	}
	if table.Attackable != nil {
		update["attackable"] = table.Attackable
	}
	if table.Hp != 0 {
		update["hp"] = table.Hp
	}
	if table.Mp != 0 {
		update["mp"] = table.Mp
	}
	if table.Attack != 0 {
		update["attack"] = table.Attack
	}
	if table.Defence != 0 {
		update["defence"] = table.Defence
	}
	if table.Dodge != 0 {
		update["dodge"] = table.Dodge
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a mob by id
func (d *mobDao) GetByID(ctx context.Context, id uint64) (*model.Mob, error) {
	// no cache
	if d.cache == nil {
		record := &model.Mob{}
		err := d.db.WithContext(ctx).Where("id = ?", id).First(record).Error
		return record, err
	}

	// get from cache
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	// get from database
	if errors.Is(err, database.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to database
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.Mob{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				if errors.Is(err, database.ErrRecordNotFound) {
					// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.MobExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Mob)
		if !ok {
			return nil, database.ErrRecordNotFound
		}
		return table, nil
	}

	if d.cache.IsPlaceholderErr(err) {
		return nil, database.ErrRecordNotFound
	}

	return nil, err
}

// GetByColumns get a paginated list of mobs by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *mobDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Mob, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.MobColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Mob{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Mob{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *mobDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Mob) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *mobDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Mob{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *mobDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Mob) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
