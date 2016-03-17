# gofence

Tool for geofencing with different algorithms for profiling

```
NAME:
   fence - Fence geojson features from stdin

USAGE:
   fence [global options] command [command options] fence_file
   
VERSION:
   0.0.0
   
COMMANDS:
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --fence "rtree"		Type of fence to use rtree|brute|qtree|qrtree|city|city-bbox|bbox
   --concurrency, -c "4"	Concurrency factor, defaults to number of cores
   --zoom, -z "18"		Some fences require a zoom level
   --help, -h			show help
   --version, -v		print the version
   --debug			generate profiling data
```

The city algorithms are special cases and both require NYC_BOROS_PATH envvar to be set to a geojson file
