/*
  (C) 2022, 2023 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/exp/slices"
)

// Measurement specification object, to be passed to the API
type MeasurementSpec struct {
	apiSpec measurementSpec
}

type measurementSpec struct {
	Definitons []measurementTargetDefinition `json:"definitions"`
	Probes     []measurementProbeDefinition  `json:"probes"`
	OneOff     bool                          `json:"is_oneoff"`
	BillTo     *string                       `json:"bill_to,omitempty"`
	Start      *uniTime                      `json:"start_time,omitempty"`
	Stop       *uniTime                      `json:"stop_time,omitempty"`
}

type measurementTargetDefinition interface {
	MarshalJSON() (b []byte, e error)
}

type measurementTargetBase struct {
	Description    string    `json:"description"`
	Target         string    `json:"target"`
	Type           string    `json:"type"`
	AddressFamily  uint      `json:"af"`
	ResolveOnProbe bool      `json:"resolve_on_probe"`
	Tags           *[]string `json:"tags,omitempty"`
	Spread         *uint     `json:"spread,omitempty"`
	SkipDNSCheck   bool      `json:"skip_dns_check,omitempty"`
}

type measurementTargetPing struct {
	measurementTargetBase
	Packets        uint  `json:"packets"`
	PacketSize     uint  `json:"packet_size"`
	Interval       *uint `json:"interval,omitempty"`
	IncludeProbeID *bool `json:"include_probe_id,omitempty"`
}

// various measuerement options
type PingOptions struct {
	Tags           []string
	Spread         uint
	SkipDNSCheck   bool
	Packets        uint // default: 3
	PacketSize     uint // default: 48 bytes
	Interval       uint // Time between packets (ms)
	IncludeProbeID bool // Include the probe ID (encoded as ASCII digits) as part of the payload
}

type measurementProbeDefinition struct {
	Type      string                          `json:"type"`
	Value     string                          `json:"value"`
	Requested int                             `json:"requested"`
	Tags      *measurementProbeDefinitionTags `json:"tags,omitempty"`
}

type measurementProbeDefinitionTags struct {
	Include *[]string `json:"include,omitempty"`
	Exclude *[]string `json:"exclude,omitempty"`
}

var areas = []string{"WW", "West", "North-Central", "South-Central", "North-East", "South-East"}

func NewMeasurementSpec() (spec *MeasurementSpec) {
	spec = new(MeasurementSpec)
	spec.apiSpec.Definitons = make([]measurementTargetDefinition, 0)
	spec.apiSpec.Probes = make([]measurementProbeDefinition, 0)
	return spec
}

func (spec *MeasurementSpec) Start(time time.Time) {
	t := uniTime(time)
	spec.apiSpec.Start = &t
}

func (spec *MeasurementSpec) Stop(time time.Time) {
	t := uniTime(time)
	spec.apiSpec.Stop = &t
}

func (spec *MeasurementSpec) OneOff(oneoff bool) {
	spec.apiSpec.OneOff = oneoff
}

func (spec *MeasurementSpec) BillTo(billto string) {
	spec.apiSpec.BillTo = &billto
}

func (spec *MeasurementSpec) addProbeSet(
	settype string,
	setvalue string,
	n int,
	tagsincl *[]string,
	tagsexcl *[]string,
) error {
	if n < -1 || n == 0 {
		return fmt.Errorf("number of probes requested should be positive")
	}
	msp := measurementProbeDefinition{
		Type:      settype,
		Value:     setvalue,
		Requested: n,
	}
	if (tagsincl != nil && len(*tagsincl) > 0) || (tagsexcl != nil && len(*tagsexcl) > 0) {
		msp.Tags = new(measurementProbeDefinitionTags)
		if tagsincl != nil && len(*tagsincl) > 0 {
			msp.Tags.Include = tagsincl
		}
		if tagsexcl != nil && len(*tagsexcl) > 0 {
			msp.Tags.Exclude = tagsexcl
		}
	}
	spec.apiSpec.Probes = append(spec.apiSpec.Probes, msp)
	return nil
}

func (spec *MeasurementSpec) AddProbesArea(area string, n int) error {
	return spec.AddProbesAreaWithTags(area, n, nil, nil)
}

func (spec *MeasurementSpec) AddProbesCountry(cc string, n int) error {
	return spec.AddProbesCountryWithTags(cc, n, nil, nil)
}

func (spec *MeasurementSpec) AddProbesList(list []uint) error {
	return spec.AddProbesListWitTags(list, nil, nil)
}

func (spec *MeasurementSpec) AddProbesReuse(msm uint, n int) error {
	return spec.AddProbesReuseWithTags(msm, n, nil, nil)
}

func (spec *MeasurementSpec) AddProbesAreaWithTags(area string, n int, tagsincl *[]string, tagsexcl *[]string) error {
	if !slices.Contains(areas, area) {
		return fmt.Errorf("invalid area: %v", area)
	}
	return spec.addProbeSet("area", area, n, tagsincl, tagsexcl)
}

func (spec *MeasurementSpec) AddProbesCountryWithTags(cc string, n int, tagsincl *[]string, tagsexcl *[]string) error {
	if len(cc) != 2 { // TODO: add proper country code validation
		return fmt.Errorf("invalid country code %v", cc)
	}
	return spec.addProbeSet("cc", cc, n, tagsincl, tagsexcl)
}

func (spec *MeasurementSpec) AddProbesListWitTags(list []uint, tagsincl *[]string, tagsexcl *[]string) error {
	n := len(list)
	if n == 0 {
		return fmt.Errorf("probe list cannot be empty")
	}
	return spec.addProbeSet("probes", makeCsv(list), n, tagsincl, tagsexcl)
}

func (spec *MeasurementSpec) AddProbesReuseWithTags(msm uint, n int, tagsincl *[]string, tagsexcl *[]string) error {
	if msm <= 1000000 {
		return fmt.Errorf("measurement ID must be >1M")
	}
	return spec.addProbeSet("msm", fmt.Sprintf("%d", msm), n, tagsincl, tagsexcl)
}

func (def *measurementTargetBase) addCommonFields(
	description string,
	typ string,
	target string,
	resolveOnProbe bool,
	af uint,
) {
	// common fields
	def.Type = "ping"
	def.Description = description
	def.Target = target
	def.ResolveOnProbe = resolveOnProbe
	def.AddressFamily = af
}

func (spec *MeasurementSpec) AddPing(
	description string,
	target string,
	resolveOnProbe bool,
	af uint,
	options *PingOptions,
) error {
	var def = new(measurementTargetPing)

	if description == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if target == "" {
		return fmt.Errorf("target cannot be empty")
	}
	if af != 4 && af != 6 {
		return fmt.Errorf("address familty must be 4 or 6")
	}

	def.addCommonFields("ping", description, target, resolveOnProbe, af)

	// defaults
	def.PacketSize = 48
	def.Packets = 3

	// options & specific fields
	if options != nil {
		def.SkipDNSCheck = options.SkipDNSCheck
		// tags spread, skipdnscheck
		if options.Tags != nil {
			def.Tags = &options.Tags
		}
		if options.Spread != 0 {
			def.Spread = &options.Spread
		}

		// ping specific fields
		if options.Packets != 0 {
			def.Packets = options.Packets
		}
		if options.PacketSize != 0 {
			def.PacketSize = options.PacketSize
		}
		if options.Interval != 0 {
			def.Interval = &options.Interval
		}
		if options.IncludeProbeID {
			def.IncludeProbeID = &options.IncludeProbeID
		}
	}

	spec.apiSpec.Definitons = append(spec.apiSpec.Definitons, def)

	return nil
}

func (target *measurementTargetPing) MarshalJSON() (b []byte, e error) {
	return json.Marshal(*target)
}

func (spec *MeasurementSpec) Submit(verbose bool) error {
	//if len(spec.apiSpec.Definitons) == 0 {
	//	return fmt.Errorf("need at least 1 measurement defintion")
	//}

	if len(spec.apiSpec.Probes) == 0 {
		return fmt.Errorf("need at least 1 probe specification")
	}

	b, err := json.Marshal(spec.apiSpec)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(b))
	return nil
}
