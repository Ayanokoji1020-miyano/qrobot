package qrobot

import "fmt"

func DeviceInfo(protocol int) string {
	cfg := `{
  "display": "MIRAI.211876.001",
  "product": "mirai",
  "device": "mirai",
  "board": "mirai",
  "model": "mirai",
  "finger_print": "mamoe/mirai/mirai:10/MIRAI.200122.001/4910920:user/release-keys",
  "boot_id": "b1cc7e1c-fc9c-f512-07e2-acbdb829f00a",
  "proc_version": "Linux version 3.0.31-DLNvZxhc (android-build@xxx.xxx.xxx.xxx.com)",
  "protocol": %v,
  "imei": "597005714727425",
  "brand": "mamoe",
  "bootloader": "unknown",
  "base_band": "",
  "version": {
    "incremental": "5891938",
    "release": "10",
    "codename": "REL",
    "sdk": 29
  },
  "sim_info": "T-Mobile",
  "os_type": "android",
  "mac_address": "00:50:56:C0:00:08",
  "ip_address": [
    10,
    0,
    1,
    3
  ],
  "wifi_bssid": "00:50:56:C0:00:08",
  "wifi_ssid": "\u003cunknown ssid\u003e",
  "imsi_md5": "8ad8f0747209c932d9d2914ebf690f7e",
  "android_id": "6aa5f4f149a7c50f",
  "apn": "wifi",
  "vendor_name": "MIUI",
  "vendor_os_name": "mirai"
}`
	return fmt.Sprintf(cfg, protocol)
}

const (
	ProtocolDefault = 0 // QQ登录方式(Default/Unset)	当前版本下默认为iPad
	ProtocolAndroid = 1 // QQ登录方式(ndroid Phone)
	ProtocolWatch   = 2 // QQ登录方式(Android Watch)
	ProtocolMacOS   = 3 // QQ登录方式(MacOS)
	ProtocolQiDian  = 4 // QQ登录方式(企点)	只能登录企点账号或企点子账号
	ProtocolIPad    = 5 // QQ登录方式(iPad)
	ProtocolAPad    = 6 // QQ登录方式(aPad)
)
