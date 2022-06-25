/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package result

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
)

type CertResult struct {
	BaseResult
	Alert           *CertAlert `json:"alert"`         //
	Method          *string    `json:"method"`        //
	ReplyTime       *float64   `json:"rt"`            //
	ServerCipher    *string    `json:"server_cipher"` //
	ConnectTime     *float64   `json:"ttc"`           //
	ProtocolVersion *string    `json:"ver"`           //
	Error           *string    `json:"err"`           //
	RawCertificates *[]string  `json:"cert"`          //
}

// CertAlert is an error could be sent by the server
type CertAlert struct {
	Level       uint //
	Description uint //
}

func (result *CertResult) ShortString() string {
	certs, _ := result.Certificates()
	ret := result.BaseShortString() +
		valueOrNA("", false, result.Error) +
		valueOrNA("", false, result.Method) +
		valueOrNA("", false, result.ProtocolVersion) +
		valueOrNA("", false, result.ReplyTime) +
		fmt.Sprintf("\t%d", len(certs))
	return ret
}

func (result *CertResult) LongString() string {
	res := result.ShortString() +
		valueOrNA("", false, result.ServerCipher)
	certs, err := result.Certificates()
	if err != nil || len(certs) == 0 {
		res += "N/A"
	} else {
		res += fmt.Sprintf("\t%s\t%s",
			certs[0].SerialNumber,
			certs[0].Subject.CommonName,
		)
	}
	return res
}

func (result *CertResult) TypeName() string {
	return "sslcert"
}

func (cert *CertResult) Parse(from string) (err error) {
	err = json.Unmarshal([]byte(from), &cert)
	if err != nil {
		return err
	}
	if cert.Type != "sslcert" {
		return fmt.Errorf("this is not a TLS/SSL certificate result (type=%s)", cert.Type)
	}
	return nil
}

func (result *CertResult) Certificates() (list []x509.Certificate, err error) {
	list = make([]x509.Certificate, 0)
	if result.RawCertificates == nil {
		return
	}
	for _, item := range *result.RawCertificates {
		block, _ := pem.Decode([]byte(item))
		if block == nil || block.Type != "CERTIFICATE" {
			log.Fatal("failed to decode PEM block containing certificate")
		}
		var cert *x509.Certificate
		cert, err = x509.ParseCertificate(block.Bytes)
		if err != nil {
			return list, err
		}
		list = append(list, *cert)
	}
	return list, nil
}
