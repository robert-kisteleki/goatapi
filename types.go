/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"fmt"
	"strings"
	"time"
)

// Geolocation type
type Geolocation struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

// Tag type
type Tag struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ErrorResponse type
type ErrorResponse struct {
	Detail ErrorDetail `json:"error"`
}

// ErrorDetails type
type ErrorDetail struct {
	Detail string `json:"details"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Code   int    `json:"code"`
}

// appendValueOrNA turns various types into a string if they have values
// (i.e. pinter is not nil) or "N/A" otherwise
func appendValueOrNA[T any](prefix string, quote bool, val *T) string {
	if val != nil {
		if quote {
			return fmt.Sprintf("\t\"%s%v\"", prefix, *val)
		} else {
			return fmt.Sprintf("\t%s%v", prefix, *val)
		}
	} else {
		return "\tN/A"
	}
}

// a datetime type that can be unmarshaled from ISO-8601 (without Z) datetime format
type isoTime time.Time

func (ut *isoTime) UnmarshalJSON(data []byte) error {
	layout := "2006-01-02T15:04:05"
	noquote := strings.ReplaceAll(string(data), "\"", "")
	noz := strings.ReplaceAll(noquote, "Z", "")
	unix, err := time.Parse(layout, noz)
	if err != nil {
		return err
	}
	*ut = isoTime(unix)
	return nil
}

// default output format for isoTime type is ISO8601Z
func (ut isoTime) String() string {
	return time.Time(ut).Format("2006-01-02T15:04:05Z")
}
