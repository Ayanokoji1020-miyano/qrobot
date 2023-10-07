package plugins

import (
	"fmt"
	"qrobot/robot"
	"strings"
)

func NewMenuPluginInstance(customerPlugins []*robot.Plugin) *robot.Plugin {
	return &robot.Plugin{
		Uid:  2,
		Name: "菜单",
		RCVMessage: func(client *robot.Client, messageInterface interface{}) bool {
			content := client.MessageContent(messageInterface)
			if strings.EqualFold("nb", content) {
				builder := strings.Builder{}
				builder.WriteString("菜单 : ")
				for i := 0; i < len(customerPlugins); i++ {
					builder.WriteString(fmt.Sprintf("\n♦️ %s", (*customerPlugins[i]).Name))
				}
				client.ReplyText(messageInterface, builder.String())
				return true
			}
			return false
		},
	}
}
