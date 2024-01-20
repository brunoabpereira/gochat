package utils

import (
	"os"
	"crypto/rsa"
	"crypto/x509"
)

func GetEnvVar(name string, dflt string) string {
	if val, ok := os.LookupEnv(name); ok {
		return val
	}
	return dflt
}

func ReadJWTKey(jwtKeyFilename string) *rsa.PublicKey {
	raw, err := os.ReadFile(jwtKeyFilename)
	if err != nil {
		panic("failed to read public key file" + err.Error())
	}
	pub, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	return pub.(*rsa.PublicKey)
}