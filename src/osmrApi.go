package main

import (
	"net/http"
	"errors"
	"strings"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
	"net/url"
	"sort"
)

const (
	OSMR_ROUTE_PATH string = "http://router.project-osrm.org/route/v1/driving/"
)

var OSMRRouteDefaultArgs = url.Values{
	"overview":{"false"},
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

/* Create Node with sorted Edges from request */
func RequestNode (r *http.Request) (*Node) {
	osmr, err := RequestOSMRRoute(r)
	if err != nil { 
		log.Println(err)
		return nil 
	}
        src, dst := GetSourceAndDestinations(r)
        node := ParseToNode(osmr, src, dst)
	node.SortEdges()
	return &node
}


func ExtractBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

/* Validate request query under OSMR restrictions:
- only one source
- multiple destinations
- source and destination cannot be empty
*/
func ValidateQuery (r *http.Request) error {
	src, dst := GetSourceAndDestinations(r)
	if len(src) == 0 || len(dst) == 0 {
		return  errors.New("Not enough arguments specified")
	} else if len(src) > 1 {
		return  errors.New("Only one source point can be provided")
	}
	return nil
}

func GetSourceAndDestinations(r *http.Request) ([]string,[]string){
        query := r.URL.Query()
        return query["src"], query["dst"]
}

/* Create OSMRRResponse struct from request. If query parameters
were incorrectly provided returns erros with code from OSMR API
*/
func RequestOSMRRoute(r *http.Request) (*OSMRResponse, error) {
        src, dst := GetSourceAndDestinations(r)
	rel := FormatOSMRRouteQuery(src, dst)
        log.Println(rel.String())

	resp, err := http.PostForm(rel.String(), OSMRRouteDefaultArgs)
	if err != nil { return nil, err	}

        body, err := ExtractBody(resp)
	if err != nil {	return nil, err	}

        var osmr OSMRResponse
        json.Unmarshal(body, &osmr)
	if osmr.Code != "Ok" {
		return nil, fmt.Errorf("OSMR API returned code: %s", osmr.Code)
	}
	return &osmr, nil
}

/* From src and dst arrays create URL correct for route request in OSMR API */
func FormatOSMRRouteQuery(src []string, dst []string) *url.URL {
        u, _ := url.Parse(OSMR_ROUTE_PATH)
        path := append(src, dst...)
        s := strings.Join(path, ";")
        rel, _ := u.Parse(s)
	return rel
}

/* Create Node from OSMRResponse (OSMR response), source and destination array.
For proper work destination array order should correspond to routes order.
*/
func ParseToNode(osmr *OSMRResponse, src []string, dst []string) (node Node){
	node = Node{Source: src[0], Edge: make([]Edge, len(osmr.Routes[0].Legs))}
	for idx, val := range osmr.Routes[0].Legs {
		node.Edge[idx] = Edge{
			Destination: dst[idx], 
			Duration: val.Duration, 
			Distance: val.Distance,
		}
	}
	return node
}

/* Sort egdes in node from smallest distance to highest in bubble sort algorithm. */
func (node *Node) SortEdges(){
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
