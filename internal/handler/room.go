package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/go-dev-frame/sponge/pkg/copier"
	"github.com/go-dev-frame/sponge/pkg/gin/middleware"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
	"github.com/go-dev-frame/sponge/pkg/logger"

	"fs/internal/cache"
	"fs/internal/dao"
	"fs/internal/database"
	"fs/internal/ecode"
	"fs/internal/model"
	"fs/internal/types"
)

var _ RoomHandler = (*roomHandler)(nil)

// RoomHandler defining the handler interface
type RoomHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
}

type roomHandler struct {
	iDao dao.RoomDao
}

// NewRoomHandler creating the handler interface
func NewRoomHandler() RoomHandler {
	return &roomHandler{
		iDao: dao.NewRoomDao(
			database.GetDB(), // db driver is mysql
			cache.NewRoomCache(database.GetCacheType()),
		),
	}
}

// Create a new room
// @Summary Create a new room
// @Description Creates a new room entity using the provided data in the request body.
// @Tags room
// @Accept json
// @Produce json
// @Param data body types.CreateRoomRequest true "room information"
// @Success 200 {object} types.CreateRoomReply{}
// @Router /api/v1/room [post]
// @Security BearerAuth
func (h *roomHandler) Create(c *gin.Context) {
	form := &types.CreateRoomRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	room := &model.Room{}
	err = copier.Copy(room, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRoom)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, room)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": room.ID})
}

// DeleteByID delete a room by id
// @Summary Delete a room by id
// @Description Deletes a existing room identified by the given id in the path.
// @Tags room
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRoomByIDReply{}
// @Router /api/v1/room/{id} [delete]
// @Security BearerAuth
func (h *roomHandler) DeleteByID(c *gin.Context) {
	id, isAbort := getRoomIDFromPath(c)
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

// UpdateByID update a room by id
// @Summary Update a room by id
// @Description Updates the specified room by given id in the path, support partial update.
// @Tags room
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRoomByIDRequest true "room information"
// @Success 200 {object} types.UpdateRoomByIDReply{}
// @Router /api/v1/room/{id} [put]
// @Security BearerAuth
func (h *roomHandler) UpdateByID(c *gin.Context) {
	id, isAbort := getRoomIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoomByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	room := &model.Room{}
	err = copier.Copy(room, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRoom)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, room)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a room by id
// @Summary Get a room by id
// @Description Gets detailed information of a room specified by the given id in the path.
// @Tags room
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRoomByIDReply{}
// @Router /api/v1/room/{id} [get]
// @Security BearerAuth
func (h *roomHandler) GetByID(c *gin.Context) {
	id, isAbort := getRoomIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	room, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RoomObjDetail{}
	err = copier.Copy(data, room)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRoom)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"room": data})
}

// List get a paginated list of rooms by custom conditions
// For more details, please refer to https://go-sponge.com/component/data/custom-page-query.html
// @Summary Get a paginated list of rooms by custom conditions
// @Description Returns a paginated list of rooms based on query filters, including page number and size.
// @Tags room
// @Accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRoomsReply{}
// @Router /api/v1/room/list [post]
// @Security BearerAuth
func (h *roomHandler) List(c *gin.Context) {
	form := &types.ListRoomsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	rooms, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRooms(rooms)
	if err != nil {
		response.Error(c, ecode.ErrListRoom)
		return
	}

	response.Success(c, gin.H{
		"rooms": data,
		"total": total,
	})
}

func getRoomIDFromPath(c *gin.Context) (string, bool) {
	idStr := c.Param("id")

	if idStr == "" {
		logger.Warn("id is empty", middleware.GCtxRequestIDField(c))
		return "", true
	}
	return idStr, false

}

func convertRoom(room *model.Room) (*types.RoomObjDetail, error) {
	data := &types.RoomObjDetail{}
	err := copier.Copy(data, room)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertRooms(fromValues []*model.Room) ([]*types.RoomObjDetail, error) {
	toValues := []*types.RoomObjDetail{}
	for _, v := range fromValues {
		data, err := convertRoom(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
