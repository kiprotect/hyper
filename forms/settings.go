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

package forms

import (
	"github.com/kiprotect/go-helpers/forms"
)

var DirectorySettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsString{},
				IsValidDirectoryType{},
			},
		},
		{
			Name: "settings",
			Validators: []forms.Validator{
				forms.IsStringMap{},
				AreValidDirectorySettings{},
			},
		},
	},
}

var ChannelForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "type",
			Validators: []forms.Validator{
				forms.IsString{},
				IsValidChannelType{},
			},
		},
		{
			Name: "settings",
			Validators: []forms.Validator{
				forms.IsStringMap{},
				AreValidChannelSettings{},
			},
		},
	},
}

var MetricsSettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "bind_address",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

var SettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "name",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "directory",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &DirectorySettingsForm,
				},
			},
		},
		{
			Name: "metrics",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{
					Form: &MetricsSettingsForm,
				},
			},
		},
		{
			Name: "signing",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringMap{
					Form: &SigningSettingsForm,
				},
			},
		},
		{
			Name: "channels",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &ChannelForm,
						},
					},
				},
			},
		},
	},
}
