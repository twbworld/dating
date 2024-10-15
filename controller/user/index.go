package user

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/twbworld/dating/dao"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/service"
	"github.com/twbworld/dating/utils"
)

type DatingApi struct {
}

// 事务(用于执行service.Match())
func (d *DatingApi) tx(datingId uint, fc func(tx *sqlx.Tx) error) error {

	if err := dao.Tx(fc); err != nil {
		return err
	}

	d.match(datingId)

	return nil

}

func (d *DatingApi) match(datingId uint) {
	service.Service.UserServiceGroup.DatingService.MatchGoroutine(datingId, func() {
		if _, ok := cliData.list.Load(datingId); !ok {
			return
		}
		if da, dr, du, err := service.Service.UserServiceGroup.DatingService.GetDating(datingId); err == nil {
			//这里懒得判断用户是否存在于会面了
			if clis, ok := cliData.list.Load(datingId); ok {
				//广播
				for cl := range clis.(map[*client]bool) {
					common.SuccessWs(cl.send, &common.DatingInfo{
						Dating: common.DatingSimple{
							CreateUserId: da.CreateUserId,
							Id:           da.Id,
							Status:       da.Status,
							Result:       *dr,
						},
						Users: du,
					})
				}
			}
		}
	})
}

// 创建会面 || 加入会面 || 手动加入会面
func (d *DatingApi) JoinDating(ctx *gin.Context) {
	var (
		err      error
		data     common.DatingPost
		dating   db.Dating
		userTime db.UserTime
		utId     uint
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln("[fijn0k]", p)
			common.Fail(ctx, `系统错误[jgfdd]`)
		}
	}()

	userId := ctx.MustGet(`userId`).(uint)
	if userId < 1 {
		common.Fail(ctx, `系统错误[thojpi]`)
		return
	}

	if ctx.ShouldBindJSON(&data) != nil {
		common.Fail(ctx, `参数错误[ddsgsj]`)
		return
	}

	if err = service.Service.UserServiceGroup.Validator.ValidatorInfo(&data.Info); err != nil {
		common.Fail(ctx, err.Error())
		return
	}

	infoJson, err := data.Info.Marshal()
	if err != nil {
		panic("[gjdas]" + err.Error())
	}

	if userId == dao.BaseUserId {
		common.Fail(ctx, `参数错误[eoifjs]`)
		return
	}

	joinUserId := userId
	if data.Id > 0 {
		if err = dao.App.DatingDb.GetDating(&dating, data.Id); err != nil {
			if err == sql.ErrNoRows {
				common.Fail(ctx, `会面不存在[jokrfjs]`)
				return
			}
			panic("[8nkuiuh]" + err.Error())
		}
		if dating.Status == 0 {
			common.Fail(ctx, `会面已结束`)
			return
		}
		if err = dao.App.DatingDb.GetUserTime(&userTime, dating.Id, userId); err != sql.ErrNoRows {
			if err == nil && userId == dating.CreateUserId {
				//"手动添加会面"
				joinUserId = dao.BaseUserId
			} else if err == nil {
				common.Fail(ctx, `您已加入会面`)
				return

			} else {
				panic("[fijsuh]" + err.Error())
			}
		}
	}

	err = dao.Tx(func(tx *sqlx.Tx) (e error) {
		if data.Id < 1 {
			//创建会面
			if id, e := dao.App.DatingDb.AddDating(userId, tx); e != nil {
				panic("[hpisd]" + e.Error())
			} else if id < 1 {
				panic("系统错误[gdsdioj]")
			} else {
				data.Id = id
			}
			if e = dao.App.DatingDb.GetDating(&dating, data.Id, tx); e != nil {
				panic("[fwije]" + e.Error())
			}
		}

		//加入会面
		utId, e = dao.App.DatingDb.JoinDating(&db.UserTime{
			DatingId: dating.Id,
			UserId:   joinUserId,
			Info:     infoJson,
		}, tx)

		return
	})

	if err != nil {
		panic("[fioasj]" + err.Error())
	}

	d.match(dating.Id)

	common.Success(ctx, gin.H{
		"id":    dating.Id,
		"ut_id": utId,
	})

}

// 更新会面时间
func (d *DatingApi) UpdateUserTime(ctx *gin.Context) {
	var (
		err            error
		data           common.UerTimePost
		datingUserJson common.DatingUserJoin
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `系统错误[iodjaso]`)
		}
	}()

	userId := ctx.MustGet(`userId`).(uint)
	if userId < 1 {
		common.Fail(ctx, `系统错误[thljkjpi]`)
		return
	}

	if ctx.ShouldBindJSON(&data) != nil {
		common.Fail(ctx, `参数错误[ddsgj]`)
		return
	}

	if err = service.Service.UserServiceGroup.Validator.ValidatorInfo(&data.Info); err != nil {
		common.Fail(ctx, err.Error())
		return
	}
	if data.UtId < 1 {
		common.Fail(ctx, `参数错误[doisgj]`)
		return
	}
	if userId == dao.BaseUserId {
		common.Fail(ctx, `参数错误[podjig]`)
		return
	}

	if dao.App.DatingDb.CheckUserTime(&datingUserJson, data.UtId) != nil {
		common.Fail(ctx, `参数错误[kkpiaojh]`)
		return
	}

	if datingUserJson.UserId != userId && (datingUserJson.UserId != dao.BaseUserId || datingUserJson.CreateUserId != userId) {
		//非会面创建者,不允许修改user_id=1数据
		common.Fail(ctx, `权限不足[dssesd]`)
		return
	}

	infoJson, err := data.Info.Marshal()
	if err != nil {
		panic(err)
	}

	err = d.tx(datingUserJson.DatingId, func(tx *sqlx.Tx) (e error) {
		return dao.App.DatingDb.JoinUpdate(data.UtId, &db.UserTime{
			Info: infoJson,
		}, tx)
	})

	if err != nil {
		panic(err)
	}

	common.Success(ctx, gin.H{
		"id":    datingUserJson.DatingId,
		"ut_id": data.UtId,
	})

}

// 获取当前用户参与的所有会面id
func (d *DatingApi) GetDatingAmount(ctx *gin.Context) {
	var (
		datingIds []uint = make([]uint, 0)
		err       error
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `系统错误[o23j4]`)
		}
	}()

	userId := ctx.MustGet(`userId`).(uint)
	if userId < 1 {
		common.Fail(ctx, `系统错误[thjytpi]`)
		return
	}

	if userId == dao.BaseUserId {
		common.Fail(ctx, `参数错误[6k57j]`)
		return
	}

	if err = dao.App.DatingDb.GetDatingAmount(&datingIds, userId); err != nil {
		panic("[ii8iujh]" + err.Error())
	}

	common.Success(ctx, gin.H{
		"ids": datingIds,
	})

}

// 获取当前用户参与的所有会面信息
func (d *DatingApi) GetDatingList(ctx *gin.Context) {
	var (
		data                common.GetDatingListPost
		limit               uint                = 5
		datingList          []common.DatingList = make([]common.DatingList, 0, int(limit))
		datingUserAvatar    []common.DatingUserAvatar
		err                 error
		datingKeyUserAvatar map[uint][]string
		datingIds           []uint
		keys                map[uint]bool
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `系统错误[fsiuh]`)
		}
	}()

	userId := ctx.MustGet(`userId`).(uint)
	if userId < 1 {
		common.Fail(ctx, `系统错误[thjteri]`)
		return
	}

	if ctx.ShouldBindJSON(&data) != nil {
		common.Fail(ctx, `参数错误[eriuth]`)
		return
	}
	if data.Page < 1 {
		common.Fail(ctx, `参数错误[456hbj]`)
		return
	}

	if userId == dao.BaseUserId {
		common.Fail(ctx, `参数错误[2ijn]`)
		return
	}

	if err = dao.App.DatingDb.GetDatingList(&datingList, userId, data.LastId, limit); err != nil {
		panic(err)
	}
	if len(datingList) < 1 {
		goto RESPONSE
	}

	datingIds = make([]uint, 0, len(datingList))
	keys = make(map[uint]bool, len(datingList))
	for _, value := range datingList {
		if _, ok := keys[value.Id]; ok {
			continue
		}
		datingIds = append(datingIds, value.Id)
		keys[value.Id] = true
	}

	if err = dao.App.DatingDb.GetUserByDatingIds(&datingUserAvatar, datingIds); err != nil {
		panic(err)
	}
	if len(datingUserAvatar) < 1 {
		goto RESPONSE
	}

	datingKeyUserAvatar = make(map[uint][]string, len(datingIds))
	for _, value := range datingUserAvatar {
		if _, ok := datingKeyUserAvatar[value.DatingId]; !ok {
			datingKeyUserAvatar[value.DatingId] = make([]string, 0)
		}
		datingKeyUserAvatar[value.DatingId] = append(datingKeyUserAvatar[value.DatingId], value.Path)
	}

	for key, val := range datingList {
		if _, ok := datingKeyUserAvatar[val.Id]; !ok {
			datingKeyUserAvatar[val.Id] = make([]string, 0)
		}
		datingList[key].AvatarUrl = datingKeyUserAvatar[val.Id]
		datingList[key].AddTimeStr = utils.TimeFormat(datingList[key].AddTime)
	}

RESPONSE:

	common.Success(ctx, gin.H{
		"list": datingList,
	})

}

// 退出/关闭 会面
func (d *DatingApi) QuitDating(ctx *gin.Context) {
	var (
		data     common.QuitDatingPost
		dating   db.Dating
		userTime db.UserTime
		isUtId   bool
		err      error
	)

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `系统错误[342ijn]`)
		}
	}()

	userId := ctx.MustGet(`userId`).(uint)
	if userId < 1 {
		common.Fail(ctx, `系统错误[thvbcri]`)
		return
	}

	if userId == dao.BaseUserId {
		common.Fail(ctx, `参数错误[43hbj]`)
		return
	}

	if ctx.ShouldBindJSON(&data) != nil {
		common.Fail(ctx, `参数错误[fd89vu]`)
		return
	}

	if data.Id > 0 {
		if dao.App.DatingDb.GetUserTime(&userTime, data.Id, userId) != nil {
			common.Fail(ctx, `找不到记录[nfg87y]`)
			return
		}
	} else if data.UtId > 0 {
		if dao.App.DatingDb.GetUserTimeById(&userTime, data.UtId) != nil {
			common.Fail(ctx, `找不到记录[7g6yasdf]`)
			return
		}
		data.Id, isUtId = userTime.DatingId, true
	} else {
		common.Fail(ctx, `参数错误[3rhbu]`)
		return
	}

	if dao.App.DatingDb.GetDating(&dating, data.Id) != nil {
		common.Fail(ctx, `参数错误[a9s8d0]`)
		return
	}

	if dating.Status == 0 {
		common.Fail(ctx, `会面已结束[a9s8d0]`)
		return
	}

	if isUtId && (dating.CreateUserId != userId || userTime.UserId != dao.BaseUserId) {
		//使用ut_id删除的, 必须为创建者同时被删是虚拟用户
		common.Fail(ctx, `参数错误[djioa]`)
		return
	}

	err = d.tx(dating.Id, func(tx *sqlx.Tx) (e error) {
		if !isUtId && dating.CreateUserId == userId {
			//关闭(创建者主动 退出或关闭 会面)
			return dao.App.DatingDb.CloseDating([]uint{dating.Id}, tx)
		}
		//退出
		return dao.App.DatingDb.QuitDating(userTime.Id, tx)
	})
	if err != nil {
		panic(err)
	}

	common.Success(ctx, data)
}
