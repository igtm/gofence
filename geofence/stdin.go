// +build ignore
// An example stdin implementation
package geofence

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/buckhx/diglet/geo"
)

func execute(in io.Reader, fence GeoFence, w int) *sync.WaitGroup {
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
				query, err := geo.GeojsonFeatureAdapter(gj)
				if err != nil {
					fmt.Println(err)
					continue
				}
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
