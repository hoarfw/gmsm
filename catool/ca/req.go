package ca

import (
	"encoding/asn1"
	//	"fmt"
	//	"io/ioutil"

	//	"errors"
	"crypto/elliptic"
	"crypto/x509/pkix"
//	"crypto/rand"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	 "github.com/hoarfw/gmsm/sm2"
	 "os"
	 "encoding/pem"

)

const requestCmdDescription = "request a channel"


func reqCmd() *cobra.Command {
	requestCmd := &cobra.Command{
		Use:   "request",
		Short: requestCmdDescription,
		Long:  requestCmdDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return request(cmd, args)
		},
	}
	// flagList := []string{
	// 	"x",
	// 	"y",
	// 	"template",
	// }
	//attachFlags(requestCmd, flagList)

	return requestCmd
}

func executeReq(prvK * sm2.PrivateKey , req *sm2.CertificateRequest) error {
	_, err := sm2.CreateCertificateRequestToPem("req.pem", req, prvK)
	if err != nil {
		logger.Fatal(err)
	}
	req, err2 := sm2.ReadCertificateRequestFromPem("req.pem")
	if err != nil {
		logger.Fatal(err2)
	}
	err = req.CheckSignature()
	if err != nil {
		logger.Fatal(err)
	} 
	return nil
}



func CreateReqForSign(pubK *sm2.PublicKey,req *sm2.CertificateRequest) (msg []byte ,err error){
	//rndR :=rand.Reader
	
	var publicKeyBytes []byte
	var publicKeyAlgorithm pkix.AlgorithmIdentifier
	
	publicKeyBytes = elliptic.Marshal(pubK.Curve, pubK.X, pubK.Y)
	oid := sm2.OidNamedCurveP256SM2
	
	publicKeyAlgorithm.Algorithm = sm2.OidPublicKeyECDSA


	paramBytes, err := asn1.Marshal(oid)
	if err != nil {
		return nil,err
	}
	publicKeyAlgorithm.Parameters.FullBytes = paramBytes


	publicKeyBytes, publicKeyAlgorithm, err = sm2.MarshalPublicKey(pubK)
	if err != nil {
		return nil, err
	}
	var extensions []pkix.Extension
	template := req
	if (len(template.DNSNames) > 0 || len(template.EmailAddresses) > 0 || len(template.IPAddresses) > 0) &&
		!sm2.OidInExtensions(sm2.OidExtensionSubjectAltName, template.ExtraExtensions) {
		sanBytes, err := sm2.MarshalSANs(template.DNSNames, template.EmailAddresses, template.IPAddresses)
		if err != nil {
			return nil, err
		}

		extensions = append(extensions, pkix.Extension{
			Id:    sm2.OidExtensionSubjectAltName,
			Value: sanBytes,
		})
	}

	extensions = append(extensions, template.ExtraExtensions...)

	var attributes []pkix.AttributeTypeAndValueSET
	attributes = append(attributes, template.Attributes...)

	if len(extensions) > 0 {
		// specifiedExtensions contains all the extensions that we
		// found specified via template.Attributes.
		specifiedExtensions := make(map[string]bool)

		for _, atvSet := range template.Attributes {
			if !atvSet.Type.Equal(sm2.OidExtensionRequest) {
				continue
			}

			for _, atvs := range atvSet.Value {
				for _, atv := range atvs {
					specifiedExtensions[atv.Type.String()] = true
				}
			}
		}

		atvs := make([]pkix.AttributeTypeAndValue, 0, len(extensions))
		for _, e := range extensions {
			if specifiedExtensions[e.Id.String()] {
				// Attributes already contained a value for
				// this extension and it takes priority.
				continue
			}

			atvs = append(atvs, pkix.AttributeTypeAndValue{
				// There is no place for the critical flag in a CSR.
				Type:  e.Id,
				Value: e.Value,
			})
		}

		// Append the extensions to an existing attribute if possible.
		appended := false
		for _, atvSet := range attributes {
			if !atvSet.Type.Equal(sm2.OidExtensionRequest) || len(atvSet.Value) == 0 {
				continue
			}

			atvSet.Value[0] = append(atvSet.Value[0], atvs...)
			appended = true
			break
		}

		// Otherwise, add a new attribute for the extensions.
		if !appended {
			attributes = append(attributes, pkix.AttributeTypeAndValueSET{
				Type: sm2.OidExtensionRequest,
				Value: [][]pkix.AttributeTypeAndValue{
					atvs,
				},
			})
		}
	}

	asn1Subject := template.RawSubject
	if len(asn1Subject) == 0 {
		asn1Subject, err = asn1.Marshal(template.Subject.ToRDNSequence())
		if err != nil {
			return
		}
	}

	rawAttributes, err := sm2.NewRawAttributes(attributes)
	if err != nil {
		return
	}

	tbsCSR := sm2.TbsCertificateRequest{
		Version: 0, // PKCS #10, RFC 2986
		Subject: asn1.RawValue{FullBytes: asn1Subject},
		PublicKey: sm2.PublicKeyInfo {
			Algorithm: publicKeyAlgorithm,
			PublicKey: asn1.BitString{
				Bytes:     publicKeyBytes,
				BitLength: len(publicKeyBytes) * 8,
			},
		},
		RawAttributes: rawAttributes,
		}

		tbsCSRContents, err := asn1.Marshal(tbsCSR)
		if err != nil {
		return nil, err
		}
		
		return tbsCSRContents,nil
	}	
	
	func CreateSignature ( tbsCSRContents []byte, priv *sm2.PrivateKey) (signs []byte ,seer error){

		signature, err := priv.Sign(nil, tbsCSRContents, nil)

		return signature, err
	}

	func SignReq(tbsCSRContents []byte, signature []byte ) (csr []byte ,serr error){

		var tbsCSR sm2.TbsCertificateRequest
		//signature, err := priv.Sign(nil, tbsCSRContents, nil)
		_, err :=asn1.Unmarshal(tbsCSRContents, &tbsCSR)
	
		if err != nil {
			return nil, err
		}

		tbsCSR.Raw = tbsCSRContents
		var sigAlgo pkix.AlgorithmIdentifier
		sigAlgo.Algorithm = sm2.OidSignatureSM2WithSM3

		return asn1.Marshal(sm2.SMCertificateRequest{
		TBSCSR:             tbsCSR,
		SignatureAlgorithm: sigAlgo,
		SignatureValue: asn1.BitString{
			Bytes:     signature,
			BitLength: len(signature) * 8,
		},
	})
	
} 

func CreateCertificateReqPem(filename string ,der []byte)(suc bool,serr error){
	
		block := &pem.Block{
			Type:  "CERTIFICATE REQUEST",
			Bytes: der,
		}
		file, err := os.Create(filename)
		if err != nil {
			return false, err
		}
		defer file.Close()
		err = pem.Encode(file, block)
		if err != nil {
			return false, err
		}
		return true, nil
	

}


func request(cmd *cobra.Command, args []string) error {
	//cmd.Name() == upgradeCmdName
	prvKeyPem := viper.GetString("req.PrivKey")
	prvKey, err := sm2.ReadPrivateKeyFromPem(prvKeyPem, nil) 
	//pubKeyPem := viper.GetString("req.PubKey")
	//pubKey, err := sm2.ReadPublicKeyFromPem(pubKeyPem, nil) // 读取公钥

	commonName := viper.GetString("req.Subject.CommonName")
	organization :=viper.GetString("req.Subject.Organization")
	templateReq := sm2.CertificateRequest{
		Subject: pkix.Name{
			CommonName: commonName,
			Organization: []string{organization},
		},
		//		SignatureAlgorithm: ECDSAWithSHA256,
		SignatureAlgorithm: sm2.SM2WithSM3,
	}

	if err != nil {
		logger.Fatal(err)
	}
	return executeReq(prvKey,&templateReq)
}
