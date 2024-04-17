package chat

import "errors"

var (
	ErrUserNotFound         = errors.New("user not loaded")
	ErrUserUnmarshalFailed  = errors.New("user could not be parsed")
	ErrRoomNotFound         = errors.New("room not loaded")
	ErrUserApiRequestFailed = errors.New("user api request failed")
	ErrRouteUnmarshalFailed = errors.New("route could not be parsed")
	ErrMessageMalformed     = errors.New("message is malformed")
)
