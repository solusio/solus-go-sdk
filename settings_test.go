package solus

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSettings_Get(t *testing.T) {
	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/settings", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		writeResponse(t, w, http.StatusCreated, fakeSettingsRawJSON)
	})
	defer s.Close()

	actual, err := createTestClient(t, s.URL).Settings.Get(context.Background())
	require.NoError(t, err)

	fakeSettingsResponse := struct {
		Data Settings `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal(fakeSettingsRawJSON, &fakeSettingsResponse))

	require.Equal(t, fakeSettingsResponse.Data, actual)
}

func TestSettings_ChangeHostname(t *testing.T) {
	host := "new-hostname.tld"
	data := struct {
		Hostname string `json:"hostname"`
	}{
		Hostname: host,
	}

	s := startTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/settings", r.URL.Path)
		assert.Equal(t, http.MethodPatch, r.Method)
		assertRequestBody(t, r, data)

		writeResponse(t, w, http.StatusCreated, fakeSettingsRawJSON)
	})
	defer s.Close()

	request := SettingsUpdateRequest{
		Hostname: &host,
	}
	actual, err := createTestClient(t, s.URL).Settings.Patch(context.Background(), request)
	require.NoError(t, err)

	fakeSettingsResponse := struct {
		Data Settings `json:"data"`
	}{}
	require.NoError(t, json.Unmarshal(fakeSettingsRawJSON, &fakeSettingsResponse))

	require.Equal(t, fakeSettingsResponse.Data, actual)
}

var fakeSettingsRawJSON = []byte(`{
	"data": {
		"uuid": "00000000-0000-0000-0000-000000000000",
		"force_autoupdate": true,
		"limit_group": {
			"id": 1,
			"name": "LimitGroup",
			"created_at": "2019-08-27T13:04:39.000000Z",
			"updated_at": "2020-12-16T08:16:11.000000Z",
			"vms": {
				"limit": 0,
				"is_enabled": false
			},
			"running_vms": {
				"limit": 5,
				"is_enabled": true
			},
			"additional_ips": {
				"limit": 2,
				"is_enabled": true
			}
		},
		"update_channels": [
			"stable",
			"mainline",
			"testing"
		],
		"update_schedule": {
			"scheduled_days": [
				6,
				1,
				3,
				4,
				2,
				5,
				7
			],
			"scheduled_time": "06:00"
		},
		"token_api": {
			"url": "https://hostname.tld:443/api/v1",
			"is_enabled": true
		},
		"non_existent_vms_remover": {
			"enabled": false,
			"interval": 1
		},
		"network_rules": {
			"arp": false,
			"dhcp": true,
			"cloud_init": true,
			"smtp": true,
			"icmp": false,
			"icmp_reply": false,
			"portmapper": true
		},
		"mail": {
			"host": "hostname.tld",
			"port": "25",
			"username": "username",
			"password": "password",
			"encryption": false,
			"from_email": "user@hostname.tld",
			"from_name": "hostname.tld",
			"test_mail": "user@hostname.tld"
		},
		"hostname": "hostname.tld",
		"registration": {
			"role": "RoleName"
		},
		"send_statistic": true,
		"compute_resource": {
			"rescue_iso_url": "http://images.prod.solus.io/rescue/rescue-latest.iso",
			"balance_strategy": "round-robin"
		},
		"dns": {
			"type": "powerdns",
			"ttl": 60,
			"server_hostname_template": "{{random-prefix}}.domain.tld",
			"reverse_dns_domain_template": "{{ip}}.example.com",
			"register_fqdn_on_server_create": false,
			"drivers": {
				"powerdns": {
					"host": "hostname.tld",
					"api_key": "00000000-0000-0000-0000-000000000000",
					"port": "443",
					"https": true
				}
			}
		},
		"features": {
			"hide_plan_name": false,
			"hide_plan_section": false,
			"hide_user_data": false,
			"hide_location_section": false,
			"allow_registration": true,
			"allow_password_recovery": true
		},
		"billing_integration": {
			"type": "whmcs",
			"drivers": {
				"whmcs": {
					"url": "https://hostname.tld/modules/addons/solusio/api/",
					"token": "token"
				}
			}
		},
		"update": {
			"method": "auto",
			"channel": "testing"
		},
		"latest-version": "1.1.19040",
		"theme": {
			"brand_name": "SOLUS IO",
			"primary_color": "#000000",
			"secondary_color": "#FFFFFF",
			"logo": "https://hostname.tld:443/public/25d664da-48a4-42a8-9ac5-c371b9036fae.png",
			"favicon": null,
			"terms_and_conditions_url": "https://hostname.tld"
		},
		"notifications": {
			"server_create": {
				"enabled": true,
				"override_templates": {
					"en_US": false
				},
				"subject_templates": {
					"en_US": "Server {{ name }} Created"
				},
				"body_templates": {
					"en_US": "<p>The server <b>{{ name }}</b> has been created."
				}
			},
			"server_reset_password": {
				"enabled": true,
				"override_templates": [],
				"subject_templates": {
					"en_US": "Server Password Reset"
				},
				"body_templates": {
					"en_US": "<p>Your server <b>{{ name }}</b> password has been reset to: <b>{{ password }}</b></p>"
				}
			},
			"user_reset_password": {
				"enabled": true,
				"override_templates": [],
				"subject_templates": {
					"en_US": "Password Reset Request"
				},
				"body_templates": {
					"en_US": "<h1>Password Reset Request</h1>\n<p>You are receiving this email because we received a"
				}
			},
			"user_verify_email": {
				"enabled": true,
				"override_templates": [],
				"subject_templates": {
					"en_US": "Email Address Verification"
				},
				"body_templates": {
					"en_US": "<h1>Email Verification</h1>\n<p>Thank you for signing up. Please verify your email"
				}
			},
			"project_user_invite": {
				"enabled": true,
				"override_templates": [],
				"subject_templates": {
					"en_US": "Project Invitation"
				},
				"body_templates": {
					"en_US": "<h1>Project Invitation</h1>\n<p>You have been invited to collaborate on the"
				}
			},
			"project_user_left": {
				"enabled": true,
				"override_templates": [],
				"subject_templates": {
					"en_US": "User Left Project"
				},
				"body_templates": {
					"en_US": "<p>The user <b>{{ email }}</b> has left the <b>{{ project }}</b> project.</p>"
				}
			},
			"server_incoming_traffic_exceeded": {
				"enabled": true,
				"override_templates": [],
				"subject_templates": {
					"en_US": "Incoming Traffic Limit Exceeded"
				},
				"body_templates": {
					"en_US": "<p>The incoming traffic limit of the server <b>{{ name }}</b> has been exceeded.</p>"
				}
			},
			"server_outgoing_traffic_exceeded": {
				"enabled": true,
				"override_templates": [],
				"subject_templates": {
					"en_US": "Outgoing Traffic Limit Exceeded"
				},
				"body_templates": {
					"en_US": "<p>The outgoing traffic limit of the server <b>{{ name }}</b> has been exceeded.</p>"
				}
			}
		},
		"latest_version": "1.1.20873"
	}
}`)
