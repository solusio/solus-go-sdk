package solus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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
		Logger: log.New(os.Stderr, "", 0),
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

func (c *Client) request(ctx context.Context, method, path string, body interface{}) ([]byte, int, error) {
	var bodyByte []byte
	var reqBody io.ReadWriter
	if body != nil {
		var err error
		bodyByte, err = json.Marshal(body)
		if err != nil {
			return nil, 0, err
		}
		reqBody = bytes.NewBuffer(bodyByte)
	}

	fullUrl, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, method, fullUrl.String(), reqBody)
	if err != nil {
		return nil, 0, err
	}

	for k, values := range c.Headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	req.Header.Set("User-Agent", c.UserAgent)

	c.Logger.Println(method, fullUrl.String(), string(bodyByte))

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Println("failed to close body", method, path)
		}
	}()

	code := resp.StatusCode

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if code >= 404 {
		return respBody, code, fmt.Errorf("HTTP %d: %s", code, respBody)
	}

	return respBody, code, nil
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
