package Domain

/* StreetViewImages represents a repository capable of retrieving images from persistence (cache in our case). */
type StreetViewImages interface {
	/*
	   Save stores the image in persistence and returns whether or not this storing was successful.

	   Why it might be unsuccessful is not the client's concern.
	*/
	Save(image StreetViewImage) bool

	/*
	   Find retrieves an image from persistence if one exists.
	*/
	Find(latitude float64, longitude float64) StreetViewImage
}
