package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type UserApiCredentials struct {
	AuthUrl  *url.URL
	Email    string `json:"email"`
	Password string `json:"password"`
	token    *string
}

func bytesReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}

func (creds UserApiCredentials) newLoginRequest() (*http.Response, error) {
	bytes, err := json.Marshal(creds)
	if err != nil {
		return nil, fmt.Errorf("could not marshal credentials: %w", err)
	}
	return http.Post(
		creds.AuthUrl.String(),
		"application/json",
		bytesReader(bytes),
	)
}

func (creds UserApiCredentials) Login() error {
	response, err := creds.newLoginRequest()
	if err != nil {
		return fmt.Errorf("could not build login request: %w", err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	t := struct {
		Token string `json:"token"`
	}{}
	if err := json.Unmarshal(body, &t); err != nil {
		return fmt.Errorf("could not unmarshall token: %w", err)
	}
	creds.token = &t.Token
	return nil
}

func (creds UserApiCredentials) Token() string {
	return *creds.token
}
