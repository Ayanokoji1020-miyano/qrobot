package service

import (
	"errors"
	"github.com/Mrs4s/MiraiGo/client"
	logger "github.com/sirupsen/logrus"
	"log"
	"os"
	"qqRobot/consts"
	"qqRobot/robot"
	"time"
)

func QrcodeLogin(c *robot.Client) error {
	rsp, err := c.QQClient.FetchQRCodeCustomSize(1, 2, 1)
	if err != nil {
		return err
	}
	_ = os.WriteFile("qrcode.png", rsp.ImageData, 0o644)
	defer func() { _ = os.Remove("qrcode.png") }()
	if c.QQClient.Uin != 0 {
		logger.Infof("请使用账号 %v 登录手机QQ扫描二维码 (qrcode.png) : ", c.QQClient.Uin)
	} else {
		logger.Infof("请使用手机QQ扫描二维码 (qrcode.png) : ")
	}
	time.Sleep(time.Second)
	printQRCode(rsp.ImageData)
	s, err := c.QQClient.QueryQRCodeStatus(rsp.Sig)
	if err != nil {
		return err
	}
	prevState := s.State

	for {
		time.Sleep(time.Second)
		s, _ = c.QQClient.QueryQRCodeStatus(rsp.Sig)
		if s == nil {
			continue
		}
		if prevState == s.State {
			continue
		}
		prevState = s.State
		switch s.State {
		case client.QRCodeCanceled:
			logger.Infof("扫码被用户取消.")
		case client.QRCodeTimeout:
			logger.Infof("二维码过期")
		case client.QRCodeWaitingForConfirm:
			logger.Infof("扫码成功, 请在手机端确认登录.")
		case client.QRCodeConfirmed:
			res, err := c.QQClient.QRCodeLogin(s.LoginInfo)
			if err != nil {
				return err
			}
			return loginResponseProcessor(c, res)
		case client.QRCodeImageFetch, client.QRCodeWaitingForScan:
			// ignore
		}
	}
}

func QrcodeLoginWithSingle(c *robot.Client, sin chan string) error {
	rsp, err := c.QQClient.FetchQRCodeCustomSize(1, 2, 1)
	if err != nil {
		return err
	}
	_ = os.WriteFile("qrcode.png", rsp.ImageData, 0o644)
	defer func() { _ = os.Remove("qrcode.png") }()
	if c.QQClient.Uin != 0 {
		logger.Infof("请使用账号 %v 登录手机QQ扫描二维码 (qrcode.png) : ", c.QQClient.Uin)
	} else {
		logger.Infof("请使用手机QQ扫描二维码 (qrcode.png) : ")
	}
	time.Sleep(time.Second)
	printQRCode(rsp.ImageData)
	s, err := c.QQClient.QueryQRCodeStatus(rsp.Sig)
	if err != nil {
		return err
	}
	prevState := s.State

	if sin == nil {
		return errors.New("sin 信号量没有初始化")
	}
	sin <- consts.ScanSuccessSingle

	for {
		time.Sleep(time.Second)
		s, _ = c.QQClient.QueryQRCodeStatus(rsp.Sig)
		if s == nil {
			continue
		}
		if prevState == s.State {
			continue
		}
		prevState = s.State
		switch s.State {
		case client.QRCodeCanceled:
			logger.Infof("扫码被用户取消.")
			return errors.New("扫码被用户取消")
		case client.QRCodeTimeout:
			logger.Infof("二维码过期")
			return errors.New("二维码过期")
		case client.QRCodeWaitingForConfirm:
			logger.Infof("扫码成功, 请在手机端确认登录.")
		case client.QRCodeConfirmed:
			res, err := c.QQClient.QRCodeLogin(s.LoginInfo)
			if err != nil {
				return err
			}
			return loginResponseProcessor(c, res)
		case client.QRCodeImageFetch, client.QRCodeWaitingForScan:
			// ignore
		}
	}
}

// Check 检测err是否为nil
func Check(err error, deleteSession bool) {
	if err != nil {
		if deleteSession && PathExists("session.token") {
			_ = os.Remove("session.token")
		}
		log.Fatalf("遇到错误: %v", err)
	}
}

// PathExists 判断给定path是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, os.ErrExist)
}
