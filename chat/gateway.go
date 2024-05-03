package chat

type Gateway[Pk any, Resource any] interface {
	Exists(Pk) bool

	Create(*Resource) error
	Read(Pk) *Resource
	ReadAll() []*Resource
	Update(Pk, *Resource) error
	Delete(Pk)
}

type GatewayManager struct {
	userGw Gateway[string, User]
	roomGw Gateway[string, Room]
	msgGw  Gateway[MessageKey, Message]
}

func (gwm *GatewayManager) UserGateway() Gateway[string, User] {
	return gwm.userGw
}

func (gwm *GatewayManager) RoomGateway() Gateway[string, Room] {
	return gwm.roomGw
}

func (gwm *GatewayManager) MessageGateway() Gateway[MessageKey, Message] {
	return gwm.msgGw
}
