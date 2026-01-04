package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateMobRequest request params
type CreateMobRequest struct {
	MobID      string `json:"mobID" binding:""`
	MobName    string `json:"mobName" binding:""`
	MobCname   string `json:"mobCname" binding:""`
	MobDesc    string `json:"mobDesc" binding:""`
	Attackable *bool  `json:"attackable" binding:""`
	Hp         int    `json:"hp" binding:""`
	Mp         int    `json:"mp" binding:""`
	Attack     int    `json:"attack" binding:""`
	Defence    int    `json:"defence" binding:""`
	Dodge      int    `json:"dodge" binding:""`
}

// UpdateMobByIDRequest request params
type UpdateMobByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	MobID      string `json:"mobID" binding:""`
	MobName    string `json:"mobName" binding:""`
	MobCname   string `json:"mobCname" binding:""`
	MobDesc    string `json:"mobDesc" binding:""`
	Attackable *bool  `json:"attackable" binding:""`
	Hp         int    `json:"hp" binding:""`
	Mp         int    `json:"mp" binding:""`
	Attack     int    `json:"attack" binding:""`
	Defence    int    `json:"defence" binding:""`
	Dodge      int    `json:"dodge" binding:""`
}

// MobObjDetail detail
type MobObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	MobID      string `json:"mobID"`
	MobName    string `json:"mobName"`
	MobCname   string `json:"mobCname"`
	MobDesc    string `json:"mobDesc"`
	Attackable *bool  `json:"attackable"`
	Hp         int    `json:"hp"`
	Mp         int    `json:"mp"`
	Attack     int    `json:"attack"`
	Defence    int    `json:"defence"`
	Dodge      int    `json:"dodge"`
}

// CreateMobReply only for api docs
type CreateMobReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteMobByIDReply only for api docs
type DeleteMobByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateMobByIDReply only for api docs
type UpdateMobByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetMobByIDReply only for api docs
type GetMobByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Mob MobObjDetail `json:"mob"`
	} `json:"data"` // return data
}

// ListMobsRequest request params
type ListMobsRequest struct {
	query.Params
}

// ListMobsReply only for api docs
type ListMobsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Mobs []MobObjDetail `json:"mobs"`
	} `json:"data"` // return data
}
