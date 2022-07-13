/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package result

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"time"
)

type BaseResult struct {
	FirmwareVersion uint        `json:"fw"`               //
	CodeVersion     string      `json:"mver"`             //
	MeasurementID   uint        `json:"msm_id"`           //
	GroupID         uint        `json:"group_id"`         //
	ProbeID         uint        `json:"prb_id"`           //
	MeasurementName string      `json:"msm_name"`         // measurement name (better use type)
	Type            string      `json:"type"`             // measurement type
	TimeStamp       uniTime     `json:"timestamp"`        // when was this result collected
	StoreTimeStamp  uniTime     `json:"stored_timestamp"` // when was this result stored
	Bundle          uint        `json:"bundle"`           // ID for a collection of related measurement results
	LastTimeSync    int         `json:"lts"`              // how long ago was the probe's clock synced
	DestinationName string      `json:"dst_name"`         //
	DestinationAddr *netip.Addr `json:"dst_addr"`         //
	SourceAddr      netip.Addr  `json:"src_addr"`         // source address used by probe
	FromAddr        netip.Addr  `json:"from"`             // IP address of the probe as known by the infra
	AddressFamily   uint        `json:"af"`               // 4 or 6
	ResolveTime     *float64    `json:"ttr"`              // only if resolve-on-probe was used
}

func (result *BaseResult) Parse(from string) (err error) {
	err = json.Unmarshal([]byte(from), &result)
	if err != nil {
		return err
	}
	return nil
}

func (result *BaseResult) String() string {
	d := "N/A"
	if result.DestinationName != "" {
		d = result.DestinationName
	}
	ret := fmt.Sprintf("%d\t%d\t%v\t%s",
		result.MeasurementID,
		result.ProbeID,
		result.TimeStamp,
		d,
	)
	ret += valueOrNA("", false, result.DestinationAddr)
	return ret
}

func (result *BaseResult) DetailString() string {
	return result.BaseDetailString() +
		fmt.Sprintf("\t%d", result.AddressFamily)
}

func (result *BaseResult) BaseString() string {
	return result.String()
}

func (result *BaseResult) BaseDetailString() string {
	return result.BaseString()
}

func (result *BaseResult) TypeName() string {
	return result.Type
}

func (result *BaseResult) GetTimeStamp() time.Time {
	return time.Time(result.TimeStamp)
}

func (result *BaseResult) GetProbeID() uint {
	return result.ProbeID
}
