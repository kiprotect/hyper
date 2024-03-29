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

package helpers

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/hyper"
	hyperForms "github.com/kiprotect/hyper/forms"
)

func InitializeDirectory(settings *hyper.Settings) (hyper.Directory, error) {
	definition := settings.Definitions.DirectoryDefinitions[settings.Directory.Type]
	return definition.Maker(settings.Name, settings.Directory.Settings)
}

// Integrates a record into the directory
func IntegrateChangeRecord(record *hyper.SignedChangeRecord, entry *hyper.DirectoryEntry) error {

	config := map[string]interface{}{
		record.Record.Section: record.Record.Data,
	}

	// we directly coerce the updated settings into the entry
	if err := hyperForms.DirectoryEntryForm.Coerce(entry, config); err != nil {
		return err
	} else {
		if entry.Records == nil {
			entry.Records = make([]*hyper.SignedChangeRecord, 0)
		}
		// we append the change record to the entry for audit logging purposes
		entry.Records = append(entry.Records, record)
	}
	return nil
}

type CertificatesList struct {
	Certificates []*hyper.OperatorCertificate `json:"certificates"`
}

var CertificatesListForm = forms.Form{
	Fields: []forms.Field{
		{
			Name: "certificates",
			Validators: []forms.Validator{
				forms.IsOptional{Default: []interface{}{}},
				forms.IsList{
					Validators: []forms.Validator{
						forms.IsStringMap{
							Form: &hyperForms.OperatorCertificateForm,
						},
					},
				},
			},
		},
	},
}

func GetRecordFingerprint(records []*hyper.SignedChangeRecord, name, keyUsage string) string {
	lastFingerprint := ""
	for _, record := range records {
		if record.Record.Section != "certificates" || record.Record.Name != name {
			continue
		}
		if params, err := CertificatesListForm.Validate(map[string]interface{}{"certificates": record.Record.Data}); err != nil {
			hyper.Log.Error(err)
			continue
		} else {
			certificatesList := &CertificatesList{}
			if err := CertificatesListForm.Coerce(certificatesList, params); err != nil {
				hyper.Log.Error(err)
				continue
			}
			for _, certificate := range certificatesList.Certificates {
				if certificate.KeyUsage == keyUsage {
					lastFingerprint = certificate.Fingerprint
				}
			}
		}
	}
	return lastFingerprint
}

func VerifyRecordHash(record *hyper.SignedChangeRecord) (bool, error) {

	submittedHash := record.Hash

	err := CalculateRecordHash(record)

	if err != nil {
		return false, err
	}

	if submittedHash != record.Hash {
		return false, nil
	}

	return true, nil
}

func VerifyRecord(record *hyper.SignedChangeRecord, verifiedRecords []*hyper.SignedChangeRecord, rootCerts []*x509.Certificate, intermediateCerts []*x509.Certificate) (bool, error) {
	signedData := &hyper.SignedData{
		Data:      record,
		Signature: record.Signature,
	}

	if ok, err := VerifyRecordHash(record); err != nil {
		return false, fmt.Errorf("error verifying record hash: %w", err)
	} else if !ok {
		return false, fmt.Errorf("invalid hash value")
	}

	// we temporarily remove the signature from the signed record
	signature := record.Signature
	record.Signature = nil
	defer func() { record.Signature = signature }()

	cert, err := LoadCertificateFromString(signature.Certificate, true)

	if err != nil {
		return false, fmt.Errorf("error loading entry certificate: %w", err)
	}

	subjectInfo, err := GetSubjectInfo(cert)

	if err != nil {
		return false, fmt.Errorf("error retrieving subject info: %w", err)
	}

	admin := false
	for _, group := range subjectInfo.Groups {
		if group == "sd-admin" {
			// service directory admins can upload its own certificate info
			// (but only if that info doesn't exist yet)
			admin = true
			break
		}
	}

	// only service-directory admins can continue
	if !admin {
		return false, nil
	}

	fingerprint := GetRecordFingerprint(verifiedRecords, subjectInfo.Name, "signing")

	if fingerprint != "" {
		if !VerifyFingerprint(cert, fingerprint) {
			// the fingerprint does not match the one we have on record
			return false, nil
		}
	}

	// finally we verify the cryptographic signature
	return Verify(signedData, rootCerts, intermediateCerts, "")
}

func CalculateRecordHash(record *hyper.SignedChangeRecord) error {

	// we always reset the hash before calculating the new one
	record.Hash = ""

	hash, err := StructuredHash(record.Record)

	if err != nil {
		return fmt.Errorf("error calculating record hash: %w", err)
	}

	record.Hash = hex.EncodeToString(hash[:])

	return nil

}
