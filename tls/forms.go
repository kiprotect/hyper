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

package tls

import (
	"github.com/kiprotect/go-helpers/forms"
)

var TLSSettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "insecure_skip_verify",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
		{
			Name: "verify_client",
			Validators: []forms.Validator{
				forms.IsOptional{Default: true},
				forms.IsBoolean{},
			},
		},
		{
			Name: "request_client_cert",
			Validators: []forms.Validator{
				forms.IsOptional{Default: false},
				forms.IsBoolean{},
			},
		},
		{
			Name: "ca_certificate_files",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsString{},
					},
				},
			},
		},
		{
			Name: "certificate_file",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
			},
		},
		{
			Name: "server_name",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{}, // to do: add URL validation
			},
		},
		{
			Name: "key_file",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsString{},
			},
		},
	},
}
