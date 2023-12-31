package qrobot

import (
	"context"
	"github.com/Ayanokoji1020-miyano/qrobot/consts"
	"github.com/Ayanokoji1020-miyano/qrobot/robot"
	"os"
	"testing"
)

var c *robot.Client

func TestMain(m *testing.M) {
	c = robot.NewClient(000000000, "")
	os.Exit(m.Run())
}

func TestLoginWithInstance(t *testing.T) {
	sin := make(chan string, 1)
	go LoginWithInstance(c, sin)

	for {
		select {
		case info, ok := <-sin:
			if !ok {
				break
			}
			switch info {
			case consts.InstanceEmptySingle:
				t.Log(info)
			case consts.ScanSuccessSingle:
				t.Log(info)
			case consts.LoginSuccessSingle:
				//todo ctx 改为子上下文
				ctx := context.Background()
				go RunQQService(ctx, c)
				t.Log("success")
			}
		}
	}
}

func TestIsOnline(t *testing.T) {
	TestLoginWithInstance(t)
	ok := IsOnline(c)
	t.Log(ok)
}
