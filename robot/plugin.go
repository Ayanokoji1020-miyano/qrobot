package robot

import (
	"github.com/Ayanokoji1020-miyano/qrobot/consts"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

type Plugin struct {
	PluginLevel consts.Level
	Uid         UID
	Name        string
	// 收到私聊消息时
	RCVPrivateMessage func(client *Client, privateMessage *message.PrivateMessage) bool
	// 收到组群消息时
	RCVGroupMessage func(client *Client, groupMessage *message.GroupMessage) bool
	// 收到临时消息时
	RCVTempMessage func(client *Client, tempMessage *message.TempMessage) bool
	// 收到消息时, 优先级低于明确类型的Message
	RCVMessage func(client *Client, messageInterface interface{}) bool
	// 收到好友请求时
	RCVNewFriendRequest func(client *Client, request *client.NewFriendRequest) bool
	// 添加了好友时
	RCVNewFriendAdded func(client *Client, event *client.NewFriendEvent) bool
	// 收到组群邀请时
	RCVGroupInvited func(client *Client, info *client.GroupInvitedRequest) bool
	// 加入组群时
	RCVJoinGroup func(client *Client, info *client.MemberJoinGroupEvent) bool
	// 离开组群时
	RCVLeaveGroup func(client *Client, event *client.GroupLeaveEvent) bool
}
