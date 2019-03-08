package common

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"log"
	"os"
)

func GetCredentialsFromEnvironment(accessKeyId, secretKey string) (*credentials.Credentials, error) {
	accessKeyVal := os.Getenv(accessKeyId)
	secretKeyVal := os.Getenv(secretKey)

	if accessKeyVal == "" || secretKeyVal == "" {
		log.Println("Could not load " + accessKeyId + " (or) " + secretKey + " from local env")
		return nil, errors.New("Could not load " + accessKeyId + " (or) " + secretKey + " from local env")
	} else {
		return credentials.NewStaticCredentials(accessKeyVal, secretKeyVal, ""), nil
	}
}
