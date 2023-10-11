package usecase

import (
	"fmt"
	"im/config"
	settingRepo "im/internal/api/setting/repo"
	"im/pkg/code"
	"im/pkg/logger"
	"im/pkg/sms"
	"im/pkg/util"
	"image/color"
	"math/rand"
	"sync"
	"time"

	captchaConfig "github.com/TestsLing/aj-captcha-go/config"
	"github.com/TestsLing/aj-captcha-go/service"
)

var (
	captchaFactory *service.CaptchaServiceFactory
	lock           = &sync.Mutex{}
)

func GetCaptchaFactory() *service.CaptchaServiceFactory {
	if captchaFactory == nil {
		lock.Lock()
		defer lock.Unlock()
		var (
			watermarkConfig = &captchaConfig.WatermarkConfig{
				FontSize: 12,
				Color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
				Text:     config.Config.Captcha.DefaultText,
			}
			clickWordConfig = &captchaConfig.ClickWordConfig{
				FontSize: 25,
				FontNum:  4,
			}
			blockPuzzleConfig = &captchaConfig.BlockPuzzleConfig{Offset: 10}
		)
		var configCap = captchaConfig.BuildConfig("redis", config.Config.Captcha.DefaultResourceRoot, watermarkConfig,
			clickWordConfig, blockPuzzleConfig, config.Config.Captcha.CacheExpireSec)
		captchaFactory = service.NewCaptchaServiceFactory(configCap)
		captchaFactory.RegisterCache("redis", service.NewConfigRedisCacheService([]string{config.Config.Redis.Address}, "", config.Config.Redis.Password, false, 0))
		captchaFactory.RegisterService("blockPuzzle", service.NewBlockPuzzleCaptchaService(captchaFactory))
		captchaFactory.RegisterService("clickWord", service.NewClickWordCaptchaService(captchaFactory))
	}
	return captchaFactory
}

func GenSmsCode(operationId, phoneNumber, countryCode, smsServer string) error {
	var (
		err error
		ok  bool
	)
	ok, err = settingRepo.SettingCache.IsAccountCodeExists(phoneNumber + "_repeat")
	if ok || err != nil {
		logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", "短信验证码1分钟后重试")
		return code.ErrBadSmsSend
	}
	rand.Seed(time.Now().UnixNano())
	smsCode := 100000 + rand.Intn(900000)
	if err = sms.NewSmsFactory().GetSms(smsServer).Send(phoneNumber, countryCode, util.IntToString(smsCode)); err != nil {
		logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("Send message, error: %v", err))
		return code.ErrSmsSend
	}

	logger.Sugar.Infow(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("phone_number: %s%s, smsCode: %d", countryCode, phoneNumber, smsCode))
	err = settingRepo.SettingCache.SetAccountCode(phoneNumber, smsCode, config.Config.Sms.ExpireTTL)
	if err != nil {
		logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("set redis,error: %v", err))
		return code.ErrSmsSend
	}
	err = settingRepo.SettingCache.SetAccountCode(phoneNumber+"_repeat", smsCode, config.Config.Sms.CodeTTL)
	if err != nil {
		logger.Sugar.Errorw(operationId, "func", util.GetSelfFuncName(), "error", fmt.Sprintf("set repeat in redis,error: %v", err))
		return code.ErrSmsSend
	}

	return nil
}
