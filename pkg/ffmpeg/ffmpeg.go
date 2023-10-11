package ffmpeg

import (
	"bytes"
	"fmt"
	"im/pkg/logger"
	"im/pkg/util"
	_ "image"
	_ "image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gofrs/uuid"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		logger.Sugar.Errorw(util.GetSelfFuncName(), "error", err.Error())
	}
	return buf
}

func GetFirstFrame(src string) string {
	reader := ExampleReadFrameAsJpeg(src, 1)
	img, err := imaging.Decode(reader)
	if err != nil {
		logger.Sugar.Errorw(util.GetSelfFuncName(), "error", err.Error())

	}
	dst := fmt.Sprintf("%s/%s%s", os.TempDir(), uuid.Must(uuid.NewV4()), ".jpeg")
	err = imaging.Save(img, dst)
	if err != nil {
		logger.Sugar.Errorw(util.GetSelfFuncName(), "error", err.Error())
	}
	return dst
}

func ResizeImage(src string) string {
	suffix := path.Ext(src)
	suffix = strings.ToLower(suffix)
	var dst string
	if suffix == ".gif" {
		return src
	} else {
		dst = fmt.Sprintf("%s/%s%s", os.TempDir(), uuid.Must(uuid.NewV4()), ".jpeg")
	}
	var newX, newY int
	img, err := imaging.Open(src, imaging.AutoOrientation(true))
	if err != nil {
		logger.Sugar.Errorw(util.GetSelfFuncName(), "error", err.Error())
		return src
	}
	x := img.Bounds().Dx()
	y := img.Bounds().Dy()
	if x <= 400 && y <= 400 {
		return src
	}
	if x > 400 && x >= y {
		newX = 400
		newY = int(float64(400) / float64(x) * float64(y))
	}
	if y > 400 && y >= x {
		newX = int(float64(400) / float64(y) * float64(x))
		newY = 400
	}
	dstImage := imaging.Resize(img, newX, newY, imaging.Lanczos)
	err = imaging.Save(dstImage, dst)
	if err != nil {
		logger.Sugar.Errorw(util.GetSelfFuncName(), "error", err.Error())
		return src
	}
	return dst
}
