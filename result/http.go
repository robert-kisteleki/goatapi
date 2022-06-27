/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package result

import (
	"encoding/json"
	"fmt"
	"net/netip"
)

type HttpResult struct {
	BaseResult
	Uri          string         `json:"uri"`    //
	RawHttpReply []RawHttpReply `json:"result"` //
}

type RawHttpReply struct {
	AddressFamily   uint            `json:"af"`         //
	BodySize        uint            `json:"bsize"`      //
	DnsError        *string         `json:"dnserr"`     //
	DestinationAddr netip.Addr      `json:"dst_addr"`   //
	Error           *string         `json:"err"`        //
	Headers         *[]string       `json:"header"`     //
	HeaderSize      uint            `json:"hsize"`      //
	Method          string          `json:"method"`     //
	ReadTiming      *HttpReadTiming `json:"readtiming"` //
	ResultCode      uint            `json:"res"`        //
	ReplyTime       float64         `json:"rt"`         //
	SubID           *uint           `json:"subid"`      //
	SubMax          *uint           `json:"submax"`     //
	Time            *uniTime        `json:"time"`       //
	TimeToConnect   float64         `json:"ttc"`        //
	TimeToFirstByte float64         `json:"ttfb"`       //
	TimeToResolve   *float64        `json:"ttr"`        //
	Version         string          `json:"ver"`
}

type HttpReadTiming struct {
	Offset    uint `json:"o"` //
	TimeSince uint `json:"t"` //
}

func (result *HttpResult) String() string {
	ret := result.BaseString() +
		fmt.Sprintf("\t%s", result.Uri)
	return ret
}

func (result *HttpResult) DetailString() string {
	res := result.String() +
		fmt.Sprintf("\t%d", len(result.RawHttpReply))
	res += "\t["
	for _, answer := range result.RawHttpReply {
		if answer.Error != nil {
			res += *answer.Error
		} else {
			res += fmt.Sprintf("%s\t%d\t%d\t%d",
				answer.Method,
				answer.ResultCode,
				answer.HeaderSize,
				answer.BodySize,
			)
		}
	}
	res += "]"
	return res
}

func (result *HttpResult) TypeName() string {
	return "http"
}

func (http *HttpResult) Parse(from string) (err error) {
	err = json.Unmarshal([]byte(from), &http)
	if err != nil {
		return err
	}
	if http.Type != "http" {
		return fmt.Errorf("this is not a HTTP result (type=%s)", http.Type)
	}
	return nil
}
