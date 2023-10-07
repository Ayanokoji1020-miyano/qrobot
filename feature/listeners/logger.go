package listeners

import (
	"github.com/Mrs4s/MiraiGo/message"
	logger "github.com/sirupsen/logrus"
	"qrobot/consts"
	"qrobot/robot"
)

func NewLogListenerInstance() *robot.ActionListener {
	return &robot.ActionListener{
		Uid:  consts.ListenerLogUID,
		Name: consts.ListenerMap[consts.ListenerLogUID],
		SendPrivateMessage: func(c *robot.Client, message *message.PrivateMessage) bool {
			buff, err := c.FormatMessageElements(message.Elements)
			if err == nil {
				logger.Info(string(buff))
			}
			return false
		},
		SendGroupMessage: func(c *robot.Client, message *message.GroupMessage) bool {
			buff, err := c.FormatMessageElements(message.Elements)
			if err == nil {
				logger.Info(string(buff))
			}
			return false
		},
		SendTempMessage: func(c *robot.Client, message *message.TempMessage, target int64) bool {
			buff, err := c.FormatMessageElements(message.Elements)
			if err == nil {
				logger.Info(string(buff))
			}
			return false
		},
	}
}
