package service

import (
	"strconv"

	"github.com/cloud-gov/cf-service-connect/launcher"
	"github.com/cloud-gov/cf-service-connect/models"
)

type mongoDB struct{}

func (p mongoDB) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("mongo")
}

func (p mongoDB) Launch(localPort int, creds models.Credentials) error {
	return launcher.StartShell("mongo", []string{
		"-u", creds.GetUsername(),
		"-p", creds.GetPassword(),
		"--port", strconv.Itoa(localPort),
		creds.GetDBName(),
	})
}

// MongoDB is the service singleton.
var MongoDB = mongoDB{}
