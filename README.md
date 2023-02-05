# goatapi - Go (RIPE) Atlas Tools - API Library

goatapi is a Go package to interact with [RIPE Atlas](https://atlas.ripe.net/) [APIs](https://atlas.ripe.net/api/v2/)
using [Golang](https://go.dev/). It is similar to [Cousteau](https://github.com/RIPE-NCC/ripe-atlas-cousteau) and
[Sagan](https://github.com/RIPE-NCC/ripe-atlas-sagan) combined.

It supports:
* finding probes
* finding anchors
* finding measurements
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

# Future Additions / TODO

* schedule a new measurement, stop existing measurements
* modify participants of an existing measurement (add/remove probes)
* check credit balance, transfer credits, ...

# Copyright, Contributing

(C) 2022, 2023 [Robert Kisteleki](https://kistel.eu/) & [RIPE NCC](https://www.ripe.net)

Contribution is possible and encouraged via the [Github repo](https://github.com/robert-kisteleki/goatapi/)

# License

See the LICENSE file
