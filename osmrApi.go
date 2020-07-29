package main

import (
	"net/http"
	"errors"
	"strings"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sort"
)

// Validate request query under OSMR restrictions
func ValidateQuery (r *http.Request) error {
	query := r.URL.Query()
	src, dst := query["src"], query["dst"]
	if len(src) == 0 || len(dst) == 0 {
		return  errors.New("Not enough arguments specified")
	} else if len(src) > 1 {
		return  errors.New("Only one source is available")
	}
	return nil
}

// Make call to OSMR API and return Node struct
func RequestNode (r *http.Request) *Node {
        query := r.URL.Query()
        src := query["src"]
        dst := query["dst"]
        u, _ := url.Parse("http://router.project-osrm.org/route/v1/driving/")
        path := append(src, dst...)
        s := strings.Join(path, ";")
        rel, _ := u.Parse(s)
        log.Println(rel.String())
	resp, _ := http.PostForm(rel.String(), url.Values{"overview":{"false"}})
        body, _ := ExtractBody(resp)
        var osmr OSMRResponse
        json.Unmarshal(body, &osmr)
        fmt.Println(osmr)
        node := ParseToNode(osmr, src, dst)
	node.sortEdges()
	return &node
}

// Sort egdes in node from smallest distance to highest
func (node *Node) sortEdges(){
	sort.Sort(ByTimeAndDistance(node.Edge)) 
}

type ByTimeAndDistance []Edge

func (a ByTimeAndDistance) Len() int { return len(a) }

func (a ByTimeAndDistance) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByTimeAndDistance) Less(i, j int) bool {
	if a[i].Duration == a[j].Duration {
		return a[i].Distance < a[j].Distance
	} else { 
		return a[i].Duration < a[j].Duration
	}
}


// Struct used to marshal data for "/routes" response
type Node struct {
	Source string   `json:"source"`
	Edge []Edge `json:"routes"`
}
type Edge struct {
	Destination string  `json:"destination"`
	Duration    float64 `json:"duration"`
	Distance    float64 `json:"distance"`
}


// Struct used to unmarshal json from OSMR
type OSMRResponse struct {
	Routes    []Routes    `json:"routes"`
	Waypoints []Waypoints `json:"waypoints"`
	Code      string      `json:"code"`
}
type Legs struct {
	Summary  string        `json:"summary"`
	Weight   int           `json:"weight"`
	Duration float64       `json:"duration"`
	Steps    []interface{} `json:"steps"`
	Distance float64       `json:"distance"`
}
type Routes struct {
	Legs       []Legs  `json:"legs"`
	WeightName string  `json:"weight_name"`
	Weight     int     `json:"weight"`
	Duration   float64 `json:"duration"`
	Distance   float64 `json:"distance"`
}
type Waypoints struct {
	Hint     string    `json:"hint"`
	Name     string    `json:"name"`
	Location []float64 `json:"location"`
}
