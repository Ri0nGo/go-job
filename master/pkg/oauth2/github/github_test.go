package github

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
		config.App.OAuth2[model.AuthTypeGithub.String()].ClientID,
		config.App.OAuth2[model.AuthTypeGithub.String()].ClientSecret,
		config.App.OAuth2[model.AuthTypeGithub.String()].RedirectURL,
	)
	authUrl := auth2Service.GetAuthUrl(context.Background(), "test-state")

	fmt.Println(authUrl)
	code := "9eeeda312d45f9a939ad"
	accessToken, err := auth2Service.getAccessToken(context.Background(), code)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(accessToken)
}
