package redismiddleware

import (
	"rummy-session/repository/models"

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

func (r *RedisRepository) GetAuthToken(userId string) models.AuthToken {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     r.RedisConfig.Addr,
		Password: r.RedisConfig.Password,
		DB:       0,
	})
	val, err := redisClient.Get(userId).Result()
	if err != nil {
		return models.AuthToken{
			Error: err,
		}
	}
	return models.AuthToken{
		AuthToken: val,
	}
}

func (r *RedisRepository) DeleteAuthToken(userId string) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     r.RedisConfig.Addr,
		Password: r.RedisConfig.Password,
		DB:       0,
	})
	err := redisClient.Del(userId).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisRepository) CreateAuthToken(userId, authToken string) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     r.RedisConfig.Addr,
		Password: r.RedisConfig.Password,
		DB:       0,
	})
	err := redisClient.Set(userId, authToken, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisRepository(c Config, ma monitoring.Agent) *RedisRepository {
	return &RedisRepository{
		RedisConfig:     c,
		monitoringAgent: ma,
	}
}
