package chat

import (
	db "github.com/pes2324q2-gei-upc/ppf-chat-engine/persist"
)

// Gateways are in charge of mapping domain types to Record types.
type Gateway interface {
	ToRecord(any) db.Record
	ToDomain(db.Record) any
}
