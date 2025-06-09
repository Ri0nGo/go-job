package qq

import (
	"context"
	"fmt"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/pkg/config"
	"testing"
)

func TestGetAccessToken(t *testing.T) {
	configPath := utils.GetMasterConfigPath()
	err := config.InitConfig(configPath)
	if err != nil {
		t.Error(err)
	}
	auth2Service := NewOAuth2Service(
		config.App.OAuth2[model.AuthTypeQQ.String()].ClientID,
		config.App.OAuth2[model.AuthTypeQQ.String()].ClientSecret,
		config.App.OAuth2[model.AuthTypeQQ.String()].RedirectURL,
	)
	authUrl := auth2Service.GetAuthUrl(context.Background(), "test-state")
	fmt.Println(authUrl)
}
