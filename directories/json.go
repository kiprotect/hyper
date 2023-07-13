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

// The JSON directory loads the service directory from a single JSON file.
// This is just for testing, for production use please use the "file" directory
// which provides support for signed service directory records.

package directories

import (
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/hyper"
	hyperForms "github.com/kiprotect/hyper/forms"
	"github.com/kiprotect/hyper/helpers"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sync"
)

var JSONDirectorySettingsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "path",
			Validators: []forms.Validator{
				forms.IsString{},
			},
		},
	},
}

var JSONRecordsForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "records",
			Validators: []forms.Validator{
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &hyperForms.ChangeRecordForm,
						},
					},
				},
			},
		},
	},
}

type JSONDirectorySettings struct {
	Path string `json:"path"`
}

type Records struct {
	Records []*hyper.ChangeRecord `json:"records"`
}

type JSONDirectory struct {
	hyper.BaseDirectory
	mutex    sync.Mutex
	settings JSONDirectorySettings
	records  []*hyper.ChangeRecord
	entries  map[string]*hyper.DirectoryEntry
}

func JSONDirectorySettingsValidator(settings map[string]interface{}) (interface{}, error) {
	if params, err := JSONDirectorySettingsForm.Validate(settings); err != nil {
		return nil, err
	} else {
		validatedSettings := &JSONDirectorySettings{}
		if err := JSONDirectorySettingsForm.Coerce(validatedSettings, params); err != nil {
			return nil, err
		}
		return validatedSettings, nil
	}
}

func MakeJSONDirectory(name string, settings interface{}) (hyper.Directory, error) {
	d := &JSONDirectory{
		BaseDirectory: hyper.BaseDirectory{
			Name_: name,
		},
		settings: settings.(JSONDirectorySettings),
	}

	return d, d.load()
}

func (f *JSONDirectory) load() error {
	if records, err := loadRecords(f.settings.Path); err != nil {
		return err
	} else {
		entries := make(map[string]*hyper.DirectoryEntry)

		for _, record := range records {
			entry, ok := entries[record.Name]
			if !ok {
				entry = hyper.MakeDirectoryEntry()
				entry.Name = record.Name
			}
			if err := helpers.IntegrateChangeRecord(&hyper.SignedChangeRecord{Record: record}, entry); err != nil {
				return err
			}
			entries[record.Name] = entry
		}

		f.entries = entries
		f.records = records

		hyper.Log.Tracef("Loaded %d directory entries from %d records...", len(f.entries), len(f.records))
		return nil
	}
}

func (f *JSONDirectory) Entries(query *hyper.DirectoryQuery) ([]*hyper.DirectoryEntry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if err := f.load(); err != nil {
		return nil, fmt.Errorf("error loading JSON directory: %w", err)
	}

	entries := make([]*hyper.DirectoryEntry, 0)

	for _, entry := range f.entries {
		entries = append(entries, entry)
	}

	return hyper.FilterDirectoryEntriesByQuery(entries, query), nil
}

func (f *JSONDirectory) EntryFor(name string) (*hyper.DirectoryEntry, error) {
	if entries, err := f.Entries(&hyper.DirectoryQuery{Operator: name}); err != nil {
		return nil, fmt.Errorf("error retrieving entry: %w", err)
	} else if len(entries) == 0 {
		return nil, hyper.NoEntryFound
	} else {
		return entries[0], nil
	}
}

func (f *JSONDirectory) OwnEntry() (*hyper.DirectoryEntry, error) {
	return f.EntryFor(f.Name())
}

func getRecordsFiles(recordsPath string) []string {
	paths := make([]string, 0)
	files, err := ioutil.ReadDir(recordsPath)
	if err != nil {
		return paths
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		r, err := regexp.MatchString(".json", file.Name())
		if err == nil && r {
			paths = append(paths, path.Join(recordsPath, file.Name()))
		}
	}
	return paths
}

func loadRecords(recordsPath string) ([]*hyper.ChangeRecord, error) {

	fi, err := os.Stat(recordsPath)
	if err != nil {
		return nil, fmt.Errorf("error retrieving records path info: %w", err)
	}
	var recordsFiles []string
	if fi.Mode().IsDir() {
		recordsFiles = getRecordsFiles(recordsPath)
	} else {
		recordsFiles = []string{recordsPath}
	}

	allRecords := make([]*hyper.ChangeRecord, 0)

	for _, recordsFile := range recordsFiles {
		hyper.Log.Tracef("Adding records from %v...", recordsFile)
		if file, err := os.Open(recordsFile); err != nil {
			return nil, fmt.Errorf("error opening records file '%s': %w", recordsFile, err)
		} else {
			if data, err := ioutil.ReadAll(file); err != nil {
				return nil, fmt.Errorf("error reading records file '%s': %w", recordsFile, err)
			} else {
				file.Close()
				rawRecords := map[string]interface{}{}
				if err := json.Unmarshal(data, &rawRecords); err != nil {
					return nil, fmt.Errorf("error parsing JSON from file '%s': %w", recordsFile, err)
				} else if params, err := JSONRecordsForm.Validate(rawRecords); err != nil {
					return nil, err
				} else {
					records := &Records{}
					if err := forms.Coerce(records, params); err != nil {
						// this should not happen if the forms are correct...
						return nil, err
					}
					// records are valid
					allRecords = append(allRecords, records.Records...)
				}
			}
		}
	}
	return allRecords, nil
}
