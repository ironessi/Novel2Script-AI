package redis

import (
	"context"
	"fmt"
	"time"

	"novel2script-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// Init 初始化 Redis 连接
func Init() error {
	cfg := config.C.Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	return nil
}

// GetClient 获取 Redis 客户端
func GetClient() *redis.Client {
	return rdb
}

// TaskProgressKey 任务进度 Key
func TaskProgressKey(taskId int64) string {
	return fmt.Sprintf("task:%d:progress", taskId)
}

// SetTaskProgress 设置任务进度
func SetTaskProgress(ctx context.Context, taskId int64, progress int, step string) error {
	key := TaskProgressKey(taskId)
	data := fmt.Sprintf(`{"progress":%d,"step":"%s","updated_at":"%s"}`,
		progress, step, time.Now().Format(time.RFC3339))
	return rdb.Set(ctx, key, data, 24*time.Hour).Err()
}

// GetTaskProgress 获取任务进度
func GetTaskProgress(ctx context.Context, taskId int64) (string, error) {
	key := TaskProgressKey(taskId)
	return rdb.Get(ctx, key).Result()
}

// DeleteTaskProgress 删除任务进度
func DeleteTaskProgress(ctx context.Context, taskId int64) error {
	key := TaskProgressKey(taskId)
	return rdb.Del(ctx, key).Err()
}
