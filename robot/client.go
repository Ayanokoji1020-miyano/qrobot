package robot

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ayanokoji1020-miyano/qrobot/consts"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strings"
	"sync"
)

type UID = int

type Client struct {
	*client.QQClient
	logger *logrus.Logger
	// 行为监听器
	actionListeners map[UID]*ActionListener
	amux            sync.RWMutex

	// plugins 接收到事件才执行
	plugins map[UID]*Plugin
	pmux    sync.RWMutex
	// 插件拦截器
	pluginBlocker func(plugin *Plugin, contactType int, contactNumber int64) bool
}

func NewClient(uin int64, password string) *Client {
	return newClientMd5(uin, md5.Sum([]byte(password)))
}

func newClientMd5(uin int64, password [16]byte) *Client {
	c := &Client{
		logger:          logrus.New(),
		QQClient:        client.NewClientMd5(uin, password),
		actionListeners: make(map[UID]*ActionListener),
		plugins:         make(map[UID]*Plugin),
	}
	c.logger.SetLevel(logrus.InfoLevel)
	c.logger.SetOutput(os.Stdout)
	c.SubscribeEventHandler(&Handler{c: c})
	// todo
	//c.OnServerUpdated(func(bot *client.QQClient, e *client.ServerUpdatedEvent) bool {
	//	if !base.UseSSOAddress {
	//		log.Infof("收到服务器地址更新通知, 根据配置文件已忽略.")
	//		return false
	//	}
	//	log.Infof("收到服务器地址更新通知, 将在下一次重连时应用. ")
	//	return true
	//})
	//if global.PathExists("address.txt") {
	//	log.Infof("检测到 address.txt 文件. 将覆盖目标IP.")
	//	addr := global.ReadAddrFile("address.txt")
	//	if len(addr) > 0 {
	//		c.SetCustomServer(addr)
	//	}
	//	log.Infof("读取到 %v 个自定义地址.", len(addr))
	//}
	return c
}

// SetActionListenersAndPlugins 设置监听器以及插件
func (c *Client) SetActionListenersAndPlugins(actionListeners []*ActionListener, plugins []*Plugin) error {
	err := c.SetActionListeners(actionListeners...)
	if err != nil {
		return err
	}
	return c.SetPlugins(plugins...)
}

// SetActionListeners 设置监听器
func (c *Client) SetActionListeners(actionListeners ...*ActionListener) error {
	c.amux.Lock()
	defer c.amux.Unlock()

	for _, action := range actionListeners {
		if action.Uid == 0 || action.Name == "" {
			return errors.New("注册监听器错误: Uid和Name不能为零值")
		}
		//_, ok := c.actionListeners[action.Uid]
		//if ok {
		//	continue
		//	c.logger.Info(fmt.Sprintf("不能重复注册监听器UID:%v", action.Uid))
		//}
		c.actionListeners[action.Uid] = action
	}
	return nil
}

// InsertActionListeners 插入或更新监听器
func (c *Client) InsertActionListeners(actionListeners ...*ActionListener) error {
	err := c.SetActionListeners(actionListeners...)
	if err != nil {
		return err
	}
	return nil
}

// RemoveActionListener 移除监听器
func (c *Client) RemoveActionListener(uid UID) {
	c.amux.Lock()
	defer c.amux.Unlock()
	delete(c.actionListeners, uid)
}

// 执行所有监听器
func (c *Client) execActionListeners(fun func(actionListener *ActionListener) bool) {
	c.amux.RLock()
	defer c.amux.RUnlock()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.logger.Error(fmt.Sprintf("action listener error: %v\n%s", err, debug.Stack()))
			}
		}()
		for id, a := range c.actionListeners {
			if fun(a) {
				c.logger.Info(fmt.Sprintf("<<< PROCESS BY MODULE(%d)", id))
				return
			}
		}
		c.logger.Info(fmt.Sprintf("<<< ACTION NOT PROCESS"))
	}()
}

// SetPlugins 设置插件
func (c *Client) SetPlugins(plugins ...*Plugin) error {
	c.pmux.Lock()
	defer c.pmux.Unlock()

	for _, p := range plugins {
		if p.PluginLevel == 0 {
			p.PluginLevel = consts.LevelFeature
		}
		if p.Uid == 0 || p.Name == "" {
			return errors.New("注册监听器错误: Uid和Name不能为零值")
		}
		c.plugins[p.Uid] = p
	}
	return nil
}

// InsertPlugins 插入或更新插件
func (c *Client) InsertPlugins(plugins ...*Plugin) error {
	err := c.SetPlugins(plugins...)
	if err != nil {
		return err
	}
	c.updateSystemMenuPlugin()
	return nil
}

// 更新系统菜单插件
func (c *Client) updateSystemMenuPlugin() {
	delete(c.plugins, consts.PluginMenuUID)
	p := &Plugin{
		PluginLevel: consts.LevelSystem,
		Uid:         consts.PluginMenuUID,
		Name:        consts.PluginMap[consts.PluginMenuUID],
		RCVMessage: func(client *Client, messageInterface interface{}) bool {
			content := client.MessageContent(messageInterface)
			if strings.EqualFold("系统功能", content) {
				builder := strings.Builder{}
				builder.WriteString("系统功能 : ")

				var s []int
				for k, p := range c.plugins {
					if p.PluginLevel == consts.LevelSystem {
						continue
					}
					s = append(s, k)
				}
				sort.Ints(s)
				var idx int
				for i := range s {
					idx++
					builder.WriteString(fmt.Sprintf("\n%v. %s", idx, c.plugins[i].Name))
				}

				client.ReplyText(messageInterface, builder.String())
				return true
			}
			return false
		},
	}
	c.SetPlugins(p)
}

// SetPluginBlocker 设置插件拦截器
func (c *Client) SetPluginBlocker(fun func(plugin *Plugin, contactType int, contactNumber int64) bool) {
	c.pluginBlocker = fun
}

// RemovePlugin 移除插件
func (c *Client) RemovePlugin(uid UID) {
	p, ok := c.plugins[uid]
	if !ok {
		return
	}

	var up bool
	up = p.PluginLevel == consts.LevelFeature

	c.pmux.Lock()
	defer c.pmux.Unlock()
	delete(c.plugins, uid)
	if up {
		c.updateSystemMenuPlugin()
	}
}

// 执行所有插件
func (c *Client) execPlugins(fun func(plugin *Plugin) bool) {
	c.pmux.RLock()
	defer c.pmux.RUnlock()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				c.logger.Error(fmt.Sprintf("event error: %v\n%s", err, debug.Stack()))
			}
		}()
		for id, p := range c.plugins {
			if fun(p) {
				c.logger.Info(fmt.Sprintf("<<< PROCESS BY MODULE(%d)", id))
				return
			}
		}
		c.logger.Info(fmt.Sprintf("<<< PLUGIN NOT PROCESS"))
	}()
}

//------------------------------------------------ 消息处理相关 ------------------------------------

// MessageElements 获取消息
func (c *Client) MessageElements(messageInterface interface{}) []message.IMessageElement {
	in := reflect.ValueOf(messageInterface).Elem().FieldByName("Elements").Interface()
	if array, ok := in.([]message.IMessageElement); ok {
		return array
	}
	return nil
}

// MessageContent 获取消息的文本内容
func (c *Client) MessageContent(messageInterface interface{}) string {
	return reflect.ValueOf(messageInterface).MethodByName("ToString").Call([]reflect.Value{})[0].String()
}

// MessageSenderUin 获取消息的发送者
func (c *Client) MessageSenderUin(source interface{}) int64 {
	if privateMessage, b := (source).(*message.PrivateMessage); b {
		return privateMessage.Sender.Uin
	} else if groupMessage, b := (source).(*message.GroupMessage); b {
		return groupMessage.Sender.Uin
	} else if tempMessage, b := (source).(*message.TempMessage); b {
		return tempMessage.Sender.Uin
	}
	return 0
}

// UploadReplyImage 上传图片用作回复
func (c *Client) UploadReplyImage(source interface{}, buffer []byte) (message.IMessageElement, error) {
	if privateMessage, b := (source).(*message.PrivateMessage); b {
		return c.UploadImage(
			message.Source{
				SourceType: message.SourcePrivate,
				PrimaryID:  privateMessage.Sender.Uin,
			}, bytes.NewReader(buffer),
		)
	} else if groupMessage, b := (source).(*message.GroupMessage); b {
		return c.UploadImage(
			message.Source{
				SourceType: message.SourcePrivate,
				PrimaryID:  groupMessage.GroupCode,
			}, bytes.NewReader(buffer),
		)
	}
	return nil, errors.New("only group message and private message")
}

// UploadReplyVideo 上传视频文件
func (c *Client) UploadReplyVideo(source interface{}, video []byte, thumb []byte) (*message.ShortVideoElement, error) {
	if groupMessage, b := (source).(*message.GroupMessage); b {
		return c.UploadShortVideo(
			message.Source{
				SourceType: message.SourceGroup,
				PrimaryID:  groupMessage.GroupCode,
			}, bytes.NewReader(video), bytes.NewReader(thumb),
		)
	}
	if privateMessage, b := (source).(*message.PrivateMessage); b {
		return c.UploadShortVideo(
			message.Source{
				SourceType: message.SourcePrivate,
				PrimaryID:  privateMessage.Sender.Uin,
			}, bytes.NewReader(video), bytes.NewReader(thumb),
		)
	}
	return nil, errors.New("only group message and private message")
}

// UploadReplyVoice 上传声音用作回复
func (c *Client) UploadReplyVoice(source interface{}, buffer []byte) (message.IMessageElement, error) {
	if privateMessage, b := (source).(*message.PrivateMessage); b {
		return c.UploadVoice(
			message.Source{
				SourceType: message.SourcePrivate,
				PrimaryID:  privateMessage.Sender.Uin,
			}, bytes.NewReader(buffer),
		)
	} else if groupMessage, b := (source).(*message.GroupMessage); b {
		return c.UploadVoice(
			message.Source{
				SourceType: message.SourceGroup,
				PrimaryID:  groupMessage.GroupCode,
			}, bytes.NewReader(buffer),
		)
	}
	return nil, errors.New("only group message and private message")
}

//------------------------------------------------ 群相关 ------------------------------------

// MessageFirstAt 第一个At的用户
func (c *Client) MessageFirstAt(groupMessage *message.GroupMessage) int64 {
	for _, element := range groupMessage.Elements {
		if element.Type() == message.At {
			if at, ok := element.(*message.AtElement); ok {
				return at.Target
			}
		}
	}
	return 0
}

// CardNameInGroup 获取成员名称
func (c *Client) CardNameInGroup(groupCode int64, uin int64) string {
	for _, group := range c.GroupList {
		if group.Code != groupCode {
			continue
		}
		for _, member := range group.Members {
			if member.Uin == uin {
				name := member.CardName
				if len(name) == 0 {
					name = member.Nickname
				}
				return name
			}
		}
	}
	return fmt.Sprintf("%d", uin)
}

// MakeReplySendingMessage 创建一个SendingMessage, 将会用于回复
func (c *Client) MakeReplySendingMessage(source interface{}) *message.SendingMessage {
	sending := message.NewSendingMessage()
	if groupMessage, b := (source).(*message.GroupMessage); b {
		sendGroupCode := groupMessage.GroupCode
		atUin := groupMessage.Sender.Uin
		return sending.Append(c.AtElement(sendGroupCode, atUin)).Append(message.NewText("\n\n"))
	}
	return sending
}

// AtElement 创建一个At
func (c *Client) AtElement(groupCode int64, uin int64) *message.AtElement {
	return message.NewAt(uin, fmt.Sprintf("@%s", c.CardNameInGroup(groupCode, uin)))
}

// ReplyText 快捷回复消息
func (c *Client) ReplyText(source interface{}, content string) {
	c.ReplyRawMessage(source, c.MakeReplySendingMessage(source).Append(message.NewText(content)))
}

// ReplyRawMessage 回复一个消息到源消息
func (c *Client) ReplyRawMessage(source interface{}, sendingMessage *message.SendingMessage) {
	if privateMessage, b := (source).(*message.PrivateMessage); b {
		c.SendPrivateMessage(privateMessage.Sender.Uin, sendingMessage)
	} else if groupMessage, b := (source).(*message.GroupMessage); b {
		c.SendGroupMessage(groupMessage.GroupCode, sendingMessage)
	} else if tempMessage, b := (source).(*message.TempMessage); b {
		c.SendGroupTempMessage(tempMessage.GroupCode, tempMessage.Sender.Uin, sendingMessage)
	}
}

//------------------------------------------------ 消息发送 ------------------------------------

// SendPrivateMsg 发送私聊消息
func (c *Client) SendPrivateMsg(target int64, msg string) {
	m := message.NewSendingMessage()
	m.Append(message.NewText(msg))
	_ = c.SendPrivateMessage(target, m)
}

// SendPrivateMessage 发送私聊消息
func (c *Client) SendPrivateMessage(target int64, m *message.SendingMessage) *message.PrivateMessage {
	privateMessage := c.QQClient.SendPrivateMessage(target, m)
	// todo 终端打印消息
	c.logMessage(privateMessage, logFlagSending)
	c.execActionListeners(func(actionListener *ActionListener) bool {
		return actionListener.SendPrivateMessage != nil && actionListener.SendPrivateMessage(c, privateMessage)
	})
	return privateMessage
}

// SendGroupMsg 发送群消息
func (c *Client) SendGroupMsg(groupCode int64, msg string) {
	m := message.NewSendingMessage()
	m.Append(message.NewText(msg))
	_ = c.SendGroupMessage(groupCode, m)
}

// SendGroupMessage 发送群消息
func (c *Client) SendGroupMessage(groupCode int64, m *message.SendingMessage) *message.GroupMessage {
	groupMessage := c.QQClient.SendGroupMessage(groupCode, m)
	c.logMessage(groupMessage, logFlagSending)
	c.execActionListeners(func(actionListener *ActionListener) bool {
		return actionListener.SendGroupMessage != nil && actionListener.SendGroupMessage(c, groupMessage)
	})
	return groupMessage
}

// SendGroupTempMsg 发送临时消息
func (c *Client) SendGroupTempMsg(groupCode, target int64, msg string) {
	m := message.NewSendingMessage()
	m.Append(message.NewText(msg))
	_ = c.SendGroupTempMessage(groupCode, target, m)
}

// SendGroupTempMessage 发送临时消息
func (c *Client) SendGroupTempMessage(groupCode, target int64, m *message.SendingMessage) *message.TempMessage {
	tempMessage := c.QQClient.SendGroupTempMessage(groupCode, target, m)
	c.logMessage(tempMessage, logFlagSending, target)
	c.execActionListeners(func(actionListener *ActionListener) bool {
		return actionListener.SendTempMessage != nil && actionListener.SendTempMessage(c, tempMessage, target)
	})
	return tempMessage
}

// SendMessage 发送消息。
//
// group == 0, qq != 0: 发送私聊消息
//
// group != 0, qq == 0: 发送群消息
//
// group != 0, qq != 0: 发送临时消息
func (c *Client) SendMessage(targetGroup, targetQQ int64, msg string) {
	switch targetGroup != 0 {
	case false:
		if targetQQ != 0 {
			c.SendPrivateMsg(targetQQ, msg)
		} else {
			return
		}
	case true:
		if targetQQ == 0 {
			c.SendGroupMsg(targetGroup, msg)
		}

		if targetQQ != 0 {
			if targetQQ == targetGroup {
				// targetGroup, targetQQ 都不为0, 且相等时发送私聊消息
				c.SendPrivateMsg(targetQQ, msg)
			} else {
				c.SendGroupTempMsg(targetGroup, targetQQ, msg)
			}
		}
	}
}

//------------------------------------------------ 开发临时使用 ------------------------------------

// 显示用的日志

const logFlagReceiving = "RECEIVING"
const logFlagSending = "SENDING"

func (c *Client) logMessage(m interface{}, logFlag string, ext ...interface{}) {

	var flag string
	var entries []message.IMessageElement

	if logFlag == logFlagSending {
		flag = "<<< Sending <<<"
	}
	if logFlag == logFlagReceiving {
		flag = ">>> Receiving >>>"
	}

	if privateMessage, ok := m.(*message.PrivateMessage); ok {
		entries = privateMessage.Elements
		flag += " PRIVATE :"
		if logFlag == logFlagSending {
			flag += fmt.Sprintf(" UID(%d) ", privateMessage.Target)
		}
		if logFlag == logFlagReceiving {
			flag += fmt.Sprintf(" UID(%d) ", privateMessage.Sender.Uin)
		}
	}
	if groupMessage, ok := m.(*message.GroupMessage); ok {
		entries = groupMessage.Elements
		flag += " GROUP :"
		if logFlag == logFlagSending {
			flag += fmt.Sprintf(" GID(%d) ", groupMessage.GroupCode)
		}
		if logFlag == logFlagReceiving {
			flag += fmt.Sprintf(" GID(%d) UID(%d) ", groupMessage.GroupCode, groupMessage.Sender.Uin)
		}
	}
	if tempMessage, ok := m.(*message.TempMessage); ok {
		entries = tempMessage.Elements
		flag += " TEMP :"
		if logFlag == logFlagSending {
			flag += fmt.Sprintf(" GID(%d) ", tempMessage.GroupCode)
			if len(ext) > 0 {
				if id, ok := ext[0].(int64); ok {
					flag += fmt.Sprintf("UID(%d) ", id)
				}
			}
		}
		if logFlag == logFlagReceiving {
			flag += fmt.Sprintf(" GID(%d) UID(%d) ", tempMessage.GroupCode, tempMessage.Sender.Uin)
		}
	}

	contentBuff, e := c.FormatMessageElements(entries)

	if e != nil {
		c.logger.Error("LOG ERROR : ", flag, " : ", e.Error())
	}
	content := string(contentBuff)

	builder := strings.Builder{}
	builder.WriteString(flag)
	builder.WriteString("\n")
	builder.WriteString(content)
	c.logger.Info(builder.String())
}

func (c *Client) FormatMessageElements(entries []message.IMessageElement) ([]byte, error) {
	var fEntries []interface{}

	for i := range entries {
		if app, b := (entries[i]).(*message.LightAppElement); b {
			fEntries = append(fEntries, map[string]string{
				"Type":    "LightAPP",
				"Content": app.Content,
			})
		} else if text, b := (entries[i]).(*message.TextElement); b {
			fEntries = append(fEntries, map[string]string{
				"Type":    "Text",
				"Content": text.Content,
			})
		} else if img, b := (entries[i]).(*message.GroupImageElement); b {
			fEntries = append(fEntries, map[string]string{
				"ImageId": img.ImageId,
				"Type":    "Image",
				"Url":     img.Url,
			})
		} else if img, b := (entries[i]).(*message.FriendImageElement); b {
			fEntries = append(fEntries, map[string]string{
				"ImageId": img.ImageId,
				"Type":    "Image",
				"Url":     img.Url,
			})
		} else if at, b := (entries[i]).(*message.AtElement); b {
			fEntries = append(fEntries, map[string]interface{}{
				"Type":    "At",
				"Target":  at.Target,
				"Display": at.Display,
			})
		} else if voice, b := (entries[i]).(*message.VoiceElement); b {
			fEntries = append(fEntries, map[string]string{
				"Type":    "Voice",
				"Name":    voice.Name,
				"Display": voice.Url,
			})
		} else if redBag, b := (entries[i]).(*message.RedBagElement); b {
			fEntries = append(fEntries, map[string]interface{}{
				"Type":   "RegBag",
				"Title":  redBag.Title,
				"RbType": int(redBag.MsgType),
			})
		} else if face, b := (entries[i]).(*message.FaceElement); b {
			fEntries = append(fEntries, map[string]interface{}{
				"Type":  "Face",
				"Name":  face.Name,
				"Index": face.Index,
			})
		} else {
			fEntries = append(fEntries, map[string]interface{}{
				"Type":    "Other",
				"SubType": int(entries[i].Type()),
			})
		}
	}
	return json.Marshal(&fEntries)
}
