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
		json, err := json.Marshal(node)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "%s", json)
	})

        if err := http.ListenAndServe(":8080",nil); err != nil {
                log.Fatal("ListenAndServe:",err)
        }  
}
