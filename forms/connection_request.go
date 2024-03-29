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

var ConnectionRequestForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "channel",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "endpoint",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
		{
			Name: "token",
			Validators: []forms.Validator{
				forms.IsBytes{Encoding: "base64"},
			},
		},
		{
			Name: "client",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &ClientInfoForm,
				},
			},
		},
	},
}
