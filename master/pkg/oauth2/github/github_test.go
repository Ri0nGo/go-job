package github

import (
	"fmt"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/config"
	"testing"
)

func TestNewOAuth2Service(t *testing.T) {
	err := config.InitConfig(utils.GetMasterConfigPath())
	if err != nil {
		panic(err)
	}
	fmt.Println(config.App.OAuth2)
}
