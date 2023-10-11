package minio

import (
	"context"
	"fmt"
	"im/config"
	"im/pkg/logger"
	"im/pkg/util"
	"io"
	"path"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var (
	MinioClient *minio.Client
	Bucket      string
	OssPrefix   string
)

func Init() {
	var err error
	cfg := config.Config.Minio
	Bucket = cfg.OssBucket
	OssPrefix = cfg.OssPrefix
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: false,
	}

	MinioClient, err = minio.New(cfg.Endpoint, opts)

	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("init minio client failed", err.Error()))
		return
	}
	opt := minio.MakeBucketOptions{
		Region:        "",
		ObjectLocking: false,
	}
	err = MinioClient.MakeBucket(context.Background(), cfg.OssBucket, opt)
	if err != nil {
		logger.Sugar.Error(zap.String("func", util.GetSelfFuncName()), zap.String("MakeBucket failed", err.Error()))
		exists, err := MinioClient.BucketExists(context.Background(), cfg.OssBucket)
		if err == nil && exists {
			logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("oss bucket ready", cfg.OssBucket))
		} else {
			if err != nil {
				logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "get bucket failed", err)
			}
			return
		}
	}
	logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("minio init ready", cfg.OssBucket))
}

func UploadToOss(filename string, fileObj io.Reader, fileType int, fileSize int64) (string, error) {
	newName, contentType := GetNewFileNameAndContentType(filename)
	savePath := GetFilePath(fileType)
	fullPath := fmt.Sprintf("%s%s", savePath, newName)
	_, err := MinioClient.PutObject(context.Background(), Bucket, fullPath, fileObj, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "PutObject failed", err)
		logger.Sugar.Errorw("", "Bucket", Bucket, "fullPath", fullPath, "fileSize", fileSize)
		return "", err
	}
	logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("Successfully PutObject", filename))
	return fmt.Sprintf("%s%s%s", OssPrefix, Bucket, fullPath), nil
}

func UploadToOssV2(filename string, fileObj io.Reader, fileType int, fileSize int64) (string, error) {
	newName, contentType := GetNewFileNameAndContentType(filename)
	savePath := GetFilePath(fileType)
	fullPath := fmt.Sprintf("%s%s", savePath, newName)
	_, err := MinioClient.PutObject(context.Background(), Bucket, fullPath, fileObj, fileSize, minio.PutObjectOptions{ContentType: contentType})

	if err != nil {
		logger.Sugar.Errorw("", "func", util.GetSelfFuncName(), "PutObject failed", err)
		logger.Sugar.Errorw("", "Bucket", Bucket, "fullPath", fullPath, "fileSize", fileSize)
		return "", err
	}
	logger.Sugar.Info(zap.String("func", util.GetSelfFuncName()), zap.String("Successfully PutObject", filename))
	return fmt.Sprintf("/%s%s", Bucket, fullPath), nil
}

func GetNewUrl(filename string,fileType int)string{
	savePath := GetNewFilePath(fileType)
	fullPath := fmt.Sprintf("%s%s", savePath, filename)
	return fmt.Sprintf("-%s%s", Bucket, fullPath)
}

func GetNewFileNameAndContentType(fileName string) (string, string) {
	suffix := path.Ext(fileName)
	newName := fmt.Sprintf("%s%s", uuid.Must(uuid.NewV4()), suffix)
	return newName, GetContentType(fileName)
}

func GetContentType(fileName string) string {
	suffix := path.Ext(fileName)
	suffix = strings.ToLower(suffix)
	if val, ok := contentTypeMap[suffix]; ok {
		return val
	}
	return "application/octet-stream"
}

func GetFilePath(fileType int) string {
	cfg := config.Config
	subPath := time.Now().Format("2006-01-02")
	if val, ok := fileTypeDirMap[fileType]; ok {
		return fmt.Sprintf("/%s/%s/%s/", cfg.Station, val, subPath)
	}
	return fmt.Sprintf("/%s/%s/%s/", cfg.Station, "file", subPath)
}

func GetNewFilePath(fileType int) string {
	cfg := config.Config
	subPath := time.Now().Format("2006-01-02")
	if val, ok := fileTypeDirMap[fileType]; ok {
		return fmt.Sprintf("-%s-%s/%s-", cfg.Station, val, subPath)
	}
	return fmt.Sprintf("-%s-%s-%s-", cfg.Station, "file", subPath)
}

func OssResizeImage(url string,x,y,fileType int) string {
	suffix := path.Ext(url)
	suffix = strings.ToLower(suffix)
	if suffix == ".gif" {
		return ""
	}
	var newX, newY int

	if x <= 400 && y <= 400 {
		return url
	}
	if x > 400 && x >= y {
		newX = 400
		newY = int(float64(400) / float64(x) * float64(y))
	}
	if y > 400 && y >= x {
		newX = int(float64(400) / float64(y) * float64(x))
		newY = 400
	}
	switch fileType {
	case 3:
		return fmt.Sprintf("?x-oss-process=image/resize,m_fill,w_%d,quality,q_60",url,newX)
	case 4:
		return fmt.Sprintf("%s?x-oss-process=video/snapshot,t_7000,f_jpg,w_%d,h_%d,m_fast",url,newX,newY)
		}
	return url
}

func GetRealURL(url string) string {
	if strings.HasPrefix(url, "http") || url == "" {
		return url
	}
	return config.Config.Minio.OssPrefix + url
}

var fileTypeDirMap map[int]string = map[int]string{
	1: "file",
	2: "audio",
	3: "image",
	4: "video",
}
var contentTypeMap map[string]string = map[string]string{
	".avi":  "video/avi",
	".aac":  "audio/aac",
	".bmp":  "image/bmp",
	".css":  "text/css",
	".csv":  "text/csv",
	".doc":  "application/msword",
	".dwg":  "application/x-dwg",
	".flv":  "video/x-flv",
	".gif":  "image/gif",
	".htm":  "text/html",
	".tiff": "image/tiff",
	".tif":  "image/tiff",
	".jfif": "image/jpeg",
	".exif": "image/jpeg",
	".pcx":  "image/x-pcx",
	".svg":  "image/svg+xml",
	".psd":  "image/vnd.adobe.photoshop",
	".cdr":  "application/vnd.corel-draw",
	".pcd":  "application/vnd.corel-draw",
	".dxf":  "image/vnd.dxf",
	".ufo":  "image/ufo",
	".eps":  "image/x-eps",
	".ai":   "application/postscript",
	".raw":  "image/raw",
	".wmf":  "image/x-wmf",
	".webp": "image/webp",
	".apng": "image/apng",
	".avif": "image/avif",
	".html": "text/html",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".tga":  "image/x-tga",
	".fpx":  "image/x-xbitmap",
	".js":   "application/x-javascript",
	".json": "application/json",
	".m1v":  "video/x-mpeg",
	".m2v":  "video/x-mpeg",
	".m3u":  "audio/mpegurl",
	".m4a":  "audio/m4a",
	".m4e":  "video/mpeg4",
	".m4v":  "video/x-m4v",
	".mp3":  "audio/mp3",
	".mp4":  "video/mp4",
	".mov":  "video/quicktime",
	".dat":  "video/x-dat",
	".mkv":  "video/mkv",
	".wmv":  "video/wmv",
	".asf":  "video/asf",
	".asx":  "video/x-ms-asf",
	".rm":   "application/vnd.rn-realmedia",
	".rmvb": "video/rmvb",
	".3gp":  "video/3gp",
	".vob":  "video/vob",
	".png":  "image/png",
	".txt":  "text/plain",
	".zip":  "application/zip",
	".gzip": "application/gzip",
}
