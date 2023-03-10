package redisUtil

import (
	"github.com/go-redis/redis/v8"
	"github.com/sunmfei/mfus/common/MFei"
	"github.com/sunmfei/mfus/config"
	"golang.org/x/net/context"
	"time"
)

type RedisPlay struct {
	myRedis *redis.Client
}

// Setup redis配置
func Setup() *RedisPlay {
	redisConf := config.NewConfig().RedisConf
	var rp = &RedisPlay{
		myRedis: redis.NewClient(&redis.Options{
			Addr:     redisConf.Host,
			Password: redisConf.Password, // no password set
			DB:       redisConf.DB,       // use default sun
		})}
	_, err := rp.myRedis.Ping(context.Background()).Result()
	if err != nil {
		MFei.LOGGER.Error("Redis connect ping failed, err:", err)
		return nil
	}
	MFei.LOGGER.Info("Redis connect succeeded")
	return rp
}

func (rp *RedisPlay) SetRedis(key string, value string, expire time.Duration) bool {
	if err := rp.myRedis.Set(context.Background(), key, value, expire).Err(); err != nil {
		MFei.LOGGER.Error("保存信息失败, err:", err)
		return false
	}
	return true
}

func (rp *RedisPlay) GetRedis(key string) string {
	result, err := rp.myRedis.Get(context.Background(), key).Result()
	if err != nil {
		MFei.LOGGER.Error("查询信息失败, err:", err)
		return ""
	}
	return result
}

func (rp *RedisPlay) DelRedis(key string) bool {
	_, err := rp.myRedis.Del(context.Background(), key).Result()
	if err != nil {
		MFei.LOGGER.Error("删除信息失败, err:", err)
		return false
	}
	return true
}

func (rp *RedisPlay) ExpireRedis(key string, t int64) bool {
	// 延长过期时间
	expire := time.Duration(t) * time.Second
	if err := rp.myRedis.Expire(context.Background(), key, expire).Err(); err != nil {
		MFei.LOGGER.Error("延长过期时间失败, err:", err)
		return false
	}
	return true
}
