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
	ResponseTime    float64 `json:"rt"`      //
	ResponseSize    uint    `json:"size"`    //
	Abuf            string  `json:"abuf"`    //
	ID              uint    `json:"id"`      //
	AnswerCount     uint    `json:"ancount"` //
	QueriesCount    uint    `json:"qdcount"` //
	NameServerCount uint    `json:"nscount"` //
	AdditionalCount uint    `json:"arcount"` //
	//Detail         *[]DnsDetail `json:"answers"` //
	//TTL6            *uint        `json:"ttl"`     //
}

func (answer *DnsAnswer) ShortString() string {
	ret := fmt.Sprintf("%d\t%d\t%d\t%d",
		answer.AnswerCount,
		answer.QueriesCount,
		answer.NameServerCount,
		answer.AdditionalCount,
	)
	return ret
}

func (answer *DnsAnswer) LongString() string {
	ret := answer.ShortString() + fmt.Sprintf("\t%s", answer.Abuf)
	return ret
}

type DnsDetail struct {
	DomainName string `json:"mname"`
	Name       string `json:"name"`
	RData      string `json:"rdata"`
	RName      string `json:"rname"`
	Serial     uint   `json:"serial"`
	TTL        uint   `json:"ttl"`
	Type       string `json:"type"`
}

type DnsError struct {
	Timeout  uint   `json:"timeout"`
	AddrInfo string `json:"getaddrinfo"`
}

func (result *DnsResult) ShortString() string {
	return result.BaseShortString()
}

func (result *DnsResult) LongString() string {
	return result.ShortString() +
		"\t" + fmt.Sprint(result.RawResult) +
		"\t" + fmt.Sprint(result.RawResultSet)
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
