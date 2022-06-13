/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

// Probe object, as it comes from the API
type Probe struct {
	ID             int           `json:"id"`
	Address4       *netip.Addr   `json:"address_v4"`
	Address6       *netip.Addr   `json:"address_v6"`
	ASN4           *int          `json:"asn_v4"`
	ASN6           *int          `json:"asn_v6"`
	CountryCode    string        `json:"country_code"`
	Description    string        `json:"description"`
	FirstConnected *uniTime      `json:"first_connected"`
	LastConnected  *uniTime      `json:"last_connected"`
	Location       Geolocation   `json:"geometry"`
	Anchor         bool          `json:"is_anchor"`
	Prefix4        *netip.Prefix `json:"prefix_v4"`
	Prefix6        *netip.Prefix `json:"prefix_v6"`
	Public         bool          `json:"is_public"`
	Status         ProbeStatus   `json:"status"`
	StatusSince    uniTime       `json:"status_since"`
	TotalUptime    int           `json:"total_uptime"`
	Type           string        `json:"type"`
	Tags           []Tag         `json:"tags"`
}

// ProbeListSortOrders lists all the allowed sort orders
var ProbeListSortOrders = []string{
	"id", "-id",
}

// various probe statuses
const (
	ProbeStatusNeverConnected = iota // 0
	ProbeStatusConnected             // 1
	ProbeStatusDisconnected          // 2
	ProbeStatusAbandoned             // 3
)

// ProbeStatusDict maps the probe status codes to human readable descriptions
var ProbeStatusDict = map[int]string{
	ProbeStatusNeverConnected: "NeverConnected",
	ProbeStatusConnected:      "Connected",
	ProbeStatusDisconnected:   "Disconnected",
	ProbeStatusAbandoned:      "Abandoned",
}

// ProbeStatus as defined by the API
type ProbeStatus struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Since *uniTime `json:"since"`
}

// ShortString produces a short textual description of the probe
func (probe *Probe) ShortString() string {
	text := fmt.Sprintf("%d\t%s", probe.ID, ProbeStatusDict[probe.Status.ID])

	if probe.CountryCode != "" {
		text += fmt.Sprintf("\t\"%s\"", probe.CountryCode)
	} else {
		text += "\tN/A"
	}

	text += appendValueOrNA("AS", false, probe.ASN4)
	text += appendValueOrNA("AS", false, probe.ASN6)

	text += fmt.Sprintf("\t%v", probe.Location.Coordinates)

	if probe.Description != "" {
		text += fmt.Sprintf("\t\"%s\"", probe.Description)
	} else {
		text += "\tN/A"
	}

	return text
}

// LongString produces a longer textual description of the probe
func (probe *Probe) LongString() string {
	text := probe.ShortString()

	text += appendValueOrNA("", false, probe.Address4)
	text += appendValueOrNA("", false, probe.Prefix4)
	text += appendValueOrNA("", false, probe.Address6)
	text += appendValueOrNA("", false, probe.Prefix6)
	text += appendValueOrNA("", false, probe.FirstConnected)
	text += appendValueOrNA("", false, probe.LastConnected)
	text += fmt.Sprintf(" %d %v %v", probe.TotalUptime, probe.Anchor, probe.Public)

	tags := make([]string, 0)
	for _, tag := range probe.Tags {
		tags = append(tags, tag.Slug)
	}
	text += fmt.Sprintf(" %v", tags)

	return text
}

// the API paginates; this describes one such page
type probeListingPage struct {
	Count    int     `json:"count"`
	Next     string  `json:"next"`
	Previous string  `json:"previous"`
	Probes   []Probe `json:"results"`
}

// ProbeFilter struct holds specified filters and other options
type ProbeFilter struct {
	params url.Values
	id     int
	limit  int
}

// NewProbeFilter prepares a new probe filter object
func NewProbeFilter() ProbeFilter {
	filter := ProbeFilter{}
	filter.params = url.Values{}
	return filter
}

// FilterID filters by a particular probe ID
func (filter *ProbeFilter) FilterID(id int) {
	filter.id = id
}

// FilterCountry filters by a country code (ISO3166-1 alpha-2)
func (filter *ProbeFilter) FilterCountry(cc string) {
	filter.params.Add("country_code", cc)
}

// FilterIDGt filters for probe IDs > some number
func (filter *ProbeFilter) FilterIDGt(n int) {
	filter.params.Add("id__gt", fmt.Sprint(n))
}

// FilterIDGte filters for probe IDs >= some number
func (filter *ProbeFilter) FilterIDGte(n int) {
	filter.params.Add("id__gte", fmt.Sprint(n))
}

// FilterIDLt filters for probe IDs < some number
func (filter *ProbeFilter) FilterIDLt(n int) {
	filter.params.Add("id__lt", fmt.Sprint(n))
}

// FilterIDLte filters for probe IDs <= some number
func (filter *ProbeFilter) FilterIDLte(n int) {
	filter.params.Add("id__lte", fmt.Sprint(n))
}

// FilterIDin filters for probe ID being one of several in the list specified
func (filter *ProbeFilter) FilterIDin(list []int) {
	filter.params.Add("id__in", makeCsv(list))
}

// FilterASN filters for an ASN in IPv4 or IPv6 space
func (filter *ProbeFilter) FilterASN(n int) {
	filter.params.Add("asn", fmt.Sprint(n))
}

// FilterASN4 filters for an ASN in IPv4 space
func (filter *ProbeFilter) FilterASN4(n int) {
	filter.params.Add("asn_v4", fmt.Sprint(n))
}

// FilterASN4in filters for an ASN on this list in IPv4 space
func (filter *ProbeFilter) FilterASN4in(list []int) {
	filter.params.Add("asn_v4__in", makeCsv(list))
}

// FilterASN6 filters for an ASN in IPv6 space
func (filter *ProbeFilter) FilterASN6(n int) {
	filter.params.Add("asn_v6", fmt.Sprint(n))
}

// FilterASN6in filters for an ASN on this list in IPv6 space
func (filter *ProbeFilter) FilterASN6in(list []int) {
	filter.params.Add("asn_v6__in", makeCsv(list))
}

// FilterStatus filters for probes that have a specific status
// See: const MeasurementStatusSpecified*
func (filter *ProbeFilter) FilterStatus(n int) {
	filter.params.Add("status", fmt.Sprint(n))
}

// FilterLatitudeGt filters for probes with latitude being greater than ('north of') this value (in degrees)
func (filter *ProbeFilter) FilterLatitudeGt(f float64) {
	filter.params.Add("latitude__gt", fmt.Sprint(f))
}

// FilterLatitudeGte filters for probes with latitude being greater than or equal to ('north of') this value (in degrees)
func (filter *ProbeFilter) FilterLatitudeGte(f float64) {
	filter.params.Add("latitude__gte", fmt.Sprint(f))
}

// FilterLatitudeLt filters for probes with latitude being greater than ('south of') this value (in degrees)
func (filter *ProbeFilter) FilterLatitudeLt(f float64) {
	filter.params.Add("latitude__lt", fmt.Sprint(f))
}

// FilterLatitudeLte filters for probes with latitude being greater than or equal to ('south of') this value (in degrees)
func (filter *ProbeFilter) FilterLatitudeLte(f float64) {
	filter.params.Add("latitude__lte", fmt.Sprint(f))
}

// FilterLongitudeGt filters for probes with longitude being greater than ('east of') this value (in degrees)
func (filter *ProbeFilter) FilterLongitudeGt(f float64) {
	filter.params.Add("longitude__gt", fmt.Sprint(f))
}

// FilterLongitudeGte filters for probes with longitude being greater than or eaual to ('east of') this value (in degrees)
func (filter *ProbeFilter) FilterLongitudeGte(f float64) {
	filter.params.Add("longitude__gte", fmt.Sprint(f))
}

// FilterLongitudeLt filters for probes with longitude being greater than ('west of') this value (in degrees)
func (filter *ProbeFilter) FilterLongitudeLt(f float64) {
	filter.params.Add("longitude__lt", fmt.Sprint(f))
}

// FilterLongitudeLte filters for probes with longitude being greater than or equal to ('west of') this value (in degrees)
func (filter *ProbeFilter) FilterLongitudeLte(f float64) {
	filter.params.Add("longitude__lte", fmt.Sprint(f))
}

// FilterAnchor filters for probes that are anchors (true) or regular probes (false)
func (filter *ProbeFilter) FilterAnchor(yesno bool) {
	filter.params.Add("is_anchor", fmt.Sprint(yesno))
}

// FilterPublic filters for probes that are public or non-public
func (filter *ProbeFilter) FilterPublic(yesno bool) {
	filter.params.Add("is_public", fmt.Sprint(yesno))
}

// FilterRadius filters for probes that are within the radius (in km) of a coordinate
// It assumes the use of FilterLatitude and FilterLongitude as well
func (filter *ProbeFilter) FilterRadius(lat, lon, radius float64) {
	filter.params.Add("radius", fmt.Sprintf("%f,%f:%f", lat, lon, radius))
}

// FilterPrefixV4 filters for probes that are in a specific IPv4 prefix
func (filter *ProbeFilter) FilterPrefixV4(prefix netip.Prefix) {
	filter.params.Add("prefix_v4", fmt.Sprint(prefix))
}

// FilterPrefixV6 filters for probes that are in a specific IPv6 prefix
func (filter *ProbeFilter) FilterPrefixV6(prefix netip.Prefix) {
	filter.params.Add("prefix_v6", fmt.Sprint(prefix))
}

// FilterTags filters for probes with the specified tags
// Speficying multiple tags results in an AND behaviour
func (filter *ProbeFilter) FilterTags(tags []string) {
	filter.params.Add("tags", strings.Join(tags, ","))
}

// Sort asks the result list to be sorted by some ordering
// See also: ProbeListSortOrders
func (filter *ProbeFilter) Sort(by string) {
	filter.params.Add("sort", by)
}

// Limit limits the number of result retrieved
func (filter *ProbeFilter) Limit(limit int) {
	filter.limit = limit
}

// Verify sanity of applied filters
func (filter *ProbeFilter) verifyFilters() error {

	// there is a finite set of possible orderings
	if filter.params.Has("sort") {
		sort := filter.params.Get("sort")
		if !ValidProbeListSortOrder(sort) {
			return fmt.Errorf("invalid sort order")
		}
	}

	if filter.params.Has("country_code") {
		cc := filter.params.Get("country_code")
		// TODO: properly verify country code
		if len(cc) != 2 {
			return fmt.Errorf("invalid country code")
		}
	}

	// tags need to be slugs
	if filter.params.Has("tags") {
		tags := strings.Split(filter.params.Get("tags"), ",")
		re, _ := regexp.Compile(`^[\w\-]+$`)
		for _, tag := range tags {
			if !re.MatchString(tag) {
				return fmt.Errorf("invalid tag: %s", tag)
			}
		}
	}

	if filter.limit < 0 {
		return fmt.Errorf("limit must not be negative")
	}
	if filter.id < 0 {
		return fmt.Errorf("ID must not be negative")
	}

	return nil
}

// GetProbeCount returns the count of probes by filtering
func (filter *ProbeFilter) GetProbeCount(
	verbose bool,
) (
	count int,
	err error,
) {
	// sanity checks - late in the process, but not too late
	err = filter.verifyFilters()
	if err != nil {
		return
	}

	// counting needs application of the specified filters
	query := apiBaseURL + "probes/?" + filter.params.Encode()

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", uaString)

	// results are paginated with next= (and previous=)
	if verbose {
		fmt.Printf("API call: GET %s\n", req.URL)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// grab and store the actual content
	var page probeListingPage
	err = json.NewDecoder(resp.Body).Decode(&page)
	if err != nil {
		return 0, err
	}

	// the only really important data point is the count
	return page.Count, nil
}

// GetProbes returns an bunch of probes by filtering
func (filter *ProbeFilter) GetProbes(
	verbose bool,
) (
	probes []Probe,
	err error,
) {
	// special case: a specific ID was "filtered"
	if filter.id != 0 {
		probe, errx := GetProbe(verbose, filter.id)
		if errx != nil {
			err = errx
			return
		}
		probes = make([]Probe, 1)
		probes[0] = *probe
		return
	}

	// sanity checks - late in the process, but not too late
	err = filter.verifyFilters()
	if err != nil {
		return
	}

	query := apiBaseURL + "probes/?" + filter.params.Encode()

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", uaString)

	// results are paginated with next= (and previoous=)
	for {
		if verbose {
			fmt.Printf("API call: GET %s\n", req.URL)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// grab and store the actual content
		var page probeListingPage
		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			return probes, err
		}

		// add elements to the list while observing the limit
		if filter.limit != 0 && len(probes)+len(page.Probes) > filter.limit {
			probes = append(probes, page.Probes[:filter.limit-len(probes)]...)
		} else {
			probes = append(probes, page.Probes...)
		}

		// no next page or got to exactly the limit => we're done
		if page.Next == "" || len(probes) == filter.limit {
			break
		}

		// just follow th enext link
		req, err = http.NewRequest("GET", page.Next, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", uaString)
	}

	return probes, nil
}

// GetProbe retrieves data for a single probe, by ID
// returns probe, nil if a probe was found
// returns nil, nil if no such probe was found
// returns nil, err on error
func GetProbe(
	verbose bool,
	id int,
) (
	*Probe,
	error,
) {
	var probe *Probe

	probeurl := fmt.Sprintf("%sprobes/%d/", apiBaseURL, id)

	req, err := http.NewRequest("GET", probeurl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
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

	if resp.StatusCode == 404 {
		return nil, nil
	}

	// grab and store the actual content
	err = json.NewDecoder(resp.Body).Decode(&probe)
	if err != nil {
		return nil, err
	}

	return probe, nil
}

// ValidProbeListSortOrder checks if a sort order is supported
func ValidProbeListSortOrder(sort string) bool {
	return slices.Contains(ProbeListSortOrders, sort)
}
