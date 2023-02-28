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
	Interval       *uint     `json:"interval,omitempty"`
	ResolveOnProbe *bool     `json:"resolve_on_probe,omitempty"`
	Tags           *[]string `json:"tags,omitempty"`
	Spread         *uint     `json:"spread,omitempty"`
	SkipDNSCheck   *bool     `json:"skip_dns_check,omitempty"`
}

type measurementTargetPing struct {
	measurementTargetBase
	Packets        *uint `json:"packets,omitempty"`
	PacketSize     *uint `json:"packet_size,omitempty"`
	PacketInterval *uint `json:"packet_interval,omitempty"`
	IncludeProbeID *bool `json:"include_probe_id,omitempty"`
}

type measurementTargetTrace struct {
	measurementTargetBase
	Protocol        string `json:"protocol"`
	ResponseTimeout *uint  `json:"response_timeout,omitempty"`
	Packets         *uint  `json:"packets,omitempty"`
	PacketSize      *uint  `json:"packet_size,omitempty"` // ?
	ParisId         uint   `json:"paris,omitempty"`
	FirstHop        *uint  `json:"first_hop,omitempty"`
	LastHop         *uint  `json:"max_hops,omitempty"`
	DestinationEH   *uint  `json:"destination_option_size,omitempty"`
	HopByHopEH      *uint  `json:"hop_by_hop_option_size,omitempty"`
	DontFragment    *bool  `json:"dont_fragment,omitempty"`
}

type measurementTargetDns struct {
	measurementTargetBase
	Protocol       string  `json:"protocol"`
	Class          string  `json:"query_class"`
	Type           string  `json:"query_type"`
	Argument       *string `json:"query_argument,omitempty"`
	UseMacros      *bool   `json:"use_macros,omitempty"`
	UseResolver    *bool   `json:"use_probe_resolver,omitempty"`
	Nsid           *bool   `json:"set_nsid_bit,omitempty"`
	UdpPayloadSize *uint   `json:"udp_payload_size,omitempty"`
	Retries        *uint   `json:"retry,omitempty"`
	IncludeQbuf    *bool   `json:"include_qbuf,omitempty"`
	IncludeAbuf    *bool   `json:"include_abuf,omitempty"`
	PrependProbeID *bool   `json:"prepend_probe_id,omitempty"`
	SetRd          *bool   `json:"set_rd_bit,omitempty"`
	SetDo          *bool   `json:"set_do_bit,omitempty"`
	SetCd          *bool   `json:"set_cd_bit,omitempty"`
	Timeout        *uint   `json:"timeout,omitempty"`
}

type measurementTargetTls struct {
	measurementTargetBase
	Port uint    `json:"port"`
	Sni  *string `json:"hostname,omitempty"`
}

type measurementTargetNtp struct {
	measurementTargetBase
	Packets *uint `json:"packets,omitempty"`
	Timeout *uint `json:"timeout,omitempty"`
}

// various measurement options
type BaseOptions struct {
	ResolveOnProbe bool
	Interval       uint
	Tags           []string
	Spread         uint
	SkipDNSCheck   bool
}
type PingOptions struct {
	Packets        uint // API default: 3
	PacketSize     uint // API default: 48 bytes
	PacketInterval uint // Time between packets (ms)
	IncludeProbeID bool // Include the probe ID (encoded as ASCII digits) as part of the payload
}
type TraceOptions struct {
	Protocol        string // default: UDP
	ResponseTimeout uint   // API default: 4000 (ms)
	Packets         uint   // API default: 3
	PacketSize      uint   // API default: 48 bytes
	ParisId         uint   // API default: 16, default: 0
	FirstHop        uint   // API default: 1
	LastHop         uint   // API default: 32
	DestinationEH   uint   // API default: 0
	HopByHopEH      uint   // API default: 0
	DontFragment    bool   // API default: false
}
type DnsOptions struct {
	Protocol       string // default: UDP
	Class          string
	Type           string
	Argument       string
	UseMacros      bool // API default: false
	UseResolver    bool // API default: false
	Nsid           bool // API default: false
	UdpPayloadSize uint // API default: 512
	Retries        uint // API default: 0
	IncludeQbuf    bool // API default: false
	IncludeAbuf    bool // API default: false
	PrependProbeID bool // API default: false
	SetRd          bool // API default: false
	SetDo          bool // API default: false
	SetCd          bool // API default: false
	Timeout        uint // API default: 5000 (ms)
}
type TlsOptions struct {
	Port uint // API default: 443
	Sni  string
}
type NtpOptions struct {
	Packets uint // API default: 3
	Timeout uint // API default: 4000 (ms)
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
var traceprotocols = []string{"ICMP", "UDP", "TCP"}
var dnsprotocols = []string{"UDP", "TCP"}
var dnsclasses = []string{"IN", "CHAOS"}
var dnstypes = []string{"A", "AAAA", "ANY", "CNAME", "DNSKEY", "DS", "MX", "NS", "NSEC", "PTR", "RRSIG", "SOA", "TXT", "SRV", "NAPTR", "TLSA"}

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
	return spec.AddProbesListWithTags(list, nil, nil)
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

func (spec *MeasurementSpec) AddProbesListWithTags(list []uint, tagsincl *[]string, tagsexcl *[]string) error {
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
	typ string,
	description string,
	target string,
	af uint,
	baseoptions *BaseOptions,
) error {
	if description == "" {
		return fmt.Errorf("description cannot be empty")
	}
	if target == "" {
		return fmt.Errorf("target cannot be empty")
	}
	if af != 4 && af != 6 {
		return fmt.Errorf("address familty must be 4 or 6")
	}

	// common fields
	def.Type = typ
	def.Description = description
	def.Target = target
	def.AddressFamily = af

	if baseoptions != nil {
		if baseoptions.ResolveOnProbe {
			def.ResolveOnProbe = &baseoptions.ResolveOnProbe
		}
		def.SkipDNSCheck = &baseoptions.SkipDNSCheck
		if baseoptions.Interval != 0 {
			def.Interval = &baseoptions.Interval
		}
		if baseoptions.Interval != 0 {
			def.Interval = &baseoptions.Interval
		}
		if baseoptions.Spread != 0 {
			def.Spread = &baseoptions.Spread
		}
		if baseoptions.Tags != nil {
			def.Tags = &baseoptions.Tags
		}
	}

	return nil
}

func (spec *MeasurementSpec) AddPing(
	description string,
	target string,
	af uint,
	baseoptions *BaseOptions,
	pingoptions *PingOptions,
) error {
	var def = new(measurementTargetPing)

	if err := def.addCommonFields("ping", description, target, af, baseoptions); err != nil {
		return err
	}

	// ping specific fields
	if pingoptions != nil {
		if pingoptions.Packets != 0 {
			def.Packets = &pingoptions.Packets
		}
		if pingoptions.PacketSize != 0 {
			def.PacketSize = &pingoptions.PacketSize
		}
		if pingoptions.PacketInterval != 0 {
			def.PacketInterval = &pingoptions.PacketInterval
		}
		if pingoptions.IncludeProbeID {
			def.IncludeProbeID = &pingoptions.IncludeProbeID
		}
	}

	spec.apiSpec.Definitons = append(spec.apiSpec.Definitons, def)

	return nil
}

func (spec *MeasurementSpec) AddTrace(
	description string,
	target string,
	af uint,
	baseoptions *BaseOptions,
	traceoptions *TraceOptions,
) error {
	var def = new(measurementTargetTrace)

	if err := def.addCommonFields("traceroute", description, target, af, baseoptions); err != nil {
		return err
	}

	// explicit defaults
	def.Protocol = "UDP"
	def.ParisId = 0

	// trace specific fields
	if traceoptions != nil {
		if traceoptions.Protocol != "" &&
			slices.Contains(traceprotocols, traceoptions.Protocol) {
			def.Protocol = traceoptions.Protocol
		}
		if traceoptions.ResponseTimeout != 0 {
			def.ResponseTimeout = &traceoptions.ResponseTimeout
		}
		if traceoptions.Packets != 0 {
			def.Packets = &traceoptions.Packets
		}
		if traceoptions.PacketSize != 0 {
			def.PacketSize = &traceoptions.PacketSize
		}
		if traceoptions.ParisId != 0 {
			def.ParisId = traceoptions.ParisId
		}
		if traceoptions.FirstHop != 0 {
			def.FirstHop = &traceoptions.FirstHop
		}
		if traceoptions.LastHop != 0 {
			def.LastHop = &traceoptions.LastHop
		}
		if traceoptions.DestinationEH != 0 {
			def.DestinationEH = &traceoptions.DestinationEH
		}
		if traceoptions.HopByHopEH != 0 {
			def.HopByHopEH = &traceoptions.HopByHopEH
		}
		if traceoptions.DontFragment {
			def.DontFragment = &traceoptions.DontFragment
		}
	}

	spec.apiSpec.Definitons = append(spec.apiSpec.Definitons, def)

	return nil
}

func (spec *MeasurementSpec) AddDns(
	description string,
	target string,
	af uint,
	baseoptions *BaseOptions,
	dnsoptions *DnsOptions,
) error {
	var def = new(measurementTargetDns)

	if err := def.addCommonFields("dns", description, target, af, baseoptions); err != nil {
		return err
	}

	// explicit defaults
	def.Protocol = "UDP"
	def.Class = "IN"
	def.Type = "A"

	// dns specific fields
	if dnsoptions != nil {
		if dnsoptions.Protocol != "" &&
			slices.Contains(dnsprotocols, dnsoptions.Protocol) {
			def.Protocol = dnsoptions.Protocol
		}
		if dnsoptions.Class != "" &&
			slices.Contains(dnsclasses, dnsoptions.Class) {
			def.Class = dnsoptions.Class
		}
		if dnsoptions.Type != "" &&
			slices.Contains(dnstypes, dnsoptions.Type) {
			def.Type = dnsoptions.Type
		}
		if dnsoptions.Argument != "" {
			def.Argument = &dnsoptions.Argument
		}
		if dnsoptions.UseMacros {
			def.UseMacros = &dnsoptions.UseMacros
		}
		if dnsoptions.UseResolver {
			def.UseResolver = &dnsoptions.UseResolver
		}
		if dnsoptions.Nsid {
			def.Nsid = &dnsoptions.Nsid
		}
		if dnsoptions.UdpPayloadSize != 0 {
			def.UdpPayloadSize = &dnsoptions.UdpPayloadSize
		}
		if dnsoptions.Retries != 0 {
			def.Retries = &dnsoptions.Retries
		}
		if dnsoptions.IncludeQbuf {
			def.IncludeQbuf = &dnsoptions.IncludeQbuf
		}
		if dnsoptions.IncludeAbuf {
			def.IncludeAbuf = &dnsoptions.IncludeAbuf
		}
		if dnsoptions.PrependProbeID {
			def.PrependProbeID = &dnsoptions.PrependProbeID
		}
		if dnsoptions.SetRd {
			def.SetRd = &dnsoptions.SetRd
		}
		if dnsoptions.SetDo {
			def.SetDo = &dnsoptions.SetDo
		}
		if dnsoptions.SetRd {
			def.SetCd = &dnsoptions.SetCd
		}
		if dnsoptions.Timeout != 0 {
			def.Timeout = &dnsoptions.Timeout
		}
	}

	spec.apiSpec.Definitons = append(spec.apiSpec.Definitons, def)

	return nil
}

func (spec *MeasurementSpec) AddTls(
	description string,
	target string,
	af uint,
	baseoptions *BaseOptions,
	tlsoptions *TlsOptions,
) error {
	var def = new(measurementTargetTls)

	if err := def.addCommonFields("sslcert", description, target, af, baseoptions); err != nil {
		return err
	}

	// explicit defaults
	def.Port = 443

	// TLS specific fields
	if tlsoptions != nil {
		if tlsoptions.Port != 0 {
			def.Port = tlsoptions.Port
		}
		if tlsoptions.Sni != "" {
			def.Sni = &tlsoptions.Sni
		}
	}

	spec.apiSpec.Definitons = append(spec.apiSpec.Definitons, def)

	return nil
}

func (spec *MeasurementSpec) AddNtp(
	description string,
	target string,
	af uint,
	baseoptions *BaseOptions,
	ntpoptions *NtpOptions,
) error {
	var def = new(measurementTargetNtp)

	if err := def.addCommonFields("ntp", description, target, af, baseoptions); err != nil {
		return err
	}

	// explicit defaults

	// NTP specific fields
	if ntpoptions != nil {
		if ntpoptions.Packets != 0 {
			def.Packets = &ntpoptions.Packets
		}
		if ntpoptions.Timeout != 0 {
			def.Timeout = &ntpoptions.Timeout
		}
	}

	spec.apiSpec.Definitons = append(spec.apiSpec.Definitons, def)

	return nil
}

func (target *measurementTargetPing) MarshalJSON() (b []byte, e error) {
	return json.Marshal(*target)
}
func (target *measurementTargetTrace) MarshalJSON() (b []byte, e error) {
	return json.Marshal(*target)
}
func (target *measurementTargetDns) MarshalJSON() (b []byte, e error) {
	return json.Marshal(*target)
}
func (target *measurementTargetTls) MarshalJSON() (b []byte, e error) {
	return json.Marshal(*target)
}
func (target *measurementTargetNtp) MarshalJSON() (b []byte, e error) {
	return json.Marshal(*target)
}

func (spec *MeasurementSpec) Submit(verbose bool) error {
	if len(spec.apiSpec.Definitons) == 0 {
		return fmt.Errorf("need at least 1 measurement defintion")
	}

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
