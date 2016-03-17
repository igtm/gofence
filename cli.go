package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/buckhx/diglet/geo"
	"github.com/buckhx/gofence/lib"
	"github.com/codegangsta/cli"
	"github.com/davecheney/profile"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

func client(args []string) {
	app := cli.NewApp()
	app.Name = "fence"
	app.Usage = "Fence geojson features from stdin"
	app.ArgsUsage = "fence.geojson"
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
	}
	app.Action = func(c *cli.Context) {
		args := c.Args()
		if len(args) < 1 || args[0] == "" {
			die(c, "fence_file required")
		}
		w := c.Int("concurrency")
		z := c.Int("zoom")
		if w < 1 || z < 0 || z > 23 {
			die(c, "-c must be > 0 && 0 <= -z <= 23")
		}
		file := args[0]
		label := c.String("fence")
		fence, err := load(file, label, z)
		if err != nil {
			die(c, err.Error())
		}
		working := execute(os.Stdin, fence, w)
		working.Wait()
	}
	app.Run(args)
}

func execute(in io.Reader, fence geofence.GeoFence, w int) *sync.WaitGroup {
	lines := make(chan string, 1<<10)
	go func() {
		defer close(lines)
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
	}()
	working := &sync.WaitGroup{}
	for i := 0; i < w; i++ {
		working.Add(1)
		go func() {
			defer working.Done()
			for line := range lines {
				gj, err := geo.UnmarshalGeojsonFeature(line)
				if err != nil {
					fmt.Println(err)
					continue
				}
				query := geo.GeojsonFeatureAdapter(gj)
				matchs := fence.Get(query.Geometry[0].Head()) // it's a point
				fences := make([]map[string]interface{}, len(matchs))
				for i, match := range matchs {
					fences[i] = match.Properties
				}
				query.Properties["fences"] = fences
				res, err := json.Marshal(query)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("%s\n", res)
			}
		}()
	}
	return working
}

func load(fenceFile, fenceType string, zoom int) (fence geofence.GeoFence, err error) {
	fence, err = geofence.GetFence(fenceType, zoom)
	if err != nil {
		return
	}
	source := geo.NewGeojsonSource(fenceFile, nil)
	features, _ := source.Publish()
	for feature := range features {
		fence.Add(feature)
	}
	return
}

func main() {
	for _, arg := range os.Args {
		if arg == "--debug" {
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
