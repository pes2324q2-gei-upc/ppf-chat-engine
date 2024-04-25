package chat

type Route struct {
	Id         string  `json:"id"`
	Name       string  `json:"destinationAlias"`
	Driver     *User   `json:"driver"`
	Passengers []*User `json:"passengers"`
}

// func (u *Route) UnmarshalJSON(data []byte) error {
// 	type Alias Route
// 	aux := &struct {
// 		Id         string  `json:"id"`
// 		Name       string  `json:"destinationAlias"`
// 		Driver     *User   `json:"driver"`
// 		Passengers []*User `json:"passengers"`
// 		*Alias
// 	}{
// 		Alias:      (*Alias)(u),
// 		Id:         u.Id,
// 		Name:       u.Name,
// 		Driver:     u.Driver,
// 		Passengers: u.Passengers,
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	return nil
// }
