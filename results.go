/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/robert-kisteleki/goatapi/result"
)

// ResultsFilter struct holds specified filters and other options
type ResultsFilter struct {
	params url.Values
	id     uint
	limit  uint
}

// NewResultsFilter prepares a new result filter object
func NewResultsFilter() ResultsFilter {
	filter := ResultsFilter{}
	filter.params = url.Values{}
	filter.params.Add("format", "txt")
	return filter
}

// FilterID filters by a particular measurement ID (mandatory)
func (filter *ResultsFilter) FilterID(id uint) {
	filter.id = id
}

// FilterStart filters for results after this timestamp
func (filter *ResultsFilter) FilterStart(t time.Time) {
	filter.params.Add("start", fmt.Sprintf("%d", t.Unix()))
}

// FilterStop filters for results before this timestamp
func (filter *ResultsFilter) FilterStop(t time.Time) {
	filter.params.Add("stop", fmt.Sprintf("%d", t.Unix()))
}

// FilterProbeIDs filters for results where the probe ID is one of several in the list specified
func (filter *ResultsFilter) FilterProbeIDs(list []uint) {
	filter.params.Add("probe_ids", makeCsv(list))
}

// FilterAnchors filters for results reported by anchors
func (filter *ResultsFilter) FilterAnchors() {
	filter.params.Add("anchors-only", "true")
}

// FilterAnchors filters for results reported by public probes
func (filter *ResultsFilter) FilterPublicProbes() {
	filter.params.Add("public-only", "true")
}

// Limit limits the number of result retrieved
func (filter *ResultsFilter) Limit(max uint) {
	filter.limit = max
}

// Verify sanity of applied filters
func (filter *ResultsFilter) verifyFilters() error {
	if filter.id == 0 {
		return fmt.Errorf("ID must be specified")
	}

	if filter.limit == 0 {
		return fmt.Errorf("limit must be positive")
	}

	return nil
}

// GetResults returns a set of results by applying all the specified filters
func (filter *ResultsFilter) GetResults(
	verbose bool,
) (
	results []result.Result,
	err error,
) {
	// sanity checks - late in the process, but not too late
	err = filter.verifyFilters()
	if err != nil {
		return
	}

	query := fmt.Sprintf("%smeasurements/%d/results/?%s", apiBaseURL, filter.id, filter.params.Encode())

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", uaString)

	if verbose {
		fmt.Printf("API call: GET %s\n", req.URL)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// we're reading one result per line, a scanner is simple enough
	read := bufio.NewScanner(bufio.NewReader(resp.Body))

	var fetched uint = 0
	typehint := ""
	for read.Scan() && fetched < filter.limit {
		line := read.Text()

		res, err := result.ParseWithTypeHint(line, typehint)
		if err != nil {
			return results, err
		}
		results = append(results, res)

		fetched++

		if typehint == "" {
			typehint = res.TypeName()
		}
	}

	return results, nil
}
