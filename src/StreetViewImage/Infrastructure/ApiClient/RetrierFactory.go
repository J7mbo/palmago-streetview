package ApiClient

import (
	"app/config"
	"github.com/j7mbo/MethodCallRetrier/v2"
	"time"
)

/* RetrierFactory is responsible for creating a Retrier (third-party lib). Can be DI'd for SoC. */
type RetrierFactory struct{}

/* Create creates a new MethodCallRetrier.Retrier given a config.StreetViewApiConfiguration. */
func (*RetrierFactory) Create(config *config.StreetViewApiConfiguration) MethodCallRetrier.Retrier {
	return MethodCallRetrier.New(time.Duration(config.GetRetryDelay())*time.Second, int64(config.GetMaxRetries()), 1)
}
