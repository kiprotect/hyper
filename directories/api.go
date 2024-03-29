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
	"crypto/x509"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/hyper"
	hyperForms "github.com/kiprotect/hyper/forms"
	"github.com/kiprotect/hyper/helpers"
	"github.com/kiprotect/hyper/jsonrpc"
	"sync"
	"time"
)

var APIDirectorySettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "endpoints",
			Validators: []forms.Validator{
				forms.IsStringList{},
			},
		},
		{
			Name: "server_names",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsStringList{},
			},
		},
		{
			Name: "cache_entries_for",
			Validators: []forms.Validator{
				forms.IsOptional{Default: 5},
				forms.IsInteger{
					HasMin: true,
					Min:    0,
					HasMax: true,
					Max:    3600,
				},
			},
		},
		{
			Name: "jsonrpc_client",
			Validators: []forms.Validator{
				forms.IsStringMap{
					Form: &jsonrpc.JSONRPCClientSettingsForm,
				},
			},
		},
		{
			Name: "ca_certificate_files",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsString{},
					},
				},
			},
		},
		{
			Name: "ca_intermediate_certificate_files",
			Validators: []forms.Validator{
				forms.IsOptional{},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsString{},
					},
				},
			},
		},
	},
}

type APIDirectorySettings struct {
	Endpoints                      []string                       `json:"endpoints"`
	ServerNames                    []string                       `json:"server_names"`
	JSONRPCClient                  *jsonrpc.JSONRPCClientSettings `json:"jsonrpc_client"`
	CACertificateFiles             []string                       `json:"ca_certificate_files"`
	CAIntermediateCertificateFiles []string                       `json:"ca_intermediate_certificate_files"`
	CacheEntriesFor                int64                          `json:"cache_entries_for"`
}

type CacheEntry struct {
	Entry     *hyper.DirectoryEntry
	FetchedAt time.Time
}

type DirectoryCache struct {
	Entries []*CacheEntry
}

type APIDirectory struct {
	hyper.BaseDirectory
	lastUpdate        time.Time
	settings          APIDirectorySettings
	jsonrpcClient     *jsonrpc.Client
	rootCerts         []*x509.Certificate
	intermediateCerts []*x509.Certificate
	entries           map[string]*hyper.DirectoryEntry
	records           []*hyper.SignedChangeRecord
	mutex             sync.Mutex
}

func APIDirectorySettingsValidator(settings map[string]interface{}) (interface{}, error) {
	if params, err := APIDirectorySettingsForm.Validate(settings); err != nil {
		return nil, err
	} else {
		validatedSettings := &APIDirectorySettings{}
		if err := APIDirectorySettingsForm.Coerce(validatedSettings, params); err != nil {
			return nil, err
		}
		return validatedSettings, nil
	}
}

func MakeAPIDirectory(name string, settings interface{}) (hyper.Directory, error) {
	apiSettings := settings.(APIDirectorySettings)

	rootCerts := make([]*x509.Certificate, 0)

	for _, certificateFile := range apiSettings.CACertificateFiles {
		cert, err := helpers.LoadCertificate(certificateFile, false)

		if err != nil {
			return nil, fmt.Errorf("error loading API directory root certificate: %w", err)
		}

		rootCerts = append(rootCerts, cert)

	}

	intermediateCerts := make([]*x509.Certificate, 0)

	for _, certificateFile := range apiSettings.CAIntermediateCertificateFiles {
		cert, err := helpers.LoadCertificate(certificateFile, false)

		if err != nil {
			return nil, fmt.Errorf("error loading API directory intermediate certificate: %w", err)
		}

		intermediateCerts = append(intermediateCerts, cert)

	}

	d := &APIDirectory{
		BaseDirectory: hyper.BaseDirectory{
			Name_: name,
		},
		jsonrpcClient:     jsonrpc.MakeClient(apiSettings.JSONRPCClient),
		entries:           make(map[string]*hyper.DirectoryEntry),
		records:           []*hyper.SignedChangeRecord{},
		rootCerts:         rootCerts,
		intermediateCerts: intermediateCerts,
		settings:          apiSettings,
	}

	// we still allow the services to start even if the API is not reachable...
	if err := d.update(); err != nil {
		hyper.Log.Error(err)
	}

	return d, nil
}

var UpdateForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "records",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &hyperForms.SignedChangeRecordForm,
						},
					},
				},
			},
		},
	},
}

type UpdateRecords struct {
	Records []*hyper.SignedChangeRecord `json:"records"`
}

func (f *APIDirectory) Entries(query *hyper.DirectoryQuery) ([]*hyper.DirectoryEntry, error) {

	f.mutex.Lock()
	lastUpdate := f.lastUpdate
	f.mutex.Unlock()

	if time.Now().Add(-time.Duration(2*f.settings.CacheEntriesFor) * time.Second).After(lastUpdate) {
		// last update was more than 2 minutes ago, we update synchronously
		if err := f.update(); err != nil {
			return nil, fmt.Errorf("error updating service directory: %w", err)
		}
	} else if time.Now().Add(-time.Duration(f.settings.CacheEntriesFor) * time.Second).After(lastUpdate) {
		// last update was more than 1 minute ago, we update in the background
		go func() {
			if err := f.update(); err != nil {
				hyper.Log.Error(err)
			}
		}()
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	entries := make([]*hyper.DirectoryEntry, len(f.entries))
	i := 0
	for _, entry := range f.entries {
		entries[i] = entry
		i++
	}
	return hyper.FilterDirectoryEntriesByQuery(entries, query), nil
}

func (f *APIDirectory) EntryFor(name string) (*hyper.DirectoryEntry, error) {
	// locking is done by Entries method
	if entries, err := f.Entries(&hyper.DirectoryQuery{Operator: name}); err != nil {
		return nil, fmt.Errorf("error retrieving service directory entry: %w", err)
	} else if len(entries) == 0 {
		return nil, hyper.NoEntryFound
	} else {
		return entries[0], nil
	}
}

func (f *APIDirectory) OwnEntry() (*hyper.DirectoryEntry, error) {
	// locking is done by Entries method
	return f.EntryFor(f.Name())
}

func (f *APIDirectory) Tip() (*hyper.SignedChangeRecord, error) {

	// to do: ensure there's always one server name and endpoint
	f.jsonrpcClient.SetServerName(f.settings.ServerNames[0])
	f.jsonrpcClient.SetEndpoint(f.settings.Endpoints[0])

	request := jsonrpc.MakeRequest("getTip", "", map[string]interface{}{})

	if result, err := f.jsonrpcClient.Call(request); err != nil {
		return nil, fmt.Errorf("error getting tip from service directory: %w", err)
	} else {
		if result.Error != nil {
			return nil, fmt.Errorf("JSON-RPC error: %s", result.Error.Message)
		}

		if result.Result == nil {
			return nil, nil
		}

		if mapResult, ok := result.Result.(map[string]interface{}); !ok {
			return nil, fmt.Errorf("expected a map as result for 'getTip' call to service directory")
		} else if params, err := hyperForms.SignedChangeRecordForm.Validate(mapResult); err != nil {
			return nil, err
		} else {
			signedChangeRecord := &hyper.SignedChangeRecord{}
			if err := hyperForms.SignedChangeRecordForm.Coerce(signedChangeRecord, params); err != nil {
				return nil, err
			} else {
				return signedChangeRecord, nil
			}
		}
	}

	return nil, nil
}

func (f *APIDirectory) Submit(signedChangeRecords []*hyper.SignedChangeRecord) error {
	// to do: ensure there's always one server name and endpoint
	f.jsonrpcClient.SetServerName(f.settings.ServerNames[0])
	f.jsonrpcClient.SetEndpoint(f.settings.Endpoints[0])

	// we tell the internal proxy about an incoming connection
	request := jsonrpc.MakeRequest("submitRecords", "", map[string]interface{}{"records": signedChangeRecords})

	if result, err := f.jsonrpcClient.Call(request); err != nil {
		return fmt.Errorf("error submitting records to service directory: %w", err)
	} else {
		if result.Error != nil {
			hyper.Log.Error(result.Error)
			return fmt.Errorf("JSON-RPC error: %s", result.Error.Message)
		}
		return nil
	}

}

func (f *APIDirectory) integrate(records []*hyper.SignedChangeRecord) error {
	for _, record := range records {
		entry, ok := f.entries[record.Record.Name]
		if !ok {
			entry = hyper.MakeDirectoryEntry()
			entry.Name = record.Record.Name
		}
		if err := helpers.IntegrateChangeRecord(record, entry); err != nil {
			return fmt.Errorf("error integrating change record: %w", err)
		}
		f.entries[record.Record.Name] = entry
	}
	return nil
}

// Updates the service directory with change records from the remote API
func (f *APIDirectory) update() error {

	f.mutex.Lock()
	defer f.mutex.Unlock()

	if time.Now().Add(-time.Duration(f.settings.CacheEntriesFor) * time.Second).Before(f.lastUpdate) {
		// last update was less than a minute ago...
		return nil
	}

	hyper.Log.Tracef("Updating service directory...")
	f.lastUpdate = time.Now()

	// to do: ensure there's always one server name and endpoint
	f.jsonrpcClient.SetServerName(f.settings.ServerNames[0])
	f.jsonrpcClient.SetEndpoint(f.settings.Endpoints[0])

	var tipHash string

	if len(f.records) > 0 {
		tipHash = f.records[len(f.records)-1].Hash
	}

	// we tell the internal proxy about an incoming connection
	request := jsonrpc.MakeRequest("getRecords", "", map[string]interface{}{"after": tipHash})

	if result, err := f.jsonrpcClient.Call(request); err != nil {
		return fmt.Errorf("error getting records from service directory: %w", err)
	} else {

		if result.Error != nil {
			return fmt.Errorf("JSON-RPC error: %s", result.Error.Message)
		}

		if result.Result == nil {
			return nil
		}

		config := map[string]interface{}{
			"records": result.Result,
		}

		if params, err := UpdateForm.Validate(config); err != nil {
			return err
		} else {
			updateRecords := &UpdateRecords{}
			if err := UpdateForm.Coerce(updateRecords, params); err != nil {
				return err
			} else {
				records := updateRecords.Records

				var fullRecords []*hyper.SignedChangeRecord
				var resetEntries bool

				if len(records) > 0 && records[0].ParentHash != tipHash {
					if records[0].ParentHash != "" {
						return fmt.Errorf("expected a new root record but got one with parent hash '%s'", records[0].ParentHash)
					}
					// seems the directory changed, we make sure the new one is actually newer than the current one
					if len(f.records) > 0 && records[0].Record.CreatedAt.Time.Before(f.records[0].Record.CreatedAt.Time) {
						return fmt.Errorf("server tried to provide an outdated service directory")
					} else {
						hyper.Log.Warning("Service directory root changed!")
						// we reset the entries
						fullRecords = records
						resetEntries = true
					}
				} else {
					fullRecords = append(f.records, records...)
				}

				// we verify all records before we integate them
				for i, record := range fullRecords {
					if ok, err := helpers.VerifyRecord(record, fullRecords[:i], f.rootCerts, f.intermediateCerts); err != nil {
						return fmt.Errorf("cannot verify service directory record: %w", err)
					} else if !ok {
						return fmt.Errorf("invalid record found")
					}
				}

				f.records = fullRecords
				if resetEntries {
					f.entries = make(map[string]*hyper.DirectoryEntry)
				}

				// we integrate the new records
				return f.integrate(records)
			}
		}
	}
}
