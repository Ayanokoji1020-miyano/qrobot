package main

import (
	"context"
	"errors"
	"github.com/Mrs4s/MiraiGo/client"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"qqRobot/feature"
	"qqRobot/robot"
	"qqRobot/service"
)

var qqClient *robot.Client

func login(qq int64, pwd string, protocol int) error {
	qqClient = robot.NewClient(qq, pwd)

	device := new(client.DeviceInfo)
	if err := device.ReadJson([]byte(deviceInfo(protocol))); err != nil {
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

func LoginWithInstance(qqClient *robot.Client, protocol int) error {
	if qqClient == nil {
		return errors.New("传入实例为 nil")
	}
	device := new(client.DeviceInfo)
	if err := device.ReadJson([]byte(deviceInfo(protocol))); err != nil {
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

func RunQQService(ctx context.Context, c *robot.Client) {
	go service.Run(ctx, c)
}

// IsOnline QQ是否在线
func IsOnline(c *robot.Client) bool {
	ok, _ := service.IFQQOnline(c)
	return ok
}
