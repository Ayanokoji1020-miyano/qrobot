package robot

import (
	"github.com/Ayanokoji1020-miyano/qrobot/consts"
	"github.com/Mrs4s/MiraiGo/message"
)

type ActionListener struct {
	ListenerLevel consts.Level
	Uid           UID
	Name          string
	// 发送了私聊消息将会执行回调
	SendPrivateMessage func(c *Client, message *message.PrivateMessage) bool
	// 发送了组群消息将会执行回调
	SendGroupMessage func(c *Client, message *message.GroupMessage) bool
	// 发送了私聊消息将会执行回调
	SendTempMessage func(c *Client, message *message.TempMessage, target int64) bool
}
