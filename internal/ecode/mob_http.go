package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// mob business-level http error codes.
// the mobNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	mobNO       = 39
	mobName     = "mob"
	mobBaseCode = errcode.HCode(mobNO)

	ErrCreateMob     = errcode.NewError(mobBaseCode+1, "failed to create "+mobName)
	ErrDeleteByIDMob = errcode.NewError(mobBaseCode+2, "failed to delete "+mobName)
	ErrUpdateByIDMob = errcode.NewError(mobBaseCode+3, "failed to update "+mobName)
	ErrGetByIDMob    = errcode.NewError(mobBaseCode+4, "failed to get "+mobName+" details")
	ErrListMob       = errcode.NewError(mobBaseCode+5, "failed to list of "+mobName)

	// error codes are globally unique, adding 1 to the previous error code
)
