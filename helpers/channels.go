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
	"fmt"
	"github.com/kiprotect/hyper"
)

func GetChannelSettingsAndDefinition(settings *hyper.Settings, name string) (*hyper.ChannelSettings, *hyper.ChannelDefinition, error) {
	for _, channel := range settings.Channels {
		if channel.Name == name {
			def := settings.Definitions.ChannelDefinitions[channel.Type]
			return channel, &def, nil
		}
	}
	return nil, nil, fmt.Errorf("channel not found")
}

func InitializeChannels(broker hyper.MessageBroker, directory hyper.Directory, settings *hyper.Settings) ([]hyper.Channel, error) {
	channels := make([]hyper.Channel, 0)
	for _, channel := range settings.Channels {
		hyper.Log.Debugf("Initializing channel '%s' of type '%s'", channel.Name, channel.Type)
		definition := settings.Definitions.ChannelDefinitions[channel.Type]
		if channelObj, err := definition.Maker(channel.Settings); err != nil {
			return nil, fmt.Errorf("error initializing channel '%s': %w", channel.Name, err)
		} else {
			if err := broker.AddChannel(channelObj); err != nil {
				return nil, fmt.Errorf("error adding channel '%s': %w", channel.Name, err)
			}
			if err := channelObj.SetDirectory(directory); err != nil {
				return nil, fmt.Errorf("error setting directory for channel '%s': %w", channel.Name, err)
			}
			channels = append(channels, channelObj)
		}
	}
	return channels, nil
}

func OpenChannels(broker hyper.MessageBroker, directory hyper.Directory, settings *hyper.Settings) ([]hyper.Channel, error) {

	channels, err := InitializeChannels(broker, directory, settings)

	if err != nil {
		return nil, fmt.Errorf("error initializing channels: %w", err)
	} else {
		for i, channel := range channels {
			if err := channel.Open(); err != nil {
				return nil, fmt.Errorf("error opening channel %d: %w", i, err)
			}
		}
	}
	return channels, nil
}

func CloseChannels(channels []hyper.Channel) error {
	var lastErr error
	for i, channel := range channels {
		if err := channel.Close(); err != nil {
			lastErr = fmt.Errorf("error closing channel %d: %w", i, err)
			hyper.Log.Error(lastErr)
		}
	}
	return lastErr
}
