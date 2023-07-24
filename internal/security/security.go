package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/nmramorov/gowatcher/internal/log"
)

func GetCertificate(path string) (*x509.Certificate, error) {
	log.InfoLog.Printf("certificate path: %s", path)
	encodedCert, err := GetCryptoKey(path)
	if err != nil {
		log.ErrorLog.Printf("error getting cert: %e", err)
		return nil, err
	}
	publicKeyBlock, _ := pem.Decode(encodedCert)
	certificate, err := x509.ParseCertificate(publicKeyBlock.Bytes)
	if err != nil {
		log.ErrorLog.Printf("error parsing public key: %e", err)
		return nil, err
	}
	return certificate, nil
}

func GetPrivateKey(path string) (*rsa.PrivateKey, error) {
	log.InfoLog.Printf("private key path: %s", path)
	encodedPrivateKey, err := GetCryptoKey(path)
	if err != nil {
		log.ErrorLog.Printf("error getting private key from file: %e", err)
		return nil, err
	}
	privateKeyBlock, _ := pem.Decode(encodedPrivateKey)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.ErrorLog.Printf("error parsing private key: %e", err)
		return nil, err
	}
	return privateKey, nil
}

func EncodeMsg(payload []byte, certificate *x509.Certificate) ([]byte, error) {
	// label := []byte("metrics")
	// ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, certificate.PublicKey.(*rsa.PublicKey), payload, label)
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, certificate.PublicKey.(*rsa.PublicKey), payload)

	if err != nil {
		log.ErrorLog.Printf("encryption error: %e", err)
		return nil, err
	}
	return ciphertext, nil
}

func DecodeMsg(msg []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// label := []byte("metrics")
	// decyphered, err := rsa.DecryptOAEP(sha256.New(), nil, privateKey, msg, label)
	decyphered, err := rsa.DecryptPKCS1v15(nil, privateKey, msg)
	if err != nil {
		log.ErrorLog.Printf("error decyphering message: %e", err)
	}
	return decyphered, nil
}

func GetCryptoKey(path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
