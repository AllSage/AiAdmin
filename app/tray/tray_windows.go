package tray

import (
	"github.com/AllSage/AiAdmin/app/tray/commontray"
	"github.com/AllSage/AiAdmin/app/tray/wintray"
)

func InitPlatformTray(icon, updateIcon []byte) (commontray.AiAdminTray, error) {
	return wintray.InitTray(icon, updateIcon)
}
