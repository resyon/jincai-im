package common

import "errors"

var (
	RoomNotExistError error = errors.New("room not exist")
	RoomNameConflict  error = errors.New("room has exist")
)
