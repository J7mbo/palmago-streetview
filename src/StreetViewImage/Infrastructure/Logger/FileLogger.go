package Logger

import (
	"bufio"
	"github.com/google/uuid"
	logrusLogger "github.com/sirupsen/logrus"
	"os"
)

/* fileHeaderLine is the line written to the top of the backup log file (also used to see if we *can* write...). */
const fileHeaderLine = "[Backup log file for when elastic search writing fails]\n\n"

/* FileLogger is the fallback logger that logs to a file when ElasticSearchLogger fails to write to elasticsearch. */
type FileLogger struct {
	logrusLogger   *logrusLogger.Logger
	requiredFields *RequiredLogFields
}

/* NewFileLogger returns a new FileLogger, or error if the file provided could not be written to. */
func NewFileLogger(
	file os.File, logger logrusLogger.Logger, requiredFields RequiredLogFields,
) (*FileLogger, error) {
	logger.Level = logrusLogger.Level(logrusLogger.TraceLevel)

	_, err := file.Write([]byte(fileHeaderLine))

	if err != nil {
		return nil, err
	}

	logger.SetOutput(bufio.NewWriter(&file))

	return &FileLogger{logrusLogger: &logger, requiredFields: &requiredFields}, nil
}

/* Debug is for the usual stuff that can be ignored. */
func (l *FileLogger) Debug(message string, fields ...map[string]interface{}) error {
	l.logrusLogger.WithFields(l.mergeFields(fields)).Debug(message)

	_ = l.logrusLogger.Out.(*bufio.Writer).Flush()

	return nil
}

/* Info means stuff that should be seen as we go about the app lifecycle. */
func (l *FileLogger) Info(message string, fields ...map[string]interface{}) error {
	l.logrusLogger.WithFields(l.mergeFields(fields)).Info(message)

	_ = l.logrusLogger.Out.(*bufio.Writer).Flush()

	return nil
}

/* Warning means something bad happened but the application can continue. */
func (l *FileLogger) Warning(message string, fields ...map[string]interface{}) error {
	l.logrusLogger.WithFields(l.mergeFields(fields)).Warning(message)

	_ = l.logrusLogger.Out.(*bufio.Writer).Flush()

	return nil
}

/* Error means stuff is majorly fucked. */
func (l *FileLogger) Error(message string, fields ...map[string]interface{}) error {
	l.logrusLogger.WithFields(l.mergeFields(fields)).Error(message)

	_ = l.logrusLogger.Out.(*bufio.Writer).Flush()

	return nil
}

/*
UpdateUuid is used to update the required field uuid with one from the request, called in middleware.

The required log fields are then re-shared with the injector so that any future injections use this uuid.
*/
func (l *FileLogger) UpdateUuid(uuid uuid.UUID) {
	l.requiredFields.CorrelationId = uuid
}

/* mergeFields merges the required fields with the provided fields for logging to elasticsearch. */
func (l *FileLogger) mergeFields(fields []map[string]interface{}) map[string]interface{} {
	finalFields := l.requiredFields.toMap()

	if len(fields) > 0 {
		for _, field := range fields {
			finalFields = mergeMaps(finalFields, field)
		}
	}

	return finalFields
}
