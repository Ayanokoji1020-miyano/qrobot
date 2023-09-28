package service

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/utils"
	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"gopkg.ilharper.com/x/isatty"
	"image"
	"image/png"
	"log"
	"os"
	"qqRobot/download"
	"qqRobot/robot"
	"strings"
	"time"
)

func printQRCode(imgData []byte) {
	const (
		black = "\033[48;5;0m  \033[0m"
		white = "\033[48;5;7m  \033[0m"
	)
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Panic(err)
	}
	data := img.(*image.Gray).Pix
	bound := img.Bounds().Max.X
	buf := make([]byte, 0, (bound*4+1)*(bound))
	i := 0
	for y := 0; y < bound; y++ {
		i = y * bound
		for x := 0; x < bound; x++ {
			if data[i] != 255 {
				buf = append(buf, white...)
			} else {
				buf = append(buf, black...)
			}
			i++
		}
		buf = append(buf, '\n')
	}
	_, _ = colorable.NewColorableStdout().Write(buf)
}

var ErrSMSRequestError = errors.New("sms request error")

func loginResponseProcessor(c *robot.Client, res *client.LoginResponse) error {
	var err error
	for {
		if err != nil {
			return err
		}
		if res.Success {
			return nil
		}
		var text string
		switch res.Error {
		case client.SliderNeededError:
			logger.Warnf("登录需要滑条验证码, 请验证后重试.")
			ticket := getTicket(res.VerifyUrl)
			if ticket == "" {
				logger.Infof("按 Enter 继续....")
				readLine()
				os.Exit(0)
			}
			res, err = c.QQClient.SubmitTicket(ticket)
			continue
		case client.NeedCaptcha:
			logger.Warnf("登录需要验证码.")
			_ = os.WriteFile("captcha.jpg", res.CaptchaImage, 0o644)
			logger.Warnf("请输入验证码 (captcha.jpg)： (Enter 提交)")
			text = readLine()
			DelFile("captcha.jpg")
			res, err = c.QQClient.SubmitCaptcha(text, res.CaptchaSign)
			continue
		case client.SMSNeededError:
			logger.Warnf("账号已开启设备锁, 按 Enter 向手机 %v 发送短信验证码.", res.SMSPhone)
			readLine()
			if !c.QQClient.RequestSMS() {
				logger.Warnf("发送验证码失败，可能是请求过于频繁.")
				return errors.WithStack(ErrSMSRequestError)
			}
			logger.Warn("请输入短信验证码： (Enter 提交)")
			text = readLine()
			res, err = c.QQClient.SubmitSMS(text)
			continue
		case client.SMSOrVerifyNeededError:
			logger.Warnf("账号已开启设备锁，请选择验证方式:")
			logger.Warnf("1. 向手机 %v 发送短信验证码", res.SMSPhone)
			logger.Warnf("2. 使用手机QQ扫码验证.")
			logger.Warn("请输入(1 - 2)：")
			text = readIfTTY("2")
			if strings.Contains(text, "1") {
				if !c.QQClient.RequestSMS() {
					logger.Warnf("发送验证码失败，可能是请求过于频繁.")
					return errors.WithStack(ErrSMSRequestError)
				}
				logger.Warn("请输入短信验证码： (Enter 提交)")
				text = readLine()
				res, err = c.QQClient.SubmitSMS(text)
				continue
			}
			fallthrough
		case client.UnsafeDeviceError:
			logger.Warnf("账号已开启设备锁，请前往 -> %v <- 验证后重启Bot.", res.VerifyUrl)
			logger.Infof("按 Enter 或等待 5s 后继续....")
			readLineTimeout(time.Second * 5)
			os.Exit(0)
		case client.OtherLoginError, client.UnknownLoginError, client.TooManySMSRequestError:
			msg := res.ErrorMessage
			logger.Warnf("登录失败: %v Code: %v", msg, res.Code)
			switch res.Code {
			case 235:
				logger.Warnf("设备信息被封禁, 请删除 device.json 后重试.")
			case 237:
				logger.Warnf("登录过于频繁, 请在手机QQ登录并根据提示完成认证后等一段时间重试")
			case 45:
				logger.Warnf("你的账号被限制登录, 请配置 SignServer 后重试")
			}
			logger.Infof("按 Enter 继续....")
			readLine()
			os.Exit(0)
		}
	}
}

var console = bufio.NewReader(os.Stdin)

func readLine() (str string) {
	str, _ = console.ReadString('\n')
	str = strings.TrimSpace(str)
	return
}

func getTicket(u string) string {
	logger.Warnf("请选择提交滑块ticket方式:")
	logger.Warnf("1. 自动提交")
	logger.Warnf("2. 手动抓取提交")
	logger.Warn("请输入(1 - 2)：")
	text := readLine()
	id := utils.RandomString(8)
	auto := !strings.Contains(text, "2")
	if auto {
		u = strings.ReplaceAll(u, "https://ssl.captcha.qq.com/template/wireless_mqq_captcha.html?", fmt.Sprintf("https://captcha.go-cqhttp.org/captcha?id=%v&", id))
	}
	logger.Warnf("请前往该地址验证 -> %v ", u)
	if !auto {
		logger.Warn("请输入ticket： (Enter 提交)")
		return readLine()
	}

	for count := 120; count > 0; count-- {
		str := fetchCaptcha(id)
		if str != "" {
			return str
		}
		time.Sleep(time.Second)
	}
	logger.Warnf("验证超时")
	return ""
}

func DelFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		// 删除失败
		logger.Error(err)
		return false
	}
	// 删除成功
	logger.Info(path + "删除成功")
	return true
}

func readLineTimeout(t time.Duration) {
	r := make(chan string)
	go func() {
		select {
		case r <- readLine():
		case <-time.After(t):
		}
	}()
	select {
	case <-r:
	case <-time.After(t):
	}
}

func readIfTTY(de string) (str string) {
	if isatty.Isatty(os.Stdin.Fd()) {
		return readLine()
	}
	logger.Warnf("未检测到输入终端，自动选择%s.", de)
	return de
}

func fetchCaptcha(id string) string {
	g, err := download.Request{URL: "https://captcha.go-cqhttp.org/captcha/ticket?id=" + id}.JSON()
	if err != nil {
		logger.Debugf("获取 Ticket 时出现错误: %v", err)
		return ""
	}
	if g.Get("ticket").Exists() {
		return g.Get("ticket").String()
	}
	return ""
}
