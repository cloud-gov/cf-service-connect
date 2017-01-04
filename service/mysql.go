package service

import "github.com/18F/cf-db-connect/models"

type MySQL struct{}

func (p MySQL) Match(si models.ServiceInstance) bool {
	return si.IsMySQLService()
}
