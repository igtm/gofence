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

## Benchmarks

| Benchmark           |   Ops    |       Time       |     Bytes     |      Mallocs       | 
|---------------------|----------|------------------|---------------|--------------------| 
| BenchmarkBrute-4    |     5000 |     251641 ns/op |       19 B/op |        1 allocs/op | 
| BenchmarkCity-4     |    30000 |      39913 ns/op |       40 B/op |        1 allocs/op | 
| BenchmarkBbox-4     |    50000 |      37085 ns/op |       11 B/op |        1 allocs/op | 
| BenchmarkCityBbox-4 |   200000 |       9484 ns/op |       13 B/op |        1 allocs/op | 
| BenchmarkQfence-4   |   300000 |       3959 ns/op |      399 B/op |       18 allocs/op | 
| BenchmarkRfence-4   |  1000000 |       2290 ns/op |      174 B/op |        9 allocs/op | 
