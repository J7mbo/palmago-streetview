package config

/* GrpcServerConfiguration contains the configuration for use in starting the gRPC server. */
type GrpcServerConfiguration struct {
	protocol string `env:"GRPC_SERVER_PROTOCOL" default:"tcp"`
	/* If you're running this on localhost, not in docker, the server host should be: "". */
	host       string `env:"GRPC_SERVER_HOST" default:""`
	port       int    `env:"GRPC_SERVER_PORT" default:"4000"`
	retryDelay int    `env:"GRPC_SERVER_RETRY_DELAY" default:"5"`
	maxRetries int    `env:"GRPC_SERVER_MAX_RETRIES" default:"10"`
}

func (c *GrpcServerConfiguration) GetProtocol() string { return c.protocol }
func (c *GrpcServerConfiguration) GetHost() string     { return c.host }
func (c *GrpcServerConfiguration) GetPort() int        { return c.port }
func (c *GrpcServerConfiguration) GetRetryDelay() int  { return c.retryDelay }
func (c *GrpcServerConfiguration) GetMaxRetries() int  { return c.maxRetries }
