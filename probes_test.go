/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"testing"
)

// Test if the filter validator does a decent job
func TestProbeFilterValidator(t *testing.T) {
	var err error
	var filter ProbeFilter

	filter = NewProbeFilter()
	filter.FilterID(-1)
	err = filter.verifyFilters()
	if err == nil {
		t.Fatalf("ID filter cannot be negative")
	}
	filter = NewProbeFilter()
	filter.FilterID(1)
	err = filter.verifyFilters()
	if err != nil {
		t.Fatalf("Correct ID filter is not allowed?")
	}

	badtag := "*"
	goodtag := "ooo"
	filter = NewProbeFilter()
	filter.FilterTags([]string{badtag})
	err = filter.verifyFilters()
	if err == nil {
		t.Fatalf("Bad tag '%s' not filtered properly", badtag)
	}
	filter = NewProbeFilter()
	filter.FilterTags([]string{"ooo"})
	err = filter.verifyFilters()
	if err != nil {
		t.Fatalf("Good tag '%s' is not allowed", goodtag)
	}

	badcc := "NED"
	goodcc := "NL"
	filter = NewProbeFilter()
	filter.FilterCountry(badcc)
	err = filter.verifyFilters()
	if err == nil {
		t.Fatalf("Bad country code '%s' not filtered properly", badcc)
	}
	filter = NewProbeFilter()
	filter.FilterCountry(goodcc)
	err = filter.verifyFilters()
	if err != nil {
		t.Fatalf("Good country code '%s' is not allowed", goodcc)
	}

	filter = NewProbeFilter()
	filter.Sort("abcd")
	err = filter.verifyFilters()
	if err == nil {
		t.Fatalf("Sort order is not filtered properly")
	}

	filter = NewProbeFilter()
	filter.Limit(-1)
	err = filter.verifyFilters()
	if err == nil {
		t.Fatalf("Limit can be negative")
	}
}
