package chat

import "errors"

var (
	ErrUserNotFound         = errors.New("user not loaded")
	ErrUserUnmarshalFailed  = errors.New("user could not be parsed")
	ErrRoomNotFound         = errors.New("room not loaded")
	ErrUserApiRequestFailed = errors.New("user api request failed")
	ErrRouteUnmarshalFailed = errors.New("route could not be parsed")
	ErrUnknownCommand       = errors.New("unknown command")
	ErrMessageMalformed     = errors.New("message is malformed")
	ErrNoMessageId          = errors.New("message id is missing")
	ErrWrongMessageSender   = errors.New("wrong message sender")
)
