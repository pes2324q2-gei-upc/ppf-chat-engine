package chat

type Route struct {
	Id         string `json:"id"`
	Name       string `json:"destinationAlias"`
	Driver     User   `json:"driver"`
	Passengers []User `json:"passengers"`
}
