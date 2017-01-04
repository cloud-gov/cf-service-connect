package service

import (
	"fmt"
	"os"

	"github.com/18F/cf-db-connect/launcher"
	"github.com/18F/cf-db-connect/logger"
	"github.com/18F/cf-db-connect/models"
)

type pSQL struct{}

func (p pSQL) Match(si models.ServiceInstance) bool {
	return si.ContainsTerms("psql", "postgres")
}

func (p pSQL) Launch(localPort int, creds models.Credentials) error {
	os.Setenv("PGPASSWORD", creds.GetPassword())
	logger.Debugf("PGPASSWORD=%s ", creds.GetPassword())

	return launcher.StartShell("psql", []string{
		"-h", "localhost",
		"-p", fmt.Sprintf("%d", localPort),
		creds.GetDBName(),
		creds.GetUsername(),
	})
}

// singleton
var PSQL = pSQL{}
