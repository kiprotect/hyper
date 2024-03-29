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

package tls

type TLSSettings struct {
	ServerName         string   `json:"server_name"`
	VerifyClient       bool     `json:"verify_client"`
	RequestClientCert  bool     `json:"request_client_cert"`
	CACertificateFiles []string `json:"ca_certificate_files"`
	CertificateFile    string   `json:"certificate_file"`
	KeyFile            string   `json:"key_file"`

	// This switch only exists to accomodate the inability of certain
	// certificate authorities to provide TLS certificates with
	// the necessary rights. Since key pinning is used to verify certificates
	// in addition to the normal TLS verification enabling this will not
	// destroy the systems' security, although it will weaken it.
	// So please do not set this to true...
	InsecureSkipVerify bool `json:"insecure_skip_verify"`
}
