package main

import (
	"encoding/json"
	"net/http"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		if err := ValidateQuery(r); err != nil {
			log.Println(err)
			return
		}
		node := RequestNode(r)
	        fmt.Println(node)
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hello, %q", "HI")
	})

	resp, _ := GetRequest("http://router.project-osrm.org/route/v1/driving/13.388860,52.517037;13.397634,52.529407?overview=false")
	body, _ := ExtractBody(resp)
	var extNode OSMRResponse
	json.Unmarshal(body, &extNode)
	fmt.Println(extNode)
        // starting WWW server
        if err := http.ListenAndServe(":8080",nil); err != nil {
                log.Fatal("ListenAndServe:",err)
        }  
}

// Create Node from OSMRResponse (OSMR response), source and destination array.
// For proper work destination array order should correspond to routes order.
func ParseToNode(osmr OSMRResponse, src []string, dst []string) (node Node){
	node = Node{Source: src[0], Edge: make([]Edge, len(osmr.Routes))}
	for idx, val := range osmr.Routes {
		node.Edge[idx] = Edge{
			Destination: dst[idx], 
			Duration: val.Duration, 
			Distance: val.Distance,
		}
	}
	return node
}
