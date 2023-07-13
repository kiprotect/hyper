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

package helpers

import (
	"encoding/json"
)

// convert an arbitrary structure to a string map via the JSON encoding
func ToStringMap(value interface{}) (map[string]interface{}, error) {
	var out map[string]interface{}
	if jsonData, err := json.Marshal(value); err != nil {
		return nil, err
	} else if err := json.Unmarshal(jsonData, &out); err != nil {
		return nil, err
	} else {
		return out, nil
	}
}
