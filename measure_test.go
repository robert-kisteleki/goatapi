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
	err = spec.AddPing("description", "www.google.com", true, 3, nil)
	if err == nil {
		t.Errorf("Invalid address family is accepted")
	}
	err = spec.AddPing("description", "www.google.com", true, 4, nil)
	if err != nil {
		t.Errorf("Valid address family is not accepted: %v", err)
	}
	err = spec.AddPing("description", "www.google.com", true, 6, nil)
	if err != nil {
		t.Errorf("Valid address family is not accepted: %v", err)
	}

	spec = MeasurementSpec{}
	err = spec.AddPing("", "", true, 4, nil)
	if err == nil {
		t.Errorf("Empty description or target is accepted")
	}
}

func TestMeasureTargetPing(t *testing.T) {
	var err error
	var spec MeasurementSpec

	spec = MeasurementSpec{}
	err = spec.AddPing("description1", "www.google.com", true, 4, nil)
	if err != nil {
		t.Errorf("Valid ping measurement target spec rejected: %v", err)
	}
	b1, err := spec.apiSpec.Definitons[0].MarshalJSON()
	if err != nil {
		t.Fatalf("Ping measurement target spec failed to marshal to JSON: %v", err)
	}
	if string(b1) != `{"description":"description1","target":"www.google.com","type":"ping","af":4,"resolve_on_probe":true,"packets":3,"packet_size":48}` {
		t.Errorf("Measurement (ping) improperly serialized: %s", string(b1))
	}

	err = spec.AddPing("description2", "www.meta.com", false, 6, &PingOptions{
		Tags:           []string{"tag1", "tag2"},
		Spread:         50,
		SkipDNSCheck:   true,
		Packets:        5,
		PacketSize:     99,
		Interval:       999,
		IncludeProbeID: true,
	})
	b2, err := spec.apiSpec.Definitons[1].MarshalJSON()
	if err != nil {
		t.Fatalf("Ping measurement target spec with options failed to marshal to JSON: %v", err)
	}
	if string(b2) != `{"description":"description2","target":"www.meta.com","type":"ping","af":6,"resolve_on_probe":false,"tags":["tag1","tag2"],"spread":50,"skip_dns_check":true,"packets":5,"packet_size":99,"interval":999,"include_probe_id":true}` {
		t.Errorf("Measurement (ping) improperly serialized: %s", string(b2))
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
	spec.AddPing("ping", "google.com", false, 4, nil)
	err = spec.Submit(false)
	if err == nil {
		t.Errorf("Measurement spec without probes is accepted")
	}
}
