package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// item business-level http error codes.
// the itemNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	itemNO       = 27
	itemName     = "item"
	itemBaseCode = errcode.HCode(itemNO)

	ErrCreateItem     = errcode.NewError(itemBaseCode+1, "failed to create "+itemName)
	ErrDeleteByIDItem = errcode.NewError(itemBaseCode+2, "failed to delete "+itemName)
	ErrUpdateByIDItem = errcode.NewError(itemBaseCode+3, "failed to update "+itemName)
	ErrGetByIDItem    = errcode.NewError(itemBaseCode+4, "failed to get "+itemName+" details")
	ErrListItem       = errcode.NewError(itemBaseCode+5, "failed to list of "+itemName)

	// error codes are globally unique, adding 1 to the previous error code
)
