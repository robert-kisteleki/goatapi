# goatAPI changelog

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
