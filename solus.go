//go:generate go run paginatorgen.go

package solus

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	BaseURL     *url.URL
	UserAgent   string
	Credentials Credentials
	Headers     http.Header
	HttpClient  *http.Client
	Logger      Logger
	Retries     int

	s service

	Account           *AccountService
	BackupNodes       *BackupNodesService
	Backups           *BackupsService
	ComputeResources  *ComputeResourcesService
	IPBlocks          *IPBlocksService
	License           *LicenseService
	Locations         *LocationsService
	OsImages          *OsImagesService
	Applications      *ApplicationsService
	Plans             *PlansService
	Projects          *ProjectsService
	Roles             *RolesService
	Storage           *StorageService
	Servers           *ServersService
	ServersMigrations *ServersMigrationsService
	StorageTypes      *StorageTypesService
	SSHKeys           *SSHKeysService
	Tasks             *TasksService
	Users             *UsersService
}

type service struct {
	client *Client
}

// Authenticator interface for client authentication.
type Authenticator interface {
	// Authenticate authenticate client on SOLUS IO and return credentials
	// which should be used for making further API calls.
	// The Client is fully initialized. Any endpoints which is not requires
	// authentication may be called.
	Authenticate(c *Client) (Credentials, error)
}

// EmailAndPasswordAuthenticator authenticate at SOLUS IO with specified email
// and password.
type EmailAndPasswordAuthenticator struct {
	Email    string
	Password string
}

var _ Authenticator = EmailAndPasswordAuthenticator{}

func (a EmailAndPasswordAuthenticator) Authenticate(c *Client) (Credentials, error) {
	authRequest := AuthLoginRequest{
		Email:    a.Email,
		Password: a.Password,
	}

	resp, err := c.authLogin(context.Background(), authRequest)
	if err != nil {
		return Credentials{}, err
	}

	return resp.Credentials, nil
}

// ApiTokenAuthenticator authenticate at SOLUS IO by provided API token.
type ApiTokenAuthenticator struct {
	Token string
}

var _ Authenticator = ApiTokenAuthenticator{}

func (a ApiTokenAuthenticator) Authenticate(*Client) (Credentials, error) {
	return Credentials{
		AccessToken: a.Token,
		TokenType:   "Bearer",
		ExpiresAt:   "",
	}, nil
}

// ClientOption represent client initialization options.
type ClientOption func(c *Client)

// AllowInsecure allow to skip certificate verify.
func AllowInsecure() ClientOption {
	return func(c *Client) {
		c.HttpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
}

// AllowInsecure allow to skip certificate verify.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.Logger = logger
	}
}

// NewClient create and initialize Client instance.
func NewClient(baseURL *url.URL, a Authenticator, opts ...ClientOption) (*Client, error) {
	client := &Client{
		BaseURL:   baseURL,
		UserAgent: "solus.io Go client",
		Headers: map[string][]string{
			"Accept":       {"application/json"},
			"Content-Type": {"application/json"},
		},
		HttpClient: &http.Client{
			Timeout:   time.Second * 35,
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		},
		Logger:  NullLogger{},
		Retries: 5,
	}

	for _, o := range opts {
		o(client)
	}

	c, err := a.Authenticate(client)
	if err != nil {
		return nil, err
	}

	client.Credentials = c
	client.Headers["Authorization"] = []string{client.Credentials.TokenType + " " + client.Credentials.AccessToken}

	client.s.client = client
	client.Account = (*AccountService)(&client.s)
	client.BackupNodes = (*BackupNodesService)(&client.s)
	client.Backups = (*BackupsService)(&client.s)
	client.ComputeResources = (*ComputeResourcesService)(&client.s)
	client.IPBlocks = (*IPBlocksService)(&client.s)
	client.License = (*LicenseService)(&client.s)
	client.Locations = (*LocationsService)(&client.s)
	client.OsImages = (*OsImagesService)(&client.s)
	client.Applications = (*ApplicationsService)(&client.s)
	client.Plans = (*PlansService)(&client.s)
	client.Projects = (*ProjectsService)(&client.s)
	client.Roles = (*RolesService)(&client.s)
	client.Storage = (*StorageService)(&client.s)
	client.Servers = (*ServersService)(&client.s)
	client.ServersMigrations = (*ServersMigrationsService)(&client.s)
	client.StorageTypes = (*StorageTypesService)(&client.s)
	client.SSHKeys = (*SSHKeysService)(&client.s)
	client.Tasks = (*TasksService)(&client.s)
	client.Users = (*UsersService)(&client.s)

	return client, nil
}

func (c *Client) authLogin(ctx context.Context, data AuthLoginRequest) (AuthLoginResponse, error) {
	const path = "auth/login"
	body, code, err := c.request(ctx, http.MethodPost, path, withBody(data))
	if err != nil {
		return AuthLoginResponse{}, err
	}

	if code != http.StatusOK {
		return AuthLoginResponse{}, newHTTPError(http.MethodPost, path, code, body)
	}

	var resp AuthLoginResponseData
	return resp.Data, unmarshal(body, &resp)
}
