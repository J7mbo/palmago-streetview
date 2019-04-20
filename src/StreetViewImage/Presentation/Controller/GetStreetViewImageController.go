package Controller

import (
	"app/api/proto/v1"
	"app/src/StreetViewImage/Application/Query"
	"app/src/StreetViewImage/Application/QueryHandler"
	"context"
)

/* GetStreetViewImageController handles the request / response of a v1.GetStreetViewRequest. */
type GetStreetViewImageController struct {
	Handler    QueryHandler.GetStreetViewImageHandler
	GrpcMapper GrpcErrorMapper
}

/* GetStreetViewImage handles the request / response of a v1.GetStreetViewRequest. */
func (c *GetStreetViewImageController) GetStreetViewImage(
	context context.Context, request *v1.GetStreetViewRequest,
) (*v1.GetStreetViewResponse, error) {
	query := Query.NewGetStreetViewImageQuery(float64(request.Latitude), float64(request.Longitude))

	imageBytes, err := c.Handler.Handle(query)

	if err != nil {
		return nil, c.GrpcMapper.MapToGrpcError(err)
	}

	response := &v1.GetStreetViewResponse{Image: imageBytes}

	return response, nil
}
