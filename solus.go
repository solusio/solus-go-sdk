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

	s service

	ComputeResources *ComputeResourcesService
	IpBlocks         *IpBlocksService
	Locations        *LocationsService
	OsImages         *OsImagesService
	Plans            *PlansService
	Projects         *ProjectsService
	Servers          *ServersService
	Tasks            *TasksService
	Users            *UsersService
}

type service struct {
	client *Client
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

	client.s.client = client
	client.ComputeResources = (*ComputeResourcesService)(&client.s)
	client.IpBlocks = (*IpBlocksService)(&client.s)
	client.Locations = (*LocationsService)(&client.s)
	client.OsImages = (*OsImagesService)(&client.s)
	client.Plans = (*PlansService)(&client.s)
	client.Projects = (*ProjectsService)(&client.s)
	client.Servers = (*ServersService)(&client.s)
	client.Tasks = (*TasksService)(&client.s)
	client.Users = (*UsersService)(&client.s)

	return client, nil
}

func (c *Client) AuthLogin(ctx context.Context, data AuthLoginRequest) (AuthLoginResponse, error) {
	opts := newRequestOpts()
	opts.body = data
	body, code, err := c.request(ctx, "POST", "auth/login", withBody(opts))
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
