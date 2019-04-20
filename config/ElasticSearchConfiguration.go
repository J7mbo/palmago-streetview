package config

/* ElasticSearchConfiguration contains the configuration for use in the application logger (ElasticSearchLogger). */
type ElasticSearchConfiguration struct {
	host       string `env:"ELASTICSEARCH_HOST" default:"localhost"`
	port       int    `env:"ELASTICSEARCH_PORT" default:"9200"`
	index      string `env:"ELASTICSEARCH_INDEX" default:"palmago"`
	retryDelay int    `env:"ELASTICSEARCH_RETRY_DELAY" default:"1"`
	maxRetries int    `env:"ELASTICSEARCH_MAX_RETRIES" default:"3"`
}

func (c *ElasticSearchConfiguration) GetHost() string    { return c.host }
func (c *ElasticSearchConfiguration) GetPort() int       { return c.port }
func (c *ElasticSearchConfiguration) GetIndex() string   { return c.index }
func (c *ElasticSearchConfiguration) GetRetryDelay() int { return c.retryDelay }
func (c *ElasticSearchConfiguration) GetMaxRetries() int { return c.maxRetries }
