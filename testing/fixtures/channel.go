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

package fixtures

import (
	"fmt"
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/helpers"
)

type Channel struct {
	Name string
}

func (c Channel) Setup(fixtures map[string]interface{}) (interface{}, error) {
	settings, ok := fixtures["settings"].(*hyper.Settings)

	if !ok {
		return nil, fmt.Errorf("settings missing")
	}

	broker, ok := fixtures["broker"].(hyper.MessageBroker)

	if !ok {
		return nil, fmt.Errorf("message broker missing")
	}

	directory, ok := fixtures["directory"].(hyper.Directory)

	if !ok {
		return nil, fmt.Errorf("directory missing")
	}

	channelSettings, definition, err := helpers.GetChannelSettingsAndDefinition(settings, c.Name)

	if err != nil {
		return nil, err
	}

	if channel, err := definition.Maker(channelSettings.Settings); err != nil {
		return nil, err
	} else if err := broker.AddChannel(channel); err != nil {
		return nil, err
	} else {
		if err := channel.SetDirectory(directory); err != nil {
			return nil, err
		}
		return channel, nil
	}
}

func (c Channel) Teardown(fixture interface{}) error {
	return nil
}
