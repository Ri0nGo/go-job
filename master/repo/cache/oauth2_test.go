package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go-job/internal/model"
	"go-job/internal/pkg/utils"
	"go-job/master/database"
	"go-job/master/pkg/config"
	"testing"
)

func TestOAuth2State(t *testing.T) {
	configPath := utils.GetMasterConfigPath()
	config.InitConfig(configPath)

	cmdable := database.NewRedisClient()
	stateCache := NewOAuth2StateCache(cmdable)
	state := "state-01"
	oldAuth2State := model.OAuth2State{
		State:        state,
		Scene:        "scene-01",
		RedirectPage: "redirect_page-01",
		Platform:     "platform-01",
		Used:         false,
	}
	err := stateCache.Set(context.Background(), state, oldAuth2State)
	assert.NoError(t, err)

	err = stateCache.MarkUsed(context.Background(), state)
	assert.NoError(t, err)

	auth2State, err := stateCache.Get(context.Background(), state)
	assert.NoError(t, err)

	assert.Equal(t, auth2State.Used, true)

	oldAuth2State.Used = true
	assert.Equal(t, oldAuth2State, auth2State)

}
