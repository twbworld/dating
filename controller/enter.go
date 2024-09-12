package controller

import "github.com/twbworld/dating/controller/user"
import "github.com/twbworld/dating/controller/admin"

var Api = new(ApiGroup)

type ApiGroup struct {
	UserApiGroup  user.ApiGroup
	AdminApiGroup admin.ApiGroup
}
