package guodiauthenticate

import (
	"crypto/md5"
	"encoding/hex"
	"guodi/src/guodiredis"
)

func FrontAuthentice(appID string, md5ID string, date string) bool {
	// 调用redis服务AppID验证模块
	isexixt, appSecretID := guodiredis.CheckAppID(appID)
	// isexixt := true
	// appSecretID := "asdsdf"
	if isexixt {
		if len(appSecretID) > 0 && len(date) > 0 {
			o := md5.New()
			srcStr := date + appSecretID
			o.Write([]byte(srcStr))
			s := o.Sum(nil)
			dstencryption := hex.EncodeToString(s)
			if dstencryption == md5ID {
				return true
			}
		}
	} else {
		return false
	}
	return false
}
