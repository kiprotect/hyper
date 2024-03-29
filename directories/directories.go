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

package directories

import (
	"github.com/kiprotect/hyper"
)

var Directories = hyper.DirectoryDefinitions{
	"api": hyper.DirectoryDefinition{
		Name:              "API-Based Service Directory",
		Description:       "Loads the service directory via a JSON-RPC API",
		Maker:             MakeAPIDirectory,
		SettingsValidator: APIDirectorySettingsValidator,
	},
	"json": hyper.DirectoryDefinition{
		Name:              "JSON-Based Service Directory",
		Description:       "Loads the service directory from a JSON file",
		Maker:             MakeJSONDirectory,
		SettingsValidator: JSONDirectorySettingsValidator,
	},
}
