package service

import "github.com/18F/cf-db-connect/models"

type PSQL struct{}

func (p PSQL) Match(si models.ServiceInstance) bool {
	return si.IsPSQLService()
}
