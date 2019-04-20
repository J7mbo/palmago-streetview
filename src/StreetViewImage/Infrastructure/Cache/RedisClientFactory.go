package Cache

import (
	"app/config"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

/* RedisClientFactory exists to delay the initialisation of a Redis client that performs connections in it's init. */
type RedisClientFactory struct {
	config *config.RedisConfiguration
}

/* NewRedisClientFactory returns a newly initialised RedisClientFactory ready to initialise a client at runtime. */
func NewRedisClientFactory(config config.RedisConfiguration) *RedisClientFactory {
	return &RedisClientFactory{config: &config}
}

/* Create returns an initialised redis.Client if the ping executes successfully, otherwise errors. */
func (f *RedisClientFactory) Create() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       f.formatAddress(f.config.GetHostname(), f.config.GetPort()),
		MaxRetries: f.config.GetMaxRetries(),
		/* Looks like jitter on the backoff could be client-specified unfortunately. Oh well. */
		MinRetryBackoff: time.Duration(float64(f.config.GetRetryDelay()) * time.Second.Seconds()),
		MaxRetryBackoff: time.Duration(float64(f.config.GetRetryDelay()) * time.Second.Seconds()),
	})

	if result := client.Ping(); result.Err() != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"unable to connect to redis on host: %s, port :%d, error: %s",
				f.config.GetHostname(), f.config.GetPort(), result.Err().Error(),
			),
		)
	}

	return client, nil
}

/* formatAddress formats the host and port into an address for the Addr field of redis.Options. */
func (*RedisClientFactory) formatAddress(hostname string, port int) string {
	return fmt.Sprintf("%s:%d", hostname, port)
}
