package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/buckhx/diglet/geo"
	"github.com/buckhx/gofence/lib"
	"github.com/codegangsta/cli"
	"github.com/davecheney/profile"
)

func client(args []string) {
	app := cli.NewApp()
	app.Name = "fence"
	app.Usage = "Fence geojson features from stdin"
	app.ArgsUsage = "Path to directory with geojson to be loaded into fences"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fence",
			Value: "rtree",
			Usage: "Type of fence to use " + strings.Join(geofence.FenceLabels, "|"),
		},
		cli.IntFlag{
			Name:  "concurrency, c",
			Value: runtime.GOMAXPROCS(0),
			Usage: "Concurrency factor, defaults to number of cores",
		},
		cli.IntFlag{
			Name:  "zoom, z",
			Value: 18,
			Usage: "Some fences require a zoom level",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "8080",
			Usage: "Port to bind to",
		},
		cli.BoolFlag{
			Name:  "profile",
			Usage: "Profiles execution via pprof",
		},
	}
	app.Action = func(c *cli.Context) {
		args := c.Args()
		if len(args) < 1 || args[0] == "" {
			die(c, "fences_path required")
		}
		w := c.Int("concurrency")
		z := c.Int("zoom")
		if w < 1 || z < 0 || z > 23 {
			die(c, "-c must be > 0 && 0 <= -z <= 23")
		}
		path := args[0]
		label := c.String("fence")
		fences, err := load(path, label, z)
		if err != nil {
			die(c, err.Error())
		}
		port := fmt.Sprintf(":%s", c.String("port"))
		err = geofence.ListenAndServe(port, fences)
		die(c, err.Error())
	}
	app.Run(args)
}

func load(dir, fenceType string, zoom int) (fences geofence.FenceIndex, err error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*json")) // .geo.json/.geojson/.json
	if err != nil {
		return
	}
	fences = geofence.NewFenceIndex()
	for _, path := range paths {
		fmt.Printf("Loading fence %s\n", path)
		fence, err := geofence.GetFence(fenceType, zoom)
		if err != nil {
			fmt.Printf("Error building fence for %s, skipping...", path)
			continue
		}
		source := geo.NewGeojsonSource(path, nil) //panics on invalid json file
		features, err := source.Publish()
		if err != nil {
			return nil, err
		}
		i := 0
		for feature := range features {
			i++
			fmt.Printf("Loading feature %d\n", i)
			if feature.Type == "Point" {
				continue // points don't have containment area
			}
			fence.Add(feature)
		}
		key := slug(path)
		fences.Set(key, fence)
	}
	return
}

func main() {
	for _, arg := range os.Args {
		//wrap all execution
		if arg == "--profile" {
			config := &profile.Config{
				MemProfile: true,
				CPUProfile: true,
			}
			defer profile.Start(config).Stop()
		}
	}
	client(os.Args)
}

func die(c *cli.Context, msg string) {
	cli.ShowAppHelp(c)
	fmt.Println(msg)
	os.Exit(1)
}

var slugger = regexp.MustCompile("[^a-z0-9]+")

// Slugs the basename of the path, removing the path and extension
// "/path/to/file_2.gz " -> "file-2"
// yoinked from diglet/util
func slug(path string) string {
	s := filepath.Base(path)
	s = strings.TrimSuffix(s, filepath.Ext(s))
	return slugged(s, "-")
}

func slugged(s, delim string) string {
	return strings.Trim(slugger.ReplaceAllString(strings.ToLower(s), delim), delim)
}
