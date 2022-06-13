/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goatapi

import (
	"testing"
)

// Test if the filter validator does a decent job
func TestAnchorFilterValidator(t *testing.T) {
	var err error
	var filter AnchorFilter

	badcc := "NED"
	goodcc := "NL"
	filter = NewAnchorFilter()
	filter.FilterCountry(badcc)
	err = filter.verifyFilters()
	if err == nil {
		t.Errorf("Bad country code '%s' not filtered properly", badcc)
	}
	filter = NewAnchorFilter()
	filter.FilterCountry(goodcc)
	err = filter.verifyFilters()
	if err != nil {
		t.Errorf("Good country code '%s' is not allowed", goodcc)
	}

	filter = NewAnchorFilter()
	filter.FilterID(-1)
	err = filter.verifyFilters()
	if err == nil {
		t.Error("ID filter cannot be negative")
	}
	filter = NewAnchorFilter()
	filter.FilterID(1)
	err = filter.verifyFilters()
	if err != nil {
		t.Error("Correct ID filter is not allowed?")
	}

	filter = NewAnchorFilter()
	filter.Limit(-1)
	err = filter.verifyFilters()
	if err == nil {
		t.Error("Limit can be negative")
	}
}
