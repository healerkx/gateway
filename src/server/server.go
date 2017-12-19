package server

import (
	"pathmap"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	//"strings"
)

type GatewayHandler struct {
	
}

func (this *GatewayHandler) doGetHead(client *http.Client, req *http.Request, abi *pathmap.ApiBindingInfo) ([]byte, error) {
	proxy, err := http.NewRequest(req.Method, abi.Url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(proxy)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
	
func (this *GatewayHandler) doPostPut(client *http.Client, req *http.Request, abi *pathmap.ApiBindingInfo) ([]byte, error) {
	proxy, err := http.NewRequest(req.Method, abi.Url, req.Body)
	if err != nil {
		return nil, err
	}

	// TODO: Post body

	resp, err := client.Do(proxy)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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

	fmt.Printf("[%q]", abi)


	if req.Method == http.MethodGet || req.Method == http.MethodHead {
		if body, err := this.doGetHead(client, req, abi); err == nil {
			fmt.Fprint(w, string(body))
		}
	} else if req.Method == http.MethodPost || req.Method == http.MethodPut {
		if body, err := this.doPostPut(client, req, abi); err == nil {
			fmt.Fprint(w, string(body))
		}
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