package private_keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"time"
)

type CA struct {
	key       *rsa.PrivateKey
	publicKey rsa.PublicKey
	ca        []byte
	keyPem    []byte
	certPem   []byte
}

func NewCA(serviceName string) (*CA, error) {
	ca := &CA{}
	err := ca.generatePrivateKey()
	if err != nil {
		return nil, err
	}
	err = ca.generateCertificate(serviceName)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (ca *CA) PublicKey() rsa.PublicKey {
	return ca.publicKey
}

func (ca *CA) CA() []byte {
	return ca.ca
}

func (ca *CA) KeyPem() []byte {
	return ca.keyPem
}

func (ca *CA) CertPem() []byte {
	return ca.certPem
}

func (ca *CA) generatePrivateKey() error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Printf("generate private key err: %s\n", err)
		return err
	}
	ca.key = privateKey
	ca.publicKey = privateKey.PublicKey
	return nil
}

func (ca *CA) generateCertificate(serviceName string) error {
	maxInt := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, maxInt)
	if err != nil {
		log.Printf("generate serial number err: %s\n", err)
		return err
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"002 Co. Ltd"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),                                 // 有限期10年
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment, // 用于数字签名和私钥加密
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
		DNSNames:              []string{serviceName},
	}
	ca.ca, err = x509.CreateCertificate(rand.Reader, &template, &template, &ca.publicKey, ca.key)
	if err != nil {
		log.Printf("create certificate err: %s\n", err)
		return err
	}
	ca.keyPem = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(ca.key)})
	ca.certPem = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: ca.ca})
	fmt.Println(string(ca.keyPem))
	fmt.Println(string(ca.certPem))
	return nil
}
