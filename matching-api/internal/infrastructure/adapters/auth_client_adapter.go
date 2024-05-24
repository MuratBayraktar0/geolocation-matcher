package adapters

import (
	"net/http"
)

type AuthClient struct {
	Client  *http.Client
	AuthURL string
}

type TokenResponse struct {
	Token string `json:"token"`
}

func NewAuthClient(authURL string) *AuthClient {
	return &AuthClient{
		Client:  &http.Client{},
		AuthURL: authURL,
	}
}

func (c *AuthClient) GetToken() (string, error) {
	// However, since this is a case study, we are skipping these steps and returning a hardcoded token instead. This allows us to simulate the behavior of the GetToken function without needing to set up and interact with an actual authentication service.
	// The commented out code below demonstrates how you might implement this function in a real-world application.

	// req, err := http.NewRequest("POST", c.AuthURL, nil)
	// if err != nil {
	// 	return "", err
	// }

	// resp, err := c.Client.Do(req)
	// if err != nil {
	// 	return "", err
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return "", err
	// }

	// var tokenResponse TokenResponse
	// err = json.Unmarshal(body, &tokenResponse)
	// if err != nil {
	// 	return "", err
	// }

	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJpc3MiOiJNYXRjaGluZy1BUEkifQ.WG2jCQsxhIGn8SANi1GLbhA0CHBCO14KacGrkngdA_Q", nil
}
