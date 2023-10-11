package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"im/config"
	"im/internal/cms_api/third/model"
	"im/pkg/code"
	"im/pkg/ffmpeg"
	"im/pkg/http"
	"im/pkg/logger"
	"im/pkg/minio"
	"im/pkg/oss"
	"im/pkg/util"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

var MinioService = new(minioService)

type minioService struct{}

func (s *minioService) Upload(c *gin.Context) {
	req := new(model.UploadReq)
	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "接口收到文件 Step 1")
	fileObj, err := file.Open()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrUnknown)
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "接口打开文件并发送到minio Step 2")
	url, err := minio.UploadToOss(file.Filename, fileObj, req.FileType, file.Size)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrUnknown)
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "文件发送minio完成 Step 3")
	newName := filepath.Base(url)
	ret := new(model.UploadResp)
	ret.OldName = file.Filename
	ret.NewName = newName
	ret.Url = url
	ret.ContentType = minio.GetContentType(ret.OldName)
	switch req.FileType {
	case 3:

		defer fileObj.Close()
		suffix := path.Ext(file.Filename)
		dst := fmt.Sprintf("%s/%s%s", os.TempDir(), uuid.Must(uuid.NewV4()), suffix)
		out, err := os.Create(dst)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrUnknown)
			return
		}
		defer os.Remove(dst)
		defer out.Close()
		_, err = io.Copy(out, fileObj)
		logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "缩略图制作 Step 4")
		thumbnailFile := ffmpeg.ResizeImage(dst)
		if thumbnailFile != dst {
			thumbnailFileObj, _ := os.Open(thumbnailFile)
			defer os.Remove(thumbnailFile)
			defer thumbnailFileObj.Close()
			info, _ := thumbnailFileObj.Stat()
			logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "缩略图上传minio Step 5")
			thumbnail, err := minio.UploadToOss(thumbnailFileObj.Name(), thumbnailFileObj, 3, info.Size())
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
				http.Failed(c, code.ErrUnknown)
				return
			}
			logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "缩略图上传完成minio Step 6")
			ret.Thumbnail = thumbnail
		} else {
			ret.Thumbnail = url
			logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "无需制作缩略图 Step 5")
		}

	case 4:

		suffix := path.Ext(file.Filename)
		dst := fmt.Sprintf("%s/%s%s", os.TempDir(), uuid.Must(uuid.NewV4()), suffix)
		out, err := os.Create(dst)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrUnknown)
			return
		}
		defer os.Remove(dst)
		defer out.Close()
		_, err = io.Copy(out, fileObj)
		frameFile := ffmpeg.GetFirstFrame(dst)
		defer os.Remove(frameFile)
		defer fileObj.Close()
		thumbnailFile := ffmpeg.ResizeImage(frameFile)

		if thumbnailFile != dst {
			thumbnailFileObj, _ := os.Open(thumbnailFile)
			defer os.Remove(thumbnailFile)
			defer thumbnailFileObj.Close()
			info, _ := thumbnailFileObj.Stat()
			thumbnail, err := minio.UploadToOss(thumbnailFileObj.Name(), thumbnailFileObj, 4, info.Size())
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
				http.Failed(c, code.ErrUnknown)
				return
			}
			ret.Thumbnail = thumbnail
		} else {
			ret.Thumbnail = url
		}
	}
	http.Success(c, ret)
	return
}

func (s *minioService) UploadV2(c *gin.Context) {
	req := new(model.UploadReq)

	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "接口收到文件 Step 1")
	fileObj, err := file.Open()
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrUnknown)
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "接口打开文件并发送到minio Step 2")
	url, err := minio.UploadToOssV2(file.Filename, fileObj, req.FileType, file.Size)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrUnknown)
		return
	}
	logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "文件发送minio完成 Step 3")
	newName := filepath.Base(url)
	ret := new(model.UploadResp)
	ret.OldName = file.Filename
	ret.NewName = newName
	ret.Url = url
	ret.ContentType = minio.GetContentType(ret.OldName)
	switch req.FileType {
	case 3:

		defer fileObj.Close()
		suffix := path.Ext(file.Filename)
		dst := fmt.Sprintf("%s/%s%s", os.TempDir(), uuid.Must(uuid.NewV4()), suffix)
		out, err := os.Create(dst)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrUnknown)
			return
		}
		defer os.Remove(dst)
		defer out.Close()
		_, err = io.Copy(out, fileObj)
		logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "缩略图制作 Step 4")
		thumbnailFile := ffmpeg.ResizeImage(dst)
		if thumbnailFile != dst {
			thumbnailFileObj, _ := os.Open(thumbnailFile)
			defer os.Remove(thumbnailFile)
			defer thumbnailFileObj.Close()
			info, _ := thumbnailFileObj.Stat()
			logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "缩略图上传minio Step 5")
			thumbnail, err := minio.UploadToOssV2(thumbnailFileObj.Name(), thumbnailFileObj, 3, info.Size())
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
				http.Failed(c, code.ErrUnknown)
				return
			}
			logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "缩略图上传完成minio Step 6")
			ret.Thumbnail = thumbnail
		} else {
			ret.Thumbnail = url
			logger.Sugar.Debugw(req.OperationID, "func", util.GetSelfFuncName(), "time", time.Now().Format("2006-01-02 15:04:05"), "msg", "无需制作缩略图 Step 5")
		}

	case 4:

		suffix := path.Ext(file.Filename)
		dst := fmt.Sprintf("%s/%s%s", os.TempDir(), uuid.Must(uuid.NewV4()), suffix)
		out, err := os.Create(dst)
		if err != nil {
			logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
			http.Failed(c, code.ErrUnknown)
			return
		}
		defer os.Remove(dst)
		defer out.Close()
		_, err = io.Copy(out, fileObj)
		frameFile := ffmpeg.GetFirstFrame(dst)
		defer os.Remove(frameFile)
		defer fileObj.Close()
		thumbnailFile := ffmpeg.ResizeImage(frameFile)

		if thumbnailFile != dst {
			thumbnailFileObj, _ := os.Open(thumbnailFile)
			defer os.Remove(thumbnailFile)
			defer thumbnailFileObj.Close()
			info, _ := thumbnailFileObj.Stat()
			thumbnail, err := minio.UploadToOssV2(thumbnailFileObj.Name(), thumbnailFileObj, 4, info.Size())
			if err != nil {
				logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
				http.Failed(c, code.ErrUnknown)
				return
			}
			ret.Thumbnail = thumbnail
		} else {
			ret.Thumbnail = url
		}
	}

	http.Success(c, ret)
	return
}

func (s *minioService) GetFileUrl(c *gin.Context) {
	req := new(model.GetUrlReq)

	err := c.ShouldBind(&req)
	if err != nil {
		logger.Sugar.Errorw(req.OperationID, "func", util.GetSelfFuncName(), "error", err.Error())
		http.Failed(c, code.ErrBadRequest)
		return
	}
	ret := new(model.UploadResp)
	url := minio.GetNewUrl(req.Filename, req.FileType)
	newName := filepath.Base(url)
	ret.OldName = req.Filename
	ret.NewName = newName
	ret.Url = url
	ret.ContentType = minio.GetContentType(ret.OldName)
	ret.Thumbnail = minio.OssResizeImage(url, req.Width, req.Height, req.FileType)
	http.Success(c, ret)
	return
}

func (s *minioService) GetSTS(c *gin.Context) {
	var (
		err         error
		OperationID string
	)
	if err = c.ShouldBindQuery(&OperationID); err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("bind json, error: %v", err))
		http.Failed(c, code.ErrBadRequest)
		return
	}
	cfg := config.Config
	sessionName := strconv.Itoa(rand.Int())
	credential, err := oss.OS.GetSTS(cfg.AliYun.OSS.BucketName, sessionName)
	if err != nil {
		logger.Sugar.Errorw(OperationID, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("GetSTS , error: %v", err))
		http.Failed(c, code.GetError(err, OperationID))
		return
	}
	http.Success(c, credential)
}
