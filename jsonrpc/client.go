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

package jsonrpc

import (
	"bytes"
	"encoding/json"
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/tls"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
	settings *JSONRPCClientSettings
}

func MakeClient(settings *JSONRPCClientSettings) *Client {
	return &Client{
		settings: settings,
	}
}

func (c *Client) SetServerName(serverName string) {
	c.settings.TLS.ServerName = serverName
}

func (c *Client) SetEndpoint(endpoint string) {
	c.settings.Endpoint = endpoint
}

func (c *Client) Call(request *Request) (*Response, error) {
	data, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	transport := &http.Transport{
		DisableKeepAlives: true, // removing this will cause connections to pile up
	}

	if c.settings.ProxyUrl != "" {
		hyper.Log.Debugf("Using proxy URL %s", c.settings.ProxyUrl)
		if proxyUrl, err := url.Parse(c.settings.ProxyUrl); err != nil {
			return nil, err
		} else {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	client.Transport = transport

	if c.settings.TLS != nil {
		tlsConfig, err := tls.TLSClientConfig(c.settings.TLS)
		if err != nil {
			return nil, err
		}
		transport.TLSClientConfig = tlsConfig
	}

	hyper.Log.Debugf("Generating request to endpoint %s...", c.settings.Endpoint)

	req, err := http.NewRequest("POST", c.settings.Endpoint, bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// to do: sanity checks...

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}

	response := &Response{}

	if err := json.Unmarshal(body, response); err != nil {
		return nil, err
	}

	return response, nil
}
