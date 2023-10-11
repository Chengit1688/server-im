package usecase

import (
	"im/internal/cms_api/config/model"
	"im/internal/cms_api/config/repo"
	"im/pkg/common"
	"im/pkg/util"
)

var ConfigUseCase = new(configUseCase)

type configUseCase struct{}

func (c *configUseCase) GetParameterConfig() (p *model.ParameterConfigResp, err error) {
	return repo.ConfigRepo.GetSystemConfig()
}

func (c *configUseCase) GetDepositConfig() (p *model.GetDepositeResp, err error) {
	var deposit model.GetDepositeResp
	configs, err := repo.ConfigRepo.GetDepositInfo()
	if err != nil {
		return nil, err
	}
	for _, config := range configs {
		switch config.Name {
		case repo.ConfigDepositHTML:
			deposit.Html = config.Value
		case repo.ConfigDepositURL:
			deposit.Url = config.Value
		case repo.ConfigDepositSWITCH:
			deposit.Switch = util.String2Int(config.Value)
		}
	}
	return &deposit, nil
}

func (c *configUseCase) GetWithdrawConfig(lang ...string) (p *model.WithdrawConfigResp, err error) {
	if len(lang) > 0{
		return repo.ConfigRepo.GetWithdrawConfig(lang[0])
	}
	return repo.ConfigRepo.GetWithdrawConfig("zh_CN")
}

func (c *configUseCase) GetJPushAuthInfo() (appKey, masterSecret string, err error) {
	configs, err := repo.ConfigRepo.GetJPushAuthInfo()
	for _, config := range configs {
		switch config.Name {
		case repo.ConfigJPushAppKey:
			appKey = config.Value
		case repo.ConfigJPushMasterSecret:
			masterSecret = config.Value
		}
	}
	return
}

func (c *configUseCase) GetFeihuConfig() (data model.GetFeihuResp, err error) {
	configs, err := repo.ConfigRepo.GetFeihuAuthInfo()
	for _, config := range configs {
		switch config.Name {
		case repo.ConfigFeihuAppKey:
			data.AppKey, _ = util.Encrypt([]byte(config.Value), common.ContentKey)
		case repo.ConfigFeihuAppSecret:
			data.AppSecret, _ = util.Encrypt([]byte(config.Value), common.ContentKey)
		}
	}
	return
}

func (c *configUseCase) GetAboutUs() (data string, err error) {
	return repo.ConfigRepo.GetAboutUs()
}

func (c *configUseCase) GetPrivacyPolicy() (data string, err error) {
	return repo.ConfigRepo.GetPrivacyPolicy()
}

func (c *configUseCase) GetUserAgreement() (data string, err error) {
	return repo.ConfigRepo.GetUserAgreement()
}

func (c *configUseCase) GetPrivilegeUserFreezeIsOpen() (status int, err error) {
	return repo.ConfigRepo.GetPrivilegeUserFreezeIsOpen()
}

func (c *configUseCase) SetPrivilegeUserFreezeIsOpen(status int) (err error) {
	return repo.ConfigRepo.SetPrivilegeUserFreezeIsOpen(status)
}
