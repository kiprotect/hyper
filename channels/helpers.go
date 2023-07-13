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

package channels

import (
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/forms"
)

func parseConnectionRequest(request *hyper.Request) (*hyper.ConnectionRequest, error) {
	var connectionRequest hyper.ConnectionRequest
	if params, err := forms.ConnectionRequestForm.Validate(request.Params); err != nil {
		return nil, err
	} else if err := forms.ConnectionRequestForm.Coerce(&connectionRequest, params); err != nil {
		return nil, err
	} else {
		return &connectionRequest, nil
	}
}

func parseRequestConnectionResponse(response map[string]interface{}) (*RequestConnectionResponse, error) {
	var requestConnectionResponse RequestConnectionResponse
	if params, err := RequestConnectionResponseForm.Validate(response); err != nil {
		return nil, err
	} else if err := RequestConnectionResponseForm.Coerce(&requestConnectionResponse, params); err != nil {
		return nil, err
	} else {
		return &requestConnectionResponse, nil
	}
}
