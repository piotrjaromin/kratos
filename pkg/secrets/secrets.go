package secrets

import (
	"fmt"

	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type Secrets struct {
	kms *kms.KMS
}

func NewSecrets() Secrets {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := kms.New(sess)

	return Secrets{
		kms: svc,
	}
}

func (s Secrets) Decrypt(encodedSecret string) (string, error) {
	blob, err := base64.StdEncoding.DecodeString(encodedSecret)
	if err != nil {
		return "", fmt.Errorf("Error while doing base64 decoding %s", err.Error())
	}

	result, err := s.kms.Decrypt(&kms.DecryptInput{CiphertextBlob: blob})
	if err != nil {
		return "", fmt.Errorf("Error while calling kms %s", err.Error())
	}

	decodedSecret := string(result.Plaintext)
	return decodedSecret, nil
}
