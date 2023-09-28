package robot

import "github.com/Mrs4s/MiraiGo/message"

type ActionListener struct {
	Uid  UID
	Name string
	// 发送了私聊消息将会执行回调
	SendPrivateMessage func(c *Client, message *message.PrivateMessage) bool
	// 发送了组群消息将会执行回调
	SendGroupMessage func(c *Client, message *message.GroupMessage) bool
	// 发送了私聊消息将会执行回调
	SendTempMessage func(c *Client, message *message.TempMessage, target int64) bool
}
