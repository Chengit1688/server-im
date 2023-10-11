package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "im/docs/api"
	chatService "im/internal/api/chat/service"
	newsService "im/internal/api/discover/service"
	friendService "im/internal/api/friend/service"
	groupService "im/internal/api/group/service"
	momentsService "im/internal/api/moments/service"
	operatorService "im/internal/api/operator/service"
	settingService "im/internal/api/setting/service"
	shoppingService "im/internal/api/shopping/service"
	userService "im/internal/api/user/service"
	walletService "im/internal/api/wallet/service"
	announcementService "im/internal/cms_api/announcement/service"
	configService "im/internal/cms_api/config/service"
	minioService "im/internal/cms_api/third/service"
	"im/pkg/middleware"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(gin.Recovery())

	router.Use(middleware.OAuth())

	router.Use(middleware.IPBlocker())
	baseRouter := router.Group("/api")

	baseRouter.Use(middleware.DeviceBlocker())

	baseRouter.GET("/swagger/*any", RegistryDoc())

	userRouter := baseRouter.Group("/user")
	{
		userRouter.POST("/register", userService.UserService.Register)
		userRouter.POST("/login", userService.AuthService.Login)
		userRouter.POST("/update_info", userService.UserService.UpdateInfo)
		userRouter.POST("/sign_today", userService.UserService.SignToday)
		userRouter.GET("/get_user_sign_info", userService.UserService.SignList)
		userRouter.GET("/user_sign_info_week", userService.UserService.SignListWeek)
		userRouter.GET("/info", userService.UserService.Info)
		userRouter.POST("/password_secure", userService.UserService.PasswordSecure)
		userRouter.POST("/get_user_info", userService.UserService.GetUserInfo)
		userRouter.POST("/get_user_base_info", userService.UserService.GetUserBaseInfo)
		userRouter.POST("/forgot_password", userService.UserService.ForgotPassword)
		userRouter.POST("/verify_code", userService.UserService.VerifyPhoneCode)
		userRouter.GET("/device_list", userService.UserService.DeviceList)
		userRouter.POST("/server_version", userService.UserService.GetServerVersion)
		userRouter.POST("/user_config_handle", userService.UserService.UserConfigHandle)
		userRouter.GET("/get_user_config", userService.UserService.GetUserConfig)
		userRouter.POST("/get_user_online_status", userService.UserService.GetUserOnlineStatus)
		userRouter.POST("/suggestion", userService.UserService.Suggestion)
		userRouter.POST("/bind_phone", userService.UserService.BindPhone)
		userRouter.GET("/favorite/image", userService.UserService.FavoriteImagePaging)
		userRouter.POST("/favorite/image_add", userService.UserService.FavoriteImageAdd)
		userRouter.POST("/favorite/image_del", userService.UserService.FavoriteImageDel)
		userRouter.POST("/set_privacy", userService.UserService.SetPrivacy)
		userRouter.POST("/get_privacy", userService.UserService.GetPrivacy)
		userRouter.POST("/real_name", userService.UserService.RealName)
		userRouter.POST("/real_name_info", userService.UserService.RealNameInfo)

		userRouter.POST("/prize_list", userService.UserService.PrizeList)
		userRouter.POST("/redeem_prize", userService.UserService.RedeemPrize)
		userRouter.POST("/redeem_prize_list", userService.UserService.RedeemPrizeList)
	}

	settingRouter := baseRouter.Group("/setting")
	{
		settingRouter.GET("/config", settingService.SettingService.RegisterAndLoginConfig)
		settingRouter.GET("/about_us", settingService.SettingService.AboutUs)
		settingRouter.GET("/privacy_policy", settingService.SettingService.PrivacyPolicy)
		settingRouter.GET("/user_agreement", settingService.SettingService.UserAgreement)
		settingRouter.GET("/version", settingService.SettingService.Version)

		settingRouter.POST("/sms", settingService.SettingService.SmsCode)
		settingRouter.POST("/shield_list", settingService.SettingService.GetShieldList)
		settingRouter.POST("/domain_list", settingService.SettingService.DomainList)
	}

	captchaRouter := baseRouter.Group("/captcha")
	{
		captchaRouter.POST("/get", userService.CaptchaService.GetCaptcha)
		captchaRouter.POST("/check", userService.CaptchaService.CheckCaptcha)
	}
	chatRouter := baseRouter.Group("/chat")
	{
		chatRouter.GET("/conversation_list", chatService.ConversationService.List)
		chatRouter.POST("/conversation_ack_seq", chatService.ConversationService.AckSeq)
		chatRouter.POST("/message_send", chatService.MessageService.Send)
		chatRouter.POST("/message_multi_send", chatService.MessageService.MultiSend)
		chatRouter.POST("/message_forward", chatService.MessageService.Forward)
		chatRouter.POST("/message_change", chatService.MessageService.Change)
		chatRouter.GET("/message_pull", chatService.MessageService.Pull)
		chatRouter.POST("/message_pull_v2", chatService.MessageService.PullV2)
		chatRouter.POST("/message_clear", chatService.MessageService.Clear)
		chatRouter.GET("/rtc_info", chatService.RTCService.RTCInfo)
		chatRouter.POST("/rtc", chatService.RTCService.RTC)
		chatRouter.POST("/rtc_operate", chatService.RTCService.RTCOperate)
		chatRouter.POST("/rtc_update", chatService.RTCService.RTCUpdate)
	}

	groupRouter := baseRouter.Group("/group")
	{
		groupRouter.POST("/search", groupService.GroupService.Search)
		groupRouter.POST("/create_group", groupService.GroupService.CreateGroup)
		groupRouter.POST("/join_group_apply", groupService.GroupService.JoinGroupApply)
		groupRouter.POST("/group_apply_list", groupService.GroupService.JoinApplyList)
		groupRouter.POST("/join_group_verify", groupService.GroupService.JoinGroupVerify)
		groupRouter.POST("/quit_group", groupService.GroupService.QuitGroup)
		groupRouter.POST("/remove_group_member", groupService.GroupService.RemoveGroupMember)
		groupRouter.POST("/invite_group_member", groupService.GroupService.InviteGroupMember)
		groupRouter.POST("/group_info", groupService.GroupService.GroupInfo)
		groupRouter.POST("/group_member_list", groupService.GroupService.GroupMemberList)
		groupRouter.POST("/update_group_member", groupService.GroupService.UpdateGroupMember)
		groupRouter.POST("/joind_group_list", groupService.GroupService.JoindGroupList)
		groupRouter.POST("/my_group_list", groupService.GroupService.MyGroupList)
		groupRouter.POST("/information", groupService.GroupService.Information)
		groupRouter.POST("/manage", groupService.GroupService.Manage)
		groupRouter.POST("/remove", groupService.GroupService.Remove)
		groupRouter.POST("/group_update", groupService.GroupService.GroupUpdate)
		groupRouter.POST("/group_update_avatar", groupService.GroupService.GroupUpdateAvatar)
		groupRouter.POST("/group_sync", groupService.GroupService.GroupSync)
		groupRouter.POST("/group_info_sync", groupService.GroupService.GroupListSync)
		groupRouter.POST("/set_admin", groupService.GroupService.GroupSetAdmin)
		groupRouter.POST("/set_owner", groupService.GroupService.GroupSetOwner)
		groupRouter.POST("/get_owner_admin", groupService.GroupService.GetOwnerAdmin)
		groupRouter.POST("/mute_member", groupService.GroupService.GroupMuteMember)
		groupRouter.POST("/mute_group", groupService.GroupService.GroupMuteAll)
		groupRouter.POST("/get_my_group_max_seq", groupService.GroupService.GetMyGroupMaxSeq)
		groupRouter.POST("/face2face_invite", groupService.GroupService.Face2FaceInvite)
		groupRouter.POST("/face2face_add", groupService.GroupService.Face2FaceAdd)
	}

	friendRouter := baseRouter.Group("/friend")
	{
		friendRouter.POST("/search", friendService.FriendService.Search)
		friendRouter.POST("/add_friend", friendService.FriendService.AddFriend)
		friendRouter.POST("/add_friend_ack", friendService.FriendService.AddFriendAck)
		friendRouter.POST("/delete_friend", friendService.FriendService.DeleteFriend)
		friendRouter.POST("/set_friend_remark", friendService.FriendService.SetFriendRemark)
		friendRouter.POST("/get_friend_remark", friendService.FriendService.GetFriendRemark)
		friendRouter.POST("/check_friend_remark", friendService.FriendService.CheckFriendRemark)
		friendRouter.POST("/get_friend", friendService.FriendService.GetFriendsInfo)
		friendRouter.POST("/get_friend_list", friendService.FriendService.GetFriendList)
		friendRouter.POST("/get_friend_apply_list", friendService.FriendService.GetFriendApplyList)
		friendRouter.POST("/get_friend_apply_all", friendService.FriendService.GetFriendApplyAll)
		friendRouter.POST("/get_self_friend_apply_list", friendService.FriendService.GetSelfFriendApplyList)
		friendRouter.POST("/get_friend_max_seq", friendService.FriendService.GetFriendMaxSeq)
		friendRouter.POST("/search_user", friendService.FriendService.SearchFriend)
		friendRouter.POST("/friend_sync", friendService.FriendService.FriendListSync)
		friendRouter.POST("/create_friend_label", friendService.FriendService.CreateFriendLabel)
		friendRouter.POST("/delete_friend_label", friendService.FriendService.DeleteFriendLabel)
		friendRouter.POST("/update_friend_label", friendService.FriendService.UpdateFriendLabel)
		friendRouter.POST("/get_friend_label", friendService.FriendService.GetFriendLabel)
		friendRouter.POST("/change_friend_label", friendService.FriendService.ChangeFriendLabel)
		friendRouter.POST("/get_friends_msg_max_seq", friendService.FriendService.GetFriendsMsgMaxSeq)
		friendRouter.POST("/add_black", friendService.FriendService.AddBlack)
		friendRouter.POST("/remove_black", friendService.FriendService.RemoveBlack)
		friendRouter.POST("/get_black_list", friendService.FriendService.GetBlackList)
		friendRouter.POST("/get_black_listv2", friendService.FriendService.GetBlackListV2)
		friendRouter.POST("/is_friend", friendService.FriendService.IsFriend)
	}
	thirdRouter := baseRouter.Group("/third")
	{
		thirdRouter.POST("/upload", minioService.MinioService.Upload)
		thirdRouter.POST("/upload/v2", minioService.MinioService.UploadV2)
		thirdRouter.GET("/sts", minioService.MinioService.GetSTS)
		thirdRouter.GET("/get_file_url", minioService.MinioService.GetFileUrl)

	}
	announcementRouter := baseRouter.Group("/announcement")
	{
		announcementRouter.GET("", announcementService.AnnouncementService.GetAnnouncement)
	}
	discoverRouter := baseRouter.Group("/discover")
	{
		discoverRouter.GET("", settingService.SettingService.GetDiscoverInfo)
	}
	walletRouter := baseRouter.Group("/wallet")
	{
		walletRouter.GET("", walletService.WalletService.GetWalletInfo)
		walletRouter.GET("/withdraw_config", configService.ConfigService.GetWithdrawConfig)
		walletRouter.POST("/withdraw", walletService.WalletService.WithdrawCommit)
		walletRouter.GET("/redpack_single", walletService.WalletService.RedpackSingleGetInfo)
		walletRouter.POST("/redpack_single/send", walletService.WalletService.RedpackSingleSend)
		walletRouter.POST("/redpack_single/recv", walletService.WalletService.RedpackSingleRecv)
		walletRouter.POST("/redpack_group", walletService.WalletService.RedpackGroupGetInfo)
		walletRouter.POST("/redpack_group/send", walletService.WalletService.RedpackGroupSend)
		walletRouter.POST("/redpack_group/recv", walletService.WalletService.RedpackGroupRecv)
		walletRouter.GET("/billing_records", walletService.WalletService.BillingRecordsList)
		walletRouter.POST("/set_paypasswd", walletService.WalletService.WalletSetPayPass)
	}

	momentsRouter := baseRouter.Group("/moments")
	{
		momentsRouter.POST("/add_moments_message", momentsService.MomentsService.AddMomentsMessage)
		momentsRouter.POST("/get_moments_message", momentsService.MomentsService.GetMomentsMessage)
		momentsRouter.POST("/del_moments_message", momentsService.MomentsService.DelMomentsMessage)
		momentsRouter.POST("/moments_comments_add", momentsService.MomentsService.MomentsCommentsAdd)
		momentsRouter.POST("/del_moments_comments", momentsService.MomentsService.DelMomentsComments)
		momentsRouter.POST("/like_moments", momentsService.MomentsService.LikeMoments)
		momentsRouter.POST("/moments_comments_list", momentsService.MomentsService.MomentsCommentsList)
		momentsRouter.POST("/moments_detail", momentsService.MomentsService.MomentsDetail)
	}
	contactsRouter := baseRouter.Group("/contacts")
	{
		contactsRouter.POST("/tag_add", momentsService.ContactsService.AddTag)
		contactsRouter.POST("/tag_add_friend", momentsService.ContactsService.AddFriendTag)
		contactsRouter.POST("/tag_check_friend", momentsService.ContactsService.CheckFriendTag)
		contactsRouter.POST("/tag_fetch_friend", momentsService.ContactsService.FetchFriendTag)
		contactsRouter.POST("/tag_update", momentsService.ContactsService.UpdateTag)
		contactsRouter.POST("/tag_delete", momentsService.ContactsService.DeleteTag)
		contactsRouter.POST("/tag_list", momentsService.ContactsService.ListTag)
		contactsRouter.POST("/tag_detail", momentsService.ContactsService.TagDetail)
	}
	shoppingRouter := baseRouter.Group("/shopping")
	{
		shoppingRouter.POST("/search", shoppingService.ShoppingService.Search)
		shoppingRouter.POST("/apply_for", shoppingService.ShoppingService.ApplyFor)
		shoppingRouter.POST("/update", shoppingService.ShoppingService.Update)
		shoppingRouter.POST("/detail", shoppingService.ShoppingService.Detail)
		shoppingRouter.POST("/team_member_list", shoppingService.ShoppingService.TeamMemberList)
		shoppingRouter.POST("/join_team", shoppingService.ShoppingService.JoinTeam)
		shoppingRouter.POST("/remove_team", shoppingService.ShoppingService.RemoveTeam)
		shoppingRouter.POST("/team_member_info", shoppingService.ShoppingService.TeamMemberInfo)
		shoppingRouter.POST("/team_leader_info", shoppingService.ShoppingService.JoinTeamInfo)
	}

	operatorRouter := baseRouter.Group("/operator")
	{
		operatorRouter.POST("/search", operatorService.OperatorService.Search)
		operatorRouter.POST("/apply_for", operatorService.OperatorService.ApplyFor)
		operatorRouter.POST("/update", operatorService.OperatorService.Update)
		operatorRouter.POST("/detail", operatorService.OperatorService.Detail)
		operatorRouter.POST("/team_member_list", operatorService.OperatorService.TeamMemberList)
		operatorRouter.POST("/join_team", operatorService.OperatorService.JoinTeam)
		operatorRouter.POST("/remove_team", operatorService.OperatorService.RemoveTeam)
		operatorRouter.POST("/team_member_info", operatorService.OperatorService.TeamMemberInfo)
		operatorRouter.POST("/team_leader_info", operatorService.OperatorService.JoinTeamInfo)
	}

	newsRouter := baseRouter.Group("/news")
	{
		newsRouter.POST("/list", newsService.NewsService.List)
		newsRouter.POST("/detail", newsService.NewsService.Detail)
	}

	return router
}

func RegistryDoc() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
