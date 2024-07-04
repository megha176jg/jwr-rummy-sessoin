package redismiddleware

import (
	"bitbucket.org/junglee_games/getsetgo/monitoring"
	"github.com/go-redis/redis"
)

type Config struct {
	Addr     string
	Password string
}
type RedisRepository struct {
	RedisConfig     Config
	monitoringAgent monitoring.Agent
}

func (r *RedisRepository) GetTitle(name string) (string, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     r.RedisConfig.Addr,
		Password: r.RedisConfig.Password,
		DB:       0,
	})
	val, err := redisClient.Get(name).Result()
	if err != nil {
		return "", err
	}

	_, err = redisClient.Set(name, val, 0).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func NewRedisRepository(c Config, ma monitoring.Agent) *RedisRepository {
	return &RedisRepository{
		RedisConfig:     c,
		monitoringAgent: ma,
	}
}
