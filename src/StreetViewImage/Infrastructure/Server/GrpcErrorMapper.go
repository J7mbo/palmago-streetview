package Server

import (
	"app/src/StreetViewImage/Application/Error"
	"app/src/StreetViewImage/Infrastructure/ApiClient"
	"app/src/StreetViewImage/Infrastructure/Logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/* grpcError is a simple grouping of the error code, grpc code to return and the error message to return. */
type grpcError struct {
	/* Code is the UserError.Code field. */
	Code string
	/* GrpcCode is the Grpc Code to be returned in the response. */
	GrpcCode codes.Code
	/* Error is the optional string to display to the user. */
	Error string
}

/* unknownError is the error returned to the user when a wild error we haven't mapped appears. */
const unknownError = "An unknown error occurred, please retry the request later."

/* okError is returned to the user when the operation completed successfully, weirdly enough. HTTP 200 equivalent. */
const okError = "The operation completed successfully."

/* Given a user error with this code, return this error string and grpc code. */
var errorMap = []grpcError{
	/* Server errors. */
	{Code: EmptyCorrelationIdCode, GrpcCode: codes.Unknown, Error: EmptyCorrelationIdErr},
	{Code: InvalidCorrelationIdCode, GrpcCode: codes.InvalidArgument, Error: EmptyCorrelationIdErr},
	/* User errors. */
	{Code: ApiClient.InvalidLocationCode, GrpcCode: codes.NotFound, Error: ApiClient.InvalidLocationCodeErr},
}

/*
grpcErrorMapper is responsible for handling the conversion of a UserError to a grpc code and message.

The application error architecture defines two error types that we care about: an ApplicationError and a UserError. An
ApplicationError is only for logging, so any error of this type received at this layer is for logging purposes only,
whereas a UserError contains an error that should be displayed to the user and so only needs to be converted to the
corresponding metadata for the correct grpc error.

ApplicationErrors are errors that ruin application flow, NOT warnings which can be logged and ignored for this request.
*/
type GrpcErrorMapper interface {
	/* MapToGrpcError maps a given UserError to an error returnable with the Grpc service (controller in our case). */
	MapToGrpcError(err error) error
}

/* grpcErrorMapper is responsible for handling the conversion of a user error to a grpc error. */
type grpcErrorMapper struct {
	logger *Logger.LoggingStrategy
}

/* NewGrpcErrorMapper returns a newly initialised GrpcErrorMapper. */
func NewGrpcErrorMapper(logger *Logger.LoggingStrategy) GrpcErrorMapper {
	return &grpcErrorMapper{logger: logger}
}

/* MapToGrpcError maps a given UserError to an error returnable with the Grpc service (controller in our case). */
func (m *grpcErrorMapper) MapToGrpcError(err error) error {
	if err == nil {
		return status.Error(codes.OK, okError)
	}

	if _, isUserErrorType := err.(Error.UserError); !isUserErrorType {
		m.logger.Error("Received non-user-error: " + err.Error())

		return status.Error(codes.Unknown, unknownError)
	}

	grpcError := m.findGrpcErrorMappingByErrorCode(err.(Error.UserError).Code)

	if grpcError == nil {
		return status.Error(codes.Unknown, unknownError)
	}

	return status.Error(grpcError.GrpcCode, grpcError.Error)
}

/* findGrpcErrorMappingByErrorCode searches the error map and returns the grpcError, if found. */
func (m *grpcErrorMapper) findGrpcErrorMappingByErrorCode(errorCode string) *grpcError {
	for _, grpcError := range errorMap {
		if grpcError.Code == errorCode {
			return &grpcError
		}
	}

	return nil
}
