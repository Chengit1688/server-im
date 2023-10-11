package oss

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/aliyun/aliyun-sts-go-sdk/sts"
	OC "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"im/config"
	"io/ioutil"
	"strings"
	"strconv"
	"net/http"
)

var OS *Client

func InitOss() {
	//初始化OSS客户端
	OSSClient, err := NewClient(Options{
		AccessId:     config.Config.AliYun.AccessKeyId,
		AccessSecret: config.Config.AliYun.AccessKeySecret,
		EndPoint:     config.Config.AliYun.OSS.EndPoint,
		RoleArn:      config.Config.AliYun.STS.RoleArn,
		ExpireTIme:   config.Config.AliYun.STS.ExpireTime,
	})
	if err != nil {
		panic("OSS存储初始化失败")
		return
	}
	OS = OSSClient
}

//初始化参数
type Options struct {
	//访问密钥
	AccessId     string
	AccessSecret string
	EndPoint     string
	RoleArn      string
	ExpireTIme   int64
}



type Client struct {
	opt Options
	oss *OC.Client
}

type Credentials struct {
	Region        string `json:"region"` //访问区域
	BucketName    string `json:"bucketName"`
	AccessId      string `json:"accessId"`
	AccessSecret  string `json:"accessSecret"`
	SecurityToken string `json:"securityToken"`
	Expiration    uint   `json:"expiration"`
	EndPoint      string `json:"endPoint"`
}

//初始化客户端
func NewClient(options Options) (client *Client, err error) {
	oss, err := OC.New(options.EndPoint, options.AccessId, options.AccessSecret)
	if err != nil {
		return nil, err
	}
	client = &Client{
		opt: options,
		oss: oss,
	}
	return client, nil
}

func MD5(file []byte) string {
	crypto := md5.New()
	crypto.Write(file)
	return hex.EncodeToString(crypto.Sum(nil))
}



//获取STS授权证书
func (c *Client) GetSTS(bucketName string, sessionName string) (cred Credentials, err error) {
	region, err := c.oss.GetBucketLocation(bucketName)
	if err != nil {
		return cred, err
	}
	client := sts.NewClient(c.opt.AccessId, c.opt.AccessSecret, c.opt.RoleArn, sessionName)
	res, err := client.AssumeRole(uint(c.opt.ExpireTIme))
	if err != nil {
		return cred, err
	}
	cred = Credentials{
		Region:        region,
		BucketName:    bucketName,
		AccessId:      res.Credentials.AccessKeyId,
		AccessSecret:  res.Credentials.AccessKeySecret,
		SecurityToken: res.Credentials.SecurityToken,
		Expiration:    uint(res.Credentials.Expiration.Unix()),
		EndPoint: c.opt.EndPoint,
	}
	return cred, nil
}

// 上传网络图片到OSS
func (c *Client) UploadNetworkImage(bucketName string, link string) (file FileDetail, err error) {
	return c.UploadNetworkFile(bucketName, link, "image/jpg")
}

// UploadNetworkAmr 上传网络MP3到OSS
func (c *Client) UploadNetworkAmr(bucketName string, link string) (file FileDetail, err error) {
	return c.UploadNetworkFile(bucketName, link, "amr")
}

// 上传网络文件到OSS
func (c *Client) UploadNetworkFile(bucketName string, link string, mimeType string) (file FileDetail, err error) {
	response, err := http.Get(link)
	if err != nil {
		return file, err
	}
	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	mimeList := strings.Split(mimeType, "/")
	if len(mimeList) == 2 {
		file.Format = mimeList[1]
	} else {
		file.Format = mimeType
	}
	//计算文件hash
	file.Name = MD5(data) + "." + file.Format

	bucket, err := c.oss.Bucket(bucketName)
	if err != nil {
		return file, err
	}

	reader := bytes.NewReader(data)
	file.Size = int(reader.Size())
	osMeta := OC.Meta("size", strconv.Itoa(int(reader.Size())))
	contentType := OC.ContentType(mimeType)

	if header, _ := bucket.GetObjectDetailedMeta(file.Name); header != nil {
		sizeStr := header.Get("X-Oss-Meta-Size")
		if sizeStr == "" {
			_ = bucket.SetObjectMeta(file.Name, osMeta)
		}
		return file, nil
	}

	err = bucket.PutObject(file.Name, reader, osMeta, contentType)
	return file, err
}

//分片上传文件
type MultipartCallback func(mediaId string, indexBuf string) (MediaData, error)

func (c *Client) UploadMultipartFile(bucketName string, fileName string, mediaId string, fn MultipartCallback) error {
	bucket, err := c.oss.Bucket(bucketName)
	if err != nil {
		return err
	}

	//如果文件已存在，则直接返回当前文件信息
	if header, _ := bucket.GetObjectDetailedMeta(fileName); header != nil {
		return nil
	}

	isFinish := false
	indexBuf := ""
	repeat := 3
	chunkIndex := 1
	parts := make([]OC.UploadPart, 0)

	storageType := OC.ObjectStorageClass(OC.StorageStandard)
	mulObj, err := bucket.InitiateMultipartUpload(fileName, storageType)
	if err != nil {
		return err
	}

	for !isFinish {
		mediaData, err := fn(mediaId, indexBuf)
		if err != nil {
			//重试3次
			if repeat > 0 {
				repeat--
				continue
			}
			return err
		}

		if mediaData.IsFinish {
			isFinish = mediaData.IsFinish
		}
		indexBuf = mediaData.OutIndexBuf

		reader := bytes.NewReader(mediaData.Data)
		part, err := bucket.UploadPart(mulObj, reader, reader.Size(), chunkIndex)
		if err != nil {
			return err
		}
		parts = append(parts, part)
		chunkIndex++
	}

	objectAcl := OC.ObjectACL(OC.ACLPublicRead)
	_, err = bucket.CompleteMultipartUpload(mulObj, parts, objectAcl)
	return err
}

// UploadFile 上传文件
func (c *Client) UploadFile(bucketName string, filePath string, mimeType string) (file FileDetail, err error) {
	data, _ := ioutil.ReadFile(filePath)
	file.Format = mimeType
	//计算文件hash
	file.Name = MD5(data) + "." + file.Format

	bucket, err := c.oss.Bucket(bucketName)
	if err != nil {
		return file, err
	}

	reader := bytes.NewReader(data)
	file.Size = int(reader.Size())
	osMeta := OC.Meta("size", strconv.Itoa(int(reader.Size())))
	contentType := OC.ContentType(file.Format)

	if header, _ := bucket.GetObjectDetailedMeta(file.Name); header != nil {
		sizeStr := header.Get("X-Oss-Meta-Size")
		if sizeStr == "" {
			_ = bucket.SetObjectMeta(file.Name, osMeta)
		}
		return file, nil
	}

	err = bucket.PutObject(file.Name, reader, osMeta, contentType)
	return file, err
}

