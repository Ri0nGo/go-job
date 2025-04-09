package utils

import "time"

func TimestampToTime(ts int64) time.Time {
	switch {
	case ts > 1e18:
		return time.Unix(0, ts).In(time.Local) // 纳秒
	case ts > 1e12:
		return time.UnixMilli(ts).In(time.Local) // 毫秒
	default:
		return time.Unix(ts, 0).In(time.Local) // 秒
	}
}
