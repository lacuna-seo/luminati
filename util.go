package luminati

import (
	"github.com/pkg/errors"
	"net/url"
)

// cleanURL removes and shebangs and query string
// parameters from a given link.
//
// Returns an error if there was an issue parsing
// the URL.
func cleanURL(link string) (string, error) {
	uri, err := url.Parse(link)
	if err != nil {
		return "", errors.Wrap(err, "error parsing luminati response url")
	}
	uri.RawQuery = ""
	uri.Fragment = ""
	return uri.String(), nil
}

// stringInSlice checks if a string exists in a slice,
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
