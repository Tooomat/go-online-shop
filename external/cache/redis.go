package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(addr string, password string, db int) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := r.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return r, nil
}

type RedisRepository interface {
	SaveRefreshToken(ctx context.Context, id, token string, ttl time.Duration) (err error)
	GetRefreshToken(ctx context.Context, id string) (res string, err error)
	DeleteKey(ctx context.Context, id string) (err error)
	BlacklistAccessToken(ctx context.Context, token string, ttl time.Duration) (err error)
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type redisRepository struct {
	rdb *redis.Client
}

func NewRedisRepository(rdb *redis.Client) RedisRepository {
	return redisRepository{
		rdb: rdb,
	}
}

func (r redisRepository) SaveRefreshToken(ctx context.Context, id, token string, ttl time.Duration) (err error) {
	return r.rdb.Set(ctx, "refresh:"+id, token, ttl).Err()
}

func (r redisRepository) GetRefreshToken(ctx context.Context, id string) (res string, err error) {
	return r.rdb.Get(ctx, "refresh:"+id).Result()
}

func (r redisRepository) DeleteKey(ctx context.Context, id string) (err error) {
	return r.rdb.Del(ctx, id).Err()
}

func (r redisRepository) BlacklistAccessToken(ctx context.Context, token string, ttl time.Duration) (err error) {
	return r.rdb.Set(ctx, "bl_"+token, "1", ttl).Err()
}

func (r redisRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	_, err := r.rdb.Get(ctx, "bl_"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
