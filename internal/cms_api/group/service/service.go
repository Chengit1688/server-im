package service

import (
	"fmt"
	"im/internal/api/group/model"
	"im/internal/api/group/repo"
	"im/internal/api/group/usecase"
	userModel "im/internal/api/user/model"
	userUseCase "im/internal/api/user/usecase"
	configUseCase "im/internal/cms_api/config/usecase"
	cmsModel "im/internal/cms_api/group/model"
	cmsRepo "im/internal/cms_api/group/repo"
	"im/pkg/code"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/util"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var GroupService = new(groupService)

type groupService struct{}

var BatchJoinChan chan bool = make(chan bool, 1)
var batchOnce sync.Once

func (s *groupService) GetLoginUserId(c *gin.Context) (string, error) {

	userId := c.GetString("o_user_id")
	return userId, nil
}

func (s *groupService) CreateGroup(c *gin.Context) {
	var (
		req  cmsModel.CreateGroupReq
		resp cmsModel.CreateGroupResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	_, err = userUseCase.UserUseCase.GetBaseInfo(req.OwnerId)
	if err != nil {
		http.Failed(c, code.ErrUserIdNotExist)
		return
	}
	groupInfo, err := usecase.GroupUseCase.CreateGroup(req.OperationID, req.Name, req.FaceUrl, req.OwnerId, req.IsTopannocuncement, req.Notification)
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	util.CopyStructFields(&resp, &groupInfo)
	http.Success(c, resp)

}

func (s *groupService) GroupInfo(c *gin.Context) {
	var (
		req  model.GroupInfoReq
		resp cmsModel.GroupInfoResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupId)
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	res := cmsModel.GroupInfo{}
	util.CopyStructFields(&res, groupInfo)
	cmsRepo.GroupRepo.GetGroupOwnerInfo(&res)
	util.CopyStructFields(&resp, res)
	http.Success(c, resp)
}

func (s *groupService) GroupUpdate(c *gin.Context) {
	var (
		req  cmsModel.GroupUpdateReq
		resp cmsModel.GroupUpdateResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if req.MuteAllMember == constant.MuteMemberPeriod {
		times := strings.Split(req.MuteAllPeriod, "-")
		if len(times) != 2 {
			http.Failed(c, code.ErrMutePeriod)
			return
		}
	}

	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupId)
	if err != nil || groupInfo.Status != 1 {
		http.Failed(c, code.ErrGroupNotExist)
		return
	}
	updateInfo := model.Group{}
	util.CopyStructFields(&updateInfo, req)
	groupInfo, err = usecase.GroupUseCase.UpdateGroup(req.OperationID, req.GroupId, &updateInfo, nil, "")
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	util.CopyStructFields(&resp, groupInfo)
	http.Success(c, resp)

}

func (s *groupService) GroupRobotUpdate(c *gin.Context) {
	var (
		req  cmsModel.GroupRobotUpdateReq
		resp cmsModel.GroupRobotUpdateResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupId)
	if err != nil || groupInfo.Status != 1 {
		http.Failed(c, code.ErrGroupNotExist)
		return
	}
	updateInfo := map[string]interface{}{
		"robot_total": req.RobotTotal,
	}

	err = db.Update(&model.Group{}, req.GroupId, updateInfo)
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}

	repo.GroupCache.UpGroupInfoCache(req.GroupId)
	repo.GroupRepo.GroupChangeMsg(req.OperationID, req.GroupId, 3)

	http.Success(c, resp)
}

func (s *groupService) Information(c *gin.Context) {
	var (
		req  model.GroupInformationReq
		resp model.GroupInformationResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if req.Name != nil && *req.Name == "" {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "name is empty")
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if resp.Group, err = usecase.GroupUseCase.UpdateInformation(req.OperationID, "", req.GroupID, &req.GroupInformationInfo); err != nil {
		http.Failed(c, err)
		return
	}
	http.Success(c, resp)
}

func (s *groupService) Manage(c *gin.Context) {
	var (
		req  model.GroupManageReq
		resp model.GroupManageResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if resp.Group, err = usecase.GroupUseCase.UpdateManage(req.OperationID, req.GroupID, &req.GroupManageInfo); err != nil {
		http.Failed(c, err)
		return
	}
	http.Success(c, resp)
}

func (s *groupService) Remove(c *gin.Context) {
	var (
		req  model.GroupRemoveReq
		resp model.GroupRemoveResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if err = usecase.GroupUseCase.Remove(req.OperationID, req.GroupID); err != nil {
		http.Failed(c, err)
		return
	}

	resp.GroupID = req.GroupID
	http.Success(c, resp)
}

func (s *groupService) GroupMemberList(c *gin.Context) {
	var (
		req  cmsModel.GroupMemberListReq
		resp cmsModel.GroupMemberListResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	groupId := req.GroupId
	members := []model.GroupMember{}
	total := int64(0)
	where := map[string]interface{}{
		"group_id": groupId,
		"status":   1,
	}

	if req.SearchKey != "" {

		where["user_id"] = fmt.Sprintf("?%s", req.SearchKey)
	}
	db.Find(model.GroupMember{}, where, "id desc", req.Page, req.PageSize, &total, &members)
	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	for _, v := range members {
		temp := cmsModel.GroupMemberInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	cmsRepo.GroupRepo.GetGroupMemberUserInfo(&resp.List)

	http.Success(c, resp)
}

func (s *groupService) RemoveGroupMember(c *gin.Context) {
	var (
		req  cmsModel.RemoveGroupMemberReq
		resp cmsModel.RemoveGroupMemberResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	err = usecase.GroupUseCase.GroupRemoveMember(req.OperationID, req.GroupId, req.UserIdList, "", model.ReasonTypeKick)
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	http.Success(c, resp)

}

func (s *groupService) GroupSetAdmin(c *gin.Context) {
	var (
		req  model.GroupSetAdminReq
		resp model.GroupSetAdminResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	memberInfo := model.GroupMember{}
	err = db.Info(&memberInfo, model.GroupMember{
		GroupId: req.GroupId,
		UserId:  req.UserId,
		Status:  1,
	})

	if err != nil {
		http.Failed(c, code.ErrNotInGroup)
		return
	}

	if memberInfo.Role == model.RoleTypeOwner {
		http.Failed(c, code.ErrUserPermissions)
		return
	}

	err = usecase.GroupUseCase.SetGroupAdmin(req.OperationID, req.GroupId, req.UserId, req.Status, "")
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	http.Success(c, resp)

}

func (s *groupService) GroupSetOwner(c *gin.Context) {
	var (
		req  model.GroupSetOwnerReq
		resp model.GroupSetOwnerResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	memberInfo := model.GroupMember{}
	err = db.Info(&memberInfo, model.GroupMember{
		GroupId: req.GroupId,
		UserId:  req.UserId,
		Status:  1,
	})

	if err != nil {
		http.Failed(c, code.ErrNotInGroup)
		return
	}

	if memberInfo.Role == model.RoleTypeOwner {
		http.Success(c, resp)
		return
	}
	err = usecase.GroupUseCase.SetGroupOwner(req.OperationID, req.GroupId, req.UserId, "")
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	http.Success(c, resp)

}

func (s *groupService) AddGroupMembers(c *gin.Context) {
	var (
		req  cmsModel.AddGroupMembersReq
		resp cmsModel.AddGroupMembersResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupId)
	if err != nil || groupInfo.Status != 1 {
		http.Failed(c, code.ErrGroupNotExist)
		return
	}

	setting, err := configUseCase.ConfigUseCase.GetParameterConfig()
	if err == nil && setting.GroupLimit > 0 {
		logger.Sugar.Debug(req.OperationID, "入群人数检测", setting.GroupLimit, groupInfo.MembersTotal, len(req.UserIdList))
		if groupInfo.MembersTotal+len(req.UserIdList) > int(setting.GroupLimit) {
			http.Failed(c, code.ErrGroupMemberOutMax)
			return
		}
	}
	_, err = usecase.GroupUseCase.BatchJoinGroup(req.OperationID, req.GroupId, req.UserIdList, "")
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	http.Success(c, resp)

}

func (s *groupService) GroupMerge(c *gin.Context) {
	var (
		req  cmsModel.GroupMergeReq
		resp cmsModel.GroupMergeResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	fromUserList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
		"group_id": req.FromGroupId,
		"status":   1,
	}, "user_id")
	if err != nil {
		http.Failed(c, code.ErrGroupNotMember)
		return
	}
	fromUserIdList := []string{}
	for _, v := range fromUserList {
		fromUserIdList = append(fromUserIdList, v.(string))
	}

	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.ToGroupId)
	if err != nil || groupInfo.Status != 1 {
		http.Failed(c, code.ErrGroupNotExist)
		return
	}

	setting, err := configUseCase.ConfigUseCase.GetParameterConfig()
	if err == nil && setting.GroupLimit > 0 {
		if groupInfo.MembersTotal+len(fromUserList) > int(setting.GroupLimit) {
			http.Failed(c, code.ErrGroupMemberOutMax)
			return
		}
	}
	_, err = usecase.GroupUseCase.BatchJoinGroup(req.OperationID, req.ToGroupId, fromUserIdList, "")
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}
	toGroup := model.Group{}
	db.Info(&toGroup, req.ToGroupId)
	toGroupInfo := cmsModel.GroupInfo{}
	util.CopyStructFields(&toGroupInfo, toGroup)
	cmsRepo.GroupRepo.GetGroupOwnerInfo(&toGroupInfo)
	util.CopyStructFields(&resp, toGroupInfo)

	http.Success(c, resp)

}

func (s *groupService) GroupList(c *gin.Context) {
	var (
		req  cmsModel.GroupListReq
		resp cmsModel.GroupListResq
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	groups := []model.Group{}
	total := int64(0)
	queryWhere := map[string]interface{}{
		"status": 1,
	}
	if req.GroupName != "" {
		queryWhere["name"] = fmt.Sprintf("?%s", req.GroupName)
	}

	if req.OwnerName != "" {

		uidList, err := db.CloumnList(&userModel.User{}, map[string]interface{}{
			"nick_name": "?" + req.OwnerName,
		}, "user_id")
		if err == nil {
			uidStringList := []string{}
			for _, v := range uidList {
				uidStringList = append(uidStringList, v.(string))
			}
			groupIdList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
				"user_id": uidStringList,
				"role":    model.RoleTypeOwner,
				"status":  1,
			}, "group_id")
			if err == nil {
				groupIdStringList := []string{}
				for _, v := range groupIdList {
					groupIdStringList = append(groupIdStringList, v.(string))
				}
				queryWhere["group_id"] = groupIdStringList
			}
		}
	}
	if req.IsDefault != 0 {
		queryWhere["is_default"] = req.IsDefault
	}
	db.Find(model.Group{}, queryWhere, "id desc", req.Page, req.PageSize, &total, &groups)
	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	for _, v := range groups {
		temp := cmsModel.GroupInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	cmsRepo.GroupRepo.GetGroupsOwnerInfo(&resp.List)
	http.Success(c, resp)
}

func (s *groupService) JoinGroups(c *gin.Context) {
	var (
		req  cmsModel.JoinGroupsReq
		resp cmsModel.JoinGroupsResp
		err  error
	)
	logger.Sugar.Debug("收到请求")
	batchOnce.Do(func() {

		BatchJoinChan <- true
		logger.Sugar.Debug("初始化通道")
	})
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	select {
	case <-BatchJoinChan:

		logger.Sugar.Debug("进入业务")
		break
	default:
		http.Failed(c, code.ErrFuncRunning)
		return
	}
	_, err = userUseCase.UserUseCase.GetBaseInfo(req.UserId)
	if err != nil {
		http.Failed(c, code.ErrUserNotFound)
		return
	}

	groupIdList, _ := db.CloumnList(&model.Group{}, model.Group{
		Status: 1,
	}, "group_id")
	if len(groupIdList) == 0 {
		http.Success(c, resp)
		return
	}

	hadGroupList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
		"user_id": req.UserId,
		"status":  1,
	}, "group_id")
	hadGroupIdMap := map[string]bool{}
	if err == nil {
		for _, v := range hadGroupList {
			groupid := v.(string)
			hadGroupIdMap[groupid] = true
		}
	}
	newGroupIdList := []string{}
	for _, v := range groupIdList {
		if had := hadGroupIdMap[v.(string)]; !had {
			newGroupIdList = append(newGroupIdList, v.(string))
		}
	}

	go func() {
		defer func() { BatchJoinChan <- true }()
		for i := 0; i < len(newGroupIdList); i++ {
			usecase.GroupUseCase.JoinGroup(req.OperationID, newGroupIdList[i], req.UserId)
			logger.Sugar.Debug("执行入群完毕", newGroupIdList[i])
			if i%10 == 0 {
				time.Sleep(1 * time.Second)
			}
		}
	}()

	http.Success(c, resp)
}

func (s *groupService) GroupSearch(c *gin.Context) {
	var (
		req  cmsModel.GroupSearchReq
		resp cmsModel.GroupSearchResq
		err  error
	)

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	groups, count, err := cmsRepo.GroupRepo.GroupSearch(req.Search)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "GroupSearch error:", err)
		http.Failed(c, code.ErrDB)
		return
	}
	resp.Count = int(count)
	for _, v := range groups {
		temp := cmsModel.GroupInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	cmsRepo.GroupRepo.GetGroupsOwnerInfo(&resp.List)
	http.Success(c, resp)
}

func (s *groupService) GroupMuteAll(c *gin.Context) {
	var (
		req  model.GroupMuteAllReq
		resp model.GroupMuteAllResp
		err  error
	)

	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if req.MuteAllMember == constant.MuteMemberPeriod {
		times := strings.Split(req.MuteAllPeriod, "-")
		if len(times) != 2 {
			http.Failed(c, code.ErrMutePeriod)
			return
		}
	}

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, code.ErrGroupNotExist)
		return
	}

	err = usecase.GroupUseCase.UpdateGroupMuteInfo(req)
	if err != nil {
		http.Failed(c, code.ErrFailRequest)
		return
	}

	http.Success(c, resp)
}
