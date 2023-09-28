package consts

const (
	DirectionSending   = "SENDING"
	DirectionReceiving = "RECEIVING"
	TypePrivate        = "PRIVATE"
	TypeGroup          = "GROUP"
	TypeTemp           = "TEMP"
)

//------------------------------------------------ action listeners ------------------------------------

const (
	ListenerLogUID = iota + 1 // 日志监听者 uid
)

var ListenerMap = map[int]string{
	ListenerLogUID: "日志",
}

//------------------------------------------------ plugins ------------------------------------

const (
	PluginLogUID = iota + 1 // 日志插件 uid
)

var PluginMap = map[int]string{
	PluginLogUID: "日志",
}
