package enroll

import (
	"errors"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
)

const PushTopicASN1 string = "0.9.2342.19200300.100.1.1"

func GetPushTopicFromPKCS12(certPath string, certPass string) (string, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", err
	}

	_, cert, err := pkcs12.Decode(certData, certPass)
	if err != nil {
		return "", err
	}

	for _, v := range cert.Subject.Names {
		if v.Type.String() == PushTopicASN1 {
			return v.Value.(string), nil
		}
	}

	return "", errors.New("Could not find Push Topic in the provided pkcs12 bundle.")
}
