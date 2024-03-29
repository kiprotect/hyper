// KIProtect Hyper
// Copyright (C) 2021-2023 KIProtect GmbH
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package proxy

import (
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	hyperForms "github.com/kiprotect/hyper/forms"
	"github.com/kiprotect/hyper/jsonrpc"
	"github.com/kiprotect/hyper/net"
	"github.com/kiprotect/hyper/tls"
	"regexp"
	"time"
)

var DirectorySettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "allowed_domains",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsString{},
						forms.MatchesRegex{Regexp: regexp.MustCompile(`^\.`)},
					},
				},
			},
		},
	},
}

var SettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "metrics",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{
					Form: &hyperForms.MetricsSettingsForm,
				},
			},
		},
		{
			Name: "public",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{
					Form: &PublicSettingsForm,
				},
			},
		},
		{
			Name: "private",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{
					Form: &PrivateSettingsForm,
				},
			},
		},
	},
}

type IsValidExpiresAtTime struct{}

func (i IsValidExpiresAtTime) Validate(value interface{}, values map[string]interface{}) (interface{}, error) {
	timeValue, ok := value.(time.Time)
	if !ok {
		return nil, fmt.Errorf("expected a time")
	}
	// we subtract 7*24 hours from the time value and make sure it's before the current time
	if timeValue.Add(-7 * 24 * time.Hour).After(time.Now()) {
		return nil, fmt.Errorf("timed announcements need to expire in 7 days or less")
	}
	return timeValue, nil
}

var InternalEndpointForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "address",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "localhost:8888"},
				forms.IsString{},
			},
		},
		{
			Name: "tls",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{
					Form: &tls.TLSSettingsForm,
				},
			},
		},
		{
			Name: "jsonrpc_client",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{
					Form: &jsonrpc.JSONRPCClientSettingsForm,
				},
			},
		},
		{
			Name: "jsonrpc_path",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "/jsonrpc"},
				forms.IsString{},
			},
		},
		{
			Name: "timeout",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 30.0},
				forms.IsFloat{HasMin: true, Min: 0, HasMax: true, Max: 3000},
			},
		},
	},
}

var PrivateSettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "datastore",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &hyperForms.DatastoreForm,
				},
			},
		},
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "internal_endpoint",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &InternalEndpointForm,
				},
			},
		},
		{
			Name: "jsonrpc_client",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &jsonrpc.JSONRPCClientSettingsForm,
				},
			},
		},
		{
			Name: "jsonrpc_server",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &jsonrpc.JSONRPCServerSettingsForm,
				},
			},
		},
	},
}

var PublicSettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "datastore",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &hyperForms.DatastoreForm,
				},
			},
		},
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		net.TCPRateLimitsField,
		{
			Name: "tls_bind_address",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "localhost:443"},
				forms.IsString{},
			},
		},
		{
			Name: "internal_bind_address",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "localhost:9999"},
				forms.IsString{},
			},
		},
		{
			Name: "internal_endpoint",
			Validators: []forms.Validator{
				forms.IsOptional{Default: "localhost:9999"},
				forms.IsString{},
			},
		},
		{
			Name: "jsonrpc_client",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &jsonrpc.JSONRPCClientSettingsForm,
				},
			},
		},
		{
			Name: "jsonrpc_server",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &jsonrpc.JSONRPCServerSettingsForm,
				},
			},
		},
		{
			Name: "accept_timeout",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 10.0},
				forms.IsFloat{HasMin: true, Min: 0, HasMax: true, Max: 3000},
			},
		},
	},
}
