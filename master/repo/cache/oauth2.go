package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go-job/internal/iface/oauth2"
	"go-job/internal/model"
	"strconv"
	"time"
)

var ErrOAuth2StateExpire = errors.New("state is expired")
var ErrOAuth2IdentityExpire = errors.New("auth2 identity is expired")

var (
	defaultStateTTL    = time.Minute * 5
	defaultStatePrefix = "oauth2:state:"
	defaultAuthPrefix  = "oauth2:auth:"
)

type OAuth2StateCache struct {
	redisCache redis.Cmdable
	ttl        time.Duration
}

func NewOAuth2StateCache(redisCache redis.Cmdable) oauth2.IOAuth2Cache {
	return &OAuth2StateCache{
		redisCache: redisCache,
		ttl:        defaultStateTTL,
	}
}

// Set 设置state
func (c *OAuth2StateCache) SetState(ctx context.Context, state string, oauth2State model.OAuth2State) error {
	return c.hset(ctx, c.getStateKey(state), c.oauth2StateToMap(oauth2State), c.ttl)
}

func (c *OAuth2StateCache) SetAuth(ctx context.Context, key string, val map[string]string, ttl time.Duration) error {
	return c.hset(ctx, c.getAuthKey(key), val, ttl)
}

// Get 获取state
func (c *OAuth2StateCache) GetState(ctx context.Context, state string) (model.OAuth2State, error) {
	result, err := c.redisCache.HGetAll(ctx, c.getStateKey(state)).Result()
	if err != nil { // redis 执行报错了
		return model.OAuth2State{}, err
	}
	if len(result) == 0 {
		return model.OAuth2State{}, ErrOAuth2StateExpire
	}
	return c.mapToOauth2State(result), nil
}

func (c *OAuth2StateCache) GetAuth(ctx context.Context, key string) (map[string]string, error) {
	result, err := c.redisCache.HGetAll(ctx, c.getAuthKey(key)).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrOAuth2IdentityExpire
	}
	return result, nil
}

// MarkUsed 将state标记为已使用
func (c *OAuth2StateCache) MarkUsed(ctx context.Context, state string) error {
	return c.redisCache.HSet(ctx, c.getStateKey(state), "used", "true").Err()
}

func (c *OAuth2StateCache) getStateKey(state string) string {
	return defaultStatePrefix + state
}

func (c *OAuth2StateCache) getAuthKey(key string) string {
	return defaultAuthPrefix + key
}

func (c *OAuth2StateCache) oauth2StateToMap(oauth2State model.OAuth2State) map[string]string {
	return map[string]string{
		"uid":      strconv.Itoa(oauth2State.Uid),
		"state":    oauth2State.State,
		"scene":    string(oauth2State.Scene),
		"platform": oauth2State.Platform,
		"used":     strconv.FormatBool(oauth2State.Used),
	}
}

func (c *OAuth2StateCache) hset(ctx context.Context, key string, value any, ttl time.Duration) error {
	if err := c.redisCache.HSet(ctx, key, value).Err(); err != nil {
		return err
	}
	return c.redisCache.Expire(ctx, key, ttl).Err()
}

func (c *OAuth2StateCache) mapToOauth2State(result map[string]string) model.OAuth2State {
	return model.OAuth2State{
		State:    result["state"],
		Scene:    model.Auth2Scene(result["scene"]),
		Platform: result["platform"],
		Used:     result["used"] == "true",
	}
}
