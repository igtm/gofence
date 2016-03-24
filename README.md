# gofence

[![Build Status](https://travis-ci.org/buckhx/gofence.svg?branch=master)](https://travis-ci.org/buckhx/gofence)

Tool for geofencing with different algorithms for profiling. 

## Installation

```
// omitting the dot installs to /usr/local/bin
curl -sSL https://raw.githubusercontent.com/buckhx/gofence/master/scripts/install.py | python - .
```

## Usage

Invoking the fence cli will start an HTTP server and read geojson files from a directory into memory for searching.
The features in the geojson will be searchable at different endpoints for each file.
The endpoints will use a url-safe slug of the file name as their identifiers.

```
NAME:
   fence - Fence geojson point features

USAGE:
   cli [global options] command [command options] path/to/geojson/dir
   
VERSION:
   0.0.0
   
COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --fence "rtree"	Type of fence to use rtree|brute|qtree|qrtree|city|city-bbox|bbox
   --zoom, -z "18"	Some fences require a zoom level
   --port, -p "8080"	Port to bind to
   --help, -h		show help
   --version, -v	print the version
```

The city algorithms are special cases and both require NYC_BOROS_PATH envvar to be set to a geojson file. 
Don't use either of them for anything besides benchmarking.
The boros and tracts data can be found on [NYC Open Data Maps](http://www1.nyc.gov/site/planning/data-maps/open-data/districts-download-metadata.page)

## HTTP Methods

    GET /fence

A list of the available fence names

    POST /fence/:name/add

Adds a geojon feature from the post body to the fence at. 
This feature will not be saved to the server and will be gone if the server is restarted. 
Features in a fence have no notion of uniqueness, so if you add the same feature twice, the searchs will return both.

    POST /fence/:name/search

Search a fence for the query in the post body. This query must be a geojson feature with a point geometry. The properties of the matched features in the fence will be returned as a list.

    GET /fence/:name/search?lat=<LAT>&lon=<LON>

Convenience method for search with GET parameters. Both lat and lon and required and must be numbers. Any other parameters in the query string will be treated as properties of the query in the result.

## Micro Benchmarks

| Benchmark           | Operations | Time (ns/op) | Bytes (b/op) | Mallocs (allocs/op) | 
|---------------------|-----------:|-------------:|-------------:|--------------------:| 
| BenchmarkBrute-4    |       5000 |       251641 |           19 |                   1 | 
| BenchmarkCity-4     |      30000 |        39913 |           40 |                   1 | 
| BenchmarkBbox-4     |      50000 |        37085 |           11 |                   1 | 
| BenchmarkCityBbox-4 |     200000 |         9484 |           13 |                   1 | 
| BenchmarkQfence-4   |     300000 |         3959 |          399 |                  18 | 
| BenchmarkRfence-4   |    1000000 |         2290 |          174 |                   9 | 

![chart link broken](https://docs.google.com/spreadsheets/d/1PYoxb7nhPA_zrh9oPFnUH0mvo8geYvEkjfe8Jtc0vvY/pubchart?oid=1486005290&format=image)

Benchmarking requires NYC_TRACTS_PATH envvar to be set. Benchmarks are ran by checking which census tract a point is in [code here](lib/fence_test.go)

HTTP profiling is done via https://golang.org/pkg/net/http/pprof/

See https://github.com/buckhx/gofence-profiling for more in depth benmarking.
