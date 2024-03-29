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

package sd

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/helpers"
	"sync"
	"time"
)

const (
	SignedChangeRecordEntry uint8 = 1
)

type RecordDirectorySettings struct {
	Datastore                      *hyper.DatastoreSettings `json:"datastore"`
	CACertificateFiles             []string                 `json:"ca_certificate_files"`
	CAIntermediateCertificateFiles []string                 `json:"ca_intermediate_certificate_files"`
}

type RecordDirectory struct {
	rootCerts         []*x509.Certificate
	intermediateCerts []*x509.Certificate
	dataStore         hyper.Datastore
	settings          *RecordDirectorySettings
	entries           map[string]*hyper.DirectoryEntry
	recordsByHash     map[string]*hyper.SignedChangeRecord
	recordChildren    map[string][]*hyper.SignedChangeRecord
	orderedRecords    []*hyper.SignedChangeRecord
	mutex             sync.Mutex
}

func MakeRecordDirectory(settings *RecordDirectorySettings, definitions *hyper.Definitions) (*RecordDirectory, error) {

	rootCerts := make([]*x509.Certificate, 0)

	for _, certificateFile := range settings.CACertificateFiles {

		cert, err := helpers.LoadCertificate(certificateFile, false)

		if err != nil {
			return nil, err
		}

		rootCerts = append(rootCerts, cert)

	}

	intermediateCerts := make([]*x509.Certificate, 0)

	for _, certificateFile := range settings.CAIntermediateCertificateFiles {

		cert, err := helpers.LoadCertificate(certificateFile, false)

		if err != nil {
			return nil, err
		}

		intermediateCerts = append(intermediateCerts, cert)

	}

	dataStore, err := helpers.InitializeDatastore(settings.Datastore, definitions)

	if err != nil {
		return nil, err
	}

	f := &RecordDirectory{
		rootCerts:         rootCerts,
		intermediateCerts: intermediateCerts,
		orderedRecords:    make([]*hyper.SignedChangeRecord, 0),
		recordsByHash:     make(map[string]*hyper.SignedChangeRecord),
		recordChildren:    make(map[string][]*hyper.SignedChangeRecord),
		settings:          settings,
		dataStore:         dataStore,
	}

	if err = f.dataStore.Init(); err != nil {
		return nil, err
	}

	_, err = f.update()

	return f, err
}

func (f *RecordDirectory) Entry(name string) (*hyper.DirectoryEntry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if entry, ok := f.entries[name]; !ok {
		return nil, nil
	} else {
		return entry, nil
	}
}

func (f *RecordDirectory) AllEntries() ([]*hyper.DirectoryEntry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	entries := make([]*hyper.DirectoryEntry, len(f.entries))
	i := 0
	for _, entry := range f.entries {
		entries[i] = entry
		i++
	}
	return entries, nil
}

func (f *RecordDirectory) Entries(*hyper.DirectoryQuery) ([]*hyper.DirectoryEntry, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return nil, nil
}

// determines whether a subject can append to the service directory
func (f *RecordDirectory) canAppend(record *hyper.SignedChangeRecord, records []*hyper.SignedChangeRecord) (bool, error) {
	return helpers.VerifyRecord(record, records, f.rootCerts, f.intermediateCerts)
}

// Appends a series of records
func (f *RecordDirectory) Append(records []*hyper.SignedChangeRecord) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if _, err := f.update(); err != nil {
		return err
	}

	for _, record := range records {

		records := f.orderedRecords

		if record.ParentHash != "" {

			tip, err := f.tip()

			if err != nil {
				return err
			}

			if (tip != nil && record.ParentHash != tip.Hash) || (tip == nil && record.ParentHash != "") {
				return fmt.Errorf("stale append, please try again")
			}
		}

		if ok, err := f.canAppend(record, records); err != nil {
			return err
		} else if !ok {
			return fmt.Errorf("you cannot append")
		}

		id, err := helpers.RandomID(16)

		if err != nil {
			return err
		}

		rawData, err := json.Marshal(record)

		if err != nil {
			return err
		}

		dataEntry := &hyper.DataEntry{
			Type: SignedChangeRecordEntry,
			ID:   id,
			Data: rawData,
		}

		if err := f.dataStore.Write(dataEntry); err != nil {
			return err
		}

		// we update the store
		if newRecords, err := f.update(); err != nil {
			return err
		} else {
			found := false
			for _, newRecord := range newRecords {
				if newRecord.Hash == record.Hash {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("new record not found")
			}
		}
	}
	return nil
}

func (f *RecordDirectory) Update() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	_, err := f.update()
	return err
}

// Returns the latest record
func (f *RecordDirectory) Tip() (*hyper.SignedChangeRecord, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.tip()
}

func (f *RecordDirectory) tip() (*hyper.SignedChangeRecord, error) {
	if len(f.orderedRecords) == 0 {
		return nil, nil
	}
	return f.orderedRecords[len(f.orderedRecords)-1], nil
}

// Returns all records after a given hash
func (f *RecordDirectory) Records(after string) ([]*hyper.SignedChangeRecord, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	relevantRecords := make([]*hyper.SignedChangeRecord, 0)
	found := false
	if after == "" {
		found = true
	}
	for _, record := range f.orderedRecords {
		if found {
			relevantRecords = append(relevantRecords, record)
		}
		if record.Hash == after {
			found = true
		}
	}
	if !found {
		// we can't find the hash, so we return all records instead
		// (as the client probably has an outdated version of the directory)
		return f.orderedRecords, nil
	}
	return relevantRecords, nil
}

// Integrates a record into the directory
func (f *RecordDirectory) integrate(record *hyper.SignedChangeRecord) error {
	entry, ok := f.entries[record.Record.Name]
	if !ok {
		entry = hyper.MakeDirectoryEntry()
		entry.Name = record.Record.Name
	}
	if err := helpers.IntegrateChangeRecord(record, entry); err != nil {
		return err
	}
	f.entries[record.Record.Name] = entry
	return nil
}

// picks the best record from a series of alternatives (based on chain length)
func (f *RecordDirectory) buildChains(records []*hyper.SignedChangeRecord, visited map[string]bool) ([][]*hyper.SignedChangeRecord, error) {

	chains := make([][]*hyper.SignedChangeRecord, 0)

	for _, record := range records {
		if _, ok := visited[record.Hash]; ok {
			continue
		} else {
			visited[record.Hash] = true
		}
		chain := make([]*hyper.SignedChangeRecord, 1)
		chain[0] = record
		children, ok := f.recordChildren[record.Hash]
		if ok {
			childChains, err := f.buildChains(children, visited)
			if err != nil {
				return nil, err
			}
			for _, childChain := range childChains {
				fullChain := append(chain, childChain...)
				chains = append(chains, fullChain)
			}
		} else {
			chains = append(chains, chain)
		}
	}

	return chains, nil

}

func (f *RecordDirectory) update() ([]*hyper.SignedChangeRecord, error) {
	if entries, err := f.dataStore.Read(); err != nil {
		return nil, err
	} else {
		changeRecords := make([]*hyper.SignedChangeRecord, 0, len(entries))
		for _, entry := range entries {
			switch entry.Type {
			case SignedChangeRecordEntry:
				record := &hyper.SignedChangeRecord{}
				if err := json.Unmarshal(entry.Data, &record); err != nil {
					return nil, fmt.Errorf("invalid record format!")
				}
				changeRecords = append(changeRecords, record)
			default:
				return nil, fmt.Errorf("unknown entry type found...")
			}
		}

		for _, record := range changeRecords {
			f.recordsByHash[record.Hash] = record
		}

		for _, record := range changeRecords {
			var parentHash string
			// if a parent exists we set the hash to it. Records without
			// parent (root records) will be stored under the empty hash...
			if parentRecord, ok := f.recordsByHash[record.ParentHash]; ok {
				parentHash = parentRecord.Hash
			}
			children, ok := f.recordChildren[parentHash]
			if !ok {
				children = make([]*hyper.SignedChangeRecord, 0)
			}
			children = append(children, record)
			f.recordChildren[parentHash] = children
		}

		rootRecords, ok := f.recordChildren[""]

		// no records present it seems
		if !ok {
			return nil, nil
		}

		chains, err := f.buildChains(rootRecords, map[string]bool{})

		hyper.Log.Infof("Found %d chains, %d root records", len(chains), len(rootRecords))

		if err != nil {
			return nil, err
		}

		verifiedChains := make([][]*hyper.SignedChangeRecord, 0)
		for i, chain := range chains {
			validRecords := make([]*hyper.SignedChangeRecord, 0)
			for j, record := range chain {
				hyper.Log.Infof("Chain %d, record %d: %s", i, j, record.Hash)
				// we verify the signature of the record
				if ok, err := helpers.VerifyRecord(record, chain[:j], f.rootCerts, f.intermediateCerts); err != nil {
					hyper.Log.Errorf("Warning, error verifying record: %v", err)
					continue
				} else if !ok {
					hyper.Log.Warning("signature does not match, ignoring this chain...")
					continue
				} else {
					validRecords = append(validRecords, record)
				}
			}
			if len(validRecords) > 0 {
				verifiedChains = append(verifiedChains, validRecords)
			}
		}

		hyper.Log.Infof("%d verified chains", len(verifiedChains))

		// the most recently created chain always wins
		var bestChain []*hyper.SignedChangeRecord
		var maxCreatedAt time.Time
		for _, chain := range verifiedChains {
			if bestChain == nil || (len(chain) > 0 && chain[0].Record.CreatedAt.Time.After(maxCreatedAt)) {
				bestChain = chain
				maxCreatedAt = chain[0].Record.CreatedAt.Time
			}
		}

		hyper.Log.Infof("Best chain created at %v with length %d", maxCreatedAt, len(bestChain))

		if bestChain == nil {
			return nil, nil
		}

		// we store the ordered sequence of records
		f.orderedRecords = bestChain

		// we regenerate the entries based on the new set of records
		f.entries = make(map[string]*hyper.DirectoryEntry)
		for _, record := range bestChain {
			f.integrate(record)
		}

		return bestChain, nil
	}
}
