// [FORK-MOD] Mock VIP user — fork 模式下把所有用户信息写为已解锁状态
package api

import (
	"github.com/88250/gulu"
	"github.com/siyuan-note/siyuan/kernel/conf"
	"github.com/siyuan-note/siyuan/kernel/model"
	"github.com/siyuan-note/siyuan/kernel/util"
)

func mockCloudUser(token string) (*conf.User, error) {
	user := &conf.User{
		UserId:                          "0",
		UserName:                        "_",
		UserAvatarURL:                   "/appearance/avatar/default-vip.svg", // [FORK-MOD] 内置金色 VIP 头像，由 /appearance/ 静态路由提供
		UserHomeBImgURL:                 "",
		UserTitles:                      []*conf.UserTitle{},
		UserIntro:                       "",
		UserNickname:                    "",
		UserCreateTime:                  "29991231 00:00:00",
		UserSiYuanProExpireTime:         -1,
		UserToken:                       "token",
		UserTokenExpireTime:             "32503593600",
		UserSiYuanRepoSize:              0,
		UserSiYuanPointExchangeRepoSize: 0,
		UserSiYuanAssetSize:             0,
		UserTrafficUpload:               0,
		UserTrafficDownload:             0,
		UserTrafficAPIGet:               0,
		UserTrafficAPIPut:               0,
		UserTrafficTime:                 0,
		UserSiYuanSubscriptionPlan:      0,
		UserSiYuanSubscriptionStatus:    0,
		UserSiYuanSubscriptionType:      1,
		UserSiYuanOneTimePayStatus:      1,
	}

	model.Conf.User = user
	data, _ := gulu.JSON.MarshalJSON(user)
	model.Conf.UserData = util.AESEncrypt(string(data))
	model.Conf.Save()

	return user, nil
}
