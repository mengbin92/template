package cache


import (
    "context"
    "errors"
    "sync"
    "time"

    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
)

// ErrLockNotAcquired 获取分布式锁失败
var ErrLockNotAcquired = errors.New("redis: lock not acquired")

type RedisConfig struct {
    Addr        string
    Password    string
    DB          int
    PoolSize    int
    DialTimeout time.Duration
    ReadTimeout time.Duration
    WriteTimeout time.Duration
}

type RedisClient struct {
    rdb   *redis.Client
    logger *zap.Logger
}

// 单例管理
var (
    mu       sync.Mutex
    instances = make(map[string]*RedisClient)
)

// GetRedisClient 获取或创建单例
func GetRedisClient(name string, cfg *RedisConfig, logger *zap.Logger) *RedisClient {
    mu.Lock()
    defer mu.Unlock()
    if c, ok := instances[name]; ok {
        return c
    }
    opts := &redis.Options{
        Addr:        cfg.Addr,
        Password:    cfg.Password,
        DB:          cfg.DB,
        PoolSize:    cfg.PoolSize,
        DialTimeout: cfg.DialTimeout,
        ReadTimeout: cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
    }
    rdb := redis.NewClient(opts)
    c := &RedisClient{rdb: rdb, logger: logger}
    instances[name] = c
    return c
}

// WithTimeout 在操作时支持自定义超时
func (c *RedisClient) WithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(context.Background(), d)
}

// Close 关闭
func (c *RedisClient) Close() error {
    return c.rdb.Close()
}

// Ping 心跳
func (c *RedisClient) Ping(ctx context.Context) error {
    return c.rdb.Ping(ctx).Err()
}

// ==================== Key 操作 ====================

// Exists 检查 key 是否存在
func (c *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
    n, err := c.rdb.Exists(ctx, key).Result()
    return n > 0, err
}

// TTL 获取剩余存活时间
func (c *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
    return c.rdb.TTL(ctx, key).Result()
}

// Expire 设置过期时间
func (c *RedisClient) Expire(ctx context.Context, key string, exp time.Duration) (bool, error) {
    return c.rdb.Expire(ctx, key, exp).Result()
}

// ==================== String ====================

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
    return c.rdb.Get(ctx, key).Result()
}

// ==================== Hash ====================
func (c *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
    return c.rdb.HSet(ctx, key, values...).Result()
}

func (c *RedisClient) HGet(ctx context.Context, key, field string) (string, error) {
    return c.rdb.HGet(ctx, key, field).Result()
}

// ==================== List ====================
func (c *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
    return c.rdb.LPush(ctx, key, values...).Result()
}

func (c *RedisClient) RPop(ctx context.Context, key string) (string, error) {
    return c.rdb.RPop(ctx, key).Result()
}

// ==================== Sorted Set ====================
// ZAdd 添加
func (c *RedisClient) ZAdd(ctx context.Context, key string, members ...redis.Z) (int64, error) {
    return c.rdb.ZAdd(ctx, key, members...).Result()
}

// ZRangeByScore 按分数区间取
func (c *RedisClient) ZRangeByScore(ctx context.Context, key, min, max string) ([]string, error) {
    return c.rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: min, Max: max}).Result()
}

// ZRem 删除成员
func (c *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
    return c.rdb.ZRem(ctx, key, members...).Result()
}

// ==================== Bitmap ====================

// SetBit 设置位
func (c *RedisClient) SetBit(ctx context.Context, key string, offset int64, value int) (int64, error) {
    return c.rdb.SetBit(ctx, key, offset, value).Result()
}

// GetBit 获取位
func (c *RedisClient) GetBit(ctx context.Context, key string, offset int64) (int64, error) {
    return c.rdb.GetBit(ctx, key, offset).Result()
}

// ==================== Geo ====================

// GeoAdd 添加地理位置
func (c *RedisClient) GeoAdd(ctx context.Context, key string, geoLocation ...*redis.GeoLocation) (int64, error) {
    return c.rdb.GeoAdd(ctx, key, geoLocation...).Result()
}

// GeoRadius 按半径查询
func (c *RedisClient) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) ([]redis.GeoLocation, error) {
    return c.rdb.GeoRadius(ctx, key, longitude, latitude, query).Result()
}

// ==================== 脚本执行 ====================

// Eval 执行 Lua 脚本
func (c *RedisClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
    return c.rdb.Eval(ctx, script, keys, args...).Result()
}

// ==================== 分布式锁 ====================

// TryLock 尝试获取锁，成功返回 unlock 函数
func (c *RedisClient) TryLock(ctx context.Context, key string, ttl time.Duration) (func(), error) {
    ok, err := c.rdb.SetNX(ctx, key, "locked", ttl).Result()
    if err != nil {
        return nil, err
    }
    if !ok {
        return nil, ErrLockNotAcquired
    }
    unlock := func() {
        if err := c.rdb.Del(ctx, key).Err(); err != nil {
            c.logger.Error("unlock failed", zap.Error(err), zap.String("key", key))
        }
    }
    return unlock, nil
}

// ==================== 事务（Pipeline & Tx） ====================

// Pipeline 管道
func (c *RedisClient) Pipeline(ctx context.Context, fn func(pipe redis.Pipeliner) error) ([]redis.Cmder, error) {
    pipe := c.rdb.Pipeline()
    if err := fn(pipe); err != nil {
        return nil, err
    }
    return pipe.Exec(ctx)
}

// Tx 执行事务
func (c *RedisClient) Tx(ctx context.Context, fn func(tx *redis.Tx) error) error {
    return c.rdb.Watch(ctx, fn)
}
