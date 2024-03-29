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

package hyper

import (
	"fmt"
)

type ChannelDefinition struct {
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	Maker             ChannelMaker      `json:"-"`
	SettingsValidator SettingsValidator `json:"-"`
}

type ChannelDefinitions map[string]ChannelDefinition
type ChannelMaker func(settings interface{}) (Channel, error)

// A channel can deliver and accept message
type Channel interface {
	Type() string
	MessageBroker() MessageBroker
	SetMessageBroker(MessageBroker) error
	CanDeliverTo(*Address) bool
	DeliverRequest(*Request) (*Response, error)
	SetDirectory(Directory) error
	Directory() Directory
	Close() error
	Open() error
}

type ProxyChannel interface {
	HandleConnectionRequest(address *Address, request *Request) (*Response, error)
}

type BaseChannel struct {
	broker    MessageBroker
	directory Directory
}

func (b *BaseChannel) OperatorEntry(name string) (*DirectoryEntry, error) {
	if entries, err := b.Directory().Entries(&DirectoryQuery{
		Operator: name,
	}); err != nil {
		return nil, fmt.Errorf("error retrieving operator entry: %w", err)
	} else if len(entries) == 0 {
		return nil, fmt.Errorf("no entry found")
	} else {
		return entries[0], nil
	}

}

func (b *BaseChannel) DirectoryEntry(operator string, channel string) (*DirectoryEntry, error) {
	if entries, err := b.Directory().Entries(&DirectoryQuery{
		Operator: operator,
		Channels: []string{channel},
	}); err != nil {
		return nil, fmt.Errorf("error retrieving directory entry: %w", err)
	} else if len(entries) > 0 {
		return entries[0], nil
	}
	return nil, nil

}

func (b *BaseChannel) Directory() Directory {
	return b.directory
}

func (b *BaseChannel) SetDirectory(directory Directory) error {
	b.directory = directory
	return nil
}

func (b *BaseChannel) MessageBroker() MessageBroker {
	return b.broker
}

func (b *BaseChannel) SetMessageBroker(broker MessageBroker) error {
	b.broker = broker
	return nil
}
