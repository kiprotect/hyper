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
	"github.com/kiprotect/hyper/jsonrpc"
)

type JSONRPCClientChannel struct {
	hyper.BaseChannel
	Settings *jsonrpc.JSONRPCClientSettings
}

func JSONRPCClientSettingsValidator(settings map[string]interface{}) (interface{}, error) {
	if params, err := jsonrpc.JSONRPCClientSettingsForm.Validate(settings); err != nil {
		return nil, err
	} else {
		validatedSettings := &jsonrpc.JSONRPCClientSettings{}
		if err := jsonrpc.JSONRPCClientSettingsForm.Coerce(validatedSettings, params); err != nil {
			return nil, err
		}
		return validatedSettings, nil
	}
}

func MakeJSONRPCClientChannel(settings interface{}) (hyper.Channel, error) {
	rpcSettings := settings.(jsonrpc.JSONRPCClientSettings)
	return &JSONRPCClientChannel{
		Settings: &rpcSettings,
	}, nil
}

func (c *JSONRPCClientChannel) Type() string {
	return "jsonrpc_client"
}

func (c *JSONRPCClientChannel) Open() error {
	return nil
}

func (c *JSONRPCClientChannel) Close() error {
	return nil
}

func (c *JSONRPCClientChannel) DeliverRequest(request *hyper.Request) (*hyper.Response, error) {

	hyper.Log.Info("Delivering request via JSON-RPC...")

	client := jsonrpc.MakeClient(c.Settings)
	jsonrpcRequest := &jsonrpc.Request{}
	jsonrpcRequest.FromHyperRequest(request)

	if groups := hyper.MethodNameRegexp.FindStringSubmatch(jsonrpcRequest.Method); groups == nil {
		return nil, fmt.Errorf("invalid method name")
	} else {
		// we remove the operator name from the method call before passing it in
		jsonrpcRequest.Method = groups[2]
	}

	jsonrpcResponse, err := client.Call(jsonrpcRequest)
	if err != nil {
		hyper.Log.Error(err)
		return nil, fmt.Errorf("error calling JSON-RPC server: %w", err)
	}

	hyper.Log.Info("Delivered request...")

	return jsonrpcResponse.ToHyperResponse(), nil
}

func (c *JSONRPCClientChannel) CanDeliverTo(address *hyper.Address) bool {

	if address.Operator == c.Directory().Name() {
		return true
	}

	return false
}
