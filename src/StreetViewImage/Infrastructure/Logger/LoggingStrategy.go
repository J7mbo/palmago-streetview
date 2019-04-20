package Logger

import (
	"fmt"
	"github.com/google/uuid"
)

/*
LoggingStrategy wraps all loggers in the application and falls back to the next logger when one returns an error.

This is especially useful when, for example, logging to elasticsearch fails - fall back to writing to a file.

LoggingStrategy is a strategy of logging that falls back to the next logger when one returns an error.

The fallback ordering is in the same order of the passed in loggers.
*/
type LoggingStrategy struct {
	loggers []Logger
}

/* NewLoggingStrategy returns a new LoggingStrategy. Goij does not yet support variadics. :-( */
func NewLoggingStrategy(esLogger ElasticSearchLogger, fLogger FileLogger) *LoggingStrategy {
	return &LoggingStrategy{loggers: []Logger{&esLogger, &fLogger}}
}

/* Usual stuff that can be ignored. */
func (s *LoggingStrategy) Debug(message string, fields ...map[string]interface{}) {
	for _, logger := range s.loggers {
		if err := logger.Debug(message, fields...); err == nil {
			return
		}
	}

	fmt.Println(fmt.Sprintf("Unable to write Debug to any logger - error: '%s', fields: '%v'", message, fields))
}

/* Stuff that should be seen as we go about. */
func (s *LoggingStrategy) Info(message string, fields ...map[string]interface{}) {
	for _, logger := range s.loggers {
		if err := logger.Info(message, fields...); err == nil {
			return
		}
	}

	fmt.Println(fmt.Sprintf("Unable to write Info to any logger - error: '%s', fields: '%v'", message, fields))
}

/* Something bad happened but the application can continue. */
func (s *LoggingStrategy) Warning(message string, fields ...map[string]interface{}) {
	for _, logger := range s.loggers {
		if err := logger.Warning(message, fields...); err == nil {
			return
		}
	}

	fmt.Println(fmt.Sprintf("Unable to write Warning to any logger - error: '%s', fields: '%v'", message, fields))
}

/* Everything is majorly fucked. */
func (s *LoggingStrategy) Error(message string, fields ...map[string]interface{}) {
	for _, logger := range s.loggers {
		if err := logger.Error(message, fields...); err == nil {
			return
		}
	}

	fmt.Println(fmt.Sprintf("Unable to write Error to any logger - error: '%s', fields: '%v'", message, fields))
}

/* Every logger uses a UUID - and it can be updated depending on the request coming in, for example. */
func (s *LoggingStrategy) UpdateUuid(uuid uuid.UUID) {
	for _, logger := range s.loggers {
		logger.UpdateUuid(uuid)
	}
}
