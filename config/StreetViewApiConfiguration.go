package config

/*
StreetViewApiConfiguration contains the configuration for use in calling the Google StreetView API.

Parameters required ripped from: https://developers.google.com/maps/documentation/streetview/intro.
*/
type StreetViewApiConfiguration struct {
	endpoint   string `env:"STREETVIEW_API_ENDPOINT" default:"https://maps.googleapis.com/maps/api/streetview"`
	height     int    `env:"STREETVIEW_API_IMAGE_HEIGHT" default:"400"`
	width      int    `env:"STREETVIEW_API_IMAGE_WIDTH" default:"400"`
	fov        int    `env:"STREETVIEW_API_IMAGE_FOV" default:"90"`
	apiKey     string `env:"STREETVIEW_API_KEY"`
	maxRetries int    `env:"STREETVIEW_API_MAX_RETRIES" default:"10"`
	retryDelay int    `env:"STREETVIEW_API_RETRY_DELAY" default:"5"`
}

func (c *StreetViewApiConfiguration) GetEndpoint() string { return c.endpoint }
func (c *StreetViewApiConfiguration) GetHeight() int      { return c.height }
func (c *StreetViewApiConfiguration) GetWidth() int       { return c.width }
func (c *StreetViewApiConfiguration) GetFov() int         { return c.fov }
func (c *StreetViewApiConfiguration) GetApiKey() string   { return c.apiKey }
func (c *StreetViewApiConfiguration) GetMaxRetries() int  { return c.maxRetries }
func (c *StreetViewApiConfiguration) GetRetryDelay() int  { return c.retryDelay }
