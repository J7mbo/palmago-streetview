package Logger

import "github.com/google/uuid"

/* Loggers are for logging! Errors can be thrown if, say, the connection to elasticsearch failed and we should retry. */
type Logger interface {
	/* Usual stuff that can be ignored. */
	Debug(message string, fields ...map[string]interface{}) error
	/* Stuff that should be seen as we go about. */
	Info(message string, fields ...map[string]interface{}) error
	/* Something bad happened but the application can continue. */
	Warning(message string, fields ...map[string]interface{}) error
	/* Everything is majorly fucked. */
	Error(message string, fields ...map[string]interface{}) error
	/* Every logger uses a UUID - and it can be updated depending on the request coming in, for example. */
	UpdateUuid(uuid uuid.UUID)
}
