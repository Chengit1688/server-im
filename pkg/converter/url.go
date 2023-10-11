package converter

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"strings"
	"time"
)

var options = launcher.New().Set("--no-sandbox").Set("--disable-setuid-sandbox").MustLaunch()
var browser = rod.New().ControlURL(options).MustConnect()

// URL2Text URL地址转文本
func URL2Text(link string) (u Url, err error) {
	//避免页面打开失败导致程序panic
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("链接采集失败：%v", err)
		}
	}()

	page := browser.Timeout(10 * time.Second).MustPage(link)
	if err = page.WaitLoad(); err != nil {
		return u, err
	}

	sleeper := func() utils.Sleeper {
		retryCount := 0
		maxRetryCount := 3

		return func(context.Context) error {
			if retryCount < maxRetryCount {
				retryCount += 1
				time.Sleep(time.Second / 10)
				return nil
			}
			return errors.New("gather timeout: dom not found")
		}
	}
	titleTarget, err := page.Sleeper(sleeper).Element("title")
	if titleTarget != nil && err == nil {
		u.Title = titleTarget.MustText()
	}
	//修复微信文章无法获取标题的问题
	if u.Title == "" || u.Title == "\n" {
		ogTitleTarget, err := page.Sleeper(sleeper).Element(`meta[property="og:title"]`)
		if ogTitleTarget != nil && err == nil {
			u.Title = *ogTitleTarget.MustAttribute("content")
		}
	}
	descTarget, err := page.Sleeper(sleeper).Element(`meta[name="description"]`)
	if descTarget != nil && err == nil {
		u.Description = *descTarget.MustAttribute("content")
	}

	faviconTarget, err := page.Sleeper(sleeper).Element(`link[rel*="icon"]`)
	if faviconTarget != nil && err == nil {
		u.Favicon = *faviconTarget.MustAttribute("href")
		isRelative := strings.HasPrefix(u.Favicon, "/")
		if isRelative && !strings.HasPrefix(u.Favicon, "//") {
			domain := faviconTarget.MustEval("document.location.origin").String()
			u.Favicon = domain + u.Favicon
		}
	}

	bodyTarget, err := page.Sleeper(sleeper).Element("body")
	if bodyTarget != nil && err == nil {
		content := bodyTarget.MustText()
		//u.Tags = label.Extract(content)
		u.BodyText = content
	}

	//检索图片地址
	imagesTarget, err := page.Sleeper(sleeper).Elements("img")
	if imagesTarget != nil && err == nil {
		maxArea := 0
		imageLink := ""

		for _, image := range imagesTarget {
			width := image.MustEval("parseInt(this.clientWidth)").Int()
			height := image.MustEval("parseInt(this.clientHeight)").Int()
			area := width * height
			link := image.MustAttribute("src")

			if area > maxArea && link != nil && !strings.HasPrefix(*link, "data:image") {
				isRelative := strings.HasPrefix(*link, "/")
				if isRelative && !strings.HasPrefix(*link, "//") {
					domain := image.MustEval("document.location.origin").String()
					imageLink = domain + *link
				} else {
					imageLink = *link
				}
				maxArea = area
			}
		}

		u.Image = imageLink
	}

	return u, err
}

type LinkParam struct {
	Url string `json:"url" form:"url"`
}

// Url 链接
type Url struct {
	Title       string
	Description string
	Favicon     string
	Image       string
	Tags        []string
	BodyText    string
}
