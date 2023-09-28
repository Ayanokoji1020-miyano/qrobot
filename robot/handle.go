package robot

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

// Handler 监听QQ事件
type Handler struct {
	c *Client
}

const (
	ContactTypePrivate = 1 // 私聊消息
	ContactTypeGroup   = 2 // 群聊消息
)

func (h *Handler) PrivateMessage(qqClient *client.QQClient, privateMessage *message.PrivateMessage) {
	h.c.logMessage(privateMessage, logFlagReceiving)
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypePrivate, privateMessage.Sender.Uin) {
				return false
			}
			return (mPoint.RCVPrivateMessage != nil && mPoint.RCVPrivateMessage(h.c, privateMessage)) ||
				(mPoint.RCVMessage != nil && mPoint.RCVMessage(h.c, privateMessage))
		})
}

func (h *Handler) GroupMessage(client *client.QQClient, groupMessage *message.GroupMessage) {
	h.c.logMessage(groupMessage, logFlagReceiving)
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypeGroup, groupMessage.GroupCode) {
				return false
			}
			return (mPoint.RCVGroupMessage != nil && mPoint.RCVGroupMessage(h.c, groupMessage)) ||
				(mPoint.RCVMessage != nil && mPoint.RCVMessage(h.c, groupMessage))
		},
	)
}

func (h *Handler) TempMessageEvent(qqClient *client.QQClient, tempMessage *client.TempMessageEvent) {
	h.c.logMessage(tempMessage, logFlagReceiving)
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypePrivate, tempMessage.Message.Sender.Uin) {
				return false
			}
			return (mPoint.RCVTempMessage != nil && mPoint.RCVTempMessage(h.c, tempMessage.Message)) ||
				(mPoint.RCVMessage != nil && mPoint.RCVMessage(h.c, tempMessage.Message))
		},
	)
}

func (h *Handler) NewFriendRequest(qqClient *client.QQClient, request *client.NewFriendRequest) {
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypePrivate, request.RequesterUin) {
				return false
			}
			return mPoint.RCVNewFriendRequest != nil && mPoint.RCVNewFriendRequest(h.c, request)
		},
	)
}

func (h *Handler) NewFriendEvent(qqClient *client.QQClient, event *client.NewFriendEvent) {
	qqClient.ReloadFriendList()
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypePrivate, event.Friend.Uin) {
				return false
			}
			return mPoint.RCVNewFriendAdded != nil && mPoint.RCVNewFriendAdded(h.c, event)
		},
	)
}

func (h *Handler) GroupInvitedRequest(qqClient *client.QQClient, request *client.GroupInvitedRequest) {
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypeGroup, request.GroupCode) {
				return false
			}
			return mPoint.RCVGroupInvited != nil && mPoint.RCVGroupInvited(h.c, request)
		},
	)
}

func (h *Handler) MemberJoinGroupEvent(qqClient *client.QQClient, info *client.MemberJoinGroupEvent) {
	qqClient.ReloadGroupList()
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypeGroup, info.Group.Code) {
				return false
			}
			return mPoint.RCVJoinGroup != nil && mPoint.RCVJoinGroup(h.c, info)
		},
	)
}

func (h *Handler) GroupLeaveEvent(qqClient *client.QQClient, event *client.GroupLeaveEvent) {
	qqClient.ReloadGroupList()
	h.c.execPlugins(
		func(mPoint *Plugin) bool {
			if h.c.pluginBlocker != nil && h.c.pluginBlocker(mPoint, ContactTypeGroup, event.Group.Code) {
				return false
			}
			return mPoint.RCVLeaveGroup != nil && mPoint.RCVLeaveGroup(h.c, event)
		},
	)
}
