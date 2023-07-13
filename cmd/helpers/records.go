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
	"encoding/json"
	"fmt"
	"github.com/kiprotect/go-helpers/forms"
	"github.com/kiprotect/hyper"
	hyperForms "github.com/kiprotect/hyper/forms"
	"github.com/kiprotect/hyper/helpers"
	"github.com/urfave/cli"
	"io/ioutil"
	"time"
)

var RecordsForm = forms.Form{
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

func getEntries(c *cli.Context, settings *hyper.Settings) error {

	directory, err := helpers.InitializeDirectory(settings)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	query := &hyper.DirectoryQuery{}
	name := c.String("name")

	if name != "" {
		query.Operator = name
	}

	entries, err := directory.Entries(query)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	jsonData, err := json.Marshal(entries)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	fmt.Println(string(jsonData))
	return nil
}

type Records struct {
	Records []*hyper.ChangeRecord `json:"records"`
}

func submitRecords(c *cli.Context, settings *hyper.Settings) error {

	reset := c.Bool("reset")

	if settings.Signing == nil {
		hyper.Log.Fatalf("Signing settings undefined!")
	}

	filename := c.Args().Get(0)

	if filename == "" {
		hyper.Log.Fatal("please specify a filename")
	}

	jsonBytes, err := ioutil.ReadFile(filename)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	records := &Records{}
	var rawRecords map[string]interface{}

	if err := json.Unmarshal(jsonBytes, &rawRecords); err != nil {
		hyper.Log.Fatal(err)
	}

	if params, err := RecordsForm.Validate(rawRecords); err != nil {
		hyper.Log.Fatal(err)
	} else if RecordsForm.Coerce(records, params); err != nil {
		hyper.Log.Fatal(err)
	}

	if err := submitChangeRecords(records.Records, settings, reset); err != nil {
		hyper.Log.Fatal(err)
	}

	return nil

}

func submitChangeRecords(changeRecords []*hyper.ChangeRecord, settings *hyper.Settings, reset bool) error {

	directory, err := helpers.InitializeDirectory(settings)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	writableDirectory, ok := directory.(hyper.WritableDirectory)

	if !ok {
		hyper.Log.Fatalf("not a writable service directory")
	}

	certificate, err := helpers.LoadCertificate(settings.Signing.CertificateFile, true)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	rootCertificate, err := helpers.LoadCertificate(settings.Signing.CACertificateFile, false)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	intermediateCertificates := []*x509.Certificate{}

	for _, certificateFile := range settings.Signing.CAIntermediateCertificateFiles {
		if cert, err := helpers.LoadCertificate(certificateFile, false); err != nil {
			hyper.Log.Fatal(err)
		} else {
			intermediateCertificates = append(intermediateCertificates, cert)
		}
	}

	// we ensure the certificate is valid for signing
	if err := helpers.VerifyCertificate(certificate, rootCertificate, intermediateCertificates, settings.Name); err != nil {
		hyper.Log.Fatal(err)
	}

	key, err := helpers.LoadPrivateKey(settings.Signing.KeyFile)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	lastRecord, err := writableDirectory.Tip()

	if err != nil {
		hyper.Log.Fatal(err)
	}

	var parentHash string

	if lastRecord != nil && !reset {
		parentHash = lastRecord.Hash
	}

	signedChangeRecords := make([]*hyper.SignedChangeRecord, 0)

	for _, changeRecord := range changeRecords {

		changeRecord.CreatedAt = hyper.HashableTime{time.Now()}

		signedChangeRecord := &hyper.SignedChangeRecord{
			ParentHash: parentHash,
			Record:     changeRecord,
		}

		if err := helpers.CalculateRecordHash(signedChangeRecord); err != nil {
			hyper.Log.Fatal(err)
		}

		signedData, err := helpers.Sign(signedChangeRecord, key, certificate)

		if err != nil {
			hyper.Log.Fatal(err)
		}

		hyper.Log.Info(signedChangeRecord.Hash)

		if ok, err := helpers.Verify(signedData, []*x509.Certificate{rootCertificate}, intermediateCertificates, settings.Name); err != nil {
			hyper.Log.Fatal(err)
		} else if !ok {
			hyper.Log.Fatalf("cannot verify signature")
		}

		signedChangeRecord.Signature = signedData.Signature
		signedChangeRecords = append(signedChangeRecords, signedChangeRecord)
		parentHash = signedChangeRecord.Hash
	}

	if err := writableDirectory.Submit(signedChangeRecords); err != nil {
		hyper.Log.Fatal(err)
	}

	return nil

}

func sign(c *cli.Context, settings *hyper.Settings) error {
	if settings.Signing == nil {
		hyper.Log.Fatalf("Signing settings undefined!")
	}

	filename := c.Args().Get(0)

	if filename == "" {
		hyper.Log.Fatal("please specify a filename")
	}

	jsonBytes, err := ioutil.ReadFile(filename)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	var jsonData map[string]interface{}

	if err := json.Unmarshal(jsonBytes, &jsonData); err != nil {
		hyper.Log.Fatal(err)
	}

	certificate, err := helpers.LoadCertificate(settings.Signing.CertificateFile, true)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	rootCertificate, err := helpers.LoadCertificate(settings.Signing.CACertificateFile, false)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	intermediateCertificates := []*x509.Certificate{}

	for _, certificateFile := range settings.Signing.CAIntermediateCertificateFiles {
		if cert, err := helpers.LoadCertificate(certificateFile, false); err != nil {
			hyper.Log.Fatal(err)
		} else {
			intermediateCertificates = append(intermediateCertificates, cert)
		}
	}

	// we ensure the certificate is valid for signing
	if err := helpers.VerifyCertificate(certificate, rootCertificate, intermediateCertificates, settings.Name); err != nil {
		hyper.Log.Fatal(err)
	}

	key, err := helpers.LoadPrivateKey(settings.Signing.KeyFile)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	signedData, err := helpers.Sign(jsonData, key, certificate)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	signedDataBytes, err := json.Marshal(signedData)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	fmt.Println(string(signedDataBytes))

	loadedSignedData, err := helpers.LoadSignedData(signedDataBytes)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	if ok, err := helpers.Verify(loadedSignedData, []*x509.Certificate{rootCertificate}, intermediateCertificates, settings.Name); err != nil {
		hyper.Log.Fatal(err)
	} else if !ok {
		hyper.Log.Fatal("Signature is not valid!")
	}

	return nil
}

func verify(c *cli.Context, settings *hyper.Settings) error {

	if settings.Signing == nil {
		hyper.Log.Fatalf("Signing settings undefined!")
	}

	filename := c.Args().Get(0)

	if filename == "" {
		hyper.Log.Fatal("please specify a filename")
	}

	name := c.Args().Get(1)

	if name == "" {
		hyper.Log.Fatal("please specify a name")
	}

	jsonBytes, err := ioutil.ReadFile(filename)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	var signedData *hyper.SignedData

	if err := json.Unmarshal(jsonBytes, &signedData); err != nil {
		hyper.Log.Fatal(err)
	}

	rootCertificate, err := helpers.LoadCertificate(settings.Signing.CACertificateFile, false)

	if err != nil {
		hyper.Log.Fatal(err)
	}

	intermediateCertificates := []*x509.Certificate{}

	for _, certificateFile := range settings.Signing.CAIntermediateCertificateFiles {
		if cert, err := helpers.LoadCertificate(certificateFile, false); err != nil {
			hyper.Log.Fatal(err)
		} else {
			intermediateCertificates = append(intermediateCertificates, cert)
		}
	}

	if ok, err := helpers.Verify(signedData, []*x509.Certificate{rootCertificate}, intermediateCertificates, name); err != nil {
		hyper.Log.Fatal(err)
	} else if !ok {
		hyper.Log.Fatal("Signature is not valid!")
	} else {
		hyper.Log.Info("Signature is ok!")
	}
	return nil
}
func RecordsCommands(settings *hyper.Settings) ([]cli.Command, error) {

	return []cli.Command{
		{
			Name:    "sd",
			Aliases: []string{"s"},
			Flags:   []cli.Flag{},
			Usage:   "Manage service-directory records.",
			Subcommands: []cli.Command{
				{
					Name: "get-entries",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name",
							Usage: "the name of the entry to retrieve",
						},
					},
					Usage:  "Get all service diectory entries and print them as JSON",
					Action: func(c *cli.Context) error { return getEntries(c, settings) },
				},
				{
					Name: "submit-records",
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  "reset",
							Usage: "reset the remote records (dangerous)",
						},
					},
					Usage:  "Submit several records at once",
					Action: func(c *cli.Context) error { return submitRecords(c, settings) },
				},
				{
					Name:   "sign",
					Flags:  []cli.Flag{},
					Usage:  "Sign a change record",
					Action: func(c *cli.Context) error { return sign(c, settings) },
				},
				{
					Name:   "verify",
					Flags:  []cli.Flag{},
					Usage:  "Verify a change record",
					Action: func(c *cli.Context) error { return verify(c, settings) },
				},
			},
		},
	}, nil
}
