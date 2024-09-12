package user

import (
	"path"

	"github.com/gin-gonic/gin"

	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/service"
	"github.com/twbworld/dating/utils"
)

type BaseApi struct{}

// 登录注册
func (b *BaseApi) Login(ctx *gin.Context) {
	var (
		data common.LoginPost
		u    db.User
		err  error
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `系统错误[onsds]`)
		}
	}()

	if ctx.ShouldBindJSON(&data) != nil {
		common.Fail(ctx, `参数错误[ddsssj]`)
		return
	}
	if err = service.Service.UserServiceGroup.Validator.ValidatorLoginPost(&data); err != nil {
		common.Fail(ctx, err.Error())
		return
	}
	if gin.Mode() == gin.TestMode {
		u.OpenId = ""
		u.UnionId = ""
		u.SessionKey = ""
	} else {
		if err = service.Service.UserServiceGroup.DatingService.GetUserByCode(data.Code, &u); err != nil {
			common.Fail(ctx, err.Error())
			return
		}
	}

	token := ""
	if u.Id > 0 {
		//userId == 0, 则前端将会请求UserAdd()进行注册
		token, err = service.Service.UserServiceGroup.BaseService.LoginToken(&u)
		if err != nil {
			global.Log.Error(err)
			common.Fail(ctx, `系统错误[ojhtgnds]`)
			return
		}
	}

	ur, f := common.UserInfo{User: u}, db.File{}
	if u.Avatar > 0 && dao.App.FileDb.GetFile(&f, u.Avatar) == nil {
		ur.AvatarUrl = f.Path
	}

	common.SuccessAuth(ctx, token, gin.H{
		`user`: ur,
	})
}

// 上传头像
func (b *BaseApi) Upload(ctx *gin.Context) {
	var (
		err error
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `系统错误[on7sds]`)
		}
	}()

	file, err := ctx.FormFile("file")
	if err != nil {
		common.Fail(ctx, `参数错误[ipjdmk]`)
		return
	}

	if err = service.Service.UserServiceGroup.Validator.ValidatorUpload(file); err != nil {
		common.Fail(ctx, err.Error())
		return
	}

	dir, f := utils.ReadyFile(path.Ext(file.Filename))

	//生成目录
	if err := utils.Mkdir(dir); err != nil {
		panic("[fojmd]" + err.Error())
	}

	//保存文件
	if err = ctx.SaveUploadedFile(file, dir+f); err != nil {
		panic("[dias78]" + err.Error())
	}

	common.Success(ctx, gin.H{
		"url": dir + f,
	})

}
