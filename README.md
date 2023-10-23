# goatapi - Go (RIPE) Atlas Tools - API Library

goatapi is a Go package to interact with [RIPE Atlas](https://atlas.ripe.net/) [APIs](https://atlas.ripe.net/api/v2/)
using [Golang](https://go.dev/). It is similar to [Cousteau](https://github.com/RIPE-NCC/ripe-atlas-cousteau) and
[Sagan](https://github.com/RIPE-NCC/ripe-atlas-sagan) combined.

It supports:
* finding probes
* finding anchors
* finding measurements
* scheduling new measurements
* downloading results of measurements and turning them into Go objects
* tuning in to result streaming and turning them into Go objects
* loading a local file containing measurement results and turning them into Go objects
* (more features to come)

The tool needs Go 1.18 to compile.

# Context

[RIPE Atlas](https://atlas.ripe.net) is an open, community based active Internet
measurement network developed by the [RIPE NCC](https://www.ripe.net/) since 2010.
It provides a number of vantage points ("probes") run by volunteers, that allow
various kinds of network measurements (pings, traceroutes, DNS queries, ...) to
be run by any user.


# Quick Start

## Finding Probes

### Count Probes Matching Some Criteria

```go
	filter := goatapi.NewProbeFilter()
	filter.FilterCountry("NL")
	filter.FilterPublic(true)
	count, err := filter.GetProbeCount(false)
	if err != nil {
		// handle the error
	}
	fmt.Println(count)
```

### Search for Probes

```go
	filter := goatapi.NewProbeFilter()
	filter.FilterCountry("NL")
	filter.Sort("-id")
	probes := make(chan goatapi.AsyncProbeResult)
	go filter.GetProbes(false, probes) // false means non-verbose
	if err != nil {
		// handle the error
	}
	for probe := range probes {
		if probe.Error != nil {
			// handle the error
		} else {
			// process the result
		}
	}
```

### Get a Particular Probe

```go
	probe, err := filter.GetProbe(false, 10001) // false means non-verbose
	if err != nil {
		// handle the error
	}
	fmt.Println(probe.ShortString())
```


## Finding Anchors

### Count Anchors Matching Some Criteria

```go
	filter := goatapi.NewAnchorFilter()
	filter.FilterCountry("NL")
	count, err := filter.GetAnchorCount(false)
	if err != nil {
		// handle the error
	}
	fmt.Println(count)
```

### Search for Anchors

```go
	filter := goatapi.NewAnchorFilter()
	filter.FilterCountry("NL")
	anchors := make(chan goatapi.AsyncAnchorResult)
	go filter.GetAnchors(false, anchors) // false means non-verbose
	if err != nil {
		// handle the error
	}
	for anchor := range anchors {
		if anchor.Error != nil {
			// handle the error
		} else {
			// process the result
		}
	}
```

### Get a Particular Anchor

```go
	anchor, err := filter.GetAnchor(false, 1080) // false means non-verbose
	if err != nil {
		// handle the error
	}
	fmt.Println(anchor.ShortString())
```

## Finding Measurements

### Count Measurements Matching Some Criteria

```go
	filter := goatapi.NewMeasurementFilter()
	filter.FilterTarget(netip.ParsePrefix("193.0.0.0/19"))
	filter.FilterType("ping")
	filter.FilterOneoff(true)
	count, err := filter.GetMeasurementCount(false)
	if err != nil {
		// handle the error
	}
	fmt.Println(count)
```

### Search for Measurements

```go
	filter := goatapi.NewMeasurementFilter()
	filter.FilterTarget(netip.ParsePrefix("193.0.0.0/19"))
	filter.FilterType("ping")
	filter.FilterOneoff(true)
	filter.Sort("-id")
	msms := make(chan goatapi.AsyncMeasurementResult)
	go filter.GetMeasurements(false, msms) // false means non-verbose
	if err != nil {
		// handle the error
	}
	for msm := range msms {
		if msm.Error != nil {
			// handle the error
		} else {
			// process the result
		}
	}
```

### Get a Particular Measurement

```go
	msm, err := filter.GetMeasurement(false, 1001)
	if err != nil {
		// handle the error
	}
	fmt.Println(msm.ShortString())
```

## Processing results

All result types are defined as object types (PingResult, TracerouteResult, DnsResult, ...). The Go types try to be more useful than what the API natively provides, i.e. there's a translation from what the API gives to objects that have more meaning and simpler to understand fields and methods.

Results can be fetched via a channel. The filter support measurement ID, start/end time, probe IDs, "latest" results and combinations of these.

**Note: this is alpha level code. Some fields are probably mis-interpreted, old result versions are not processed correctly, a lot of corner cases are likely not handled properly and some objects should be passed as pointers instead. There's possibly a lot of work to be done here.**

An example of retrieving and processing results from the data API:

```go
	filter := goatapi.NewResultsFilter()
	filter.FilterID(10001)
	filter.FilterLatest()

	results := make(chan result.AsyncResult)
	go filter.GetResults(false, results) // false means verbose mode is off

	for result := range results {
		// do something with a result
	}
```

An example of retrieving and processing results from result streaming:

```go
	filter := goatapi.NewResultsFilter()
	filter.FilterID(10001)
	filter.Stream(true)

	results := make(chan result.AsyncResult)
	go filter.GetResults(false, results) // false means verbose mode is off

	for result := range results {
		// do something with a result
	}
```

An example of retrieving and processing results from a file:

```go
	filter := goatapi.NewResultsFilter()
	filter.FilterFile("-") // stdin, one can also use a proper file name
	// note: other filters can be added (namely start, stop and probe)

	results := make(chan result.AsyncResult)
	go filter.GetResults(false, results) // false means verbose mode is off

	for result := range results {
		// do something with a result
	}
```

## Result types

The `result` package contains various types to hold corresponding measurement result types:
* `BaseResult` is the basis of all and contains the basic fields such as `MeasurementID`, `ProbeId`, `TimeStamp`, `Type` and such
* `PingResult`, `TracerouteResult`, `DnsResult` etc. contain the type-specific fields

## Measurement Scheduling

You can schedule measuements with virtually all available API options. A quick example:

```go
	spec := goatapi.NewMeasurementSpec()
	spec.ApiKey(myapikey)

	include := []string{"system-v4"}
	exclude := []string{"badtag", "anotherbad"}
	spec.AddProbesAreaWithTags("ww", 10, &include, &exclude)
	spec.AddProbesList([]uint{1, 99, 999})
	spec.AddProbesCountry("NL", 15)
	spec.AddProbesPrefix(netip.MustParsePrefix("192.0.0.0/8"), 5)

	spec.Start(tomorrownoon)
	spec.OneOff(true)

	spec.AddTrace(
		"my traceroute measurement",
		"ping.ripe.net",
		4,
		&goatapi.BaseOptions{ResolveOnProbe: true},
		&goatapi.TraceOptions{FirstHop: 4, ParisId: 9},
	)

	msmid, err := spec.Submit()
	if err != nil {
		// use msmid
	}
```

### Basics

A new measuement object can be created with `NewMeasurementSpec()`. In order to successfully submit this to the API, you need to add an API key using `ApiKey()`. It also needs to contain at least one probe definition and at least one measurement definition.

### Probe Definitions

All variations of probe selection are supported:
* `AddProbesArea()` and `AddProbesAreaWithTags()` to add probes from an area (_WW_, _West_, ...)
* `AddProbesCountry()` and `AddProbesCountryWithTags()` to add probes from an country specified by its country code (ISO 3166-1 alpha-2)
* `AddProbesReuse()` and `AddProbesReuseWithTags()` to reuse probes from a previous measurement
* `AddProbesAsn()` and `AddProbesAsnWithTags()` to add probes from an ASN
* `AddProbesPrefix()` and `AddProbesPrefixWithTags()` to add probes from a prefix (IPv4 or IPv6)
* `AddProbesList()` and `AddProbesListWithTags()` to add probes with explicit probe IDs

Probe tags can be specified to include or exclude ones that have those specific tags.

### Time Definitions

You can specify whether you want a one-off or an ongoing measurement using `Oneoff()`.

Each measurement can have an explicit start time defined with `Start()`. Ongoing measurements can also have a predefined stop time with `Stop()`. These have to be sane regarding the current time (they need to be in the future) and to each other (stop needs to happen after start). By default start time is as soon as possible with an undefined end time.

### Measurement Definitions

Various measurements can be added with `AddPing()`, `AddTrace()`, `AddDns()`, `AddTls()`, `AddNtp()` and `AddHttp()`. Multiple measurements can be added to one specification; in this case they will use the same probes and timing.

All measurement types support setting common options using a `BaseOptions{}` structure. You can set the measurement interval, spread, the resolve-on-probe flag and so on here. If you are ok with the API defaults then you can leave this parameter to be `nil`.

All measurement types also accept type-specific options via the structures `PingOptions{}`, `TraceOptions{}`, `DnsOptions{}` and so on. If you are ok with the API defaults then you can leave this parameter to be `nil` as well.

### Submitting a Measurement Specification to the API

The `Submit()` function POSTs the whole specifiaton to the API. It either returns with an `error` or a list of recently created measurement IDs. In case you're only intrested in the API-compatible JSON structure without submitting it, then `GetApiJson()` should be called instead.

# Future Additions / TODO

* stop existing measurements
* modify participants of an existing measurement (add/remove probes)
* check credit balance, transfer credits, ...

# Copyright, Contributing

(C) 2022, 2023 [Robert Kisteleki](https://kistel.eu/) & [RIPE NCC](https://www.ripe.net)

Contribution is possible and encouraged via the [Github repo](https://github.com/robert-kisteleki/goatapi/)

# License

See the LICENSE file
