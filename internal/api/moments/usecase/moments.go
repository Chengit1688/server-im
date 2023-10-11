package usecase

import (
	"im/internal/api/moments/model"
	momentsRepo "im/internal/api/moments/repo"
	userusecase "im/internal/api/user/usecase"
	"im/pkg/response"
	"strings"
)

var MomentsUseCase = new(momentsUseCase)

type momentsUseCase struct{}

func (r *momentsUseCase) MomentsDetail(v model.MomentsMessage, lang string, offset, limit int) (info model.IssueInfo, err error) {
	info.MomentsID = v.ID
	user, err := userusecase.UserUseCase.GetBaseInfo(v.UserId)
	if err != nil {
		return model.IssueInfo{}, response.GetError(response.ErrDB, lang)
	}
	info.User.UserId = user.UserId
	info.User.FaceURL = user.FaceURL
	info.User.NickName = user.NickName

	momentsInfo := v
	info.MomentsInfo = &momentsInfo

	var momentsCount int64
	momentsCommentsList, err := momentsRepo.MomentsRepo.MomentsCommentsList(v.ID, &momentsCount, offset, limit)
	if err != nil {
		return model.IssueInfo{}, response.GetError(response.ErrDB, lang)
	}
	momentsComments := make([]model.CommentsInfo, 0)
	if len(momentsCommentsList) > 0 {
		for _, vv := range momentsCommentsList {
			momentsCommentsInfo := model.CommentsInfo{}
			user, err = userusecase.UserUseCase.GetBaseInfo(vv.UserID)
			if err != nil {
				return model.IssueInfo{}, response.GetError(response.ErrDB, lang)
			}
			momentsCommentsInfo.User.UserId = user.UserId
			momentsCommentsInfo.User.FaceURL = user.FaceURL
			momentsCommentsInfo.User.NickName = user.NickName
			momentsCommentsInfo.MomentsID = vv.MomentsID
			momentsCommentsInfo.ID = vv.ID
			momentsCommentsInfo.IsOwnComment = vv.IsOwnComment
			momentsCommentsInfo.ReplyToUser.UserId = vv.ReplyUser.UserID
			momentsCommentsInfo.ReplyToUser.FaceURL = vv.ReplyUser.FaceURL
			momentsCommentsInfo.ReplyToUser.NickName = vv.ReplyUser.NickName
			momentsCommentsInfo.Content = vv.Content
			momentsCommentsInfo.Images = strings.Split(vv.Image, ";")
			momentsComments = append(momentsComments, momentsCommentsInfo)
		}
	}
	info.Comments = momentsComments
	info.CommentsCount = momentsCount

	var likeCount int64
	likeList, err := momentsRepo.MomentsRepo.MomentsCommentsLikeList(v.ID, &likeCount)
	if err != nil {
		return model.IssueInfo{}, response.GetError(response.ErrDB, lang)
	}
	likes := make([]model.MomentsCommentsLikeInfo, 0)
	for _, l := range likeList {
		like := model.MomentsCommentsLikeInfo{}
		userInfo, err := userusecase.UserUseCase.GetBaseInfo(l.FriendUserID)
		if err != nil {
			return model.IssueInfo{}, response.GetError(response.ErrDB, lang)
		}
		like.MomentsID = l.MomentsID
		like.ID = l.ID
		like.User.NickName = userInfo.NickName
		like.User.UserId = userInfo.UserId
		like.User.FaceURL = userInfo.FaceURL
		likes = append(likes, like)
	}
	info.Likes = likes
	info.LikesCount = likeCount
	return info, nil
}
