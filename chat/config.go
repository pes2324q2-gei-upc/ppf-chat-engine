package chat

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Configuration struct {
	debug       bool
	userApiUrl  string
	routeApiUrl string
	creds       Credentials
}

func NewConfiguration(debug bool, userApiUrl, routeApiUrl string, email string, pass string) *Configuration {
	return &Configuration{
		debug:       debug,
		userApiUrl:  userApiUrl,
		routeApiUrl: routeApiUrl,
		creds: Credentials{
			Email:    email,
			Password: pass,
			Token:    "",
		},
	}
}
