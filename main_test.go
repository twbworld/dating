package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"mime/multipart"

	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/twbworld/dating/global"
	initGlobal "github.com/twbworld/dating/initialize/global"
	"github.com/twbworld/dating/initialize/system"
	"github.com/twbworld/dating/model/common"
	"github.com/twbworld/dating/model/db"
	"github.com/twbworld/dating/router"
	"github.com/twbworld/dating/service"
)

func TestMain(t *testing.T) {
	gin.SetMode(gin.TestMode)
	initGlobal.New("config.example.yaml").Start()
	if err := system.DbStart(); err != nil {
		t.Fatal("数据库连接失败[fsj09]", err)
	}
	defer func() {
		time.Sleep(time.Second * 1) //给足够时间处理数据
		system.DbClose()
	}()

	ginServer := gin.Default()
	router.Start(ginServer)

	n := time.Now()
	ni := n.In(global.Tz)

	token, err := getToken()
	if err != nil {
		t.Fatal("jwt错误[isodfji]", err)
	}

	uploadBody, contentType, err := createForm("test.png", nil)
	if err != nil {
		t.Fatal("[isodadfji]", err)
	}
	uploadBodyFail, contentTypeFail, err := createForm("test.ico", nil)
	if err != nil {
		t.Fatal("[isoerji]", err)
	}
	// common.UserInfoPost
	uploadBody3, contentType3, err := createForm("test2.png", map[string]string{"code": "aaaaaa", "nick_name": "测试"})
	if err != nil {
		t.Fatal("[isodfsdfji]", err)
	}
	uploadBody4, contentType4, err := createForm("test3.png", map[string]string{"code": "aaaaaa", "nick_name": "更名后"})
	if err != nil {
		t.Fatal("[ishdfji]", err)
	}
	uploadBody5, contentType5, err := createForm("test4.png", map[string]string{"desc": "测试带图反馈"})
	if err != nil {
		t.Fatal("[is89dfji]", err)
	}

	//以下是有执行顺序的, 并且库提前有必要数据
	testCases := [...]struct {
		method      string
		postRes     common.Response
		getRes      string
		url         string
		status      int
		postData    interface{}
		contentType string
	}{
		{url: "/upload", postData: uploadBody, contentType: contentType},
		{url: "/upload", postData: uploadBodyFail, contentType: contentTypeFail, postRes: common.Response{Code: 1}},
		{url: "/login", postData: common.LoginPost{Code: ""}, postRes: common.Response{Code: 1}},
		{url: "/login", postData: common.LoginPost{Code: "aaaaaa"}},
		//插入后, userId=3
		{url: "/userAdd", postData: uploadBody3, contentType: contentType3},
		//插入已存在用户(不会新增用户)
		{url: "/userAdd", postData: uploadBody4, contentType: contentType4},
		//加入会面(库已存在id为1,创建者为3的会面)
		{url: "/joinDating", postData: common.DatingPost{
			Id: 1,
			Info: common.InfoPost{Time: [][2]string{
				{ni.Format(time.DateOnly) + ` 9:00:00`, ni.Format(time.DateOnly) + ` 23:00:00`},
				{time.Now().In(global.Tz).AddDate(0, 0, 5).Format(time.DateOnly) + ` 11:00:00`, time.Now().In(global.Tz).AddDate(0, 0, 5).Format(time.DateOnly) + ` 14:00:00`},
			}},
		}},
		//会面重复加入, 失败
		{url: "/joinDating", postData: common.DatingPost{
			Id: 1,
			Info: common.InfoPost{Time: [][2]string{
				{ni.Format(time.DateOnly) + ` 9:00:00`, ni.Format(time.DateOnly) + ` 23:00:00`},
				{time.Now().In(global.Tz).AddDate(0, 0, 5).Format(time.DateOnly) + ` 11:00:00`, time.Now().In(global.Tz).AddDate(0, 0, 5).Format(time.DateOnly) + ` 14:00:00`},
			}},
		}, postRes: common.Response{Code: 1}},
		//创建会面
		{url: "/joinDating", postData: common.DatingPost{
			Info: common.InfoPost{Time: [][2]string{
				{ni.Format(time.DateOnly) + ` 10:00:00`, ni.Format(time.DateOnly) + ` 22:00:00`},
				{time.Now().In(global.Tz).AddDate(0, 0, 10).Format(time.DateOnly) + ` 11:00:00`, time.Now().In(global.Tz).AddDate(0, 0, 10).Format(time.DateOnly) + ` 14:00:00`},
			}},
		}},
		//加入会面(手动虚拟)
		{url: "/joinDating", postData: common.DatingPost{
			Id: 2,
			Info: common.InfoPost{Time: [][2]string{
				{ni.Format(time.DateOnly) + ` 10:00:00`, ni.Format(time.DateOnly) + ` 20:00:00`},
				{time.Now().In(global.Tz).AddDate(0, 0, 20).Format(time.DateOnly) + ` 13:00:00`, time.Now().In(global.Tz).AddDate(0, 0, 20).Format(time.DateOnly) + ` 14:00:00`},
			}},
		}},
		{url: "/getDatingList", postData: common.GetDatingListPost{Page: 1, LastId: 0}},
		{url: "/getDating", postData: common.GetDatingPost{Id: 1}},
		{url: "/getDating", postData: common.GetDatingPost{Id: 2}},
		{url: "/getDatingAmount", postData: ""},
		//用UtId只能退出虚拟用户
		//退出不属于自己的虚拟用户(失败)
		{url: "/quitDating", postData: common.QuitDatingPost{UtId: 1}, postRes: common.Response{Code: 1}},
		//退出不存在会面(失败)
		{url: "/quitDating", postData: common.QuitDatingPost{Id: 3}, postRes: common.Response{Code: 1}},
		//退出会面
		{url: "/quitDating", postData: common.QuitDatingPost{Id: 1}},
		//退出虚拟用户
		{url: "/quitDating", postData: common.QuitDatingPost{UtId: 7}},
		//关闭会面
		{url: "/quitDating", postData: common.QuitDatingPost{Id: 2}},
		//反馈
		{url: "/feedback", postData: common.FeedbackPost{Desc: ""}, postRes: common.Response{Code: 1}},
		{url: "/feedback", postData: common.FeedbackPost{Desc: "测试纯文本反馈"}},
		//带图反馈
		{url: "/feedback", postData: uploadBody5, contentType: contentType5},
	}

	//非web请求的测试========================begin

	t.Run("Match", func(t *testing.T) {
		if res, err := service.Service.UserServiceGroup.DatingService.Match(1); err != nil || len(res.Date) < 1 {
			t.Fatal("[erscai2]", err)
		}
	})
	t.Run("MatchFatal", func(t *testing.T) {
		//不存在会面
		if res, err := service.Service.UserServiceGroup.DatingService.Match(2); err == nil || !strings.Contains(err.Error(), "iudha09") {
			t.Fatal("[erk0ai2]", res, err)
		}
	})
	//非web请求的测试========================end

	for k, value := range testCases {
		t.Run(strconv.FormatInt(int64(k+1), 10)+value.url, func(t *testing.T) {
			if value.method == "" {
				value.method = http.MethodPost
			}
			if value.status == 0 {
				value.status = 200
			}
			if value.method == http.MethodPost {
				if value.contentType == "" {
					value.contentType = "application/json"
				}
				if value.postRes == (common.Response{}) {
					value.postRes.Code = 0
				}
			}

			requestBody := new(bytes.Buffer)
			if value.postData != nil {
				if v, ok := value.postData.(*bytes.Buffer); ok {
					requestBody = v
				} else {
					jsonVal, err := json.Marshal(value.postData)
					if err != nil {
						t.Fatal("json出错[godjg]", err)
					}
					requestBody = bytes.NewBuffer(jsonVal)
				}
			}

			b := time.Now().UnixMilli()

			//向注册的路有发起请求
			req, err := http.NewRequest(value.method, value.url, requestBody)
			if err != nil {
				t.Fatal("请求出错[godkojg]", err)
			}
			req.Header.Set("Authorization", token)
			if value.method == http.MethodPost {
				req.Header.Set("content-type", value.contentType)
			}

			res := httptest.NewRecorder() // 构造一个记录
			ginServer.ServeHTTP(res, req) //模拟http服务处理请求

			result := res.Result() //response响应

			fmt.Printf("^^^^^^处理用时%d毫秒^^^^^^\n", time.Now().UnixMilli()-b)

			assert.Equal(t, value.status, result.StatusCode)

			body, err := io.ReadAll(result.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer result.Body.Close()

			switch value.method {
			case http.MethodPost:
				var response common.Response
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatal("返回错误", err, string(body))
				}
				assert.Equal(t, value.postRes.Code, response.Code)
			case http.MethodGet:
				assert.Contains(t, string(body), value.getRes)
			}

			// fmt.Println("request!!!!!!!!!!", string(jsonVal))
			fmt.Println("response!!!!!!!!!!", string(body))

			time.Sleep(time.Millisecond * 500) //!!!!!!!!!!!!!!!!!!

		})

	}

}

func getToken() (string, error) {
	u := db.User{}
	u.Id = 2 //!!!!!!!!!!!!!!!!测试用户
	return service.Service.UserServiceGroup.BaseService.LoginToken(&u)
}

func createForm(fileName string, data map[string]string) (postData *bytes.Buffer, contentType string, err error) {

	uploadBody := &bytes.Buffer{}
	writer := multipart.NewWriter(uploadBody)
	//构造form-file的字段和文件名数据;
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return
	}

	//往文件写入
	if _, err = io.Copy(part, strings.NewReader("hi")); err != nil {
		return
	}

	// 写入表单字段
	for k, v := range data {
		if err = writer.WriteField(k, v); err != nil {
			return
		}
	}

	// 关闭 writer
	if err = writer.Close(); err != nil {
		return nil, "", err
	}

	return uploadBody, writer.FormDataContentType(), nil
}
