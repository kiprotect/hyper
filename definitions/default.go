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

package definitions

import (
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/channels"
	"github.com/kiprotect/hyper/cmd"
	"github.com/kiprotect/hyper/datastores"
	"github.com/kiprotect/hyper/directories"
)

var Default = hyper.Definitions{
	DatastoreDefinitions: datastores.Definitions,
	DirectoryDefinitions: directories.Directories,
	CommandsDefinitions:  cmd.Commands,
	ChannelDefinitions:   channels.Channels,
}
