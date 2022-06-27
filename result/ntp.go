/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package result

import (
	"encoding/json"
	"fmt"
)

type NtpResult struct {
	BaseResult
	Protocol           string  `json:"proto"`           //
	Version            uint    `json:"version"`         //
	LeapIndicator      string  `json:"li"`              //
	Mode               string  `json:"mode"`            //
	Stratum            uint    `json:"stratum"`         //
	PollInterval       uint    `json:"poll"`            //
	Precision          float64 `json:"precision"`       //
	RootDelay          float64 `json:"root-delay"`      //
	RootDispersion     float64 `json:"root-dispersion"` //
	ReferenceID        string  `json:"ref-id"`          //
	ReferenceTimestamp float64 `json:"ref-ts"`          //
	RawResult          []any   `json:"result"`          //
}

// one successful ntp reply
type NtpReply struct {
	OriginTimestamp   float64 //
	TransmitTimestamp float64 //
	ReceiveTimestamp  float64 //
	FinalTimestamp    float64 //
	Offset            float64 //
	Rtt               float64 //
}

func (result *NtpReply) String() string {
	return fmt.Sprintf("[%f\t%f\t%f\t%f\t%f\t%f]",
		result.Offset,
		result.Rtt,
		result.OriginTimestamp,
		result.TransmitTimestamp,
		result.ReceiveTimestamp,
		result.FinalTimestamp,
	)
}

func (result *NtpResult) String() string {
	ret := result.BaseString() +
		fmt.Sprintf("\t%s\t%d\t%d\t%d",
			result.ReferenceID, result.Stratum, len(result.Replies()), len(result.Errors()),
		)
	return ret
}

func (result *NtpResult) LongString() string {
	return result.String() +
		fmt.Sprintf("\t%s\t%v", result.Protocol, result.Replies())
}

func (result *NtpResult) TypeName() string {
	return "ntp"
}

func (ntp *NtpResult) Parse(from string) (err error) {
	err = json.Unmarshal([]byte(from), &ntp)
	if err != nil {
		return err
	}
	if ntp.Type != "ntp" {
		return fmt.Errorf("this is not an NTP result (type=%s)", ntp.Type)
	}
	return nil
}

func (result *NtpResult) Replies() []NtpReply {
	r := make([]NtpReply, 0)
	for _, item := range result.RawResult {
		mapitem := item.(map[string]any)
		if rtt, ok := mapitem["rtt"]; ok {
			// fill in other fields of a reply struct
			ntpr := NtpReply{Rtt: rtt.(float64)}
			if offset, ok := mapitem["offset"]; ok {
				ntpr.Offset = offset.(float64)
			}
			if origints, ok := mapitem["origin-ts"]; ok {
				ntpr.OriginTimestamp = origints.(float64)
			}
			if transmitts, ok := mapitem["transmit-ts"]; ok {
				ntpr.TransmitTimestamp = transmitts.(float64)
			}
			if receivets, ok := mapitem["receive-ts"]; ok {
				ntpr.ReceiveTimestamp = receivets.(float64)
			}
			if finalts, ok := mapitem["final-ts"]; ok {
				ntpr.FinalTimestamp = finalts.(float64)
			}

			r = append(r, ntpr)
		}
	}
	return r
}

func (result *NtpResult) Errors() []string {
	r := make([]string, 0)
	for _, item := range result.RawResult {
		mapitem := item.(map[string]any)
		if err, ok := mapitem["x"]; ok {
			r = append(r, err.(string))
		}
	}
	return r
}
