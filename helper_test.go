package solus

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var fakeApplication = Application{
	Id:   1,
	Name: "fake application",
	Icon: Icon{
		Id:   2,
		Name: "Fake application icon",
		URL:  "http://example.com/image.png",
		Type: IconTypeApplication,
	},
	Url:              "http://example.com/app",
	CloudInitVersion: "v2",
	UserData:         "fake user data",
	LoginLink: LoginLink{
		Type: LoginLinkTypeNone,
	},
	JsonSchema: "fake json schema",
	IsDefault:  true,
	IsVisible:  true,
	IsBuiltin:  true,
}

var fakeComputeResource = ComputeResource{
	Id:        1,
	Name:      "fake compute resource",
	Host:      "192.0.2.1",
	AgentPort: 1337,
	Status: ComputerResourceStatus{
		Id:   2,
		Name: "fake status",
	},
	Locations: []Location{
		fakeLocation,
	},
}

var fakeComputeResourceInstallStep = ComputeResourceInstallStep{
	Id:                1,
	ComputeResourceId: 2,
	Title:             "fake CR install step",
	Status:            ComputeResourceInstallStepStatusError,
	StatusText:        "fake status text",
	Progress:          57,
}

var fakeIpBlock = IpBlock{
	Id:      1,
	Name:    "fake ip block",
	Type:    IPv4,
	Gateway: "192.0.2.254",
	Netmask: "255.255.0.0",
	Ns1:     "8.8.8.8",
	Ns2:     "8.8.4.4",
	From:    "192.0.2.1",
	To:      "192.0.2.100",
	Subnet:  16,
	ComputeResources: []ComputeResource{
		fakeComputeResource,
	},
	Ips: []IpBlockIpAddress{
		{
			Id: 3,
			Ip: "192.0.2.1",
		},
	},
}

var fakeIpBlockIpAddress = IpBlockIpAddress{
	Id:      1,
	Ip:      "192.0.2.2",
	IpBlock: fakeIpBlock,
}

var fakeLicense = License{
	CpuCores:       1,
	CpuCoresInUse:  2,
	IsActive:       true,
	Key:            "fake key",
	KeyType:        "fake key type",
	Product:        "fake product",
	ExpirationDate: "fake expiration date",
	UpdateDate:     "fake update date",
}

var fakeLocation = Location{
	Id:   1,
	Name: "fake location",
	Icon: Icon{
		Id:   2,
		Name: "Fake location icon",
		URL:  "http://example.com/image.png",
		Type: IconTypeFlags,
	},
	Description: "fake description",
	IsDefault:   true,
	IsVisible:   true,
	ComputeResources: []ComputeResource{
		{Id: 1},
	},
}

var fakeOsImage = OsImage{
	Id:   1,
	Name: "fake os image",
	Icon: Icon{
		Id:   2,
		Name: "fake os image icon",
		URL:  "http://example.com/image.png",
		Type: IconTypeOS,
	},
	Versions: []OsImageVersion{
		fakeOsImageVersion,
	},
	IsDefault: false,
}

var fakeOsImageVersion = OsImageVersion{
	Id:               1,
	Position:         100,
	Version:          "1337",
	Url:              "http://example.com/os.qcow2",
	CloudInitVersion: "v2",
}

var fakePlan = Plan{
	Id:   1,
	Name: "fake plan",
	Params: PlanParams{
		Disk: 42,
		RAM:  1337,
		VCPU: 100500,
	},
	StorageType:         "fb",
	ImageFormat:         "qcow2",
	IsDefault:           true,
	IsSnapshotAvailable: true,
	IsSnapshotsEnabled:  true,
	Limits: PlanLimits{
		TotalBytes: PlanLimit{
			IsEnabled: true,
			Limit:     2,
		},
		TotalIops: PlanLimit{
			IsEnabled: true,
			Limit:     5,
		},
	},
	TokenPerHour:  7,
	TokenPerMonth: 8,
	Position:      11,
}

var fakeProject = Project{
	Id:          1,
	Name:        "fake project",
	Description: "fake descriptions",
	Members:     42,
	IsOwner:     true,
	IsDefault:   true,
	Owner:       fakeUser,
	Servers: []Server{
		fakeServer,
	},
}

var fakeUser = User{
	Id:        1,
	Email:     "fake@example.com",
	Password:  "fake password",
	CreatedAt: time.Now().String(),
	Status:    UserStatusActive,
	Roles: []Role{
		fakeRole,
	},
}

var fakeRole = Role{
	Id:         1,
	Name:       "fake role",
	IsDefault:  true,
	UsersCount: 42,
}

var fakeServer = Server{
	Id:          1,
	Name:        "fake server",
	Description: "fake description",
	UUID:        "123e4567-e89b-12d3-a456-426655440000",
	Status:      "running",
	Ips: []IpBlockIpAddress{
		fakeIpBlockIpAddress,
	},
}

var fakeTask = Task{
	Id:                1,
	ComputeResourceId: 2,
	Queue:             "fake queue",
	Action:            ServerActionCreate,
	Status:            TaskStatusDone,
	Output:            "fake output",
	Progress:          42,
	Duration:          23,
}

func startTestServer(t *testing.T, h http.HandlerFunc) *httptest.Server {
	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	s := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Scheme = "http"
		r.URL.Host = listener.Addr().String()
		h(w, r)
	}))

	err = s.Listener.Close()
	require.NoError(t, err)
	s.Listener = listener

	s.Start()

	return s
}

type authenticator struct{}

func (authenticator) Authenticate(*Client) (Credentials, error) { return Credentials{}, nil }

func createTestClient(t *testing.T, addr string) *Client {
	u, err := url.Parse(addr)
	require.NoError(t, err)

	c, err := NewClient(u, authenticator{})
	require.NoError(t, err)
	return c
}

func assertRequestQuery(t *testing.T, r *http.Request, expected url.Values) {
	assert.Equal(t, expected.Encode(), r.URL.Query().Encode())
}

func assertRequestBody(t *testing.T, r *http.Request, expected interface{}) {
	b, err := ioutil.ReadAll(r.Body)
	require.NoError(t, err)

	d := reflect.New(reflect.TypeOf(expected)).Interface()
	err = json.Unmarshal(b, d)
	require.NoError(t, err)

	assert.Equal(t, expected, reflect.ValueOf(d).Elem().Interface())
}

func writeJSON(t *testing.T, w http.ResponseWriter, statusCode int, r interface{}) {
	data, err := json.Marshal(r)
	require.NoError(t, err)

	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	require.NoError(t, err)
}

func writeResponse(t *testing.T, w http.ResponseWriter, statusCode int, r interface{}) {
	writeJSON(t, w, statusCode, struct {
		Data interface{} `json:"data"`
	}{r})
}
