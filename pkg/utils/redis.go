package utils

import "github.com/go-redis/redis/v9"

func MakeRedisClient(uri string, isClusterMode bool) (redis.UniversalClient, error) {
	opt, err := redis.ParseURL(uri)
	if err != nil {
		return nil, err
	}

	if isClusterMode {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{opt.Addr},
			Password: opt.Password,
		}), nil
	}

	rdb := redis.NewClient(opt)
	return rdb, err
}
