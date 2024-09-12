**Dating**
===========
[![](https://github.com/twbworld/dating/workflows/ci/badge.svg?branch=main)](https://github.com/twbworld/dating/actions)
[![](https://img.shields.io/github/tag/twbworld/dating?logo=github)](https://github.com/twbworld/dating)
![](https://img.shields.io/badge/language-golang-cyan)
[![](https://img.shields.io/github/license/twbworld/dating)](https://github.com/twbworld/dating/blob/main/LICENSE)

## 简介
**小程序 `我们何时见` 源码**

## 目录结构 : 
``` sh
├── client/             #前端源码, uni-app框架
├── config.example.yaml #以其为例, 自行创建config.yaml
├── controller/         #MVC模式的C
│   ├── admin/          #后台api
│   ├── enter.go
│   └── user/           #前台api
├── dao/            #sql
├── Dockerfile      #存放GitHub-Actions的工作流文件
├── .editorconfig
├── .gitattributes
├── .github/
│   └── workflows/ #存放GitHub-Actions的工作流文件
├── .gitignore
├── global/
│   └── global.go   #全局变量的初始化
├── go.mod
├── go.sum
├── initialize/     #服务初始化相关
│   ├── global/
│   ├── server.go   #gin服务
│   └── system/
├── LICENSE
├── log/
│   ├── gin.log     #gin日志
│   ├── .gitkeep
│   └── run.log     #业务日志
├── main.go         #入口
├── main_test.go    #测试
├── middleware/     #路由中间件以及验参
├── model/
│   ├── common/     #业务要用的结构体
│   ├── config/     #配置文件的结构体
│   └── db/         #数据库模型结构体
├── README.md
├── router/         #gin路由
├── service/
│   ├── admin/
│   ├── enter.go
│   └── user/
├── static/         #静态资源
├── task/
│   └── dating.go   #定时任务
└── utils/
    └── tool.go
```

## 后端运行

### docker-compose
``` yaml
version: "3"
services:
    dating:
        image: ghcr.io/twbworld/dating:latest
        ports:
            - 80:80
        volumes:
            - ${PWD}/config.yaml:config.yaml:ro
            - ${PWD}/dating.db:dating.db:rw
```

### 打包本地运行
```sh
$ cp config.example.yaml config.yaml

$ go mod download && go mod tidy

$ go build -o server main.go

$ ./server
或
$ ./server -c config.yaml
```

## 前端运行

1. 前端源码在`client`目录下, 或直达项目[传送门](https://github.com/twbworld/dating-client/)
    ```sh
    $ git submodule update --remote
    $ cd client/
    $ npm install
    ```
2. 下载"微信开发者工具"并登陆
3. 本前端项目用`uni-app`框架, 下载专用的`HBuilderX`编辑器
4. 配置`manifest.json`
5. 运行: 菜单-运行-运行到小程序模拟器-微信开发者工具

## 展示

![](https://cdn.jsdelivr.net/gh/twbworld/hosting@main/img/202409111424755.png)
![](https://cdn.jsdelivr.net/gh/twbworld/hosting@main/img/202409111425133.png)
![](https://cdn.jsdelivr.net/gh/twbworld/hosting@main/img/202409111425284.png)
![](https://cdn.jsdelivr.net/gh/twbworld/hosting@main/img/202409111425518.png)

## 计划中:

1.  广告位
2.  头像屏蔽层
3.  已加入列表左滑
4.  显示匹配失败的用户(屏蔽层)
5.  会面详情添加日历显示
6.  根据天气状况推荐
7.  长连接实时更新会面详情
8.  预设会面天数
9.  修改用户头像等信息
10. 后台
11. 无登录使用
