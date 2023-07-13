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
)

var Channels = hyper.ChannelDefinitions{
	"stdout": hyper.ChannelDefinition{
		Name:              "Stdout Channel",
		Description:       "Prints messages to stdout (just for testing and debugging)",
		Maker:             MakeStdoutChannel,
		SettingsValidator: StdoutSettingsValidator,
	},
	"jsonrpc_client": hyper.ChannelDefinition{
		Name:              "JSONRPC Client Channel",
		Description:       "Creates outgoing JSONRPC connections to deliver and receive messages",
		Maker:             MakeJSONRPCClientChannel,
		SettingsValidator: JSONRPCClientSettingsValidator,
	},
	"grpc_client": hyper.ChannelDefinition{
		Name:              "gRPC Client Channel",
		Description:       "Creates outgoing gRPC connections to deliver and receive messages",
		Maker:             MakeGRPCClientChannel,
		SettingsValidator: GRPCClientSettingsValidator,
	},
	"jsonrpc_server": hyper.ChannelDefinition{
		Name:              "JSONRPC Server Channel",
		Description:       "Accepts incoming JSONRPC connections to deliver and receive messages",
		Maker:             MakeJSONRPCServerChannel,
		SettingsValidator: JSONRPCServerSettingsValidator,
	},
	"grpc_server": hyper.ChannelDefinition{
		Name:              "gRPC Server Channel",
		Description:       "Accepts incoming gRPC connections to deliver and receive messages",
		Maker:             MakeGRPCServerChannel,
		SettingsValidator: GRPCServerSettingsValidator,
	},
}
