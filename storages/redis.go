package storages

import (
	"log"
	"sync"

	"github.com/go-redis/redis"
)

// RedisClient is a pointer to the redis client
type RedisClient struct {
	*redis.Client
}

// key name for set in redis
const key = "drivers"

// For a singleton pattern, use Once.Do of sync package
var once sync.Once
var redisClient *RedisClient

// GetRedisClient connects and returns a redis client instance
func GetRedisClient() *RedisClient {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		redisClient = &RedisClient{client}
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to redis %v", err)
	}
	return redisClient
}


// AddDriverLocation upserts a driver location in redis
func (c *RedisClient) AddDriverLocation(lon, lat float64, id string) {
	c.GeoAdd(
		key,
		&redis.GeoLocation{Longitude: lon, Latitude: lat, Name: id},
	)
}

// RemoveDriverLocation removes a driver location from redis
func (c *RedisClient) RemoveDriverLocation(id string) {
	c.ZRem(key, id)
}

// SearchDrivers returns a list of redis geolocations
// from a specific lat lon within r radius
func (c *RedisClient) SearchDrivers(limit int, lat, lon, r float64) []redis.GeoLocation {
	res, _ := c.GeoRadius(key, lon, lat, &redis.GeoRadiusQuery{
		Radius:      r,
		Unit:        "m",
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
		Count:       limit,
		Sort:        "ASC",
	}).Result()
	return res

}
