/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package result

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"sort"
)

type PingResult struct {
	BaseResult
	Min        float64 `json:"min"`    // -1 if N/A
	Avg        float64 `json:"avg"`    // -1 if N/A
	Max        float64 `json:"max"`    // -1 if N/A
	Sent       uint    `json:"sent"`   //
	Received   uint    `json:"rcvd"`   //
	Duplicates uint    `json:"dup"`    //
	Size       uint    `json:"size"`   //
	Protocol   string  `json:"proto"`  //
	TTL        *uint   `json:"ttl"`    //
	Step       uint    `json:"step"`   //
	RawResult  []any   `json:"result"` //
}

// one successful ping reply; there are at most Sent ones in a ping result
type PingReply struct {
	Rtt       float64
	Source    *netip.Addr
	Ttl       uint
	Duplicate bool
}

func (result *PingResult) ShortString() string {
	ret := result.BaseShortString() +
		fmt.Sprintf("\t%d/%d/%d\t%f/%f/%f",
			result.Sent, result.Received, result.Duplicates,
			result.Min, result.Avg, result.Max,
		)

	return ret
}

func (result *PingResult) LongString() string {
	return result.ShortString() +
		fmt.Sprintf("\t%s\t%v", result.Protocol, result.ReplyRtts())
}

func (result *PingResult) TypeName() string {
	return "ping"
}

func (ping *PingResult) Parse(from string) (err error) {
	err = json.Unmarshal([]byte(from), &ping)
	if err != nil {
		return err
	}
	if ping.Type != "ping" {
		return fmt.Errorf("this is not a ping result (type=%s)", ping.Type)
	}
	return nil
}

func (result *PingResult) Replies() ([]PingReply, error) {
	r := make([]PingReply, 0)
	for _, item := range result.RawResult {
		mapitem := item.(map[string]any)
		if rtt, ok := mapitem["rtt"]; ok {
			// fill in other fields of a reply struct
			pr := PingReply{Rtt: rtt.(float64)}
			if src, ok := mapitem["src_addr"]; ok {
				src, err := netip.ParseAddr(src.(string))
				if err != nil {
					return r, err
				}
				pr.Source = &src
			} else {
				pr.Source = result.DestinationAddr // TODO: is this correct?
			}
			if ttl, ok := mapitem["ttl"]; ok {
				pr.Ttl = ttl.(uint)
			} else {
				pr.Ttl = *result.TTL
			}
			_, pr.Duplicate = mapitem["dup"]

			r = append(r, pr)
		}
	}
	return r, nil
}

func (result *PingResult) ReplyRtts() []float64 {
	r := make([]float64, 0)
	for _, item := range result.RawResult {
		mapitem := item.(map[string]any)
		if rtt, ok := mapitem["rtt"]; ok {
			r = append(r, rtt.(float64))
		}
	}
	return r
}

func (result *PingResult) Errors() []string {
	r := make([]string, 0)
	for _, item := range result.RawResult {
		mapitem := item.(map[string]any)
		if err, ok := mapitem["error"]; ok {
			r = append(r, err.(string))
		}
	}
	return r
}

func (result *PingResult) Timeouts() int {
	n := 0
	for _, item := range result.RawResult {
		mapitem := item.(map[string]any)
		if _, ok := mapitem["x"]; ok {
			n++
		}
	}
	return n
}

func (result *PingResult) MedianRTT() (float64, error) {
	vals := result.ReplyRtts()
	return median(vals)
}

func median(vals []float64) (float64, error) {
	n := len(vals)
	if n == 0 {
		return 0.0, fmt.Errorf("zero values")
	}
	slice := vals[:]
	sort.Float64s(slice)

	// follow the definition of median
	if n%2 == 0 {
		return (vals[n/2-1] + vals[n/2]) / 2, nil
	} else {
		return vals[n/2], nil
	}
}
