package Query

/* GetStreetViewImage represents a query used for retrieving an image from Google StreetView. */
type GetStreetViewImage interface {
	GetLatitude() float64
	GetLongitude() float64
}

/* getStreetViewImage represents a query used for retrieving an image from Google StreetView. */
type getStreetViewImage struct {
	latitude  float64
	longitude float64
}

/* NewGetStreetViewImageQuery returns a new GetStreetViewImage. */
func NewGetStreetViewImageQuery(latitude float64, longitude float64) GetStreetViewImage {
	return &getStreetViewImage{latitude: latitude, longitude: longitude}
}

/* GetLatitude retrieves the Latitude from the GetStreetViewImage query object. */
func (q *getStreetViewImage) GetLatitude() float64 {
	return q.latitude
}

/* GetLongitude retrieves the Longitude from the GetStreetViewImage query object. */
func (q *getStreetViewImage) GetLongitude() float64 {
	return q.longitude
}
