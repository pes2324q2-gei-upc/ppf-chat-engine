package chat

type GatewayManager struct {
	userGw UserGateway
	roomGw RoomGateway
	msgGw  MessageGateway
}

func (gwm *GatewayManager) RoomGateway() *RoomGateway {
	return &gwm.roomGw
}

func (gwm *GatewayManager) UserGateway() *UserGateway {
	return &gwm.userGw
}

func (gwm *GatewayManager) MessageGateway() *MessageGateway {
	return &gwm.msgGw
}
