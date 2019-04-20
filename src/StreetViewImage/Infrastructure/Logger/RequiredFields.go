package Logger

import (
	"github.com/google/uuid"
)

/* These fields MUST be consistently present in the logs. They are merged in with LogFields when logged. */
type RequiredLogFields struct {
	/* The environment the application is running on ("dev", "prod" etc). */
	Env string

	/* The correlation id for the request. */
	CorrelationId uuid.UUID
}

/* Converts the struct to a map for logging; the key and value pairs in the struct tag become the KvPs in the map. */
func (lf *RequiredLogFields) toMap() map[string]interface{} {
	return map[string]interface{}{
		"env":            lf.Env,
		"correlation_id": lf.CorrelationId.String(),
	}
}
