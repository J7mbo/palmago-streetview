package Domain

import (
	"errors"
	"fmt"
	"net/http"
)

/* expectedImageType is the content type the StreetView API specifies it will return; used for validation purposes. */
const expectedImageType = "image/jpeg"

/* StreetViewImage contains the raw data of an image from Google StreetView. */
type StreetViewImage interface {
	GetUuid() string
	GetLatitude() float64
	GetLongitude() float64
	GetBytes() []byte

	/* Save saves an image for future use. Technically this is caching it. Here's your DDD-style stuff -.-. */
	Save(images StreetViewImages)
}

/* streetViewImage contains the raw data of an image from Google StreetView. */
type streetViewImage struct {
	uuid       *ImageUuid
	latitude   float64
	longitude  float64
	imageBytes []byte
}

/* NewStreetViewImage returns an initialised StreetViewImage or an error if the image was considered invalid. */
func NewStreetViewImage(latitude float64, longitude float64, byteArray []byte) (StreetViewImage, error) {
	if err := validateImage(byteArray); err != nil {
		return nil, err
	}

	uuid := NewImageUuid(latitude, longitude)

	return &streetViewImage{uuid: uuid, latitude: latitude, longitude: longitude, imageBytes: byteArray}, nil
}

/* GetUuid returns the uuid for this image as a string. */
func (i *streetViewImage) GetUuid() string {
	return i.uuid.String()
}

/* GetBytes retrieves the bytes of the image. */
func (i *streetViewImage) GetBytes() []byte {
	return i.imageBytes
}

/* GetLatitude retrieves the latitude. */
func (i *streetViewImage) GetLatitude() float64 {
	return i.latitude
}

/* GetLongitude retrieves the longitude. */
func (i *streetViewImage) GetLongitude() float64 {
	return i.longitude
}

/* Save saves an image for future use. */
func (i *streetViewImage) Save(images StreetViewImages) {
	/* It doesn't really matter if this fails, this is optional and is already logged. */
	_ = images.Save(i)
}

/* validateImage uses http.DetectContentType to ensure that the image is of type image/jpeg as the API specifies. */
func validateImage(imageBytes []byte) error {
	contentType := http.DetectContentType(imageBytes)

	if contentType != expectedImageType {
		return errors.New(
			fmt.Sprintf("image data used to create a StreetViewImage should be an image/jpeg, got: '%s'", contentType),
		)
	}

	return nil
}
