package chat

import "errors"

var ErrUserNotFound = errors.New("user not loaded")
var ErrUserUnmarshalFailed = errors.New("user could not be parsed")
var ErrRoomNotFound = errors.New("room not loaded")
var ErrUserApiRequestFailed = errors.New("user api request failed")
