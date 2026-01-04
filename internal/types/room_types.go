package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRoomRequest request params
type CreateRoomRequest struct {
	Title string `json:"title" binding:""`
	Desc  string `json:"desc" binding:""`
	Way   string `json:"way" binding:""`
	Mobs  string `json:"mobs" binding:""`
}

// UpdateRoomByIDRequest request params
type UpdateRoomByIDRequest struct {
	ID    string `json:"id" binding:""`
	Title string `json:"title" binding:""`
	Desc  string `json:"desc" binding:""`
	Way   string `json:"way" binding:""`
	Mobs  string `json:"mobs" binding:""`
}

// RoomObjDetail detail
type RoomObjDetail struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Way   string `json:"way"`
	Mobs  string `json:"mobs"`
}

// CreateRoomReply only for api docs
type CreateRoomReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID string `json:"id"`
	} `json:"data"` // return data
}

// DeleteRoomByIDReply only for api docs
type DeleteRoomByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateRoomByIDReply only for api docs
type UpdateRoomByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetRoomByIDReply only for api docs
type GetRoomByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Room RoomObjDetail `json:"room"`
	} `json:"data"` // return data
}

// ListRoomsRequest request params
type ListRoomsRequest struct {
	query.Params
}

// ListRoomsReply only for api docs
type ListRoomsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Rooms []RoomObjDetail `json:"rooms"`
	} `json:"data"` // return data
}
