package service

import "github.com/twbworld/dating/service/user"
import "github.com/twbworld/dating/service/admin"

var Service = new(ServiceGroup)

type ServiceGroup struct {
	UserServiceGroup  user.ServiceGroup
	AdminServiceGroup admin.ServiceGroup
}
