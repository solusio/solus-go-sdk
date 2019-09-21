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
)

type Client struct {
	BaseURL    *url.URL
	UserAgent  string
	Email      string
	Password   string
	Token      string
	Headers    http.Header
	HttpClient *http.Client
	Logger     *log.Logger
}

func NewClient(baseURL *url.URL, email, password string) *Client {
	return &Client{
		BaseURL:   baseURL,
		UserAgent: "solus.io Go client",
		Email:     email,
		Password:  password,
		Headers: map[string][]string{
			"Accept":       {"application/json"},
			"Content-Type": {"application/json"},
		},
		HttpClient: &http.Client{
			Timeout: 35,
		},
		Logger: log.New(os.Stderr, "", 0),
	}

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

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL.String()+path, reqBody)
	if err != nil {
		return nil, 0, err
	}

	for k, values := range c.Headers {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.Logger.Println("failed to close body", method, path)
		}
	}()

	c.Logger.Printf("%#v", req)
	c.Logger.Println(method, c.BaseURL.String()+path, string(bodyByte))

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
	body, code, err := c.request(ctx, "POST", "/auth/login", data)
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
