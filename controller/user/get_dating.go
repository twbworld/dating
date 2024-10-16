package user

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/middleware"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/service"
)

type client struct {
	uId  uint
	conn *websocket.Conn
	send chan *common.Response
	once sync.Once
}

// websocket在线用户数据
type clients struct {
	list sync.Map //map[uint]map[*client]bool (多个Goroutine读写共享内存, 要上锁)
}

var cliData = &clients{}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *client) updateClientList(datingId uint, add bool) {
	if clis, ok := cliData.list.Load(datingId); ok {
		clis := clis.(map[*client]bool)
		if add {
			clis[c] = true
			cliData.list.Store(datingId, clis)
		} else {
			delete(clis, c)
			if len(clis) == 0 {
				cliData.list.Delete(datingId)
			} else {
				cliData.list.Store(datingId, clis)
			}
		}
	} else if add {
		cliData.list.Store(datingId, map[*client]bool{c: true})
	}
}

// websocket断开
func (c *client) close() {
	c.once.Do(func() {
		//避免重复close
		close(c.send)
		c.conn.Close()
		// fmt.Println("已下线")
	})
}

// websocket发送数据
func (c *client) writePump() {
	defer c.close()

	for res := range c.send {
		c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if err := c.conn.WriteJSON(res); err != nil {
			return
		}
	}
}

// websocket读取数据
func (c *client) readPump() {
	var data common.GetDatingPost

	defer func() {
		c.close()
		if data.Id != 0 {
			c.updateClientList(data.Id, false)
		}
	}()

	for {
		c.conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			//下线
			break
		}
		if len(msg) == 0 {
			common.FailWs(c.send, `参数错误[nv2gnb8]`)
			continue
		}

		//测试代码
		// fmt.Println("收到=", string(msg), "=")
		// if c.uId == 0 && string(msg) == "ok\n" {
		// 	c.uId = 4
		// 	common.SuccessAuthWs(c.send, "")
		// 	continue
		// }

		if c.uId == 0 {
			if userId, newToken, err := middleware.JWTAuth(string(msg)); err != nil {
				common.FailAuthWs(c.send, err.Error())
			} else if userId == 0 {
				common.FailAuthWs(c.send, `参数错误[nvb8]`)
			} else {
				c.uId = userId
				common.SuccessAuthWs(c.send, newToken)
			}
			continue
		}

		if json.Unmarshal(msg, &data) != nil {
			common.FailWs(c.send, `参数错误[n789dab8]`)
			continue
		}

		da, dr, du, err := service.Service.UserServiceGroup.DatingService.GetDating(data.Id)
		if err != nil {
			common.FailWs(c.send, err.Error())
			continue
		}

		isset := false
		for _, value := range du {
			if value.Id == c.uId {
				isset = true
				break
			}
		}
		if !isset {
			//用户不存在于会面
			common.SuccessWs(c.send, &common.DatingInfo{
				Dating: common.DatingSimple{Status: da.Status},
				Users:  []common.DatingUser{},
			})
			continue
		}

		c.updateClientList(data.Id, true)
		common.SuccessWs(c.send, &common.DatingInfo{
			Dating: common.DatingSimple{
				CreateUserId: da.CreateUserId,
				Id:           da.Id,
				Status:       da.Status,
				Result:       dr,
			},
			Users: du,
		})
	}
}

// 获取会面详情(websocket)
func (d *DatingApi) GetDatingWs(ctx *gin.Context) {
	if !ctx.IsWebsocket() {
		return
	}
	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}

	clien := &client{
		conn: conn,
		send: make(chan *common.Response),
	}

	go clien.writePump()
	go clien.readPump()

}

// 获取会面详情
func (d *DatingApi) GetDating(ctx *gin.Context) {
	var data common.GetDatingPost

	defer func() {
		if p := recover(); p != nil {
			global.Log.Errorln(p)
			common.Fail(ctx, `系统错误[lksdfj]`)
		}
	}()

	if ctx.ShouldBindJSON(&data) != nil {
		common.Fail(ctx, `参数错误[j7n65]`)
		return
	}

	userId := ctx.MustGet(`userId`).(uint)
	if userId == 0 {
		common.Fail(ctx, `系统错误[thojpi]`)
		return
	}

	da, dr, du, err := service.Service.UserServiceGroup.DatingService.GetDating(data.Id)
	if err != nil {
		common.Fail(ctx, err.Error())
		return
	}

	isset := false
	for _, value := range du {
		if value.Id == userId {
			isset = true
			break
		}
	}
	if !isset {
		//用户不存在于会面 或 已退出
		common.Success(ctx, &common.DatingInfo{
			Dating: common.DatingSimple{Status: da.Status},
			Users:  []common.DatingUser{},
		})
		return
	}

	common.Success(ctx, &common.DatingInfo{
		Dating: common.DatingSimple{
			CreateUserId: da.CreateUserId,
			Id:           da.Id,
			Status:       da.Status,
			Result:       dr,
		},
		Users: du,
	})
}
