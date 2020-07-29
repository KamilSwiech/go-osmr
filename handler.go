package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		src := query["src"]
		dst := query["dst"]
		if len(src) == 0 || len(dst) == 0 {
			log.Println("No arguments")
			return
		}
		u, _ := url.Parse("http://router.project-osrm.org/route/v1/driving/")
		path := append(src, dst...)
		s := strings.Join(path, ";")
		rel, _ := u.Parse(s)
		log.Println(rel.String())
		resp, _ := http.PostForm(rel.String(), url.Values{"overview":{"false"}})
		body, _ := ExtractBody(resp)
		log.Println(string(body))
		w.WriteHeader(200)
		fmt.Fprintf(w, "Hello, %q", src)
	})

	resp, _ := GetRequest("http://router.project-osrm.org/route/v1/driving/13.388860,52.517037;13.397634,52.529407?overview=false")
	body, _ := ExtractBody(resp)
	var node Node
	json.Unmarshal(body, &node)
	fmt.Println(node)
        // starting WWW server
        if err := http.ListenAndServe(":8080",nil); err != nil {
                log.Fatal("ListenAndServe:",err)
        }  
}

func ParseToNode(ext ExtendedNode) (node Node){
	
}
