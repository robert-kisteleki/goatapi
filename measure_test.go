/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"reflect"
	"testing"
)

// Test if the measurement probe spec part works
func TestMeasureProbes(t *testing.T) {
	var err error
	var spec MeasurementSpec

	spec = MeasurementSpec{}
	err = spec.AddProbesArea("X", 5)
	if err == nil {
		t.Errorf("Invalid probe area is accepted")
	}

	for _, area := range areas {
		spec = MeasurementSpec{}
		err = spec.AddProbesArea(area, 5)
		if err != nil {
			t.Errorf("Valid probe area %s is not accepted: %v", area, err)
		}
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesCountry("NED", 5)
	if err == nil {
		t.Errorf("Invalid country is accepted")
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesCountry("NL", 5)
	if err != nil {
		t.Errorf("Valid country is not accepted: %v", err)
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesList([]uint{})
	if err == nil {
		t.Errorf("Empty probe list is accepted")
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesList([]uint{1, 2, 3})
	if err != nil {
		t.Errorf("Good probe list is not accepted: %v", err)
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesReuse(10, 5)
	if err == nil {
		t.Errorf("Bad msm reuse ID is accepted")
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesReuse(1000001, 5)
	if err != nil {
		t.Errorf("Good msm reuse ID is not accepted: %v", err)
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesReuse(10, 5)
	if err == nil {
		t.Errorf("Bad msm reuse ID is accepted")
	}

	spec = MeasurementSpec{}
	err = spec.AddProbesReuseWithTags(1000001, 5, &[]string{"itag1", "itag2"}, &[]string{"etag1", "etag2"})
	if err != nil {
		t.Errorf("Good probe tag filter list is not accepted: %v", err)
	}
	pspec := spec.apiSpec.Probes[0]
	if pspec.Tags.Include == nil ||
		pspec.Tags.Exclude == nil ||
		!reflect.DeepEqual(*pspec.Tags.Include, []string{"itag1", "itag2"}) ||
		!reflect.DeepEqual(*pspec.Tags.Exclude, []string{"etag1", "etag2"}) {
		t.Errorf("Probe filter tag list is not registered properly")
	}
}

// Test if the measurement target spec part works
func TestMeasureTargetGeneral(t *testing.T) {
	var err error
	var spec MeasurementSpec

	spec = MeasurementSpec{}
	err = spec.AddPing("description", "www.google.com", 3, nil, nil)
	if err == nil {
		t.Errorf("Invalid address family is accepted")
	}
	err = spec.AddPing("description", "www.google.com", 4, nil, nil)
	if err != nil {
		t.Errorf("Valid address family is not accepted: %v", err)
	}
	err = spec.AddPing("description", "www.google.com", 6, nil, nil)
	if err != nil {
		t.Errorf("Valid address family is not accepted: %v", err)
	}

	spec = MeasurementSpec{}
	err = spec.AddPing("", "", 4, nil, nil)
	if err == nil {
		t.Errorf("Empty description or target is accepted")
	}

	spec = MeasurementSpec{}
	err = spec.AddPing("description1", "www.google.com", 4, nil, nil)
	if err != nil {
		t.Errorf("Valid ping measurement target spec rejected: %v", err)
	}
}

func TestMeasureTargetBase(t *testing.T) {
	var err error
	var spec MeasurementSpec

	err = spec.AddPing("description1", "www.meta.com", 6, &BaseOptions{
		ResolveOnProbe: true,
		Interval:       999,
		Tags:           []string{"tag1", "tag2"},
		Spread:         50,
		SkipDNSCheck:   true,
	}, nil)
	b2, err := spec.apiSpec.Definitons[0].MarshalJSON()
	if err != nil {
		t.Fatalf("Ping measurement target spec with options failed to marshal to JSON: %v", err)
	}
	if string(b2) != `{"description":"description1","target":"www.meta.com","type":"ping","af":6,"interval":999,"resolve_on_probe":true,"tags":["tag1","tag2"],"spread":50,"skip_dns_check":true}` {
		t.Errorf("Measurement with base options improperly serialized: %s", string(b2))
	}

}

func TestMeasureTargetPing(t *testing.T) {
	var err error
	var spec MeasurementSpec

	err = spec.AddPing("description2", "www.meta.com", 6, nil, &PingOptions{
		Packets:        5,
		PacketSize:     9,
		PacketInterval: 99,
		IncludeProbeID: true,
	})
	b2, err := spec.apiSpec.Definitons[0].MarshalJSON()
	if err != nil {
		t.Fatalf("Ping measurement target spec with options failed to marshal to JSON: %v", err)
	}
	if string(b2) != `{"description":"description2","target":"www.meta.com","type":"ping","af":6,"packets":5,"packet_size":9,"packet_interval":99,"include_probe_id":true}` {
		t.Errorf("Measurement (ping) with options improperly serialized: %s", string(b2))
	}

}

func TestMeasureTargetTrace(t *testing.T) {
	var err error
	var spec MeasurementSpec

	err = spec.AddTrace("description2", "www.meta.com", 6, nil, &TraceOptions{
		Protocol:        "ICMP",
		ResponseTimeout: 19,
		Packets:         5,
		PacketSize:      9,
		ParisId:         8,
		FirstHop:        2,
		LastHop:         3,
		DestinationEH:   4,
		HopByHopEH:      7,
		DontFragment:    true,
	})
	b2, err := spec.apiSpec.Definitons[0].MarshalJSON()
	if err != nil {
		t.Fatalf("Trace measurement target spec with options failed to marshal to JSON: %v", err)
	}
	if string(b2) != `{"description":"description2","target":"www.meta.com","type":"traceroute","af":6,"protocol":"ICMP","response_timeout":19,"packets":5,"packet_size":9,"paris":8,"first_hop":2,"max_hops":3,"destination_option_size":4,"hop_by_hop_option_size":7,"dont_fragment":true}` {
		t.Errorf("Measurement (trace) with options improperly serialized: %s", string(b2))
	}

}

// Test if the measurement spec generator works
func TestMeasureSpec(t *testing.T) {
	var err error
	var spec MeasurementSpec

	spec = MeasurementSpec{}
	err = spec.Submit(false)
	if err == nil {
		t.Errorf("Measurement spec without probes or targets is accepted")
	}

	spec = MeasurementSpec{}
	spec.AddProbesArea("WW", 3)
	err = spec.Submit(false)
	if err == nil {
		t.Errorf("Measurement spec without targets is accepted")
	}

	spec = MeasurementSpec{}
	spec.AddPing("ping", "google.com", 4, nil, nil)
	err = spec.Submit(false)
	if err == nil {
		t.Errorf("Measurement spec without probes is accepted")
	}
}
