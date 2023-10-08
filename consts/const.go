package consts

const (
	DirectionSending   = "SENDING"
	DirectionReceiving = "RECEIVING"
	TypePrivate        = "PRIVATE"
	TypeGroup          = "GROUP"
	TypeTemp           = "TEMP"
)

type Level = int

const (
	LevelSystem  Level = 1 // 系统级别
	LevelFeature       = 2 // 功能级别
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
	PluginMenuUID = iota + 1 // 菜单插件 uid
	PluginLogUID             // 日志插件 uid
)

var PluginMap = map[int]string{
	PluginMenuUID: "菜单",
	PluginLogUID:  "日志",
}

//------------------------------------------------ chan single ------------------------------------

const (
	LoginSuccessSingle  = "login success"
	ScanSuccessSingle   = "请扫码登录"
	InstanceEmptySingle = "传入实例为空"
)
