# goatapi - Go (RIPE) Atlas Tools - API Library

goatapi is a Go package to interact with [RIPE Atlas](https://atlas.ripe.net/) [APIs](https://atlas.ripe.net/api/v2/)
using [Golang](https://go.dev/). It is similar to [Cousteau](https://github.com/RIPE-NCC/ripe-atlas-cousteau).

It supports:
* finding probes
* finding anchors
* finding measurements
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

```
	filter := goatapi.NewProbeFilter()
	filter.FilterCountry("NL")
	filter.FilterPublic(true)
	count, err := filter.GetProbeCount(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(count)
```

### Search for Probes

```
	filter := goatapi.NewProbeFilter()
	filter.FilterCountry("NL")
	filter.Sort("-id")
	probes, err := filter.GetProbes(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	for _, probe := range probes {
		fmt.Println(probe.ShortString())
	}
```

### Get a Particular Probe

```
	filter := goatapi.NewProbeFilter()
	filter.FilterID(10001)
	probe, err := filter.GetProbe(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(probe.ShortString())
```


## Finding Anchors

### Count Anchors Matching Some Criteria

```
	filter := goatapi.NewAnchorFilter()
	filter.FilterCountry("NL")
	count, err := filter.GetAnchorCount(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(count)
```

### Search for Anchors

```
	filter := goatapi.NewAnchorFilter()
	filter.FilterCountry("NL")
	probes, err := filter.GetAnchors(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	for _, anchor := range anchors {
		fmt.Println(anchor.ShortString())
	}
```

### Get a Particular Anchor

```
	filter := goatapi.NewAnchorFilter()
	filter.FilterID(1)
	anchor, err := filter.GetAnchor(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(anchor.ShortString())
```

## Finding Measurements

### Count Measurements Matching Some Criteria

```
	filter := goatapi.NewMeasurementFilter()
	filter.FilterTarget(netip.ParsePrefix("193.0.0.0/19"))
	filter.FilterType("ping")
	filter.FilterOneoff(true)
	count, err := filter.GetMeasurementCount(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(count)
```

### Search for Measurements

```
	filter := goatapi.NewMeasurementFilter()
	filter.FilterTarget(netip.ParsePrefix("193.0.0.0/19"))
	filter.FilterType("ping")
	filter.FilterOneoff(true)
	filter.Sort("-id")
	msms, err := filter.GetMeasurements(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	for _, msm := range msms {
		fmt.Println(msm.ShortString())
	}
```

### Get a Particular Measurement

```
	filter := goatapi.NewMeasurementFilter()
	filter.FilterID(1001)
	msm, err := filter.GetMeasurement(false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(msm.ShortString())
```

# Future Additions / TODO

* schedule a new measurement, stop existing measurements
* modify participants of an existing measurement (add/remove probes)
* fetch results for, or listen to real-time result stream of, an already scheduled measurement
* check credit balance, transfer credits, ...

# Copyright, Contributing

(C) 2022, [Robert Kisteleki](https://kistel.eu/) & [RIPE NCC](https://www.ripe.net)

Contribution is possible and encouraged via the [Github repo](https://github.com/robert-kisteleki/goatapi/)

# License

See the LICENSE file
