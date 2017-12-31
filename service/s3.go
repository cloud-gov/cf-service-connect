package service

import (
	"fmt"
	"os"

	"github.com/18F/cf-service-connect/launcher"
	"github.com/18F/cf-service-connect/models"
)

type s3 struct{}

func (p s3) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("s3")
}

func (p s3) Launch(localPort int, creds models.Credentials) error {
	// note that the port/tunnel isn't actually used

	// shoved support for the AWS credentials style into the database style
	os.Setenv("AWS_DEFAULT_REGION", creds.GetRegion())
	os.Setenv("AWS_ACCESS_KEY_ID", creds.GetUsername())
	os.Setenv("AWS_SECRET_ACCESS_KEY", creds.GetPassword())
	fmt.Printf("Bucket: %s\n", creds.GetDBName())

	return launcher.StartShell("aws-shell", []string{})
}

// S3 is the service singleton.
var S3 = s3{}
