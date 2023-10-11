package service

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	apiChatModel "im/internal/api/chat/model"
	apiChatRepo "im/internal/api/chat/repo"
	apiChatUseCase "im/internal/api/chat/usecase"
	apiFriendUseCase "im/internal/api/friend/usecase"
	apiGroupUseCase "im/internal/api/group/usecase"
	apiUserUseCase "im/internal/api/user/usecase"
	cmsModel "im/internal/cms_api/chat/model"
	cmsRepo "im/internal/cms_api/chat/repo"
	"im/pkg/code"
	"im/pkg/common"
	"im/pkg/db"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/mqtt"
	"im/pkg/util"
	http2 "net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var MessageService = new(messageService)

type messageService struct{}

func (s *messageService) GetLoginUserId(c *gin.Context) (string, error) {

	userId := c.GetString("o_user_id")
	return userId, nil
}

func (s *messageService) HistoryList(c *gin.Context) {
	var (
		req            cmsModel.MessageHistoryListReq
		resp           cmsModel.MessageHistoryListResp
		conversationID string
		err            error
	)

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	req.Check()
	resp.Pagination = req.Pagination

	if req.Export {
		req.Offset = 0
		req.Limit = 99999999
	}

	conversationID = apiChatUseCase.ConversationUseCase.GetConversationID(req.ConversationType, req.SendID, req.RecvID)
	if conversationID == "" && req.RecvID != "" {

		if user, err2 := apiUserUseCase.UserUseCase.GetBaseInfo(req.RecvID); err2 != nil || user == nil {
			http.Success(c, resp)
			return
		}
		conversationID = req.RecvID
	}

	switch req.ConversationType {
	case apiChatModel.ConversationTypeSingle:
		if resp.List, resp.Count, err = apiChatRepo.MessageRepo.ListSingle(conversationID, req.SendID, req.RecvID, req.Type, req.Content, req.StartTime, req.EndTime, req.Offset, req.Limit); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("db list error, error: %v", err))
			http.Failed(c, code.ErrDB)
			return
		}

	case apiChatModel.ConversationTypeGroup:
		if resp.List, resp.Count, err = apiChatRepo.MessageRepo.List(conversationID, req.SendID, req.Type, req.Content, req.StartTime, req.EndTime, req.Offset, req.Limit); err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("db list error, error: %v", err))
			http.Failed(c, code.ErrDB)
			return
		}

	default:
		http.Failed(c, code.ErrBadRequest)
		return
	}

	for i, message := range resp.List {
		if message.SendID != "" {
			if user, err2 := apiUserUseCase.UserUseCase.GetBaseInfo(message.SendID); err2 == nil && user != nil {
				message.SendNickname = user.NickName
				message.SendFaceUrl = user.FaceURL
			}

			message.RecvID = apiChatUseCase.ConversationUseCase.GetRecvID(message.SendID, message.ConversationType, message.ConversationID)

			if message.ConversationType == apiChatModel.ConversationTypeSingle {
				if user, err2 := apiUserUseCase.UserUseCase.GetBaseInfo(message.RecvID); err2 == nil && user != nil {
					message.RecvName = user.NickName
				}
			}
		}
		resp.List[i] = message
	}

	if !req.Export {
		http.Success(c, resp)
	} else {
		s.HistoryExport(c, req.OperationID, req.ConversationType, resp.List)
	}
}

func (s *messageService) HistoryExport(c *gin.Context, operationID string, conversationType apiChatModel.ConversationType, list []apiChatModel.MessageInfo) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Sheet1"
	f.NewSheet(sheetName)

	var sheetHeader []interface{}
	switch conversationType {
	case apiChatModel.ConversationTypeSingle:
		sheetHeader = []interface{}{"发送时间", "发送者", "发送者ID", "接收者", "接收者ID", "消息类型", "聊天内容"}

	case apiChatModel.ConversationTypeGroup:
		sheetHeader = []interface{}{"发送时间", "发送者", "发送者ID", "消息类型", "聊天内容"}
	}
	f.SetSheetRow(sheetName, "A1", &sheetHeader)

	var row []interface{}
	for i := range list {
		data := list[i]

		sendTime := util.FormatTime(util.UnunixMilliTime(data.SendTime))
		contentType := apiChatModel.MessageTypeString(data.Type)
		content := apiChatUseCase.MessageUseCase.ComposeMessageContent(operationID, data.Content)

		switch conversationType {
		case apiChatModel.ConversationTypeSingle:
			row = []interface{}{sendTime, data.SendNickname, data.SendID, data.RecvName, data.RecvID, contentType, content}

		case apiChatModel.ConversationTypeGroup:
			row = []interface{}{sendTime, data.SendNickname, data.SendID, contentType, content}
		}
		f.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &row)
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		logger.Sugar.Errorw(operationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("write to buffer error, error: %v", err))
		http.Failed(c, code.ErrUnknown)
		return
	}

	c.Writer.WriteHeader(http2.StatusOK)
	filename := url.QueryEscape("聊天记录.xlsx")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=utf-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Writer.Write(buf.Bytes())
}

func (s *messageService) Change(c *gin.Context) {
	var (
		req         apiChatModel.MessageChangeReq
		resp        apiChatModel.MessageChangeResp
		messageType apiChatModel.MessageType
		err         error
	)

	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if len(req.MsgIDList) == 0 {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "msg id list size 0")
		http.Failed(c, code.ErrBadRequest)
		return
	}

	var msg *apiChatModel.Message
	if msg, err = apiChatRepo.MessageRepo.Get(req.MsgIDList[0]); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Errorf("db get error, msg id list: %v, error: %v", req.MsgIDList, err))
		http.Failed(c, code.ErrDB)
		return
	}

	switch req.Status {
	case apiChatModel.MessageStatusTypeDelete:
		messageType = apiChatModel.MessageDelete

	default:
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", "message status type error")
		http.Failed(c, code.ErrBadRequest)
		return
	}

	if _, err = apiChatUseCase.MessageUseCase.UpdateStatus(req.OperationID, "", msg.ConversationType, msg.ConversationID, messageType, req.Status, req.MsgIDList); err != nil {
		http.Failed(c, err)
		return
	}
	http.Success(c, resp)
}

func (s *messageService) Clear(c *gin.Context) {
	var (
		req  cmsModel.MessageClearReq
		resp cmsModel.MessageClearResp
		err  error
	)

	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}

	l := util.NewLock(db.RedisCli, common.LockMessageClear)
	if l.IsLock() {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err)
		http.Failed(c, code.ErrTaskBusy)
		return
	}
	if err = l.Lock(); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("lock error, error: %v", err))
		http.Failed(c, code.ErrDB)
		return
	}
	defer l.Unlock()

	util.CopyStructFields(&resp, &req)
	switch req.Type {
	case cmsModel.ClearTypeUser:
		go func() {
			if err = apiChatUseCase.MessageUseCase.ClearClient(req.OperationID, req.TargetID); err != nil {
				http.Failed(c, err)
				return
			}

			mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageAdminClearPush, nil, req.TargetID)
		}()
	case cmsModel.ClearTypeGroupMember:
		go func() {
			var maxSeq int64

			groupMemberIDList := apiGroupUseCase.GroupUseCase.GroupMemberIdList(req.TargetID)
			for _, memberID := range groupMemberIDList {
				if err = apiChatUseCase.MessageUseCase.ClearConversation(req.OperationID, memberID, apiChatModel.ConversationTypeGroup, req.TargetID, maxSeq); err != nil {
					continue
				}

				mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageClearPush, apiChatModel.MessageClearResp{
					ConversationType: apiChatModel.ConversationTypeGroup,
					ConversationID:   req.TargetID,
					MaxSeq:           maxSeq,
				}, memberID)
			}
		}()
	case cmsModel.ClearTypeAll:
		go func() {
			var userIDList []string
			if userIDList, err = apiUserUseCase.UserUseCase.GetAllUserIDList(); err != nil {
				http.Failed(c, code.ErrDB)
				return
			}

			for _, userID := range userIDList {
				if err = apiChatUseCase.MessageUseCase.ClearClient(req.OperationID, userID); err != nil {
					continue
				}

				mqtt.SendMessageToUsers(req.OperationID, common.ChatMessageAdminClearPush, nil, userID)
			}
		}()
	}
	http.Success(c, resp)
}

func (s *messageService) MultiSendList(c *gin.Context) {
	var (
		req  cmsModel.GetMultiSendPagingReq
		resp cmsModel.GetMultiSendPagingResp
		err  error
	)

	if err = c.ShouldBindQuery(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind query error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	items, count, err := cmsRepo.MessageRepo.MultiSendPaging(req)
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "query db error:", err)
		http.Failed(c, code.ErrDB)
		return
	}
	resp.Count = count
	resp.Page = req.Page
	resp.PageSize = req.PageSize

	var recordUserIDMap map[int][]string
	recordUserIDMap = make(map[int][]string)
	var recordUserNicknameMap map[int][]string
	recordUserNicknameMap = make(map[int][]string)
	var recordContentMap map[int]string
	recordContentMap = make(map[int]string)
	for _, item := range items {
		if _, ok := recordUserIDMap[item.ID]; !ok {
			resp.List = append(resp.List, cmsModel.GetMultiSendPagingItem{ID: item.ID, Operate: item.Username, CreatedAt: item.CreatedAt})
			recordContentMap[item.ID] = item.Content
		}
		recordUserIDMap[item.ID] = append(recordUserIDMap[item.ID], item.SenderID)
		recordUserNicknameMap[item.ID] = append(recordUserNicknameMap[item.ID], item.SenderNickname)
	}

	for i, record := range resp.List {
		arrayID := recordUserIDMap[record.ID]
		arrayNickname := recordUserNicknameMap[record.ID]
		resp.List[i].SenderIDs = strings.Join(arrayID, ",")
		resp.List[i].SenderNicknames = strings.Join(arrayNickname, ",")
		content, _ := util.Encrypt([]byte(recordContentMap[record.ID]), common.ContentKey)
		resp.List[i].Content = content
	}

	http.Success(c, resp)
}

func (s *messageService) MultiSend(c *gin.Context) {
	var (
		req            cmsModel.MultiSendReq
		err            error
		decryptContent string
	)

	if err = c.ShouldBind(&req); err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "bind query error:", err)
		http.Failed(c, code.ErrBadRequest)
		return
	}
	timeNow := time.Now()
	timestamp := timeNow.Unix()
	operate_id := c.GetString("o_user_id")

	var userIdsArray []string
	if strings.Contains(req.SenderIDs, ",") {
		userIdsArray = strings.Split(req.SenderIDs, ",")
	} else {
		userIdsArray = append(userIdsArray, req.SenderIDs)
	}

	checkSameAccountMap := map[string]bool{}
	var users []string
	for _, checkAccount := range userIdsArray {
		if _, ok := checkSameAccountMap[checkAccount]; !ok {
			users = append(users, checkAccount)
		}
		checkSameAccountMap[checkAccount] = true
	}

	if decryptContent, err = util.Decrypt(req.Content, common.ContentKey); err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("decrypt error, error: %v", err))
		http.Failed(c, code.ErrBadRequest)
		return
	}

	for _, user := range users {
		_, err = apiUserUseCase.UserUseCase.GetBaseInfo(user)
		if err != nil {
			logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "user_id not exist error:", err)
			http.Failed(c, code.ErrUserIdNotExist)
			return
		}
	}

	running, err := cmsRepo.MessageCache.GetMultiSendLock()
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "redis error:", err)
		http.Failed(c, code.ErrDB)
		return
	}
	if running {
		http.Failed(c, code.ErrMsgSendRunning)
		return
	}
	err = cmsRepo.MessageCache.SetMultiSendLock()
	if err != nil {
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "redis error:", err)
		http.Failed(c, code.ErrDB)
		return
	}

	err = cmsRepo.MessageRepo.MultiSendAdd(cmsModel.MultiSendRecord{OperateID: operate_id, Content: decryptContent, CreatedAt: timestamp}, users)
	if err != nil {
		cmsRepo.MessageCache.DelMultiSendLock()
		logger.Sugar.Error(req.OperationID, util.GetSelfFuncName(), "db error:", err)
		http.Failed(c, code.ErrDB)
		return
	}

	go sendMessage(req.OperationID, req.Content, users)

	http.Success(c, "消息群发任务创建成功，后台发送中")
}

func sendMessage(OperationID, Content string, users []string) (err error) {

	var msgType apiChatModel.MessageType = 1
	for _, user_id := range users {

		cmsRepo.MessageCache.SetMultiSendLock()
		friends := apiFriendUseCase.FriendUseCase.GetUserFriendIdList(user_id)

		total := len(friends)
		cur := 0
		for _, recv_id := range friends {
			conversationID := apiChatUseCase.ConversationUseCase.GetConversationID(apiChatModel.ConversationTypeSingle, user_id, recv_id)
			ClientMsgID := apiChatUseCase.MessageUseCase.GetMsgID(apiChatModel.ConversationTypeSingle, conversationID, util.UnixMilliTime(time.Now()))
			if _, err = apiChatUseCase.MessageUseCase.SendMessageToUsers(OperationID, ClientMsgID, apiChatModel.ConversationTypeSingle, conversationID, user_id, msgType, Content, user_id, recv_id); err != nil {
				logger.Sugar.Error(OperationID, util.GetSelfFuncName(), "send mqtt error:", err)
				return
			}

			time.Sleep(time.Millisecond * 50)
			cur += 1
			if total >= 50 && cur == 50 {
				time.Sleep(time.Second)
				cur = 0
				cmsRepo.MessageCache.SetMultiSendLock()
			}
		}
	}
	cmsRepo.MessageCache.DelMultiSendLock()
	return
}
