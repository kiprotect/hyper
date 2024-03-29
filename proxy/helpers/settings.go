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
	"github.com/kiprotect/go-helpers/settings"
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/proxy"
	"io/fs"
)

func Settings(settingsPaths []string, fs fs.FS, definitions *hyper.Definitions) (*proxy.Settings, error) {
	if rawSettings, err := settings.MakeSettings(settingsPaths, fs); err != nil {
		return nil, err
	} else if params, err := proxy.SettingsForm.ValidateWithContext(rawSettings.Values, map[string]interface{}{"definitions": definitions}); err != nil {
		return nil, err
	} else {
		settings := &proxy.Settings{
			Definitions: definitions,
		}
		if err := proxy.SettingsForm.Coerce(settings, params); err != nil {
			// this should not happen if the forms are correct...
			return nil, err
		}
		// settings are valid
		return settings, nil
	}
}
