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

package proxy

import (
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/jsonrpc"
	"github.com/kiprotect/hyper/net"
	"github.com/kiprotect/hyper/tls"
	"time"
)

const (
	PublicAnnouncementType  uint8 = 1
	PrivateAnnouncementType uint8 = 2
)

type Settings struct {
	Definitions *hyper.Definitions     `json:"definitions"`
	Metrics     *hyper.MetricsSettings `json:"metrics"`
	Private     *PrivateServerSettings `json:"private"`
	Public      *PublicServerSettings  `json:"public"`
}

type DirectorySettings struct {
	AllowedDomains []string `json:"allowed_domains"`
}

type PublicServerSettings struct {
	Datastore           *hyper.DatastoreSettings       `json:"datastore"`
	Name                string                         `json:"name"`
	TLSBindAddress      string                         `json:"tls_bind_address"`
	InternalBindAddress string                         `json:"internal_bind_address"`
	InternalEndpoint    string                         `json:"internal_endpoint"`
	JSONRPCClient       *jsonrpc.JSONRPCClientSettings `json:"jsonrpc_client"`
	JSONRPCServer       *jsonrpc.JSONRPCServerSettings `json:"jsonrpc_server`
	TCPRateLimits       []*net.RateLimit               `json:"tcp_rate_limits"`
	AcceptTimeout       float64                        `json:"accept_timeout"`
}

type PublicAnnouncement struct {
	// When the announcement expires
	ExpiresAt *time.Time `json:"expires_at"`
	// the name of the operator to forward the connection to
	Operator string `json:"operator"`
	// the name of the domain to forward
	Domain string `json:"domain"`
}

type PrivateAnnouncement struct {
	// When the announcement expires
	ExpiresAt *time.Time `json:"expires_at"`
	// the name of the public proxy to announce this to
	Proxy string `json:"proxy"`
	// the pattern to announce, as a regexp
	Domain string `json:"domain"`
}

type InternalEndpointSettings struct {
	Address       string                         `json:"address"`
	TLS           *tls.TLSSettings               `json:"tls"`
	JSONRPCClient *jsonrpc.JSONRPCClientSettings `json:"jsonrpc_client"`
	JSONRPCPath   string                         `json:"jsonrpc_path"`
	Timeout       float64                        `json:"timeout"`
}

type PrivateServerSettings struct {
	Datastore        *hyper.DatastoreSettings       `json:"datastore"`
	Name             string                         `json:"name"`
	Announcements    []*PrivateAnnouncement         `json:"announcements"`
	InternalEndpoint *InternalEndpointSettings      `json:"internal_endpoint"`
	JSONRPCClient    *jsonrpc.JSONRPCClientSettings `json:"jsonrpc_client"`
	JSONRPCServer    *jsonrpc.JSONRPCServerSettings `json:"jsonrpc_server`
}
