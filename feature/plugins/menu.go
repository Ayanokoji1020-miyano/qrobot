package plugins

import (
	"fmt"
	"github.com/Ayanokoji1020-miyano/qrobot/consts"
	"github.com/Ayanokoji1020-miyano/qrobot/robot"
	"strings"
)

func NewMenuPluginInstance(customerPlugins []*robot.Plugin) *robot.Plugin {
	return &robot.Plugin{
		PluginLevel: consts.LevelSystem,
		Uid:         consts.PluginMenuUID,
		Name:        consts.PluginMap[consts.PluginMenuUID],
		RCVMessage: func(client *robot.Client, messageInterface interface{}) bool {
			content := client.MessageContent(messageInterface)
			if strings.EqualFold("系统功能", content) {
				builder := strings.Builder{}
				builder.WriteString("系统功能 : ")
				for i := 0; i < len(customerPlugins); i++ {
					builder.WriteString(fmt.Sprintf("\n%v. %s", i+1, (*customerPlugins[i]).Name))
				}
				client.ReplyText(messageInterface, builder.String())
				return true
			}
			return false
		},
	}
}
