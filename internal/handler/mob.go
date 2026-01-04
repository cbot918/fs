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

var _ MobHandler = (*mobHandler)(nil)

// MobHandler defining the handler interface
type MobHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type mobHandler struct {
	iDao dao.MobDao
}

// NewMobHandler creating the handler interface
func NewMobHandler() MobHandler {
	return &mobHandler{
		iDao: dao.NewMobDao(
			database.GetDB(), // db driver is mysql
			cache.NewMobCache(database.GetCacheType()),
		),
	}
}

// Create a new mob
// @Summary Create a new mob
// @Description Creates a new mob entity using the provided data in the request body.
// @Tags mob
// @Accept json
// @Produce json
// @Param data body types.CreateMobRequest true "mob information"
// @Success 200 {object} types.CreateMobReply{}
// @Router /api/v1/mob [post]
// @Security BearerAuth
func (h *mobHandler) Create(c *gin.Context) {
	form := &types.CreateMobRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	mob := &model.Mob{}
	err = copier.Copy(mob, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateMob)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, mob)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": mob.ID})
}

// DeleteByID delete a mob by id
// @Summary Delete a mob by id
// @Description Deletes a existing mob identified by the given id in the path.
// @Tags mob
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteMobByIDReply{}
// @Router /api/v1/mob/{id} [delete]
// @Security BearerAuth
func (h *mobHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getMobIDFromPath(c)
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

// UpdateByID update a mob by id
// @Summary Update a mob by id
// @Description Updates the specified mob by given id in the path, support partial update.
// @Tags mob
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateMobByIDRequest true "mob information"
// @Success 200 {object} types.UpdateMobByIDReply{}
// @Router /api/v1/mob/{id} [put]
// @Security BearerAuth
func (h *mobHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getMobIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateMobByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	mob := &model.Mob{}
	err = copier.Copy(mob, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDMob)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, mob)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a mob by id
// @Summary Get a mob by id
// @Description Gets detailed information of a mob specified by the given id in the path.
// @Tags mob
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetMobByIDReply{}
// @Router /api/v1/mob/{id} [get]
// @Security BearerAuth
func (h *mobHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getMobIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	mob, err := h.iDao.GetByID(ctx, id)
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

	data := &types.MobObjDetail{}
	err = copier.Copy(data, mob)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDMob)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"mob": data})
}

// List get a paginated list of mobs by custom conditions
// @Summary Get a paginated list of mobs by custom conditions
// @Description Returns a paginated list of mob based on query filters, including page number and size.
// @Tags mob
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListMobsReply{}
// @Router /api/v1/mob/list [post]
// @Security BearerAuth
func (h *mobHandler) List(c *gin.Context) {
	form := &types.ListMobsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	mobs, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertMobs(mobs)
	if err != nil {
		response.Error(c, ecode.ErrListMob)
		return
	}

	response.Success(c, gin.H{
		"mobs":  data,
		"total": total,
	})
}

func getMobIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertMob(mob *model.Mob) (*types.MobObjDetail, error) {
	data := &types.MobObjDetail{}
	err := copier.Copy(data, mob)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertMobs(fromValues []*model.Mob) ([]*types.MobObjDetail, error) {
	toValues := []*types.MobObjDetail{}
	for _, v := range fromValues {
		data, err := convertMob(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
