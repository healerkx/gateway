
package main

import (
	"pathmap"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"strings"
)

type GatewayHandler struct {

}

func (this *GatewayHandler) serve(w http.ResponseWriter, req *http.Request) {
	abi := pathmap.GetApiBindingInfo(req.URL.Path)
	
	client := &http.Client{}
	if abi == nil {
		if req.URL.Path == "/favicon.ico" {
			return
		}
		fmt.Printf("[%s] not found abi", req.URL)
	}

	fmt.Printf("[%q]", abi)
	proxy, err := http.NewRequest("GET", abi.Url, strings.NewReader(req.URL.RawQuery))
    if err != nil {
        // handle error
    }
 
	resp, err := client.Do(proxy)
	if err != nil {
		
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		
	}
	
	fmt.Fprint(w, string(body))
}


func initialize(handler *GatewayHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.serve)
	return mux
}


func main() {
	pathmap.Initialize()
	handler := GatewayHandler{}
	mux := initialize(&handler)

    err := http.ListenAndServe(":8080", mux)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}