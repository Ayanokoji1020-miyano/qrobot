package plugins

import (
	"github.com/Ayanokoji1020-miyano/qrobot/consts"
	"github.com/Ayanokoji1020-miyano/qrobot/robot"
	"github.com/Mrs4s/MiraiGo/message"
	logger "github.com/sirupsen/logrus"
)

func NewLogPluginInstance() *robot.Plugin {
	return &robot.Plugin{
		PluginLevel: consts.LevelSystem,
		Uid:         consts.PluginLogUID,
		Name:        consts.ListenerMap[consts.PluginLogUID],
		RCVMessage: func(client *robot.Client, messageInterface interface{}) bool {

			if privateMessage, b := (messageInterface).(*message.PrivateMessage); b {
				buff, err := client.FormatMessageElements(privateMessage.Elements)
				if err == nil {
					logger.Info(string(buff))
				}
			} else if groupMessage, b := (messageInterface).(*message.GroupMessage); b {
				buff, err := client.FormatMessageElements(groupMessage.Elements)
				if err == nil {
					logger.Info(string(buff))
				}
			} else if tempMessage, b := (messageInterface).(*message.TempMessage); b {
				buff, err := client.FormatMessageElements(tempMessage.Elements)
				if err == nil {
					logger.Info(string(buff))
				}
			}

			return false
		},
	}
}
