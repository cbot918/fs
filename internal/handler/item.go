package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"fs/internal/cache"
	"fs/internal/dao"
	"fs/internal/database"
	"fs/internal/ecode"
	"fs/internal/model"
	"fs/internal/types"
)

var _ ItemHandler = (*itemHandler)(nil)

// ItemHandler defining the handler interface
type ItemHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type itemHandler struct {
	iDao dao.ItemDao
}

// NewItemHandler creating the handler interface
func NewItemHandler() ItemHandler {
	return &itemHandler{
		iDao: dao.NewItemDao(
			database.GetDB(), // db driver is mysql
			cache.NewItemCache(database.GetCacheType()),
		),
	}
}

// Create a new item
// @Summary Create a new item
// @Description Creates a new item entity using the provided data in the request body.
// @Tags item
// @Accept json
// @Produce json
// @Param data body types.CreateItemRequest true "item information"
// @Success 200 {object} types.CreateItemReply{}
// @Router /api/v1/item [post]
// @Security BearerAuth
func (h *itemHandler) Create(c *gin.Context) {
	form := &types.CreateItemRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	item := &model.Item{}
	err = copier.Copy(item, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateItem)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, item)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": item.ID})
}

// DeleteByID delete a item by id
// @Summary Delete a item by id
// @Description Deletes a existing item identified by the given id in the path.
// @Tags item
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteItemByIDReply{}
// @Router /api/v1/item/{id} [delete]
// @Security BearerAuth
func (h *itemHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getItemIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update a item by id
// @Summary Update a item by id
// @Description Updates the specified item by given id in the path, support partial update.
// @Tags item
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateItemByIDRequest true "item information"
// @Success 200 {object} types.UpdateItemByIDReply{}
// @Router /api/v1/item/{id} [put]
// @Security BearerAuth
func (h *itemHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getItemIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateItemByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	item := &model.Item{}
	err = copier.Copy(item, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDItem)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, item)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a item by id
// @Summary Get a item by id
// @Description Gets detailed information of a item specified by the given id in the path.
// @Tags item
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetItemByIDReply{}
// @Router /api/v1/item/{id} [get]
// @Security BearerAuth
func (h *itemHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getItemIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	item, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.ItemObjDetail{}
	err = copier.Copy(data, item)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDItem)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"item": data})
}

// List get a paginated list of items by custom conditions
// @Summary Get a paginated list of items by custom conditions
// @Description Returns a paginated list of item based on query filters, including page number and size.
// @Tags item
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListItemsReply{}
// @Router /api/v1/item/list [post]
// @Security BearerAuth
func (h *itemHandler) List(c *gin.Context) {
	form := &types.ListItemsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	items, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertItems(items)
	if err != nil {
		response.Error(c, ecode.ErrListItem)
		return
	}

	response.Success(c, gin.H{
		"items": data,
		"total": total,
	})
}

func getItemIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertItem(item *model.Item) (*types.ItemObjDetail, error) {
	data := &types.ItemObjDetail{}
	err := copier.Copy(data, item)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertItems(fromValues []*model.Item) ([]*types.ItemObjDetail, error) {
	toValues := []*types.ItemObjDetail{}
	for _, v := range fromValues {
		data, err := convertItem(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
