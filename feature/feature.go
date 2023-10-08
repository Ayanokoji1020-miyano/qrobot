package feature

import (
	"github.com/Ayanokoji1020-miyano/qrobot/feature/listeners"
	"github.com/Ayanokoji1020-miyano/qrobot/feature/plugins"
	"github.com/Ayanokoji1020-miyano/qrobot/robot"
	logger "github.com/sirupsen/logrus"
)

func RegisterFeature(c *robot.Client) {
	// 事件监听器
	actionsListeners := []*robot.ActionListener{
		listeners.NewLogListenerInstance(),
	}
	c.SetActionListeners(actionsListeners...)
	// 自定义组件
	cPlugins := []*robot.Plugin{
		plugins.NewLogPluginInstance(),
	}
	// 系统组件
	sPlugins := []*robot.Plugin{
		plugins.NewMenuPluginInstance(cPlugins),
	}
	err := c.SetPlugins(
		append(sPlugins, cPlugins...)...,
	)
	if err != nil {
		logger.Error(err)
	}
	// 插件过滤器 true为阻止该插件
	c.SetPluginBlocker(func(plugin *robot.Plugin, contactType int, contactNumber int64) bool {
		return false
	})
}
