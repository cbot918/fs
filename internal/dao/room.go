package dao

import (
	"context"
	"errors"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"

	"fs/internal/cache"
	"fs/internal/database"
	"fs/internal/model"
)

var _ RoomDao = (*roomDao)(nil)

// RoomDao defining the dao interface
type RoomDao interface {
	Create(ctx context.Context, table *model.Room) error
	DeleteByID(ctx context.Context, id string) error
	UpdateByID(ctx context.Context, table *model.Room) error
	GetByID(ctx context.Context, id string) (*model.Room, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.Room, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Room) (string, error)
	DeleteByTx(ctx context.Context, tx *gorm.DB, id string) error
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Room) error
}

type roomDao struct {
	db    *gorm.DB
	cache cache.RoomCache     // if nil, the cache is not used.
	sfg   *singleflight.Group // if cache is nil, the sfg is not used.
}

// NewRoomDao creating the dao interface
func NewRoomDao(db *gorm.DB, xCache cache.RoomCache) RoomDao {
	if xCache == nil {
		return &roomDao{db: db}
	}
	return &roomDao{
		db:    db,
		cache: xCache,
		sfg:   new(singleflight.Group),
	}
}

func (d *roomDao) deleteCache(ctx context.Context, id string) error {
	if d.cache != nil {
		return d.cache.Del(ctx, id)
	}
	return nil
}

// Create a new room, insert the record and the id value is written back to the table
func (d *roomDao) Create(ctx context.Context, table *model.Room) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID delete a room by id
func (d *roomDao) DeleteByID(ctx context.Context, id string) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Room{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByID update a room by id
func (d *roomDao) UpdateByID(ctx context.Context, table *model.Room) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}

func (d *roomDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.Room) error {
	if table.ID == "" {
		return errors.New("id cannot be empty")
	}

	update := map[string]interface{}{}

	if table.Title != "" {
		update["title"] = table.Title
	}
	if table.Desc != "" {
		update["desc"] = table.Desc
	}
	if table.Way != "" {
		update["way"] = table.Way
	}
	if table.Mobs != "" {
		update["mobs"] = table.Mobs
	}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a room by id
func (d *roomDao) GetByID(ctx context.Context, id string) (*model.Room, error) {
	// no cache
	if d.cache == nil {
		record := &model.Room{}
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
		val, err, _ := d.sfg.Do(id, func() (interface{}, error) {

			table := &model.Room{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// set placeholder cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, database.ErrRecordNotFound) {
					if err = d.cache.SetPlaceholder(ctx, id); err != nil {
						logger.Warn("cache.SetPlaceholder error", logger.Err(err), logger.Any("id", id))
					}
					return nil, database.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			if err = d.cache.Set(ctx, id, table, cache.RoomExpireTime); err != nil {
				logger.Warn("cache.Set error", logger.Err(err), logger.Any("id", id))
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.Room)
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

// GetByColumns get a paginated list of rooms by custom conditions.
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
func (d *roomDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.Room, int64, error) {
	if params.Sort == "" {
		params.Sort = "-id"
	}
	queryStr, args, err := params.ConvertToGormConditions(query.WithWhitelistNames(model.RoomColumnNames))
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.Room{}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.Room{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *roomDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.Room) (string, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.ID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *roomDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id string) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Room{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.deleteCache(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *roomDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.Room) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	_ = d.deleteCache(ctx, table.ID)

	return err
}
