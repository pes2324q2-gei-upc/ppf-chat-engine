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
	UserGw Gateway[string, User]
	RoomGw Gateway[string, Room]
	MsgGw  Gateway[MessageKey, Message]
}

func (gwm *GatewayManager) UserGateway() Gateway[string, User] {
	return gwm.UserGw
}

func (gwm *GatewayManager) RoomGateway() Gateway[string, Room] {
	return gwm.RoomGw
}

func (gwm *GatewayManager) MessageGateway() Gateway[MessageKey, Message] {
	return gwm.MsgGw
}
