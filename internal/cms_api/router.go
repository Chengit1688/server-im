package cms_api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "im/docs/cms"
	adminService "im/internal/cms_api/admin/service"
	announcementService "im/internal/cms_api/announcement/service"
	chatService "im/internal/cms_api/chat/service"
	cmsuiService "im/internal/cms_api/cmsui/service"
	configService "im/internal/cms_api/config/service"
	dashboardService "im/internal/cms_api/dashboard/service"
	discoverService "im/internal/cms_api/discover/service"
	friendService "im/internal/cms_api/friend/service"
	groupService "im/internal/cms_api/group/service"
	ipblacklistService "im/internal/cms_api/ipblacklist/service"
	ipwhitelistService "im/internal/cms_api/ipwhitelist/service"
	logService "im/internal/cms_api/log/service"
	menuService "im/internal/cms_api/menu/service"
	operationService "im/internal/cms_api/operation/service"
	operatorService "im/internal/cms_api/operator/service"
	roleService "im/internal/cms_api/role/service"
	shoppingService "im/internal/cms_api/shopping/service"
	minioService "im/internal/cms_api/third/service"
	userService "im/internal/cms_api/user/service"
	walletService "im/internal/cms_api/wallet/service"
	"im/pkg/middleware"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(middleware.JWTAuth())

	router.Use(middleware.CmsOperateLog())

	router.Use(middleware.IPByPass())
	baseRouter := router.Group("/cms_api")

	baseRouter.GET("/swagger/*any", RegistryDoc())
	adminRouter := baseRouter.Group("/admin")
	{
		adminRouter.POST("/login", adminService.AdminService.Login)
		adminRouter.POST("/refresh_token", adminService.AdminService.RefreshToken)
		adminRouter.GET("/get_info", adminService.AdminService.GetInfo)
		adminRouter.GET("", adminService.AdminService.List)
		adminRouter.POST("", adminService.AdminService.Add)
		adminRouter.PUT("/update_info", adminService.AdminService.UpdateInfo)
		adminRouter.PUT("/update_password", adminService.AdminService.UpdatePassword)
		adminRouter.DELETE("", adminService.AdminService.Delete)
		adminRouter.GET("/get_google_secret", adminService.AdminService.GetGoogleCodeSecret)
	}
	ipWhiteListRouter := baseRouter.Group("/admin/ip_white_list")
	{
		ipWhiteListRouter.GET("", ipwhitelistService.IPWhiteListService.IPWhiteListPaging)
		ipWhiteListRouter.POST("", ipwhitelistService.IPWhiteListService.IPWhiteListAdd)
		ipWhiteListRouter.PUT("/:id", ipwhitelistService.IPWhiteListService.IPWhiteListUpdate)
		ipWhiteListRouter.DELETE("/:id", ipwhitelistService.IPWhiteListService.IPWhiteistDelete)
	}
	roleRouter := baseRouter.Group("/role")
	{
		roleRouter.DELETE("/:id", roleService.RoleService.RoleDelete)
		roleRouter.PUT("/:id", roleService.RoleService.RoleUpdate)
		roleRouter.GET("/:id", roleService.RoleService.RoleGet)
		roleRouter.POST("", roleService.RoleService.RoleAdd)
		roleRouter.GET("", roleService.RoleService.RoleList)
	}
	settingRouter := baseRouter.Group("/config")
	{
		settingRouter.POST("/invite_add", configService.InviteCodeService.Add)
		settingRouter.POST("/invite_delete", configService.InviteCodeService.Delete)
		settingRouter.POST("/invite_update", configService.InviteCodeService.Update)
		settingRouter.POST("/invite_update_friends", configService.InviteCodeService.UpdateFriend)
		settingRouter.POST("/invite_update_groups", configService.InviteCodeService.UpdateGroup)
		settingRouter.POST("/invite_list", configService.InviteCodeService.GetList)
		settingRouter.POST("/invite_update_status", configService.InviteCodeService.UpdateStatus)

		settingRouter.POST("/version_list", configService.AppVersionService.GetList)
		settingRouter.POST("/version_add", configService.AppVersionService.Add)
		settingRouter.POST("/version_delete", configService.AppVersionService.Delete)
		settingRouter.POST("/version_update", configService.AppVersionService.Update)
		settingRouter.POST("/version_update_status", configService.AppVersionService.UpdateStatus)

		settingRouter.POST("/default_account_add", configService.DefaultAccountService.AddOrUpdateFriend)
		settingRouter.POST("/default_account_update", configService.DefaultAccountService.AddOrUpdateFriend)
		settingRouter.POST("/default_account_delete", configService.DefaultAccountService.DeleteFriend)
		settingRouter.POST("/default_account_list", configService.DefaultAccountService.GetFriendList)

		settingRouter.POST("/shield_add", configService.ShieldService.Add)
		settingRouter.POST("/shield_update", configService.ShieldService.Update)
		settingRouter.POST("/shield_list", configService.ShieldService.GetList)
		settingRouter.POST("/shield_delete", configService.ShieldService.Delete)
	}
	configRouter := baseRouter.Group("/config")
	{
		configRouter.GET("/info", configService.ConfigService.GetCmsConfig)
		configRouter.GET("/login_config", configService.ConfigService.GetLoginConfig)
		configRouter.GET("/register_config", configService.ConfigService.GetRegisterConfig)
		configRouter.POST("/login_config", configService.ConfigService.UpdateLoginConfig)
		configRouter.POST("/register_config", configService.ConfigService.UpdateRegisterConfig)
		configRouter.POST("/sign_config_handle", configService.ConfigService.HandleSignConfig)
		configRouter.GET("/get_sign_config", configService.ConfigService.GetSignConfig)
		configRouter.POST("/google_code_is_open", configService.ConfigService.SetGoogleCodeIsOpen)
		configRouter.POST("/site_ui", cmsuiService.CmsUIService.SetCmsSiteUI)
		configRouter.POST("/jpush", configService.ConfigService.SetJPushConfig)
		configRouter.GET("/jpush", configService.ConfigService.GetJPushConfig)
		configRouter.POST("/feihu", configService.ConfigService.SetFeihuConfig)
		configRouter.GET("/feihu", configService.ConfigService.GetFeihuConfig)
		configRouter.GET("/parameter_config", configService.ConfigService.GetSystemConfig)
		configRouter.POST("/parameter_config_update", configService.ConfigService.UpdateSystemConfig)
		configRouter.GET("/deposite", configService.ConfigService.GetDepositeConfig)
		configRouter.POST("/deposite", configService.ConfigService.SetDepositeConfig)
		configRouter.GET("/withdraw", configService.ConfigService.GetWithdrawConfig)
		configRouter.POST("/withdraw", configService.ConfigService.SetWithdrawConfig)
		configRouter.GET("/about_us", configService.ConfigService.GetAboutUs)
		configRouter.POST("/about_us", configService.ConfigService.SetAboutUs)
		configRouter.GET("/privacy_policy", configService.ConfigService.GetPrivacyPolicy)
		configRouter.POST("/privacy_policy", configService.ConfigService.SetPrivacyPolicy)
		configRouter.GET("/user_agreement", configService.ConfigService.GetUserAgreement)
		configRouter.POST("/user_agreement", configService.ConfigService.SetUserAgreement)
		configRouter.GET("/ip_white_list_is_open", configService.ConfigService.GetIPWhiteListIsOpen)
		configRouter.POST("/ip_white_list_is_open", configService.ConfigService.SetIPWhiteListIsOpen)
		configRouter.POST("/default_is_open", configService.ConfigService.SetDefaultIsOpen)
		configRouter.GET("/default_is_open", configService.ConfigService.GetDefaultIsOpen)

	}
	menuRouter := baseRouter.Group("/menu")
	{
		menuRouter.GET("", menuService.MenuService.MenuList)
	}
	groupRouter := baseRouter.Group("/group")
	{
		groupRouter.POST("/create_group", groupService.GroupService.CreateGroup)
		groupRouter.POST("/remove_group_member", groupService.GroupService.RemoveGroupMember)
		groupRouter.POST("/group_info", groupService.GroupService.GroupInfo)
		groupRouter.POST("/group_member_list", groupService.GroupService.GroupMemberList)
		groupRouter.POST("/group_list", groupService.GroupService.GroupList)
		groupRouter.POST("/group_update", groupService.GroupService.GroupUpdate)
		groupRouter.POST("/information", groupService.GroupService.Information)
		groupRouter.POST("/manage", groupService.GroupService.Manage)
		groupRouter.POST("/remove", groupService.GroupService.Remove)
		groupRouter.POST("/set_admin", groupService.GroupService.GroupSetAdmin)
		groupRouter.POST("/set_owner", groupService.GroupService.GroupSetOwner)
		groupRouter.POST("/set_robot", groupService.GroupService.GroupRobotUpdate)
		groupRouter.POST("/add_group_members", groupService.GroupService.AddGroupMembers)
		groupRouter.POST("/group_merge", groupService.GroupService.GroupMerge)
		groupRouter.POST("/batch_join_group", groupService.GroupService.JoinGroups)
		groupRouter.GET("/search", groupService.GroupService.GroupSearch)
		groupRouter.POST("/mute_group", groupService.GroupService.GroupMuteAll)
	}
	userRouter := baseRouter.Group("/user")
	{
		userRouter.GET("", userService.UserService.UserList)
		userRouter.GET("/export", userService.UserService.UserListExport)
		userRouter.GET("/search", userService.UserService.UserSearch)
		userRouter.GET("/login_history", userService.UserService.LoginHistory)
		userRouter.GET("/user_for_group", userService.UserService.FindUserToGroupList)
		userRouter.POST("", userService.UserService.UserBatchAdd)
		userRouter.GET("/details", userService.UserService.UserDetails)
		userRouter.POST("/updates", userService.UserService.UserInfoUpdate)
		userRouter.POST("/freeze", userService.UserService.FreezeUser)
		userRouter.POST("/unfreeze", userService.UserService.UnFreezeUser)
		userRouter.POST("/set_password", userService.UserService.SetUserPassword)
		userRouter.POST("/sign_log_list", userService.UserService.SignLogList)
		userRouter.POST("/privilege_user_list", userService.UserService.PrivilegeUserList)
		userRouter.POST("/privilege_user_add", userService.UserService.PrivilegeUserAdd)
		userRouter.POST("/privilege_user_remove", userService.UserService.PrivilegeUserRemove)
		userRouter.POST("/privilege_user/set_freeze", userService.UserService.SetPrivilegeUserFreeze)
		userRouter.GET("/disabled/user", userService.UserService.DisabledManagermentUser)
		userRouter.GET("/disabled/device", userService.UserService.DisabledManagermentDevice)
		userRouter.POST("/disabled/device/disable", userService.UserService.DMDeviceLock)
		userRouter.POST("/disabled/device/enable", userService.UserService.DMDeviceUnLock)
		userRouter.GET("/disabled/ip", userService.UserService.DisabledManagermentIP)
		userRouter.POST("/agent_level", userService.UserService.AgentLevel)
		userRouter.POST("/operator_level", userService.UserService.OperatorLevel)
		userRouter.POST("/real_name_auth", userService.UserService.RealNameAuth)
		userRouter.POST("/real_name_list", userService.UserService.RealNameList)

	}
	thirdRouter := baseRouter.Group("/third")
	{
		thirdRouter.POST("/upload", minioService.MinioService.Upload)
		thirdRouter.POST("/upload/v2", minioService.MinioService.UploadV2)
	}
	announcementRouter := baseRouter.Group("/announcement")
	{
		announcementRouter.GET("", announcementService.AnnouncementService.GetAnnouncement)
		announcementRouter.POST("", announcementService.AnnouncementService.UpdateAnnouncement)
	}
	discoverRouter := baseRouter.Group("/discover")
	{
		discoverRouter.GET("", discoverService.DiscoverService.GetDiscoverInfo)
		discoverRouter.POST("", discoverService.DiscoverService.AddDiscover)
		discoverRouter.PUT("/:id", discoverService.DiscoverService.UpdateDiscover)
		discoverRouter.DELETE("/:id", discoverService.DiscoverService.DeleteDiscover)
		discoverRouter.GET("/status", discoverService.DiscoverService.GetDiscoverOpenStatus)
		discoverRouter.POST("/status", discoverService.DiscoverService.SetDiscoverOpenStatus)
		discoverRouter.POST("/news/add", discoverService.NewsService.Add)
		discoverRouter.POST("/news/list", discoverService.NewsService.List)

		discoverRouter.POST("/prize_add", discoverService.DiscoverService.AddPrize)
		discoverRouter.POST("/prize_update", discoverService.DiscoverService.UpdatePrize)
		discoverRouter.POST("/prize_delete", discoverService.DiscoverService.DeletePrize)
		discoverRouter.POST("/prize_list", discoverService.DiscoverService.ListPrize)
		discoverRouter.POST("/list_redeem_prize", discoverService.DiscoverService.RedeemPrizeLog)
		discoverRouter.POST("/set_redeem_prize", discoverService.DiscoverService.SetRedeemPrize)
	}
	newsRouter := baseRouter.Group("/news")
	{
		newsRouter.POST("/add", discoverService.NewsService.Add)
		newsRouter.POST("/update", discoverService.NewsService.Update)
		newsRouter.POST("/delete", discoverService.NewsService.Delete)
		newsRouter.POST("/list", discoverService.NewsService.List)
	}
	friendRouter := baseRouter.Group("/friend")
	{
		friendRouter.POST("/user_friend_list", friendService.FriendService.UserFriendList)
		friendRouter.POST("/user_add_friend", friendService.FriendService.UserAddFriend)
		friendRouter.POST("/user_remove_friend", friendService.FriendService.DeleteFriend)
	}
	operationRouter := baseRouter.Group("/operation")
	{
		operationRouter.POST("/reg_statistics", operationService.OperationService.RegistrationStatistics)
		operationRouter.GET("/online_statistics", dashboardService.DashboardService.OnlineUserDaily)
		operationRouter.GET("/single_msg_statistics", dashboardService.DashboardService.SingleMessageDaily)
		operationRouter.GET("/group_msg_statistics", dashboardService.DashboardService.GroupMessageDaily)
		operationRouter.POST("/invite_code_statistics_list", operationService.OperationService.InviteCodeStatistics)
		operationRouter.POST("/invite_code_statistics_details", operationService.OperationService.InviteCodeStatisticsDetails)
		operationRouter.GET("/suggestion_list", operationService.OperationService.GetSuggestionList)
		operationRouter.GET("/online_users", operationService.OperationService.OnlineUsers)
	}
	ipBlackListRouter := baseRouter.Group("/operation/ip_black_list")
	{
		ipBlackListRouter.GET("", ipblacklistService.IPBlackListService.IPBlackListPaging)
		ipBlackListRouter.POST("", ipblacklistService.IPBlackListService.IPBlackListAdd)
		ipBlackListRouter.POST("/batch", ipblacklistService.IPBlackListService.AddInBatch)
		ipBlackListRouter.PUT("/:id", ipblacklistService.IPBlackListService.IPBlackListUpdate)
		ipBlackListRouter.DELETE("/:id", ipblacklistService.IPBlackListService.IPBlackListDelete)
		ipBlackListRouter.DELETE("/batch", ipblacklistService.IPBlackListService.RemoveInBatch)
	}
	dashboardRouter := baseRouter.Group("/dashboard")
	{
		dashboardRouter.GET("", dashboardService.DashboardService.Info)
	}
	chatRouter := baseRouter.Group("/chat")
	{
		chatRouter.GET("/message_history", chatService.MessageService.HistoryList)
		chatRouter.POST("/message_change", chatService.MessageService.Change)
		chatRouter.POST("/message_clear", chatService.MessageService.Clear)
		chatRouter.POST("/multi_send", chatService.MessageService.MultiSend)
		chatRouter.GET("/multi_send_records", chatService.MessageService.MultiSendList)
	}
	walletRouter := baseRouter.Group("/wallet")
	{
		walletRouter.GET("/billing_records", walletService.WalletService.BillingRecordsList)
		walletRouter.GET("/billing_records/export", walletService.WalletService.BillingRecordsExport)
		walletRouter.POST("/change_amount", walletService.WalletService.WalletChangeAmount)
		walletRouter.GET("/redpack_single_records", walletService.WalletService.RedpackSingleRecordsList)
		walletRouter.POST("/redpack_group_records", walletService.WalletService.RedpackGroupRecordsList)
		walletRouter.GET("/redpack_single_records/export", walletService.WalletService.RedpackSingleRecordsExport)
		walletRouter.POST("/redpack_group_records/export", walletService.WalletService.RedpackGroupRecordsExport)
		walletRouter.GET("/withdraw_records", walletService.WalletService.WithdrawRecordsList)
		walletRouter.GET("/withdraw_records/count_pending", walletService.WalletService.WithdrawRecordsCountPending)
		walletRouter.POST("/withdraw_records/set_status", walletService.WalletService.WithdrawRecordsStatusSet)
		walletRouter.GET("/withdraw_records/:id", walletService.WalletService.WithdrawRecordsDescribe)
		walletRouter.POST("/set_paypasswd", walletService.WalletService.WalletSetPayPass)
	}
	logRouter := baseRouter.Group("/log")
	{
		logRouter.GET("/operate_log", logService.LogService.OperateLogList)
	}
	shoppingRouter := baseRouter.Group("/shopping")
	{
		shoppingRouter.POST("/list", shoppingService.ShoppingService.ShopList)
		shoppingRouter.POST("/approve", shoppingService.ShoppingService.Approve)
		shoppingRouter.POST("/update", shoppingService.ShoppingService.Update)
		shoppingRouter.POST("/member_list", shoppingService.ShoppingService.MemberList)
	}

	operatorRouter := baseRouter.Group("/operator")
	{
		operatorRouter.POST("/list", operatorService.OperatorService.ShopList)
		operatorRouter.POST("/approve", operatorService.OperatorService.Approve)
		operatorRouter.POST("/update", operatorService.OperatorService.Update)
		operatorRouter.POST("/member_list", operatorService.OperatorService.MemberList)
	}

	return router
}

func RegistryDoc() gin.HandlerFunc {
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
