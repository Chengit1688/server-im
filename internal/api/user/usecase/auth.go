package usecase

import (
	"fmt"
	"im/config"
	"im/pkg/mqtt"
)

var AuthUseCase = new(authUseCase)

type authUseCase struct{}

func (c *authUseCase) CreateAuthUsername(userID string) (err error) {
	station := config.Config.Station
	username := fmt.Sprintf("%s_%s", station, userID)
	password := "root1234"
	return mqtt.CreateAuthUsername(username, password)
}
