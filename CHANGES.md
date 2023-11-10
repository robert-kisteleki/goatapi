# goatAPI changelog

## next

* FIX: traceroute hop details can have 'error' instad of actual data
* CHANGED: firmware version 1 results are not supported
* FIX: improved handling of some old results which encode firmware as string
* CHANGED: update to go 1.21 and module updates to newest

## 0.5.0

* CHANGED: renamed `Start()` and `Stop()` to `StartTime()` and `EndTime()`
* CHANGED: renamed `Submit()` to `Schedule()`
* NEW: support for stopping measurements via `MeasurementSpec.Stop()`
* NEW: support for adding or removing probes to/from existing measurements via `MeasurementSpec.ParticipationRequest()`

## 0.4.2

* no changes, adminstrative release only

## 0.4.1

* FIX: handle dnserrors in TLS results
* FIX: traceroute numerical errors in hops were not handled properly
* CHANGED: fix typos
* CHANGED: improve stream EOF handling

## 0.4.0

* FIX: do not blow up if connection to stream fails
* FIX: defaulted to file read if stream was turned on but no measurement ID was set
* FIX: (ping parser) blew up if the generic Rtt field was missing
* FIX: (ping parser) min/avg/max were not reported correctly if they were not present but could otherwise be calculated
* FIX: typo in ErrorDetail JSON "detail"
* NEW: support for measurement status checks
* CHANGED: verbosity is now a setting, not a parameter
* CHANGED: ErrorDetail can have embedded error messages
* NEW: support for measurement scheduling with probe selection, timing, all measurement types and options

## 0.3.0

* NEW: add support to stream results from the result stream instead of the data API
* NEW: add support to save all results obtained (from data API, stream API or even
  from a file) to a file
* CHANGED: change default limit on downloading results to 0 (no limit)

## v0.2.1

* CHANGED: better error handling for non-200 responses
* CHANGED: rdata is exposed in DNS results
* CHANGED: API GET calls are not handles in one function
* CHANGED: GetMeasurement() can also use an API key
* CHANGED: all responses (probes, acnhors, ...) are async/channel based now

## v0.2.0

* NEW: support for result processing (API call or local file)
  * including downloading "latest" results
  * all result types are supported (need a lot of work for testig / older resuls though)

## v0.1.0

* support listing probes, anchors, measurements with virtually all filtering options
* support counting items, retrieveing all matching ones or just a specific one by ID
* support for "list_measurements" API key
