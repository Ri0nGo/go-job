package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go-job/internal/model"
	"strconv"
	"time"
)

var ErrOAuth2StateExpire = errors.New("state is expired")

var (
	defaultStateTTL    = time.Minute * 5
	defaultStatePrefix = "oauth2:state:"
)

type OAuth2StateCache struct {
	redisCache redis.Cmdable
	prefix     string
	ttl        time.Duration
}

func NewOAuth2StateCache(redisCache redis.Cmdable) *OAuth2StateCache {
	return &OAuth2StateCache{
		redisCache: redisCache,
		prefix:     defaultStatePrefix,
		ttl:        defaultStateTTL,
	}
}

// Set 设置state
func (c *OAuth2StateCache) Set(ctx context.Context, state string, oauth2State model.OAuth2State) error {
	return c.hset(ctx, c.getKey(state), c.oauth2StateToMap(oauth2State), c.ttl)
}

// Get 获取state
func (c *OAuth2StateCache) Get(ctx context.Context, state string) (model.OAuth2State, error) {
	result, err := c.redisCache.HGetAll(ctx, c.getKey(state)).Result()
	if err != nil { // redis 执行报错了
		return model.OAuth2State{}, err
	}
	if len(result) == 0 {
		return model.OAuth2State{}, ErrOAuth2StateExpire
	}
	return c.mapToOauth2State(result), nil
}

// MarkUsed 将state标记为已使用
func (c *OAuth2StateCache) MarkUsed(ctx context.Context, state string) error {
	return c.redisCache.HSet(ctx, c.getKey(state), "used", "true").Err()
}

func (c *OAuth2StateCache) getKey(state string) string {
	return c.prefix + state
}

func (c *OAuth2StateCache) oauth2StateToMap(oauth2State model.OAuth2State) map[string]string {
	return map[string]string{
		"state":         oauth2State.State,
		"scene":         string(oauth2State.Scene),
		"redirect_page": oauth2State.RedirectPage,
		"platform":      oauth2State.Platform,
		"used":          strconv.FormatBool(oauth2State.Used),
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
		State:        result["state"],
		Scene:        model.Auth2Scene(result["scene"]),
		RedirectPage: result["redirect_page"],
		Platform:     result["platform"],
		Used:         result["used"] == "true",
	}
}
