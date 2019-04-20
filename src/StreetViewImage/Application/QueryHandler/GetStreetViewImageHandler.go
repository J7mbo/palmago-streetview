package QueryHandler

import (
	"app/src/StreetViewImage/Application/Error"
	"app/src/StreetViewImage/Application/Query"
	"app/src/StreetViewImage/Domain"
	"app/src/StreetViewImage/Infrastructure/ApiClient"
	"app/src/StreetViewImage/Infrastructure/Logger"
	"fmt"
)

/* GetStreetViewImageHandler handles a query to retrieve an image from Google StreetView. */
type GetStreetViewImageHandler interface {
	/* Handle takes in a Query and returns an array of bytes containing an image / an error. */
	Handle(query Query.GetStreetViewImage) ([]byte, error)
}

/* getStreetViewImage handles a query to retrieve an image from Google StreetView. */
type getStreetViewImageHandler struct {
	repository Domain.StreetViewImages
	apiClient  ApiClient.StreetViewApiClient
	logger     Logger.LoggingStrategy
}

/* NewGetStreetViewImageHandler returns a new GetStreetViewImageHandler. */
func NewGetStreetViewImageHandler(
	repository Domain.StreetViewImages, apiClient ApiClient.StreetViewApiClient, logger Logger.LoggingStrategy,
) GetStreetViewImageHandler {
	return &getStreetViewImageHandler{repository: repository, apiClient: apiClient, logger: logger}
}

/* Handle takes in a Query and returns an array of bytes containing an image / an error. */
func (h *getStreetViewImageHandler) Handle(query Query.GetStreetViewImage) ([]byte, error) {
	lat, lon := query.GetLatitude(), query.GetLongitude()

	img := h.repository.Find(lat, lon)

	if img != nil {
		h.logger.Debug(fmt.Sprintf("Cache already contains image for lat: '%f', lon: '%f', returning...", lat, lon))

		return img.GetBytes(), nil
	}

	responseBytes, err := h.apiClient.Request(lat, lon)

	if err != nil {
		/* Assertion here as we only know an image doesn't exist in streetview (user's fault) when in the domain. */
		if _, isUserError := err.(Error.UserError); isUserError {
			return nil, err
		}

		return nil, Error.NewApplicationError(
			fmt.Sprintf("Unable to perform request to streetview API, error: %s", err.Error()),
		)
	}

	image, err := Domain.NewStreetViewImage(lat, lon, responseBytes)

	if err != nil {
		return nil, Error.NewApplicationError(
			fmt.Sprintf("Looks like the StreetView api image was not considered a valid image. Error: %s", err.Error()),
		)
	}

	fmt.Println(h.repository)

	image.Save(h.repository)

	return image.GetBytes(), nil
}
