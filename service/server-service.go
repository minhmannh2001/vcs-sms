package service

import (
	"github.com/minhmannh2001/sms/database"
	"github.com/minhmannh2001/sms/entity"
)

type ServerService interface {
	CreateServer(server *entity.Server) error
	// CreateServers()
	ViewServer(id int) (entity.Server, error)
	ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error)
	UpdateServer(server *entity.Server) error
	DeleteServer(id int) error
	CheckServerExistence(ip string) bool
	CheckServerName(name string) bool
}

type serverService struct {
	serverDatabase database.SMSDatabase
}

func NewServerService(db database.SMSDatabase) ServerService {
	return &serverService{
		serverDatabase: db,
	}
}

func (service *serverService) CreateServer(server *entity.Server) error {
	err := service.serverDatabase.CreateServer(server)
	if err != nil {
		return err
	}
	return nil
}

func (service *serverService) ViewServer(id int) (entity.Server, error) {
	server, err := service.serverDatabase.ViewServer(id)
	return server, err
}

func (service *serverService) ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error) {
	return service.serverDatabase.ViewServers(from, to, perpage, sortby, order, filter)
}

func (service *serverService) UpdateServer(server *entity.Server) error {
	err := service.serverDatabase.UpdateServer(server)
	if err != nil {
		return err
	}

	return nil
}

func (service *serverService) DeleteServer(id int) error {
	err := service.serverDatabase.DeleteServer(id)
	return err
}

func (service *serverService) CheckServerExistence(ip string) bool {
	return service.serverDatabase.CheckServerExistence(ip)
}

func (service *serverService) CheckServerName(name string) bool {
	return service.serverDatabase.CheckServerName(name)
}
