package server

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
	
func (this *GatewayHandler) makeRequest(client *http.Client, req *http.Request, abi *pathmap.ApiBindingInfo) ([]byte, error) {
	proxy, err := http.NewRequest(req.Method, abi.Url, nil)
	if err != nil {
		return nil, err
	}
	
	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		bodyRead := ioutil.NopCloser(strings.NewReader(string(bodyBytes)))
		defer bodyRead.Close()
		proxy.Body = bodyRead
	}

	proxy.Header = req.Header
	fmt.Printf("%s [%s]\n", req.Method, abi.Url)
	resp, err := client.Do(proxy)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}
	

func (this *GatewayHandler) serve(w http.ResponseWriter, req *http.Request) {
	
	abi := pathmap.GetApiBindingInfo(req.Method, req.URL.Path, req.URL.RawQuery)
	
	client := &http.Client{}
	if abi == nil {
		if req.URL.Path == "/favicon.ico" {
			return
		}
		fmt.Printf("[%s] not found abi", req.URL)
	}
	pathmap.Handle(abi)

	if body, err := this.makeRequest(client, req, abi); err == nil {
		fmt.Fprint(w, string(body))
	} else {
		fmt.Fprint(w, "")
	}	
}

func initialize(handler *GatewayHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.serve)

	// TODO: Other routes
	return mux
}

func Run() {
	
	pathmap.Initialize()
	handler := GatewayHandler{}
	mux := initialize(&handler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}