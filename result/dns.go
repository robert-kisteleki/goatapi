/*
  (C) 2022 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package result

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/netip"
	"strings"
)

type DnsResult struct {
	BaseResult
	RawResult    *DnsAnswer `json:"result"`    //
	RawResultSet []DnsReply `json:"resultset"` //
}

type DnsReply struct {
	Time            uniTime    `json:"time"`     //
	LastTimeSync    uint       `json:"lts"`      //
	SourceAddr      netip.Addr `json:"src_addr"` //
	DestinationAddr netip.Addr `json:"dst_addr"` //
	DestinationPort string     `json:"dst_port"` // (should be uint)
	AddressFamily   uint       `json:"af"`       //
	Protocol        string     `json:"proto"`    //
	Error           *DnsError  `json:"error"`    //
	RetryCount      *uint      `json:"retry"`    //
	SubID           uint       `json:"subid"`    //
	SubMax          uint       `json:"submax"`   //
	RawQBuf         *string    `json:"qbuf"`     //
	Answer          DnsAnswer  `json:"result"`   //
}

type DnsAnswer struct {
	ResponseTime    float64      `json:"rt"`      //
	ResponseSize    uint         `json:"size"`    //
	Abuf            string       `json:"abuf"`    //
	ID              uint         `json:"id"`      //
	AnswerCount     uint         `json:"ancount"` //
	QueriesCount    uint         `json:"qdcount"` //
	NameServerCount uint         `json:"nscount"` //
	AdditionalCount uint         `json:"arcount"` //
	Details         *[]DnsDetail `json:"answers"` //
	//TTL6            *uint        `json:"ttl"`     //
}

func (answer *DnsAnswer) String() string {
	ret := fmt.Sprintf("%d\t%d\t%d\t%d",
		answer.AnswerCount,
		answer.QueriesCount,
		answer.NameServerCount,
		answer.AdditionalCount,
	)
	return ret
}

func (answer *DnsAnswer) DetailString() string {
	res := answer.String() + "\t["
	if answer.Details != nil {
		s := make([]string, 0)
		for _, detail := range *answer.Details {
			s = append(s, detail.DetailString())
		}
		res += strings.Join(s, " ")
	}
	res += "]"
	return res
}

type DnsDetail struct {
	DomainName string `json:"mname"`
	Name       string `json:"name"`
	RData      string `json:"rdata"`
	RName      string `json:"rname"`
	Serial     uint   `json:"serial"`
	Ttl        uint   `json:"ttl"`
	Type       string `json:"type"`
}

func (detail *DnsDetail) DetailString() string {
	return fmt.Sprintf("%s %s %d", detail.Type, detail.RName, detail.Serial)
}

type DnsError struct {
	Timeout  uint   `json:"timeout"`
	AddrInfo string `json:"getaddrinfo"`
}

func (result *DnsResult) String() string {
	return result.BaseString() +
		"\t" + fmt.Sprint(len(result.Answers()))
}

func (result *DnsResult) DetailString() string {
	res := result.String() + "\t["
	s := make([]string, 0)
	for _, ans := range result.Answers() {
		s = append(s, ans.DetailString())
	}
	res += strings.Join(s, " ")
	res += "]"
	return res
}

func (result *DnsResult) TypeName() string {
	return "dns"
}

func (dns *DnsResult) Parse(from string) (err error) {
	err = json.Unmarshal([]byte(from), &dns)
	if err != nil {
		return err
	}
	if dns.Type != "dns" {
		return fmt.Errorf("this is not a DNS result (type=%s)", dns.Type)
	}
	return nil
}

func (result *DnsReply) QBuf() ([]byte, error) {
	if result.RawQBuf == nil {
		return nil, fmt.Errorf("qbuf is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(*result.RawQBuf)
	if err != nil {
		return nil, fmt.Errorf("error decoding qbuf:  %s", err.Error())
	}

	return decoded, nil
}

func (result *DnsResult) Answers() (answers []DnsAnswer) {
	if result.RawResult != nil {
		answers = append(answers, *result.RawResult)
	}
	if result.RawResultSet != nil {
		for _, rs := range result.RawResultSet {
			answers = append(answers, rs.Answer)
		}
	}
	return answers
}
