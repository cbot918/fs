package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// room business-level http error codes.
// the roomNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	roomNO       = 1
	roomName     = "room"
	roomBaseCode = errcode.HCode(roomNO)

	ErrCreateRoom     = errcode.NewError(roomBaseCode+1, "failed to create "+roomName)
	ErrDeleteByIDRoom = errcode.NewError(roomBaseCode+2, "failed to delete "+roomName)
	ErrUpdateByIDRoom = errcode.NewError(roomBaseCode+3, "failed to update "+roomName)
	ErrGetByIDRoom    = errcode.NewError(roomBaseCode+4, "failed to get "+roomName+" details")
	ErrListRoom       = errcode.NewError(roomBaseCode+5, "failed to list of "+roomName)

	// error codes are globally unique, adding 1 to the previous error code
)
