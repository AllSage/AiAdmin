//go:build !windows

package tray

import (
	"fmt"

	"github.com/AllSage/AiAdmin/app/tray/commontray"
)

func InitPlatformTray(icon, updateIcon []byte) (commontray.AiAdminTray, error) {
	return nil, fmt.Errorf("NOT IMPLEMENTED YET")
}
