package qrobot

import (
	"context"
	"github.com/Ayanokoji1020-miyano/qrobot/consts"
	"github.com/Ayanokoji1020-miyano/qrobot/feature"
	"github.com/Ayanokoji1020-miyano/qrobot/robot"
	"github.com/Ayanokoji1020-miyano/qrobot/service"
	"github.com/Mrs4s/MiraiGo/client"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var qqClient *robot.Client

func login(qq int64, pwd string, protocol int) error {
	qqClient = robot.NewClient(qq, pwd)

	device := new(client.DeviceInfo)
	if err := device.ReadJson([]byte(DeviceInfo(protocol))); err != nil {
		logger.Infof("加载设备信息失败: %v", err)
	}
	qqClient.UseDevice(device)

	var hasLogin bool
	buff, err := ioutil.ReadFile("session.token")
	if err != nil {
		logger.Info("加载登录信息失败,扫码重新登录....")
	} else {
		err = qqClient.TokenLogin(buff)
		hasLogin = err == nil
		if !hasLogin {
			hasLogin = err.Error() == "already online"
		}
	}

	if !hasLogin {
		err = service.QrcodeLogin(qqClient)
	}
	if err != nil {
		return err
	}

	feature.RegisterFeature(qqClient)

	ioutil.WriteFile("session.token", qqClient.GenToken(), os.FileMode(0600))
	logger.Infof("登录成功 欢迎使用: %v", qqClient.Nickname)
	logger.Info("开始加载好友列表...")
	service.Check(qqClient.ReloadFriendList(), true)
	logger.Infof("共加载 %v 个好友.", len(qqClient.FriendList))
	logger.Infof("开始加载群列表...")
	service.Check(qqClient.ReloadGroupList(), true)
	logger.Infof("共加载 %v 个群.", len(qqClient.GroupList))
	logger.Info("加载完成")
	return nil
}

func LoginWithInstance(qqClient *robot.Client, sin chan string) {
	if sin == nil {
		sin = make(chan string, 1)
	}
	defer func() {
		close(sin)
	}()

	if qqClient == nil {
		sin <- consts.InstanceEmptySingle
		return
	}
	device := new(client.DeviceInfo)
	if err := device.ReadJson([]byte(DeviceInfo(ProtocolWatch))); err != nil {
		logger.Infof("加载设备信息失败: %v", err)
	}
	qqClient.UseDevice(device)

	var hasLogin bool
	buff, err := ioutil.ReadFile("session.token")
	if err != nil {
		logger.Info("加载登录信息失败,扫码重新登录....")
	} else {
		err = qqClient.TokenLogin(buff)
		hasLogin = err == nil
		if !hasLogin {
			hasLogin = err.Error() == "already online"
		}
	}

	if !hasLogin {
		err = service.QrcodeLoginWithSingle(qqClient, sin)
	} else {
		sin <- consts.LoginSuccessSingle
	}
	if err != nil {
		sin <- err.Error()
		return
	}

	feature.RegisterFeature(qqClient)

	ioutil.WriteFile("session.token", qqClient.GenToken(), os.FileMode(0600))
	logger.Infof("登录成功 欢迎使用: %v", qqClient.Nickname)
	logger.Info("开始加载好友列表...")
	service.Check(qqClient.ReloadFriendList(), true)
	logger.Infof("共加载 %v 个好友.", len(qqClient.FriendList))
	logger.Infof("开始加载群列表...")
	service.Check(qqClient.ReloadGroupList(), true)
	logger.Infof("共加载 %v 个群.", len(qqClient.GroupList))
	logger.Info("加载完成")
	return
}

func RunQQService(ctx context.Context, c *robot.Client) {
	go service.Run(ctx, c)
}

// IsOnline QQ是否在线
func IsOnline(c *robot.Client) bool {
	ok, _ := service.IFQQOnline(c)
	return ok
}
