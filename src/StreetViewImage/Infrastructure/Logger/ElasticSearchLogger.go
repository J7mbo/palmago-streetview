package Logger

import (
	"app/config"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/j7mbo/MethodCallRetrier/v2"
	"github.com/j7mbo/go-multierror"
	"github.com/olivere/elastic"
	logrusLogger "github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v3"
	"time"
)

/* Currently http-only. */
const elasticUrlFormat = "http://%s:%d"

/* This log level and above will be logged. */
const minLogLevel = logrusLogger.DebugLevel

/* ElasticSearchLogger is the default application logger which logs to elastic search. */
type ElasticSearchLogger struct {
	/* The logrus library we are decorating. */
	logrusLogger logrusLogger.Logger

	/* Configuration object containing host, port etc. */
	config *config.ElasticSearchConfiguration

	/* Every log must have a correlation id for the request associated with it. */
	requiredFields *RequiredLogFields

	/* To retry calls to elasticsearch. */
	retrier MethodCallRetrier.Retrier
}

/* Create initialises a new instance of an ElasticSearchLogger. */
func NewElasticSearchLogger(
	config config.ElasticSearchConfiguration,
	logger logrusLogger.Logger,
	requiredFields RequiredLogFields,
	retrierFactory RetrierFactory,
) *ElasticSearchLogger {
	retrier := retrierFactory.Create(&config)

	logger.SetLevel(minLogLevel)

	return &ElasticSearchLogger{
		config: &config, logrusLogger: logger, requiredFields: &requiredFields, retrier: retrier,
	}
}

/*
UpdateUuid is used to update the required field uuid with one from the request, called in middleware.

The required log fields are then re-shared with the injector so that any future injections use this uuid.
*/
func (l *ElasticSearchLogger) UpdateUuid(uuid uuid.UUID) {
	l.requiredFields.CorrelationId = uuid
}

/* Debug is for the usual stuff that can be ignored. */
func (l *ElasticSearchLogger) Debug(message string, fields ...map[string]interface{}) error {
	if err := l.addHooksWithRetry(); err != nil {
		return err
	}

	l.logrusLogger.WithFields(l.mergeFields(fields)).Debug(message)

	return nil
}

/* Info means stuff that should be seen as we go about the app lifecycle. */
func (l *ElasticSearchLogger) Info(message string, fields ...map[string]interface{}) error {
	if err := l.addHooksWithRetry(); err != nil {
		fmt.Println(err)
		return err
	}

	l.logrusLogger.WithFields(l.mergeFields(fields)).Info(message)

	return nil
}

/* Warning means something bad happened but the application can continue. */
func (l *ElasticSearchLogger) Warning(message string, fields ...map[string]interface{}) error {
	if err := l.addHooksWithRetry(); err != nil {
		return err
	}

	l.logrusLogger.WithFields(l.mergeFields(fields)).Warning(message)

	return nil
}

/* Error means stuff is majorly fucked. */
func (l *ElasticSearchLogger) Error(message string, fields ...map[string]interface{}) error {
	if err := l.addHooksWithRetry(); err != nil {
		return err
	}

	l.logrusLogger.WithFields(l.mergeFields(fields)).Error(message)

	return nil
}

/* Repeatedly retries the addHooks() call and returns concatenated errors if one is provided. */
func (l *ElasticSearchLogger) addHooksWithRetry() error {
	errs, wasSuccessful := l.retrier.ExecuteFuncWithRetry(func() error {
		err := l.addHooks()

		return err
	})

	if wasSuccessful {
		return nil
	}

	return errors.New(multierror.AppendList(errs...).Error())
}

/* addHooks adds hooks to Logrus, which calls a ping method internally, so we don't do this on Logger initialisation. */
func (l *ElasticSearchLogger) addHooks() error {
	logger := l.logrusLogger

	/* We've already added the hooks, no need to do it again */
	if len(l.logrusLogger.Hooks) > 0 {
		return nil
	}

	client, err := elastic.NewClient(
		elastic.SetURL(fmt.Sprintf(elasticUrlFormat, l.config.GetHost(), l.config.GetPort())),
		/* First time spinning up docker, it can take some time for elasticsearch to start running. */
		elastic.SetHealthcheckTimeoutStartup(time.Duration(5)),
		elastic.SetSnifferTimeout(time.Duration(5)),
		elastic.SetSnifferTimeoutStartup(time.Duration(5)),
		elastic.SetHealthcheckTimeout(time.Duration(5)),
		/* Alternative: https://github.com/olivere/elastic-with-docker. */
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)

	if err != nil {
		return err
	}

	/* This call also adds all the other Hooks as well (Fatal, Error, Warn etc). */
	debugHook, err := elogrus.NewAsyncElasticHook(
		client, l.config.GetHost(), logrusLogger.DebugLevel, l.config.GetIndex(),
	)

	if err != nil {
		return err
	}

	logger.Hooks.Add(debugHook)

	return nil
}

/* mergeFields merges the required fields with the provided fields for logging to elasticsearch. */
func (l *ElasticSearchLogger) mergeFields(fields []map[string]interface{}) map[string]interface{} {
	finalFields := l.requiredFields.toMap()

	if len(fields) > 0 {
		for _, field := range fields {
			finalFields = mergeMaps(finalFields, field)
		}
	}

	return finalFields
}

/* mergeMaps merges two maps, assuming that the maps are single-dimensional, and overwrites duplicate keys. */
func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}

	return result
}
