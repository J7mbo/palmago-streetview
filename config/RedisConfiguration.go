package config

/* RedisConfiguration contains the configuration for use when connecting to redis. */
type RedisConfiguration struct {
	hostname   string `env:"REDIS_HOST" default:"palmago-redis"`
	port       int    `env:"REDIS_port" default:"6379"`
	retryDelay int    `env:"REDIS_RETRY_DELAY" default:"10"`
	maxRetries int    `env:"REDIS_MAX_RETRIES" default:"5"`
}

func (c *RedisConfiguration) GetHostname() string { return c.hostname }
func (c *RedisConfiguration) GetPort() int        { return c.port }
func (c *RedisConfiguration) GetRetryDelay() int  { return c.retryDelay }
func (c *RedisConfiguration) GetMaxRetries() int  { return c.maxRetries }
