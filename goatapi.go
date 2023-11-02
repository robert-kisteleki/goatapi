/*
  (C) 2022, 2023 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const version = "v0.5.0"

var uaString = "goatAPI " + version
var apiBaseURL = "https://atlas.ripe.net/api/v2/"
var streamBaseURL = "wss://atlas-stream.ripe.net/stream/"

// UserAgent returns the user agent used by the package as a string
func UserAgent() string {
	return uaString
}

// ModifyUserAgent allows the caller to refine the user agent to include
// some extra piece of text
func ModifyUserAgent(addition string) {
	uaString += " (" + addition + ")"
}

// SetAPIBase allows the caller to modify the API to talk to
// This is really only useful to developers who have access to compatible APIs
func SetAPIBase(newAPIBaseURL string) {
	// TODO: check sanity of new API base URL
	apiBaseURL = newAPIBaseURL
}

// SetStreamBase allows the caller to modify the stream to talk to
// This is really only useful to developers who have access to compatible APIs
func SetStreamBase(newStreamBaseURL string) {
	// TODO: check sanity of new API base URL
	streamBaseURL = newStreamBaseURL
}

// Turn a slice of ints to a comma CSV string
func makeCsv(list []uint) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), ","), "[]")
}

func apiGetRequest(
	verbose bool,
	url string,
	key *uuid.UUID,
) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", uaString)
	if key != nil {
		req.Header.Set("Authorization", "Key "+(*key).String())
	}

	if verbose {
		msg := fmt.Sprintf("# API call: GET %s", url)
		if key != nil {
			msg += fmt.Sprintf(" (using API key %s...)", (*key).String()[:8])
		}
		fmt.Println(msg)
	}
	client := &http.Client{}

	return client.Do(req)
}
