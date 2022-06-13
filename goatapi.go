/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"fmt"
	"strings"
)

const version = "v0.1.0"

var uaString = "goatAPI " + version
var apiBaseURL = "https://atlas.ripe.net/api/v2/"

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

// Turn a slice of ints to a comma CSV string
func makeCsv(list []int) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(list)), ","), "[]")
}
