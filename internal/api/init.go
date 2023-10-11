package api

import (
	"encoding/json"
	chatModel "im/internal/api/chat/model"
	"im/internal/api/chat/repo"
	discoveryModel "im/internal/api/discover/model"
	friendModel "im/internal/api/friend/model"
	groupModel "im/internal/api/group/model"
	momentsModel "im/internal/api/moments/model"
	operatorMoel "im/internal/api/operator/model"
	settingModel "im/internal/api/setting/model"
	shoppingModel "im/internal/api/shopping/model"
	userModel "im/internal/api/user/model"
	cmsConfigModel "im/internal/cms_api/config/model"
	cmsUserModel "im/internal/cms_api/user/model"
	"im/pkg/common"
	"im/pkg/converter"
	"im/pkg/db"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/tokenbucket"
	"strings"
	"time"
)

var TokenBucket *tokenbucket.TokenBucket

func Init() {

	if err := db.DB.AutoMigrate(new(userModel.User)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(settingModel.SettingConfig)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(groupModel.Group)); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(new(groupModel.GroupMember)); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(new(groupModel.GroupMemberApply)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(&friendModel.Friend{}); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(&friendModel.FriendRequest{}); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(&friendModel.Black{}); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(&friendModel.FriendLabel{}); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(cmsConfigModel.AppVersion)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(cmsConfigModel.InviteCode)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(cmsConfigModel.DefaultFriend)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(cmsUserModel.SignLog)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(userModel.UserConfig)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(chatModel.Conversation)); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(new(chatModel.Message)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(userModel.UserDevice)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(userModel.UserIp)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(userModel.LoginHistory)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(userModel.FavoriteImage)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(momentsModel.MomentsMessage)); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(new(momentsModel.MomentsComments)); err != nil {
		panic(err)
	}
	if err := db.DB.AutoMigrate(new(momentsModel.MomentsCommentsLike)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(momentsModel.MomentsInbox)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(shoppingModel.Shop)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(shoppingModel.ShopTeam)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(discoveryModel.News)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(chatModel.Link)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(momentsModel.ContactsTag)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(operatorMoel.Operator)); err != nil {
		panic(err)
	}

	if err := db.DB.AutoMigrate(new(operatorMoel.OperatorTeam)); err != nil {
		panic(err)
	}
}

func LinkTask() {
	for {

		linkList, err := repo.LinkRepo.GetListByStatus(0)
		if err != nil {
			logger.Sugar.Errorf("内容解析服务数据查询失败：%v \n", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		for _, linkInfo := range linkList {

			TokenBucket.Get()
			go func(linkInfo chatModel.Link) {

				defer TokenBucket.Set()
				if err := DealLinkData(linkInfo); err != nil {
					logger.Sugar.Errorf("内容解析服务数据解析失败：%v \n", err.Error())
				}

				_ = repo.LinkRepo.BatchUpdateStatus(linkInfo.MsgID, 2)
			}(linkInfo)

			_ = repo.LinkRepo.BatchUpdateStatus(linkInfo.MsgID, 1)
		}
	}
}

func init() {

	TokenBucket = tokenbucket.NewTokenBucket(tokenbucket.Options{
		TokenNum: 9999,
	})
}

func DealLinkData(linkInfo chatModel.Link) error {
	options := make(map[string]interface{})
	keyword := ""

	url, err := converter.URL2Text(linkInfo.Link)
	if err != nil {
		return err
	}
	if linkInfo.Title == "" {
		options["title"] = url.Title
	}
	if linkInfo.Desc == "" {
		options["desc"] = url.Description
	}
	if linkInfo.Cover == "" {
		options["cover"] = url.Image
	}
	if linkInfo.Favicon == "" {
		options["favicon"] = url.Favicon
	}
	if linkInfo.Keyword == "" {
		keyword = strings.Join(url.Tags, ",")
		options["keyword"] = keyword
	}
	if linkInfo.Content == "" {
		options["content"] = url.BodyText
	}
	options["status"] = 2
	if err := repo.LinkRepo.Updates(options, linkInfo.MsgID); err != nil {
		return err
	}

	options["msg_id"] = linkInfo.MsgID
	options["link"] = linkInfo.Link

	data, _ := json.Marshal(options)
	switch linkInfo.ConversationType {
	case chatModel.ConversationTypeSingle:
		if err = mqtt.SendMessageToUsers("", common.LinkUpdatePush, data, strings.Split(linkInfo.ConversationID, "_")...); err != nil {
			return err
		}

	case chatModel.ConversationTypeGroup:
		if err = mqtt.SendMessageToGroups("", common.LinkUpdatePush, data, linkInfo.ConversationID); err != nil {
			return err
		}
	}
	return nil
}
