package ApiClient

import (
	"app/config"
	"app/src/StreetViewImage/Application/Error"
	"app/src/StreetViewImage/Infrastructure/Logger"
	"errors"
	"fmt"
	"github.com/j7mbo/MethodCallRetrier/v2"
	"github.com/j7mbo/go-multierror"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	/* Request timeout - standard resiliency problem. */
	requestTimeout = time.Duration(5 * time.Second)

	/* Max sizes from: https://developers.google.com/maps/documentation/streetview/usage-and-billing. */
	maxWidth  = 640
	maxHeight = 640

	/* Error constants. */
	InvalidLocationCode    = "InvalidLocationCode"
	InvalidLocationCodeErr = "invalid location provided: the coordinates do not correspond to a valid street view image"
)

/* StreetViewApiClient handles requests to Google's Street View API. */
type StreetViewApiClient interface {
	/* Request performs a request to the street view api with the runtime provided latitude and longitude. */
	Request(latitude float64, longitude float64) ([]byte, error)
}

/* streetViewApiClient handles requests to Google's Street View API. */
type streetViewApiClient struct {
	config  *config.StreetViewApiConfiguration
	retrier MethodCallRetrier.Retrier
	logger  Logger.LoggingStrategy
}

/* NewStreetViewApiClient returns a new StreetViewApiClient. */
func NewStreetViewApiClient(
	config config.StreetViewApiConfiguration, retrierFactory RetrierFactory, logger Logger.LoggingStrategy,
) StreetViewApiClient {
	retrier := retrierFactory.Create(&config)

	return &streetViewApiClient{config: &config, retrier: retrier, logger: logger}
}

/* Request performs a request to the street view api with the runtime provided latitude and longitude. */
func (c *streetViewApiClient) Request(latitude float64, longitude float64) ([]byte, error) {
	uri, err := url.Parse(c.config.GetEndpoint())

	/* throw new DevRetardationException. */
	if err != nil {
		return nil, Error.NewApplicationError(
			fmt.Sprintf(
				"Unable to build url for request. Endpoint: '%s', error: '%s'", c.config.GetEndpoint(), err.Error(),
			),
		)
	}

	uri = c.addQueryToUrl(*uri, latitude, longitude)

	if !c.streetviewImageExistsInGoogle(uri) {
		return nil, Error.UserError{Code: InvalidLocationCode, Err: InvalidLocationCodeErr}
	}

	/* Remove the key from the url so it can be logged without sharing anything dangerous. */
	uriStringForLogging := regexp.MustCompile(`key=[^&]*`).ReplaceAllString(uri.String(), "${1}")

	c.logger.Debug(fmt.Sprintf("Making request to: %s", uriStringForLogging))

	var res *http.Response

	defer func() {
		if res != nil {
			_ = res.Body.Close()
		}
	}()

	errs, wasSuccessful := c.retrier.ExecuteFuncWithRetry(func() error {
		response, err := (&http.Client{Timeout: requestTimeout}).Get(uri.String())

		if err != nil {
			return err
		}

		if response == nil {
			return Error.NewApplicationError("No response from StreetView api...")
		}

		if response.StatusCode != 200 {
			return Error.NewApplicationError(
				fmt.Sprintf("response status code: '%d', full response: '%v'", response.StatusCode, res),
			)
		}

		res = response

		return nil
	})

	if !wasSuccessful {
		return nil, Error.NewApplicationError(
			fmt.Sprintf(
				"Error making request to: '%s', errors: '%s'", uri.String(), multierror.AppendList(errs...).Error(),
			),
		)
	}

	resBytes, _ := ioutil.ReadAll(res.Body)

	c.logger.Debug(
		fmt.Sprintf(
			"Received Streetview api response, status: '%d', bytes: '%d'", res.StatusCode, len(resBytes),
		),
	)

	return resBytes, nil
}

/*
streetviewImageExistsInGoogle performs a metadata endpoint call to check that google has this image ($$$ free).

See: https://developers.google.com/maps/documentation/streetview/metadata#response-format
*/
func (c *streetViewApiClient) streetviewImageExistsInGoogle(uri *url.URL) bool {
	metadataUri := strings.Replace(uri.String(), "/streetview?", "/streetview/metadata?", 1)

	uriStringForLogging := regexp.MustCompile(`key=[^&]*`).ReplaceAllString(uri.String(), "${1}")

	c.logger.Debug(fmt.Sprintf("Making request for metadata to: %s", uriStringForLogging))

	var res *http.Response

	defer func() {
		if res != nil {
			_ = res.Body.Close()
		}
	}()

	errs, wasSuccessful := c.retrier.ExecuteFuncWithRetry(func() error {
		response, err := (&http.Client{Timeout: requestTimeout}).Get(uri.String())

		if err != nil {
			return err
		}

		if response == nil {
			return errors.New("no response from streetView api")
		}

		if response.StatusCode != 200 {
			return errors.New(
				fmt.Sprintf("response status code: '%d', full response: '%v'", response.StatusCode, res),
			)
		}

		res = response

		return nil
	})

	if !wasSuccessful {
		c.logger.Error(
			fmt.Sprintf(
				"Error making request to: '%s', errors: '%s'", metadataUri, multierror.AppendList(errs...).Error(),
			),
		)

		return false
	}

	resBytes, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(resBytes))

	if strings.Contains(string(resBytes), "ZERO_RESULTS") {
		return false
	}

	return true
}

/* addQueryToUrl builds the GET query string from config vars and returns the newly appended url. */
func (c *streetViewApiClient) addQueryToUrl(url url.URL, latitude float64, longitude float64) *url.URL {
	q := url.Query()

	queryMap := map[string]string{
		"size":     c.buildSizeString(),
		"location": c.buildLocationString(latitude, longitude),
		"key":      c.config.GetApiKey(),
	}

	for key, value := range queryMap {
		q.Add(key, value)
	}

	url.RawQuery = q.Encode()

	return &url
}

/* buildSizeString builds the size GET parameter. */
func (c *streetViewApiClient) buildSizeString() string {
	width, height := c.config.GetWidth(), c.config.GetHeight()

	if width > maxWidth {
		width = maxWidth
	}

	if height > maxHeight {
		height = maxHeight
	}

	return fmt.Sprintf("%dx%d", width, height)
}

/* buildLocationString builds the location GET parameter. */
func (c *streetViewApiClient) buildLocationString(latitude float64, longitude float64) string {
	return fmt.Sprintf("%f,%f", latitude, longitude)
}
