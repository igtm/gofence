## HTTP geofencing service
Forked from https://github.com/buckhx/gofence

## Additional Changes from original
- remove implementations except rtree
- support go module
- add japan.geojson file from https://github.com/dataofjapan/land/blob/master/japan.geojson

```bash
# example usage
$ go build -o output
$ ./output --p "8899" ./

# call 
curl "http://localhost:8890/fence/japan/search?lat=35.6599017&lon=139.7169006"

{
    "query": {
        "type": "Feature",
        "properties": {},
        "geometry": {
            "type": "Point",
            "coordinates": [
                139.7169006,
                35.6599017
            ]
        }
    },
    "fences": [
        {
            "id": null,
            "nam": "Tokyo To",
            "nam_ja": "東京都"
        }
    ]
}
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
   --fence "rtree"	Type of fence to use rtree
   --zoom, -z "18"	Some fences require a zoom level
   --port, -p "8080"	Port to bind to
   --profile		Mounts profiling endpoints
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

Convenience method for search with GET parameters. Both lat and lon and required and must be numbers. Any other parameters in the query string will be treated as properties of the query in the result. Will be more performant since json unmarshalling isn't necessary.
