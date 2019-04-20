package Controller

/*
GrpcErrorMapper is a controller-specific interface that Server.GrpcErrorMapper implements. The reason it exists here is
to avoid cyclic dependencies mainly, otherwise the concrete Server.GrpcErrorMapper would be injected directly.

Note that app.go now delegates the initialisation of this interface the implementation from Server.GrpcErrorMapper.
*/
type GrpcErrorMapper interface {
	/* MapToGrpcError maps a given UserError to an error returnable with the Grpc service (controller in our case). */
	MapToGrpcError(err error) error
}
