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
	"fmt"
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/forms"
	"github.com/kiprotect/hyper/jsonrpc"
)

type JSONRPCServerChannel struct {
	hyper.BaseChannel
	Settings *jsonrpc.JSONRPCServerSettings
	Server   *jsonrpc.JSONRPCServer
}

func JSONRPCServerSettingsValidator(settings map[string]interface{}) (interface{}, error) {
	if params, err := jsonrpc.JSONRPCServerSettingsForm.Validate(settings); err != nil {
		return nil, err
	} else {
		validatedSettings := &jsonrpc.JSONRPCServerSettings{}
		if err := jsonrpc.JSONRPCServerSettingsForm.Coerce(validatedSettings, params); err != nil {
			return nil, err
		}
		return validatedSettings, nil
	}
}

func MakeJSONRPCServerChannel(settings interface{}) (hyper.Channel, error) {
	rpcSettings := settings.(jsonrpc.JSONRPCServerSettings)

	s := &JSONRPCServerChannel{
		Settings: &rpcSettings,
	}

	if server, err := jsonrpc.MakeJSONRPCServer(&rpcSettings, s.handler); err != nil {
		return nil, fmt.Errorf("error creating JSON-RPC server: %w", err)
	} else {
		s.Server = server
		return s, nil
	}
}

func (c *JSONRPCServerChannel) handler(context *jsonrpc.Context) *jsonrpc.Response {

	request := &hyper.Request{}

	// we make sure the parameters are well-formed
	if params, err := forms.RequestForm.Validate(map[string]interface{}{
		"method": context.Request.Method,
		"params": context.Request.Params,
		"id":     context.Request.ID,
	}); err != nil {
		hyper.Log.Debug(err)
		return context.Error(400, fmt.Sprintf("invalid request: %v", err), err)
	} else if err := forms.RequestForm.Coerce(request, params); err != nil {
		hyper.Log.Error(err)
		return context.InternalError()
	}

	// we replace the ID with an addressable ID that we can use to reconstruct
	// the sender of the request later
	request.ID = fmt.Sprintf("%s(%s)", request.Method, request.ID)

	// this request comes from the server itself
	clientInfo := &hyper.ClientInfo{
		Name: c.Directory().Name(),
	}

	if entry, err := c.Directory().OwnEntry(); err != nil {
		hyper.Log.Errorf("Error retrieving own directory entry: %v", err)
		return context.InternalError()
	} else {
		clientInfo.Entry = entry
	}

	if response, err := c.MessageBroker().DeliverRequest(request, clientInfo); err != nil {
		return context.Error(1, err.Error(), err)
	} else {
		if response == nil {
			return context.Result(map[string]interface{}{"message": "submitted"})
		}
		jsonrpcResponse := jsonrpc.FromHyperResponse(response)
		if jsonrpcResponse.Error != nil {
			return context.Error(jsonrpcResponse.Error.Code, jsonrpcResponse.Error.Message, jsonrpcResponse.Error.Data)
		} else {
			return context.Result(jsonrpcResponse.Result)
		}
	}
}

func (c *JSONRPCServerChannel) Type() string {
	return "jsonrpc_server"
}

func (c *JSONRPCServerChannel) Open() error {
	return c.Server.Start()
}

func (c *JSONRPCServerChannel) Close() error {
	return c.Server.Stop()
}

func (c *JSONRPCServerChannel) DeliverRequest(request *hyper.Request) (*hyper.Response, error) {
	return nil, nil
}

func (c *JSONRPCServerChannel) CanDeliverTo(address *hyper.Address) bool {
	return false
}
