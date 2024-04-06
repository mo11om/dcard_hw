package controllers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config stores options for connecting to Redis
type Config struct {
	Addr     string
	Password string
	DB       int
	PoolSize int // Optional: pool size for connection pooling
}

// Client represents a connection to a Redis server
type Client struct {
	pool *redis.ClusterClient // Cluster client for cluster connections
	cli  *redis.Client        // Single client for single node connections
	m    sync.Mutex           // Mutex for thread-safe access
}

// NewClient creates a new Redis client based on the provided configuration
func NewClient(cfg *Config) (*Client, error) {
	var err error
	c := &Client{}

	// if cfg.PoolSize > 0 { // Use connection pool if pool size is specified
	// 	c.pool = redis.NewClusterClient(&redis.ClusterOptions{
	// 		Addrs:    []string{cfg.Addr},
	// 		Password: cfg.Password,
	// 		DB:       cfg.DB,
	// 		PoolSize: cfg.PoolSize,
	// 	})
	// } else { // Use single client otherwise

	// }
	c.cli = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if _, err := c.cli.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}
	return c, err
}

// Set stores a key-value pair in Redis with an optional expiration time
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var err error
	c.m.Lock()
	defer c.m.Unlock()

	if c.pool != nil {
		err = c.pool.Set(ctx, key, value, expiration).Err()
	} else {
		err = c.cli.Set(ctx, key, value, expiration).Err()
	}
	return err
}

// Get retrieves a value for a given key from Redis
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	var err error
	var val string
	c.m.Lock()
	defer c.m.Unlock()

	if c.pool != nil {
		val, err = c.pool.Get(ctx, key).Result()
	} else {
		val, err = c.cli.Get(ctx, key).Result()
	}
	if err == redis.Nil {
		return "", nil // Key not found (handled as nil value)
	}
	return val, err
}

// FlushAll removes all keys from the current database
func (c *Client) FlushAll(ctx context.Context) error {
	c.m.Lock()
	defer c.m.Unlock()

	if c.pool != nil {
		return c.pool.FlushAll(ctx).Err()
	} else {
		return c.cli.FlushAll(ctx).Err()
	}
}

// ... Add similar functions for other Redis operations (HSet, HGetAll, etc.)

// Close closes the Redis connection(s)
func (c *Client) Close() error {
	var err error
	c.m.Lock()
	defer c.m.Unlock()

	if c.pool != nil {
		err = c.pool.Close()
	} else {
		err = c.cli.Close()
	}
	return err
}
