package user

import (
	"database/sql"
	"mime/multipart"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/service"
	"github.com/twbworld/dating/utils"
)

type UserApi struct{}

// 用户注册(用户存在, 则修改)
func (b *UserApi) UserAdd(ctx *gin.Context) {
	var (
		data  common.UserInfoPost
		u     db.User
		token string
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `出错, 请重新授权[oinds]`)
		}
	}()

	data.Code = ctx.DefaultPostForm("code", "")
	data.NickName = strings.TrimSpace(ctx.DefaultPostForm("nick_name", "微信用户"))

	if err := service.Service.UserServiceGroup.Validator.ValidatorUserAddPost(&data); err != nil {
		common.Fail(ctx, err.Error())
		return
	}

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
	fp := dir + f

	//判断新图片是否已存在数据库
	id, err := dao.App.FileDb.GetFileId(fp)
	if id > 0 {
		common.Fail(ctx, `出错, 请重新授权[ddokssj]`)
		return
	} else if err != nil && err != sql.ErrNoRows {
		panic("系统出错[tweiojlkmsdf]")
	}

	// 用Code向微信换取用户信息
	if err := service.Service.UserServiceGroup.DatingService.GetUserByCode(data.Code, &u); err != nil {
		panic(err)
	}

	if data.NickName == "" {
		data.NickName = "微信用户"
	} else if u.Id == 0 {
		if utils.WxCheckContent(u.OpenId, data.NickName) != nil {
			common.Fail(ctx, `昵称不合法, 申请重新授权`)
			return
		}
	}

	if u.Id > 0 {
		//已存在用户

		ur, f := common.UserInfo{User: u}, db.File{}
		if u.Avatar > 0 && dao.App.FileDb.GetFile(&f, u.Avatar) == nil {
			ur.AvatarUrl = f.Path
		}

		if u.NickName != data.NickName || f.Path != fp {
			//修改用户信息
			err := dao.Tx(func(tx *sqlx.Tx) (e error) {

				if f.Path != fp {
					u.Avatar, e = dao.App.FileDb.AddFile(&db.File{
						Path: fp,
						Ext:  strings.TrimLeft(filepath.Ext(fp), "."),
					}, tx)
					if e != nil {
						panic("[gdf7e6f]" + e.Error())
					}
					if u.Avatar < 1 {
						panic("[shihd78]系统错误")
					}

				}
				u.NickName = data.NickName

				e = dao.App.UserDb.UpdateUser(&u, tx)
				if e != nil {
					panic("[gdf7e6f]" + e.Error())
				}

				//生成目录
				if e = utils.Mkdir(dir); e != nil {
					panic("[fojmfd]" + e.Error())
				}
				//保存头像
				if e = ctx.SaveUploadedFile(file, fp); e != nil {
					panic("[di6s78]" + e.Error())
				}

				return
			})

			if err != nil {
				global.Log.Errorln(err)
			} else {
				//使用修改后的数据
				ur.User = u
				ur.AvatarUrl = fp
			}
		}

		common.Success(ctx, gin.H{"user": ur})
		return
	}

	//用户注册
	err = dao.Tx(func(tx *sqlx.Tx) (e error) {
		u.Avatar, e = dao.App.FileDb.AddFile(&db.File{
			Path: fp,
			Ext:  strings.TrimLeft(filepath.Ext(fp), "."),
		}, tx)
		if e != nil {
			panic("[gdf76f]" + e.Error())
		}
		if u.Avatar < 1 {
			panic("[shid78]系统错误")
		}

		u.Id = 0
		u.NickName = data.NickName

		id, e := dao.App.UserDb.AddUser(&u, tx)
		if e != nil {
			panic("[gdfjinf]" + e.Error())
		}

		if dao.App.UserDb.GetUserById(&u, id, tx) == sql.ErrNoRows {
			panic(`系统错误[o8s6nm]`)
		}

		//生成目录
		if e = utils.Mkdir(dir); e != nil {
			panic("[fojmfd]" + e.Error())
		}
		//保存头像
		if e = ctx.SaveUploadedFile(file, fp); e != nil {
			panic("[di6s78]" + e.Error())
		}

		token, e = service.Service.UserServiceGroup.BaseService.LoginToken(&u)

		return
	})
	if err != nil {
		panic(err)
	}

	common.SuccessAuth(ctx, token, gin.H{
		"user": common.UserInfo{
			User:      u,
			AvatarUrl: fp,
		},
	})
}

// 用户反馈
func (b *UserApi) Feedback(ctx *gin.Context) {

	var (
		data     common.FeedbackPost
		imgExist bool
		fp       string
		dir      string
		file     *multipart.FileHeader
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `出错, 请重新授权[osdds]`)
		}
	}()

	userId := ctx.MustGet(`userId`).(uint)
	if userId < 1 {
		common.Fail(ctx, `系统错误[th9pi]`)
		return
	}

	if userId == dao.BaseUserId {
		common.Fail(ctx, `参数错误[6f7j]`)
		return
	}

	if ctx.ContentType() == "application/json" {
		//没有图片的反馈
		ctx.ShouldBindJSON(&data)
	} else {
		// 带图反馈
		data.Desc = strings.TrimSpace(ctx.DefaultPostForm("desc", ""))
		imgExist = true
	}

	if err := service.Service.UserServiceGroup.Validator.ValidatorFeedbackPost(&data); err != nil {
		common.Fail(ctx, err.Error())
		return
	}

	if imgExist {
		var (
			err error
			f   string
		)
		file, err = ctx.FormFile("file")
		if err != nil {
			common.Fail(ctx, `参数错误[ipjdmk]`)
			return
		}

		if err = service.Service.UserServiceGroup.Validator.ValidatorUpload(file); err != nil {
			common.Fail(ctx, err.Error())
			return
		}
		dir, f = utils.ReadyFile(path.Ext(file.Filename))
		fp = dir + f
		//判断新图片是否已存在数据库
		id, err := dao.App.FileDb.GetFileId(fp)
		if id > 0 {
			common.Fail(ctx, `出错, 请重新上传[ddtyhkssj]`)
			return
		} else if err != nil && err != sql.ErrNoRows {
			panic("系统出错[tweiosdf]")
		}

	}

	err := dao.Tx(func(tx *sqlx.Tx) (e error) {
		var fileId string
		if imgExist {
			fid, e := dao.App.FileDb.AddFile(&db.File{
				Path: fp,
				Ext:  strings.TrimLeft(filepath.Ext(fp), "."),
			}, tx)
			if e != nil {
				panic("[gdf76f]" + e.Error())
			}
			fileId = strconv.Itoa(int(fid))
		}

		id, e := dao.App.FeedbackDb.AddFeedback(&db.Feedback{
			Desc:   data.Desc,
			UserId: userId,
			FileId: fileId,
		}, tx)
		if id < 1 {
			panic(`系统错误[ond0sm]`)
		}

		if imgExist {
			//生成目录
			if e = utils.Mkdir(dir); e != nil {
				panic("[fjghfd]" + e.Error())
			}
			//保存头像
			if e = ctx.SaveUploadedFile(file, fp); e != nil {
				panic("[ddgf8]" + e.Error() + fp)
			}

		}

		return
	})
	if err != nil {
		panic(err)
	}

	var str strings.Builder
	str.WriteString("用户反馈通知:\n")
	str.WriteString(data.Desc)
	if imgExist {
		str.WriteString("\n")
		str.WriteString(global.Config.Domain)
		str.WriteString(`/`)
		str.WriteString(fp)
	}

	go service.Service.UserServiceGroup.TgService.TgSend(str.String())

	common.SuccessOk(ctx, `成功`)
}
