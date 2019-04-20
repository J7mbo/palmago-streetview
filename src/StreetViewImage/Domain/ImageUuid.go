package Domain

import (
	"fmt"
	"strings"
)

/* uuidString is the format for the unique identifier of a StreetViewImage. */
const uuidString = "street_view_image:{image.latitude}:{image.longitude}"

/* ImageUuid is a unique identifier for a StreetViewImage. It may be used for persistence and is re-constructable. */
type ImageUuid struct {
	uuidString string
}

/*
NewImageUuid creates a new uuid given a latitude and longitude.

We deliberately don't pass a StreetViewImage here for easy reconstruction purposes (no domain object required).
*/
func NewImageUuid(latitude float64, longitude float64) *ImageUuid {
	latitudeString := fmt.Sprintf("%f", latitude)
	longitudeString := fmt.Sprintf("%f", longitude)

	replacer := strings.NewReplacer("{image.latitude}", latitudeString, "{image.longitude}", longitudeString)

	return &ImageUuid{uuidString: replacer.Replace(uuidString)}
}

/* String returns the uuid as a string, useful for persistence. */
func (i *ImageUuid) String() string {
	return i.uuidString
}
