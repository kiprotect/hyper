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

type DatastoreDefinition struct {
	Name              string            `json:"name"`
	Description       string            `json:"description"`
	Maker             DatastoreMaker    `json:"-"`
	SettingsValidator SettingsValidator `json:"-"`
}

type DatastoreDefinitions map[string]DatastoreDefinition
type DatastoreMaker func(settings interface{}) (Datastore, error)

type Datastore interface {
	// Write data to the store
	Write(*DataEntry) error
	// Read data from the store
	Read() ([]*DataEntry, error)
	Init() error
}

const (
	NullType = 0
)

type DataEntry struct {
	Type uint8
	ID   []byte
	Data []byte
}
