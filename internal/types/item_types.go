package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateItemRequest request params
type CreateItemRequest struct {
	ItemID     string `json:"itemID" binding:""`
	ItemName   string `json:"itemName" binding:""`
	ItemCname  string `json:"itemCname" binding:""`
	ItemDesc   string `json:"itemDesc" binding:""`
	Hp         int    `json:"hp" binding:""`
	Mp         int    `json:"mp" binding:""`
	Attack     int    `json:"attack" binding:""`
	Defence    int    `json:"defence" binding:""`
	Dodge      int    `json:"dodge" binding:""`
	Str        int    `json:"str" binding:""`
	Cor        int    `json:"cor" binding:""`
	Inte       int    `json:"inte" binding:""`
	Dex        int    `json:"dex" binding:""`
	Con        int    `json:"con" binding:""`
	Kar        int    `json:"kar" binding:""`
	Classifier string `json:"classifier" binding:""`
}

// UpdateItemByIDRequest request params
type UpdateItemByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	ItemID     string `json:"itemID" binding:""`
	ItemName   string `json:"itemName" binding:""`
	ItemCname  string `json:"itemCname" binding:""`
	ItemDesc   string `json:"itemDesc" binding:""`
	Hp         int    `json:"hp" binding:""`
	Mp         int    `json:"mp" binding:""`
	Attack     int    `json:"attack" binding:""`
	Defence    int    `json:"defence" binding:""`
	Dodge      int    `json:"dodge" binding:""`
	Str        int    `json:"str" binding:""`
	Cor        int    `json:"cor" binding:""`
	Inte       int    `json:"inte" binding:""`
	Dex        int    `json:"dex" binding:""`
	Con        int    `json:"con" binding:""`
	Kar        int    `json:"kar" binding:""`
	Classifier string `json:"classifier" binding:""`
}

// ItemObjDetail detail
type ItemObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	ItemID     string `json:"itemID"`
	ItemName   string `json:"itemName"`
	ItemCname  string `json:"itemCname"`
	ItemDesc   string `json:"itemDesc"`
	Hp         int    `json:"hp"`
	Mp         int    `json:"mp"`
	Attack     int    `json:"attack"`
	Defence    int    `json:"defence"`
	Dodge      int    `json:"dodge"`
	Str        int    `json:"str"`
	Cor        int    `json:"cor"`
	Inte       int    `json:"inte"`
	Dex        int    `json:"dex"`
	Con        int    `json:"con"`
	Kar        int    `json:"kar"`
	Classifier string `json:"classifier"`
}

// CreateItemReply only for api docs
type CreateItemReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteItemByIDReply only for api docs
type DeleteItemByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateItemByIDReply only for api docs
type UpdateItemByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetItemByIDReply only for api docs
type GetItemByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Item ItemObjDetail `json:"item"`
	} `json:"data"` // return data
}

// ListItemsRequest request params
type ListItemsRequest struct {
	query.Params
}

// ListItemsReply only for api docs
type ListItemsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Items []ItemObjDetail `json:"items"`
	} `json:"data"` // return data
}
