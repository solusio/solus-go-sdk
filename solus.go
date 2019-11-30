package solus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	BaseURL     *url.URL
	UserAgent   string
	Email       string
	Password    string
	Credentials Credentials
	Headers     http.Header
	HttpClient  *http.Client
	Logger      *log.Logger
	Retries     int
}

func NewClient(baseURL *url.URL, email, password string) (*Client, error) {
	client := &Client{
		BaseURL:   baseURL,
		UserAgent: "solus.io Go client",
		Email:     email,
		Password:  password,
		Headers: map[string][]string{
			"Accept":       {"application/json"},
			"Content-Type": {"application/json"},
		},
		HttpClient: &http.Client{
			Timeout: time.Second * 35,
		},
		Logger:  log.New(os.Stderr, "", 0),
		Retries: 5,
	}

	authRequest := AuthLoginRequest{
		Email:    email,
		Password: password,
	}

	resp, err := client.AuthLogin(context.Background(), authRequest)
	if err != nil {
		return nil, err
	}

	client.Credentials = resp.Credentials
	client.Headers["Authorization"] = []string{client.Credentials.TokenType + " " + client.Credentials.AccessToken}

	return client, nil
}

func (c *Client) AuthLogin(ctx context.Context, data AuthLoginRequest) (AuthLoginResponse, error) {
	body, code, err := c.request(ctx, "POST", "auth/login", data)
	if err != nil {
		return AuthLoginResponse{}, err
	}

	if code != 200 {
		return AuthLoginResponse{}, fmt.Errorf("HTTP %d: %s", code, body)
	}

	var resp AuthLoginResponseData
	if err := json.Unmarshal(body, &resp); err != nil {
		return AuthLoginResponse{}, fmt.Errorf("failed to decode '%s': %s", body, err)
	}

	return resp.Data, nil
}
