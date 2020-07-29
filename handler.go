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
		w.WriteHeader(200)
		b, _ := json.Marshal(node)
		fmt.Fprintf(w, "%s", b)
	})

        if err := http.ListenAndServe(":8080",nil); err != nil {
                log.Fatal("ListenAndServe:",err)
        }  
}
