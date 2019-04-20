package Server

import (
	"app/api/proto/v1"
	"app/config"
	"app/src/StreetViewImage/Application/Error"
	"app/src/StreetViewImage/Infrastructure/Logger"
	"app/src/StreetViewImage/Presentation/Controller"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/j7mbo/MethodCallRetrier/v2"
	"github.com/j7mbo/go-multierror"
	"github.com/j7mbo/goij"
	"google.golang.org/grpc"
	"net"
)

/* addressRegex example: localhost:8080. Used for listening on the Grpc port. */
const addressRegex = "%s:%d"

/* GrpcServer represents an application-specific abstraction around google's Grpc server. */
type GrpcServer interface {
	/* Run runs the GrpcServer. Great that we have to start docblocks with the method name isn't it!? */
	Run()
}

/*
grpcServer encapsulates the initialisation, configuration and execution of google's google.golang.org/grpc.Server.
*/
type grpcServer struct {
	config       *config.GrpcServerConfiguration
	logger       Logger.LoggingStrategy
	retrier      MethodCallRetrier.Retrier
	interceptors *RequestInterceptorGroup

	/* injector allows the grpcServer to route to on-the-fly-injector-initialised controllers. */
	injector Goij.Injector
}

/* Create creates a new GrpcServer. */
func New(
	config *config.GrpcServerConfiguration,
	retrierFactory RetrierFactory,
	logger Logger.LoggingStrategy,
	interceptors *RequestInterceptorGroup,
	injector Goij.Injector,
) GrpcServer {
	retrier := retrierFactory.Create(config)

	return &grpcServer{config: config, logger: logger, retrier: retrier, interceptors: interceptors, injector: injector}
}

/* Run runs the GrpcServer. Great that we have to start docblocks with the method name isn't it? */
func (s *grpcServer) Run() {
	listener, err := s.createListener()

	if err != nil {
		s.logger.Error(
			fmt.Sprintf("Unable to listen on port %d, what are you doing? Error: %s", s.config.GetPort(), err.Error()),
		)

		return
	}

	middlewareInterceptors := s.interceptors.GetInterceptors()

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(middlewareInterceptors...)))

	s.registerControllers(grpcServer)

	_, rpcFatalErrors, wasSuccessful := s.retrier.ExecuteWithRetry(grpcServer, "Serve", listener)

	if !wasSuccessful {
		s.logger.Error(fmt.Sprintf("Fatal errors for rpc server: %s", multierror.AppendList(rpcFatalErrors...).Error()))
	}
}

/* registerControllers registers the relevant controller endpoint with the server. */
func (s *grpcServer) registerControllers(server *grpc.Server) {
	v1.RegisterStreetviewServiceServer(
		server,
		v1.StreetviewServiceServer(
			s.injector.Make("GetStreetViewImageController").(*Controller.GetStreetViewImageController),
		),
	)
}

/* createListener creates a listener or returns an error if, for example, the port is taken. */
func (s *grpcServer) createListener() (net.Listener, error) {
	var listener net.Listener

	errs, wasSuccessful := s.retrier.ExecuteFuncWithRetry(func() error {
		s.logger.Info(
			fmt.Sprintf(
				"Will start listening on host: %s, port: %d for GRPC calls", s.config.GetHost(), s.config.GetPort(),
			),
		)

		lis, err := net.Listen(s.config.GetProtocol(), s.createListenerAddress())

		if err != nil {
			return err
		}

		/* Makes goland happy. Jetbrains <3s PHPStorm, but they're like "screw Go developers". */
		listener = lis

		return nil
	})

	if !wasSuccessful {
		return nil, Error.NewApplicationError(multierror.AppendList(errs...).Error())
	}

	return listener, nil
}

/* createListenerAddress uses addressRegex and the host and port to create the address to listen on for net.Listener. */
func (s *grpcServer) createListenerAddress() string {
	return fmt.Sprintf(addressRegex, s.config.GetHost(), s.config.GetPort())
}
