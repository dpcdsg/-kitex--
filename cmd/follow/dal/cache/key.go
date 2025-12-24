package cache

import "strconv"

const (
	followKey   = "follow:"
	followerKey = "follower:"
	LimiterKey  = "limiter"
)

func FollowListKey(uid int64) string {
	return followKey + strconv.FormatInt(uid, 10)
}

func FollowerListKey(uid int64) string {
	return followerKey + strconv.FormatInt(uid, 10)
}

//基于git，方便更新
