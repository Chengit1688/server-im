package service

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	chatModel "im/internal/api/chat/model"
	chatUseCase "im/internal/api/chat/usecase"
	"im/internal/api/group/model"
	"im/internal/api/group/repo"
	"im/internal/api/group/usecase"
	permissionUseCase "im/internal/api/permission/usecase"
	userRepo "im/internal/api/user/repo"
	"im/pkg/common"
	"im/pkg/common/constant"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/response"
	"im/pkg/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var GroupService = new(groupService)

type groupService struct{}

func (s *groupService) GetLoginUserId(c *gin.Context) (string, error) {

	user_id := c.GetString("user_id")
	return user_id, nil
}

func (s *groupService) Search(c *gin.Context) {
	var (
		req  model.GroupSearchReq
		resp model.GroupSearchResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	req.Check()
	resp.Pagination = req.Pagination
	userID := c.GetString("user_id")

	if resp.List, resp.Count, err = usecase.GroupUseCase.Search(userID, req.Keyword, req.Offset, req.Limit); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("search error, error: %v", err))
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	if len(resp.List) == 0 {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	http.Success(c, resp)
}

func (s *groupService) CreateGroup(c *gin.Context) {
	var (
		req  model.CreateGroupReq
		resp model.CreateGroupResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)

	if err = permissionUseCase.PermissionUseCase.CheckCreateGroupPermission(req.OperationID, loginUserId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	newGroup, err := usecase.GroupUseCase.CreateGroup(req.OperationID, req.Name, req.FaceUrl, loginUserId, 2)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}

	util.CopyStructFields(&resp, newGroup)
	resp.Role = model.RoleTypeOwner
	http.Success(c, resp)
}

func (s *groupService) Face2FaceInvite(c *gin.Context) {
	var (
		req    model.Face2FaceAddReq
		resp   model.Face2FaceInviteResp
		result map[string]string
		err    error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	userKey := "user_" + loginUserId
	if _, err = db.RedisCli.HSet(c, req.GroupNumber, userKey, loginUserId).Result(); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "err", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}

	ttl, err1 := db.RedisCli.TTL(c, req.GroupNumber).Result()
	if err1 != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "err", err1)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}
	if ttl == 0 {
		if err = db.RedisCli.Expire(c, req.GroupNumber, model.InviteGroupExpire).Err(); err != nil {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "err", err)
			http.Failed(c, response.GetError(response.ErrUnknown, lang))
			return
		}
	}
	if result, err = db.RedisCli.HGetAll(c, req.GroupNumber).Result(); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "err", err)
		http.Failed(c, response.GetError(response.ErrUnknown, lang))
		return
	}

	for _, val := range result {
		userInfo, err2 := userRepo.UserRepo.GetByUserID(userRepo.WhereOption{
			UserId: val,
		})
		if err2 != nil {
			logger.Sugar.Warnw(req.OperationID, util.GetSelfFuncName(), "err", err2)
			continue
		}
		resp.Users = append(resp.Users, model.UserInfo{
			UserID:     userInfo.UserID,
			Account:    userInfo.Account,
			NickName:   userInfo.NickName,
			FaceURL:    userInfo.FaceURL,
			BigFaceURL: userInfo.BigFaceURL,
		})
	}

	http.Success(c, resp)
}

func (s *groupService) Face2FaceAdd(c *gin.Context) {
	var (
		req           model.Face2FaceAddReq
		resp          model.Face2FaceAddResp
		err           error
		isGenLock, ok bool
		groupID       string
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	userKey := "user_" + loginUserId
	if ok, err = db.RedisCli.HExists(c, req.GroupNumber, userKey).Result(); err != nil && err != redis.Nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("HExists error, key: %s, error: %v", req.GroupNumber, err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if !ok {
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	groupKey := "group_" + req.GroupNumber
	if groupID, err = db.RedisCli.Get(c, groupKey).Result(); err != nil && err != redis.Nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("redis get error, key: %s, error: %v", groupKey, err))
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	joinFunc := func() error {
		if usecase.GroupMemberUseCase.CheckMember(groupID, loginUserId) {
			return nil
		}
		err = usecase.GroupUseCase.JoinGroup(req.OperationID, groupID, loginUserId)
		return err
	}
	tempGroupKey := "genGroup_" + req.GroupNumber
	if isGenLock, err = db.RedisCli.SetNX(c, tempGroupKey, loginUserId, model.InviteGroupExpire).Result(); err != nil && err != redis.Nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("redis setnx error, key: %s, error: %v", req.GroupNumber, err))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	defer func() {
		_, _ = db.RedisCli.HDel(c, req.GroupNumber, userKey).Result()
		if isGenLock {
			_, _ = db.RedisCli.Del(c, tempGroupKey).Result()
		}
	}()
	if false == isGenLock {
		time.Sleep(2 * time.Second)
		if joinFunc() != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, response.GetError(response.ErrFailRequest, lang))
			return
		}
		http.Success(c)
		return
	}
	if groupID == "" {
		newGroup, err1 := usecase.GroupUseCase.CreateGroup(req.OperationID, req.GroupNumber, "", loginUserId, 2)
		if err1 != nil {
			http.Failed(c, response.GetError(response.ErrFailRequest, lang))
			return
		}
		if _, err = db.RedisCli.Set(c, groupKey, newGroup.GroupId, model.InviteGroupExpire).Result(); err != nil && err != redis.Nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("redis set error, key: %s, error: %v", groupKey, err))
			http.Failed(c, response.GetError(response.ErrFailRequest, lang))
			return
		}
		groupID = newGroup.GroupId
	} else {
		if joinFunc() != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
			http.Failed(c, response.GetError(response.ErrFailRequest, lang))
			return
		}
	}
	groupInfo := model.Group{}
	err = db.Info(&groupInfo, groupID)
	if err != nil || groupInfo.Status == 2 {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	_ = util.CopyStructFields(&resp, groupInfo)
	http.Success(c, resp)
}

func (s *groupService) JoinGroupApply(c *gin.Context) {
	var (
		req  model.JoinGroupApplyReq
		resp model.JoinGroupApplyResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)

	err = db.Info(&model.GroupMemberApply{}, map[string]interface{}{
		"group_id": req.GroupId,
		"user_id":  loginUserId,
		"status":   0,
	})
	if err == nil {
		http.Success(c, resp)
		return
	}

	groupInfo := model.Group{}
	if err = db.Info(&groupInfo, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	if err = db.Info(&model.GroupMember{}, model.GroupMember{
		GroupId: req.GroupId,
		UserId:  loginUserId,
		Status:  1,
	}); err == nil {
		http.Failed(c, response.GetError(response.ErrAlreadyInGroup, lang))
		return
	}

	var needVerify bool
	if needVerify, err = permissionUseCase.PermissionUseCase.CheckJoinGroupPermission(req.OperationID, loginUserId, req.GroupId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	if !needVerify {
		err = usecase.GroupUseCase.JoinGroup(req.OperationID, req.GroupId, loginUserId)
		if err != nil {
			http.Failed(c, response.GetError(response.ErrFailRequest, lang))
			return
		}
		http.Success(c, resp)
		return
	}
	apply := model.GroupMemberApply{
		GroupId:    req.GroupId,
		UserId:     loginUserId,
		Remark:     req.Remark,
		CreateTime: time.Now().Unix(),
	}

	err = db.Insert(&apply)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}

	repo.GroupRepo.GroupMemberApplyChangeMsg(req.OperationID, apply.Id, 1)
	http.Success(c, resp)
}

func (s *groupService) JoinGroupVerify(c *gin.Context) {
	var (
		req  model.JoinGroupVerifyReq
		resp model.JoinGroupVerifyResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	applyInfo := model.GroupMemberApply{}
	err = db.Info(&applyInfo, req.ApplyId)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrApplyNotFound, lang))
		return
	}

	if applyInfo.Status != 0 {
		http.Failed(c, response.GetError(response.ErrApplyDone, lang))
		return
	}

	loginMemberInfo := model.GroupMember{}
	err = db.Info(&loginMemberInfo, &model.GroupMember{
		GroupId: applyInfo.GroupId,
		UserId:  loginUserId,
		Status:  1,
	})
	if err != nil || loginMemberInfo.Role == model.RoleTypeUser {
		http.Failed(c, response.GetError(response.ErrUserPermissions, lang))
		return
	}

	tx := db.DB.Begin()
	err = db.UpdateTx(tx, &model.GroupMemberApply{}, model.GroupMemberApply{
		Id: req.ApplyId,
	}, model.GroupMemberApply{
		Status:         req.Status,
		OperationTime:  time.Now().Unix(),
		OperatorUserId: loginUserId,
	})
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	if req.Status == 1 {
		err := usecase.GroupUseCase.JoinGroup(req.OperationID, applyInfo.GroupId, applyInfo.UserId)
		if err != nil {
			tx.Rollback()
			http.Failed(c, response.GetError(response.ErrFailRequest, lang))
			return
		}
		tx.Commit()

	} else {
		repo.GroupRepo.GroupMemberApplyChangeMsg(req.OperationID, applyInfo.Id, 2)
		tx.Commit()
	}

	http.Success(c, resp)
}

func (s *groupService) JoinApplyList(c *gin.Context) {
	var (
		req  model.JoinApplyListReq
		resp model.JoinApplyListResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	loginUserId, _ := s.GetLoginUserId(c)

	groupIdList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
		"user_id": loginUserId,
		"status":  1,
		"role":    []model.RoleType{model.RoleTypeOwner, model.RoleTypeAdmin},
	}, "group_id")
	if err != nil {
		http.Success(c, resp)
		return
	}
	applys := []model.GroupMemberApply{}
	total := int64(0)
	db.Find(model.GroupMemberApply{}, map[string]interface{}{
		"group_id": groupIdList,
	}, "id desc", req.Page, req.PageSize, &total, &applys)
	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	for _, v := range applys {
		temp := model.ApplyInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	repo.GroupRepo.GetApplyUserInfo(&resp.List)
	http.Success(c, resp)
}

func (s *groupService) JoindGroupList(c *gin.Context) {
	var (
		req  model.JoindGroupListReq
		resp model.JoindGroupListResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	loginUserId, _ := s.GetLoginUserId(c)
	wheres := map[string]interface{}{
		"user_id": loginUserId,
		"status":  1,
		"role":    []string{"admin", "user", "staff"},
	}
	if req.Version > 0 {
		wheres["version"] = fmt.Sprintf(">=%d", req.Version)
	}

	groupIdList, err := db.CloumnList(&model.GroupMember{}, wheres, "group_id")
	if err != nil {
		http.Success(c, resp)
		return
	}
	groups := []model.Group{}
	total := int64(0)
	db.Find(model.Group{}, map[string]interface{}{
		"group_id": groupIdList,
	}, "", req.Page, req.PageSize, &total, &groups)
	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	for _, v := range groups {
		temp := model.GroupInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	repo.GroupRepo.GetGroupRole(&resp.List, loginUserId)
	http.Success(c, resp)
}

func (s *groupService) GetMyGroupMaxSeq(c *gin.Context) {
	var (
		req  model.MyGroupListReq
		resp model.MyGroupListMaxSeqResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	loginUserId, _ := s.GetLoginUserId(c)

	groupIdList, err := db.CloumnList(&model.GroupMember{}, map[string]interface{}{
		"user_id": loginUserId,
		"status":  1,
		"role":    []string{"admin", "user", "owner"},
	}, "group_id")
	if err != nil {
		http.Success(c, resp)
		return
	}
	for _, group := range groupIdList {
		friendsSeq := model.MaxGroupSeq{}
		friendsSeq.ConversationID = group.(string)
		resp.List = append(resp.List, friendsSeq)
	}
	resp.Count = len(groupIdList)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	http.Success(c, resp)
}

func (s *groupService) MyGroupList(c *gin.Context) {
	var (
		req  model.MyGroupListReq
		resp model.MyGroupListResq
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	loginUserId, _ := s.GetLoginUserId(c)
	wheres := map[string]interface{}{
		"user_id": loginUserId,
		"status":  1,
		"role":    model.RoleTypeOwner,
	}
	if req.Version > 0 {
		wheres["version"] = fmt.Sprintf(">=%d", req.Version)
	}

	groupIdList, err := db.CloumnList(&model.GroupMember{}, wheres, "group_id")
	if err != nil || len(groupIdList) == 0 {
		http.Success(c, resp)
		return
	}
	logger.Sugar.Debugw("加入的群id", groupIdList)
	groups := []model.Group{}
	total := int64(0)
	db.Find(model.Group{}, map[string]interface{}{
		"group_id": groupIdList,
	}, "", req.Page, req.PageSize, &total, &groups)
	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	logger.Sugar.Debugw("群组真实列表", len(groups))
	for _, v := range groups {
		temp := model.GroupInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	repo.GroupRepo.GetGroupRole(&resp.List, loginUserId)
	http.Success(c, resp)
}

func (s *groupService) GroupInfo(c *gin.Context) {
	var (
		req  model.GroupInfoReq
		resp model.GroupInfoResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)
	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupId)
	if err != nil || groupInfo.Status == 2 {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	util.CopyStructFields(&resp, groupInfo)
	resp.ShowQrcodeByNormalMemberV2 = groupInfo.ShowQrcodeByNormalMemberV2
	resp.ShowQrcodeByNormalMember = groupInfo.ShowQrcodeByNormalMember
	groups := []model.GroupInfo{model.GroupInfo(resp)}
	repo.GroupRepo.GetGroupRole(&groups, loginUserId)
	resp = model.GroupInfoResp(groups[0])

	http.Success(c, resp)
}

func (s *groupService) Information(c *gin.Context) {
	var (
		req  model.GroupInformationReq
		resp model.GroupInformationResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.Name != nil && *req.Name == "" {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "name is empty")
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	if err = permissionUseCase.PermissionUseCase.CheckGroupInformationPermission(req.OperationID, userID, req.GroupID, lang); err != nil {
		http.Failed(c, err)
		return
	}

	if resp.Group, err = usecase.GroupUseCase.UpdateInformation(req.OperationID, userID, req.GroupID, &req.GroupInformationInfo); err != nil {
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
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	if err = permissionUseCase.PermissionUseCase.CheckGroupManagePermission(req.OperationID, userID, req.GroupID, lang); err != nil {
		http.Failed(c, err)
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
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	userID := c.GetString("user_id")
	if err = permissionUseCase.PermissionUseCase.CheckGroupRemovePermission(req.OperationID, userID, req.GroupID, lang); err != nil {
		http.Failed(c, err)
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
		req  model.GroupMemberListReq
		resp model.GroupMemberListResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	groupId := req.GroupId
	members := []model.GroupMember{}
	total := int64(0)
	searchWhere := map[string]interface{}{
		"group_id": groupId,
		"status":   1,
	}
	if req.IsMute != 0 {
		searchWhere["mute_end_time"] = fmt.Sprintf(">%d", time.Now().Unix())
	}

	if req.SearchKey != "" {
		uidList := repo.GroupRepo.GetUserIdListByNickName(req.SearchKey)
		searchWhere["user_id"] = uidList
	}
	db.Find(model.GroupMember{}, searchWhere, "role_index asc", req.Page, req.PageSize, &total, &members)
	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	for _, v := range members {
		temp := model.GroupMemberInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	repo.GroupRepo.GetGroupMemberUserInfo(&resp.List)
	http.Success(c, resp)
}

func (s *groupService) UpdateGroupMember(c *gin.Context) {
	var (
		req  model.UpdateGroupMemberReq
		resp model.GroupInfoResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.UserId == "" {
		req.UserId = c.GetString("user_id")
	}
	if err = repo.GroupMemberRepo.UpdateMember(req.GroupID, req.UserId, map[string]interface{}{
		"group_nick_name": req.GroupNickName,
	}); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "error:", err)
		http.Failed(c, response.GetError(response.ErrDB, lang))
		return
	}

	repo.GroupMemberCache.DeleteMember(req.GroupID, req.UserId)

	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupID)
	if err != nil || groupInfo.Status != 1 {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)

	res := model.GroupInfo{}
	_ = util.CopyStructFields(&res, groupInfo)
	res.ShowQrcodeByNormalMember = groupInfo.ShowQrcodeByNormalMember
	res.ShowQrcodeByNormalMemberV2 = groupInfo.ShowQrcodeByNormalMemberV2

	repo.GroupRepo.GetGroupRoleOne(&res, loginUserId)
	_ = util.CopyStructFields(&resp, res)

	http.Success(c, resp)
}

func (s *groupService) QuitGroup(c *gin.Context) {
	var (
		req  model.QuitGroupReq
		resp model.QuitGroupResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if err = permissionUseCase.PermissionUseCase.CheckGroupQuitPermission(req.OperationID, loginUserId, req.GroupId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	err = usecase.GroupUseCase.GroupRemoveMember(req.OperationID, req.GroupId, []string{loginUserId}, loginUserId, model.ReasonTypeQuit)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	http.Success(c, resp)
}

func (s *groupService) InviteGroupMember(c *gin.Context) {
	var (
		req  model.InviteGroupMemberReq
		resp model.InviteGroupMemberResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	logger.Sugar.Debug(req.OperationID, util.GetSelfFuncName(), "接到邀请入群请求")

	groupInfo := model.Group{}
	if err = db.Info(&groupInfo, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if err = permissionUseCase.PermissionUseCase.CheckGroupMemberInvitePermission(req.OperationID, loginUserId, req.GroupId, len(req.UserIdList), lang); err != nil {
		http.Failed(c, err)
		return
	}

	members, err := usecase.GroupUseCase.BatchJoinGroup(req.OperationID, req.GroupId, req.UserIdList, loginUserId)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	groupMembers := []model.GroupMemberInfo{}
	util.CopyStructFields(&groupMembers, members)
	repo.GroupRepo.GetGroupMemberUserInfo(&groupMembers)
	resp.List = groupMembers
	resp.Count = len(groupMembers)
	http.Success(c, resp)
}

func (s *groupService) RemoveGroupMember(c *gin.Context) {
	var (
		req  model.RemoveGroupMemberReq
		resp model.RemoveGroupMemberResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	logger.Sugar.Debug(req.OperationID, util.GetSelfFuncName(), "收到请求")

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	for _, userID := range req.UserIdList {
		if err = permissionUseCase.PermissionUseCase.CheckGroupMemberKickPermission(req.OperationID, loginUserId, userID, req.GroupId, lang); err != nil {
			http.Failed(c, err)
			return
		}
	}

	err = usecase.GroupUseCase.GroupRemoveMember(req.OperationID, req.GroupId, req.UserIdList, loginUserId, model.ReasonTypeKick)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	http.Success(c, resp)
}

func (s *groupService) GroupUpdate(c *gin.Context) {
	var (
		req  model.GroupUpdateReq
		resp model.GroupUpdateResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	loginMemberInfo := model.GroupMember{}
	err = db.Info(&loginMemberInfo, model.GroupMember{
		GroupId: req.GroupId,
		UserId:  loginUserId,
		Status:  1,
	})

	if err != nil || loginMemberInfo.Role == "user" {
		http.Failed(c, response.GetError(response.ErrUserPermissions, lang))
		return
	}

	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupId)
	if err != nil || groupInfo.Status != 1 {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	updateInfo := model.Group{}
	util.CopyStructFields(&updateInfo, req)
	updateInfo.ShowQrcodeByNormalMemberV2 = req.ShowQrcodeByNormalMemberV2
	updateInfo.ShowQrcodeByNormalMember = req.ShowQrcodeByNormalMember

	groupInfo, err = usecase.GroupUseCase.UpdateGroup(req.OperationID, req.GroupId, &updateInfo, nil, loginUserId)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	res := model.GroupInfo{}
	util.CopyStructFields(&res, groupInfo)
	res.ShowQrcodeByNormalMember = groupInfo.ShowQrcodeByNormalMember
	res.ShowQrcodeByNormalMemberV2 = groupInfo.ShowQrcodeByNormalMemberV2

	repo.GroupRepo.GetGroupRoleOne(&res, loginUserId)
	util.CopyStructFields(&resp, res)

	http.Success(c, resp)
}

func (s *groupService) GroupUpdateAvatar(c *gin.Context) {
	var (
		req  model.GroupUpdateAvatarReq
		resp model.GroupUpdateAvatarResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	loginMemberInfo := model.GroupMember{}
	err = db.Info(&loginMemberInfo, model.GroupMember{
		GroupId: req.GroupId,
		UserId:  loginUserId,
		Status:  1,
	})

	if err != nil || loginMemberInfo.Role == "user" {
		http.Failed(c, response.GetError(response.ErrUserPermissions, lang))
		return
	}

	groupInfo := model.Group{}
	err = db.Info(&groupInfo, req.GroupId)
	if err != nil || groupInfo.Status != 1 {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	updateInfo := model.Group{}
	_ = util.CopyStructFields(&updateInfo, req)

	groupInfo, err = usecase.GroupUseCase.UpdateGroup(req.OperationID, req.GroupId, &updateInfo, nil, loginUserId)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}
	res := model.GroupInfo{}
	_ = util.CopyStructFields(&res, groupInfo)
	res.ShowQrcodeByNormalMember = groupInfo.ShowQrcodeByNormalMember
	res.ShowQrcodeByNormalMemberV2 = groupInfo.ShowQrcodeByNormalMemberV2

	repo.GroupRepo.GetGroupRoleOne(&res, loginUserId)
	_ = util.CopyStructFields(&resp, res)

	http.Success(c, resp)
}

func (s *groupService) GroupSync(c *gin.Context) {
	var (
		req  model.GroupSyncReq
		resp model.GroupSyncResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)

	memberInfo := model.GroupMember{}
	err = db.Info(&memberInfo, &model.GroupMember{
		UserId:  loginUserId,
		GroupId: req.GroupId,
		Status:  1,
	})
	if err != nil {
		http.Failed(c, response.GetError(response.ErrUserPermissions, lang))
		return
	}

	total := int64(0)
	members := []model.GroupMember{}
	_ = db.Find(&model.GroupMember{}, map[string]interface{}{
		"group_id": req.GroupId,
		"version":  fmt.Sprintf(">%d", req.LocalVersion),
	}, "version asc", req.Page, req.PageSize, &total, &members)
	resp.Count = int(total)
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	for _, v := range members {
		temp := model.GroupMemberInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	repo.GroupRepo.GetGroupMemberUserInfo(&resp.List)
	http.Success(c, resp)
}

func (s *groupService) GroupSetAdmin(c *gin.Context) {
	var (
		req  model.GroupSetAdminReq
		resp model.GroupSetAdminResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if loginUserId == req.UserId {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error:", fmt.Sprintf("user self"))
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	if err = permissionUseCase.PermissionUseCase.CheckGroupSetAdminPermission(req.OperationID, loginUserId, req.UserId, req.GroupId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	err = usecase.GroupUseCase.SetGroupAdmin(req.OperationID, req.GroupId, req.UserId, req.Status, loginUserId)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
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
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if err = permissionUseCase.PermissionUseCase.CheckGroupSetOwnerPermission(req.OperationID, loginUserId, req.UserId, req.GroupId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	err = usecase.GroupUseCase.SetGroupOwner(req.OperationID, req.GroupId, req.UserId, loginUserId)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}

	http.Success(c, resp)
}

func (s *groupService) GetOwnerAdmin(c *gin.Context) {
	var (
		req     model.GetOwnerAdminReq
		resp    model.GetAdminOwnerResp
		members []model.GroupMember
		err     error
	)
	lang := c.GetHeader("Locale")
	if err = c.ShouldBindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	if members, err = repo.GroupMemberRepo.GetAdminOwner(req.GroupId); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotMember, lang))
		return
	}
	for _, v := range members {
		temp := model.GroupMemberInfo{}
		_ = util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	repo.GroupRepo.GetGroupMemberUserInfo(&resp.List)

	http.Success(c, resp)
}

func (s *groupService) GroupMuteMember(c *gin.Context) {
	var (
		req  model.GroupMuteMemberReq
		resp model.GroupMuteMemberResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)

	if err = permissionUseCase.PermissionUseCase.CheckGroupMemberMutePermission(req.OperationID, loginUserId, req.UserId, req.GroupId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	memberInfo := model.GroupMember{}
	err = db.Info(&memberInfo, model.GroupMember{
		GroupId: req.GroupId,
		UserId:  req.UserId,
		Status:  1,
	})

	if err != nil {
		http.Failed(c, response.GetError(response.ErrNotInGroup, lang))
		return
	}

	endTime := time.Now().Unix() + int64(req.MuteSec)
	if req.MuteSec == 0 {
		endTime = 0
	}

	tx := db.DB.Begin()
	err = db.UpdateTx(tx, &model.GroupMember{}, memberInfo.Id, map[string]interface{}{
		"mute_end_time": endTime,
	})
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}

	tx.Commit()

	repo.GroupMemberCache.DeleteMember(req.GroupId, req.UserId)

	var content chatModel.MessageContent
	content.OperatorID = loginUserId
	content.BeOperatorList = append(content.BeOperatorList, chatModel.MessageBeOperator{BeOperatorID: req.UserId})
	switch req.MuteSec {
	case 60 * 60:
		content.TimeType = chatModel.MessageContentTimeTypeOneHour

	case 24 * 60 * 60:
		content.TimeType = chatModel.MessageContentTimeTypeOneDay

	case 365 * 24 * 60 * 60:
		content.TimeType = chatModel.MessageContentTimeTypeForever
	}

	if content.TimeType == 0 {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(req.OperationID, req.GroupId, chatModel.MessageGroupOneUnmuteNotify, &content)
	} else {
		chatUseCase.MessageUseCase.SendSystemMessageToGroup(req.OperationID, req.GroupId, chatModel.MessageGroupOneMuteNotify, &content)
	}

	repo.GroupRepo.GroupMemberChangeMsg(req.OperationID, req.GroupId, []string{req.UserId}, common.GroupMemberChangePush)
	http.Success(c, resp)
}

func (s *groupService) GroupListSync(c *gin.Context) {
	var (
		req  model.GroupListSyncReq
		resp model.GroupListSyncResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}
	loginUserId, _ := s.GetLoginUserId(c)

	groups := []model.Group{}
	total := int64(0)
	db.Find(model.Group{}, map[string]interface{}{
		"group_id": req.GroupIdList,
	}, "", 1, 999, &total, &groups)
	resp.Count = int(total)
	resp.Page = 1
	resp.PageSize = 0
	for _, v := range groups {
		temp := model.GroupInfo{}
		util.CopyStructFields(&temp, v)
		resp.List = append(resp.List, temp)
	}
	repo.GroupRepo.GetGroupRole(&resp.List, loginUserId)
	http.Success(c, resp)
}

func (s *groupService) GroupMuteAll(c *gin.Context) {
	var (
		req  model.GroupMuteAllReq
		resp model.GroupMuteAllResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	if req.MuteAllMember == constant.MuteMemberPeriod {
		times := strings.Split(req.MuteAllPeriod, "-")
		if len(times) != 2 {
			http.Failed(c, response.GetError(response.ErrMutePeriod, lang))
			return
		}
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if err = db.Info(&model.Group{}, model.Group{
		GroupId: req.GroupId,
		Status:  1,
	}); err != nil {
		http.Failed(c, response.GetError(response.ErrGroupNotExist, lang))
		return
	}

	if err = permissionUseCase.PermissionUseCase.CheckGroupMutePermission(req.OperationID, loginUserId, req.GroupId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	err = usecase.GroupUseCase.UpdateGroupMuteInfo(req)
	if err != nil {
		http.Failed(c, response.GetError(response.ErrFailRequest, lang))
		return
	}

	http.Success(c, resp)
}

func (s *groupService) SetGroupNickname(c *gin.Context) {
	var (
		req  model.GroupNickNameReq
		resp model.GroupNickNameResp
		err  error
	)
	lang := c.GetHeader("Locale")
	if err = c.BindJSON(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind json error:", err)
		http.Failed(c, response.GetError(response.ErrBadRequest, lang))
		return
	}

	loginUserId, _ := s.GetLoginUserId(c)

	if err = permissionUseCase.PermissionUseCase.CheckGroupBaseInfo(req.OperationID, loginUserId, req.GroupId, lang); err != nil {
		http.Failed(c, err)
		return
	}

	if err = repo.GroupMemberRepo.SetMemberNickName(req.GroupId, loginUserId, req.GroupNickName); err != nil {
		http.Failed(c, err)
		return
	}

	repo.GroupMemberCache.DeleteMember(req.GroupId, loginUserId)

	repo.GroupRepo.GroupMemberChangeMsg(req.OperationID, req.GroupId, []string{loginUserId}, common.GroupMemberChangePush)
	http.Success(c, resp)
}
