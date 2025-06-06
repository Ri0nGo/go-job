package consts

import "time"

const (
	DefaultLoginJwtExpireTime     = time.Hour * 2
	DefaultAuthStateJwtExpireTime = time.Minute * 10
)
