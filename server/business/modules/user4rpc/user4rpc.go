// Copyright (C) 2020 WuPeng <wup364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 用户管理模块

package user4rpc

import (
	"fileservice/business/constants"
	"fileservice/business/service"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/fileutil"
	"github.com/wup364/pakku/utils/logs"
)

// User4RPC 用户管理模块
type User4RPC struct {
	Signature
	us *UserStory
	c  ipakku.AppConfig `@autowired:"AppConfig"`
}

// AsModule 模块加载器接口实现, 返回模块信息&配置
func (umg *User4RPC) AsModule() ipakku.Opts {
	return ipakku.Opts{
		Name:        "User4RPC",
		Version:     1.0,
		Description: "用户信息",
		OnReady: func(mctx ipakku.Loader) {
			umg.us = &UserStory{}
			if err := mctx.AutoWired(&umg.Signature); nil != err {
				logs.Panicln(err)
			} else {
				// 注册 accesstoken 缓存库, x分钟过期
				if err := umg.ch.RegLib(service.Cachelib_UserAccessToken, 60*60*24); nil != err {
					logs.Panicln(err)
				}
			}
			deftDataSource := "./.datas/" + mctx.GetParam(ipakku.PARAMKEY_APPNAME).ToString("app") + ".db?cache=shared"
			confDataSource := umg.c.GetConfig("authuser.sotre.datasource").ToString(deftDataSource)
			if confDataSource == deftDataSource {
				umg.mkSqliteDIR() // 创建sqlite文件存放目录
			}
			if err := umg.us.Initial(constants.DBSetting{
				DriverName:     umg.c.GetConfig("authuser.sotre.driver").ToString("sqlite3"),
				DataSourceName: confDataSource,
			}); nil != err {
				logs.Panicln(err)
			}
		},
		OnSetup: func() {
			if err := umg.us.Install(); nil != err {
				logs.Panicln(err)
			}
		},
	}
}

// ListAllUsers 列出所有用户数据, 无分页
func (umg *User4RPC) ListAllUsers() ([]service.UserInfoDto, error) {
	if users, err := umg.us.ListAllUsers(); nil == err && len(users) > 0 {
		dtos := make([]service.UserInfoDto, len(users))
		for i := 0; i < len(users); i++ {
			dtos[i] = *users[i].ToDto()
		}
		return dtos, nil
	} else {
		return make([]service.UserInfoDto, 0), err
	}
}

// QueryUser 根据用户ID查询详细信息
func (umg *User4RPC) QueryUser(userID string) (*service.UserInfoDto, error) {
	if user, err := umg.us.QueryUser(userID); nil == err && nil != user {
		return user.ToDto(), nil
	} else {
		return nil, err
	}
}

// Clear 清空用户
func (umg *User4RPC) Clear() error {
	if users, err := umg.us.ListAllUsers(); nil == err {
		if len(users) > 0 {
			for _, val := range users {
				if err = umg.us.DelUser(val.UserID); nil != err {
					return err
				}
			}
		}
	} else {
		return err
	}
	return nil
}

// AddUser 添加用户
func (umg *User4RPC) AddUser(user *service.UserInfo) error {
	return umg.us.AddUser(user)
}

// UpdatePWD 修改用户密码
func (umg *User4RPC) UpdatePWD(userID, pwd string) error {
	userOld, err := umg.us.QueryUser(userID)
	if nil != err {
		return err
	}
	if nil == userOld {
		return service.ErrorUserNotExist
	}
	userOld.UserPWD = pwd
	return umg.us.UpdatePWD(userOld)
}

// UpdateUserName 修改用户昵称
func (umg *User4RPC) UpdateUserName(userID, userName string) error {
	if len(userName) == 0 {
		return service.ErrorUserNameIsNil
	}
	userOld, err := umg.us.QueryUser(userID)
	if nil != err {
		return err
	}
	if nil == userOld {
		return service.ErrorUserNotExist
	}
	userOld.UserName = userName
	return umg.us.UpdateUser(userOld)
}

// DelUser 根据userID删除用户
func (umg *User4RPC) DelUser(userID string) error {
	return umg.us.DelUser(userID)
}

// CheckPwd 校验密码是否一致
func (umg *User4RPC) CheckPwd(userID, pwd string) bool {
	return umg.us.CheckPwd(userID, pwd)
}

// GetAuthFilterFunc 获取过滤器实现
func (umg *User4RPC) GetAuthFilterFunc() ipakku.FilterFunc {
	return umg.RestfulAPIFilter
}

// AskAccess 获取access
func (umg *User4RPC) AskAccess(userID, pwd string) (*service.UserAccessDto, error) {
	if !umg.CheckPwd(userID, pwd) {
		return nil, service.ErrorAuthentication
	}
	if user, err := umg.us.QueryUser(userID); nil != err {
		return nil, err
	} else {
		return umg.Signature.AskAccess(*user)
	}
}

// mkSqliteDIR 创建sqlite文件存放目录
func (umg *User4RPC) mkSqliteDIR() {
	if !fileutil.IsExist("./.datas") {
		if err := fileutil.MkdirAll("./.datas"); nil != err {
			logs.Panicln(err)
		}
	}
}
