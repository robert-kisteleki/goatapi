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
	"os"
	"time"

	"github.com/robert-kisteleki/goatapi/result"
	"golang.org/x/exp/slices"
)

// ResultsFilter struct holds specified filters and other options
type ResultsFilter struct {
	params url.Values
	id     uint
	file   string
	limit  uint
	start  *time.Time
	stop   *time.Time
	probes []uint
}

// NewResultsFilter prepares a new result filter object
func NewResultsFilter() ResultsFilter {
	filter := ResultsFilter{}
	filter.params = url.Values{}
	filter.params.Add("format", "txt")
	filter.probes = make([]uint, 0)
	return filter
}

// FilterID filters by a particular measurement ID
func (filter *ResultsFilter) FilterID(id uint) {
	filter.id = id
}

// FilterFile "filters" results from a particular file
func (filter *ResultsFilter) FilterFile(filename string) {
	filter.file = filename
}

// FilterStart filters for results after this timestamp
func (filter *ResultsFilter) FilterStart(t time.Time) {
	filter.start = &t
	filter.params.Add("start", fmt.Sprintf("%d", t.Unix()))
}

// FilterStop filters for results before this timestamp
func (filter *ResultsFilter) FilterStop(t time.Time) {
	filter.stop = &t
	filter.params.Add("stop", fmt.Sprintf("%d", t.Unix()))
}

// FilterProbeIDs filters for results where the probe ID is one of several in the list specified
func (filter *ResultsFilter) FilterProbeIDs(list []uint) {
	filter.probes = list
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
	if filter.id == 0 && filter.file == "" {
		return fmt.Errorf("ID or filename must be specified")
	}

	if filter.limit == 0 {
		return fmt.Errorf("limit must be positive")
	}

	return nil
}

// GetResults returns results in a slice by applying the specified filters
func (filter *ResultsFilter) GetResults(
	verbose bool,
) (
	results []result.Result,
	err error,
) {
	// prepare to read results
	read, err := filter.openNetworkResults(verbose)
	if err != nil {
		return nil, err
	}

	var fetched uint = 0
	typehint := ""
	for read.Scan() && fetched < filter.limit {
		line := read.Text()

		res, err := result.ParseWithTypeHint(line, typehint)
		if err != nil {
			return results, err
		}

		// add the result to the result set
		results = append(results, res)

		fetched++

		// a type hint makes parsing much faster
		if typehint == "" {
			typehint = res.TypeName()
		}
	}

	return results, nil
}

// GetResultsAsync returns results from the API or a file via a channel
func (filter *ResultsFilter) GetResultsAsync(
	verbose bool,
	results chan result.AsyncResult,
) {
	if filter.id != 0 {
		filter.GetNetworkResultsAsync(verbose, results)
	} else {
		if filter.file != "" {
			filter.GetFileResultsAsync(verbose, results)
		} else {
			results <- result.AsyncResult{Result: nil, Error: fmt.Errorf("neither ID nor input file were specified")}
			close(results)
		}
	}
}

// GetNetworkResultsAsync returns results from the API
// via a channel by applying the specified filters
func (filter *ResultsFilter) GetNetworkResultsAsync(
	verbose bool,
	results chan result.AsyncResult,
) {
	defer close(results)

	// prepare to read results
	read, err := filter.openNetworkResults(verbose)
	if err != nil {
		results <- result.AsyncResult{Result: nil, Error: err}
		return
	}

	filter.getResultsAsync(verbose, read, results)
}

// GetFileResultsAsync returns results from a file via a channel
// If the file is "-" then it reads from stdin
func (filter *ResultsFilter) GetFileResultsAsync(
	verbose bool,
	results chan result.AsyncResult,
) {
	defer close(results)

	var file *os.File
	if filter.file == "-" {
		file = os.Stdin
		if verbose {
			fmt.Printf("# Reading results from stdin\n")
		}
	} else {
		var err error
		file, err = os.Open(filter.file)
		if err != nil {
			results <- result.AsyncResult{Result: nil, Error: err}
			return
		}
		defer file.Close()

		if verbose {
			fmt.Printf("# Reading results from file: %s\n", filter.file)
		}
	}

	read := bufio.NewScanner(bufio.NewReader(file))

	filter.getResultsAsync(verbose, read, results)
}

func (filter *ResultsFilter) getResultsAsync(
	verbose bool,
	read *bufio.Scanner,
	results chan result.AsyncResult,
) {
	var fetched uint = 0
	typehint := ""
	for read.Scan() && fetched < filter.limit {
		line := read.Text()

		res, err := result.ParseWithTypeHint(line, typehint)
		if err != nil {
			results <- result.AsyncResult{Result: nil, Error: err}
			continue
		}

		// check if time interval and probe constraints match (applicable if we're
		// reading from a file), and if so, put the result on the channel
		ts := time.Time(res.GetTimeStamp())
		if (filter.start == nil || filter.start.Before(ts.Add(time.Duration(1)))) &&
			(filter.stop == nil || filter.stop.After(ts.Add(time.Duration(-1)))) &&
			(len(filter.probes) == 0 || slices.Contains(filter.probes, res.GetProbeID())) {
			results <- result.AsyncResult{Result: res, Error: nil}
			fetched++
		}

		// a type hint makes parsing much faster
		if typehint == "" {
			typehint = res.TypeName()
		}
	}
}

// prepare fecthing results, i.e. verify parameters, connect to the API, etc.
func (filter *ResultsFilter) openNetworkResults(
	verbose bool,
) (
	read *bufio.Scanner,
	err error,
) {
	// sanity checks - late in the process, but not too late
	err = filter.verifyFilters()
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf("%smeasurements/%d/results/?%s", apiBaseURL, filter.id, filter.params.Encode())

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", uaString)

	if verbose {
		fmt.Printf("# API call: GET %s\n", req.URL)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// we're reading one result per line, a scanner is simple enough
	return bufio.NewScanner(bufio.NewReader(resp.Body)), nil
}
