package service

import (
	"strconv"

	"github.com/18F/cf-db-connect/launcher"
	"github.com/18F/cf-db-connect/models"
)

type mySQL struct{}

func (p mySQL) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("mysql")
}

func (p mySQL) Launch(localPort int, creds models.Credentials) error {
	return launcher.StartShell("mysql", []string{
		"-u", creds.GetUsername(),
		"-h", "0",
		"-p" + creds.GetPassword(),
		"-D", creds.GetDBName(),
		"-P", strconv.Itoa(localPort),
	})
}

// singleton
var MySQL = mySQL{}
