package Server

import (
	"app/api/proto/v1"
	"app/src/StreetViewImage/Application/Error"
	"app/src/StreetViewImage/Infrastructure/Logger"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/j7mbo/goij"
	"github.com/kazegusuri/grpc-panic-handler"
	"google.golang.org/grpc"
	"strings"
)

/* Generic error string constants for all requests. */
const (
	EmptyCorrelationIdCode   = "EmptyCorrelationId"
	EmptyCorrelationIdErr    = "invalid correlation id provided, it must not be empty nor a blank string"
	InvalidCorrelationIdCode = "InvalidCorrelationId"
	InvalidCorrelationIdErr  = "invalid non-version-4 uuid provided, example v4 format: acca4678-fbbd-43b9-9d8a-83f8794935cb"
)

/* RequestInterceptorGroup returns user-defined middleware functions used for intercepting grpc requests. */
type RequestInterceptorGroup struct {
	Logger   Logger.LoggingStrategy
	Injector Goij.Injector
}

/* GetInterceptors retrieves all the user-defined middleware functions used for intercepting grpc requests. */
func (ri *RequestInterceptorGroup) GetInterceptors() []grpc.UnaryServerInterceptor {
	panichandler.InstallPanicHandler(func(r interface{}) {
		ri.Logger.Error(fmt.Sprintf("Panic from grpc: %v", r))
	})

	return []grpc.UnaryServerInterceptor{
		/* Add interceptors here. */
		grpc.UnaryServerInterceptor(ri.addUuidToInjector),
		grpc.UnaryServerInterceptor(panichandler.UnaryPanicHandler),
	}
}

/*
addUuidToInjector is a middleware function to retrieve the Uuid from the request and share it with the injector.

The UUID should be a version 4 UUID (easily found on google).

Note that at this point if an error is returned it is shown to the user.
*/
func (ri *RequestInterceptorGroup) addUuidToInjector(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	correlationId := req.(*v1.GetStreetViewRequest).CorrelationId

	if strings.Trim(correlationId, " ") == "" {
		return nil, Error.UserError{Code: EmptyCorrelationIdCode, Err: EmptyCorrelationIdErr}
	}

	newUuid, err := uuid.Parse(correlationId)

	if err != nil {
		return nil, Error.UserError{Code: InvalidCorrelationIdCode, Err: InvalidCorrelationIdErr}
	}

	/* Overwrite injected object's correlation id with the request value and re-share with the injector. */
	ri.Logger.UpdateUuid(newUuid)

	ri.Injector.Share(ri.Logger)

	ri.Logger.Info(fmt.Sprintf("Request received: %v", req))

	return handler(ctx, req)
}
