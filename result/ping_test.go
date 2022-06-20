/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package result

import (
	"fmt"
	"net/netip"
	"reflect"
	"testing"
)

// Test if the ping parser does a decent job
func TestProbeParser(t *testing.T) {
	var ping PingResult
	err := ping.Parse(`
{"fw":5040,
"mver":"2.4.1",
"lts":20,
"dst_name":"example.com",
"ttr":1.234,
"af":4,
"dst_addr":"10.1.2.3",
"src_addr":"10.2.3.4",
"proto":"ICMP",
"ttl":54,
"size":64,
"result":[
	{"rtt":10.000},
	{"rtt":15.000},
	{"rtt":4.750},
	{"rtt":5.250},
	{"rtt":25.000}
],
"dup":0,
"rcvd":4,
"sent":4,
"min":4.75,
"max":25.0,
"avg":12.0,
"msm_id":1234567,
"prb_id":2345678,
"timestamp":1655443320,
"msm_name":"Ping",
"from":"192.168.1.1",
"type":"ping",
"group_id":34567890,
"step":10,
"stored_timestamp":1655443322
}
`)
	if err != nil {
		t.Fatalf("Error parsing ping result: %s", err)
	}

	assertEqual(t, ping.FirmwareVersion, uint(5040), "error parsing ping field value for fw")
	assertEqual(t, ping.LastTimeSync, 20, "error parsing ping field value for lts")
	assertEqual(t, ping.DestinationName, "example.com", "error parsing ping field value for dst_name")
	dst, _ := netip.ParseAddr("10.1.2.3")
	assertEqual(t, *ping.DestinationAddr, dst, "error parsing ping field value for dst_addr")
	src, _ := netip.ParseAddr("10.2.3.4")
	assertEqual(t, ping.SourceAddr, src, "error parsing ping field value for src_addr")
	from, _ := netip.ParseAddr("192.168.1.1")
	assertEqual(t, ping.FromAddr, from, "error parsing ping field value for from")
	assertEqual(t, *ping.ResolveTime, 1.234, "error parsing ping field value for ttr")
	assertEqual(t, ping.AddressFamily, uint(4), "error parsing ping field value for af")
	assertEqual(t, ping.Protocol, "ICMP", "error parsing ping field value for proto")
	assertEqual(t, *ping.TTL, uint(54), "error parsing ping field value for ttl")
	assertEqual(t, ping.Size, uint(64), "error parsing ping field value for size")
	assertEqual(t, ping.Duplicates, uint(0), "error parsing ping field value for dup")
	assertEqual(t, ping.Received, uint(4), "error parsing ping field value for rcvd")
	assertEqual(t, ping.Sent, uint(4), "error parsing ping field value for sent")
	assertEqual(t, ping.MeasurementID, uint(1234567), "error parsing ping field value for msm_id")
	assertEqual(t, ping.ProbeID, uint(2345678), "error parsing ping field value for prb_id")
	assertEqual(t, ping.GroupID, uint(34567890), "error parsing ping field value for group_id")
	assertEqual(t, ping.TimeStamp.String(), "2022-06-17T05:22:00Z", "error parsing ping field value for timestamp")
	assertEqual(t, ping.StoreTimeStamp.String(), "2022-06-17T05:22:02Z", "error parsing ping field value for stored_timestamp")
	assertEqual(t, ping.MeasurementName, "Ping", "error parsing ping field value for msm_name")
	assertEqual(t, ping.Type, "ping", "error parsing ping field value for type")
	assertEqual(t, ping.Step, uint(10), "error parsing ping field value for step")

	assertEqual(t, ping.Min, 4.75, "error parsing ping field value for min")
	assertEqual(t, ping.Avg, 12.0, "error parsing ping field value for avg")
	assertEqual(t, ping.Max, 25.0, "error parsing ping field value for max")

	replies, err := ping.Replies()
	if err != nil {
		t.Error(err)
	}
	rtts := make([]float64, 0)
	for _, reply := range replies {
		rtts = append(rtts, reply.RTT)
	}
	assertEqual(t, fmt.Sprint(rtts), "[10 15 4.75 5.25 25]", "error parsing RTTs")
	med, _ := ping.MedianRTT()
	medexp := 10.0
	assertEqual(t, med, medexp, fmt.Sprintf("median is incorrect, got %f, expected %f", med, medexp))

	// verify errors
	// verify timeouts
}

func TestMedian(t *testing.T) {
	var list []float64

	list = []float64{}
	_, err := median(list)
	if err == nil {
		t.Errorf("median calculation of emty list should return error")
	}

	list = []float64{10, 40, 30, 20}
	med2, _ := median(list)
	assertEqual(t, med2, 25.0, "error in median calculation for even list")

	list = []float64{10, 40, 20}
	med1, _ := median(list)
	assertEqual(t, med1, 20.0, "error in median calculation for odd list")
}

// almost one-liner to reduce boiler plate
func assertEqual(t *testing.T, val1 interface{}, val2 interface{}, msg string) {
	if val1 == val2 {
		return
	}
	t.Errorf("%s: received %v (type %v), expected %v (type %v)", msg, val1, reflect.TypeOf(val1), val2, reflect.TypeOf(val2))
}
