package service

import (
	"context"
	"github.com/Ayanokoji1020-miyano/qrobot/robot"
	"os"
	"time"
)

// IFQQOnline QQ是否在线
func IFQQOnline(c *robot.Client) (bool, error) {
	if !PathExists("session.token") {
		return false, nil
	}

	token, err := os.ReadFile("session.token")
	if err != nil {
		return false, err
	}

	err = c.TokenLogin(token)
	if err != nil && err.Error() == "already online" {
		err = nil
	}
	return err == nil, err
}

func Run(ctx context.Context, c *robot.Client) {
	go Listen(ctx, c)
	go Online(ctx, c)
}

// Listen 维持QQ监听
func Listen(ctx context.Context, c *robot.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func Online(ctx context.Context, c *robot.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ok, err := IFQQOnline(c)
			if err != nil || !ok {
				_, cancel := context.WithCancel(ctx)
				cancel()
			}
		}
		time.Sleep(time.Second * 60)
	}
}
