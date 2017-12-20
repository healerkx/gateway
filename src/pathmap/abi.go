package pathmap

import (
	"fmt"
	"strconv"
	"middleware"
	"net/http"
	"strings"
)

// Global
var gGetHeadPathMap = NewPathNode()

var gPostPutPathMap = NewPathNode()

func Handle(abi *ApiBindingInfo) {
	
}


type ApiBindingInfo struct {
	Url string;
	WarningLevel int32;
	LogLevel int32;
	CheckConfig int32;
	counter *middleware.Counter;
}

func NewApiBindingInfo(url string) *ApiBindingInfo {
	return &ApiBindingInfo{
		url, 0, 0, 0, middleware.NewCounter(),
	}
}

type PathNode struct {
	bindId int32
	subNode map[string]*PathNode
	abi *ApiBindingInfo
	pathParamName string
}

func NewUrlPathNode(url string, pathParamName string, bindId int32) *PathNode {
	return &PathNode{
		subNode: make(map[string]*PathNode),
		abi: NewApiBindingInfo(url),
		pathParamName: pathParamName,
		bindId: bindId,
	}
}

func NewPathNode() *PathNode {
	return &PathNode{
		subNode: make(map[string]*PathNode),
		abi: nil,
		pathParamName: "",
		bindId: 0,
	}
}

func (this *PathNode) Update(url string, pathParamName string, bindId int32) {
	this.abi = NewApiBindingInfo(url)
	this.bindId = bindId
	this.pathParamName = pathParamName
}

func GetPathNode(method, path string) (*PathNode, map[string]string) {
	parts := strings.Split(strings.Trim(path, "/ "), "/")
	currentNode := GetPathMap(method)
	
	pathParamMap := make(map[string]string)
	for _, part := range parts {
		if node, ok := currentNode.subNode[part]; ok {
			currentNode = node
		} else {
			if node, ok := currentNode.subNode["$"]; ok {
				if node.pathParamName != "" {
					pathParamMap[node.pathParamName] = part
				}
				currentNode = node
			} else {
				return nil, pathParamMap
			}
		}
	}
	return currentNode, pathParamMap
}

func GetUrl(pathNode *PathNode, pathParamMap map[string]string, query string) string {
	url := pathNode.abi.Url
	for key, value := range pathParamMap {
		url = strings.Replace(url, fmt.Sprintf("{{%s}}", key), value, -1)
	}
	return strings.TrimRight(url, "? ") + "?" + query
}

func GetApiBindingInfo(method, path, query string) *ApiBindingInfo {
	pathNode, pathParamMap := GetPathNode(method, path)
	if pathNode != nil {
		url := GetUrl(pathNode, pathParamMap, query)
		
		// TODO: ? how to add ?
		return NewApiBindingInfo(url)
	} else {
		return nil
	}
}

func GetPathMap(method string) *PathNode {
	if method == http.MethodGet || method == http.MethodHead {
		return gGetHeadPathMap
	} else if method == http.MethodPost || method == http.MethodPut {
		return gPostPutPathMap
	} else {
		return nil
	}
}

func addRoute(method string, apiBinding map[string]string) {
	path := apiBinding["gateway_api"]
	url := apiBinding["service_api"]
	bindId, _ := strconv.Atoi(apiBinding["bind_id"])
	parts := strings.Split(strings.Trim(path, "/ "), "/")
	
	count := len(parts)
	currentNode := GetPathMap(method)
	
	for index, part := range parts {
		pathParam := ""
		if strings.HasPrefix(part, "{{") && strings.HasSuffix(part, "}}") {
			pathParam = strings.Trim(part, "{} ")
			part = "$"
		}
		if node, ok := currentNode.subNode[part]; ok {
			// Update a Non-Leaf PathNode as a Bind PathNode 
			if index + 1 == count {
				node.Update(url, pathParam, int32(bindId))
			}
			currentNode = node
		} else {
			var newNode *PathNode
			if index + 1 == count {
				newNode = NewUrlPathNode(url, pathParam, int32(bindId))
			} else {
				newNode = NewPathNode()
			}
			currentNode.subNode[part] = newNode
			currentNode = newNode
		}
	}
}

func printRoutes(pathNode *PathNode, level int) {
	currentNode := pathNode
	for key, node := range currentNode.subNode {
		if node.abi != nil {
			fmt.Printf("|_%s[%s] (%d)%q\n", strings.Repeat("_", level * 4), key, node.bindId, node.abi.Url)
		} else {
			fmt.Printf("|_%s[%s]\n", strings.Repeat("_", level * 4), key)
		}
		
		printRoutes(node, level + 1)
	}
}

var gHttpMethodMap = map[int]string {
	1: 		http.MethodGet,
	2: 		http.MethodHead,
	4: 		http.MethodPost,
	8: 		http.MethodPut,
	16: 	http.MethodPatch,
	32: 	http.MethodDelete,
	64: 	http.MethodConnect,
	128: 	http.MethodOptions,
	256: 	http.MethodTrace,
}

// TODO:
func getHttpMethod(httpMethod string) []string {
	httpMethods, _ := strconv.Atoi(httpMethod)
	methods := []string{}
	for b, m := range gHttpMethodMap {
		if httpMethods & b == b {
			methods = append(methods, m)
		}
			
	}
	return methods
}

func addRoutes(apiBindings []map[string]string) {
	for _, apiBinding := range apiBindings {
		fmt.Printf("%+v\n", apiBinding)
		methods := getHttpMethod(apiBinding["http_method"])
		for _, method := range methods {
			addRoute(method, apiBinding)
		}
		
	}
}

func Initialize() bool {

	var apiBindings []map[string]string
	var err error
	if apiBindings, err = LoadApiBindingInfo(); err != nil {
		return false
	}

	addRoutes(apiBindings)

	// For test	
	// addRoute(http.MethodGet, "/api/thsamples", "http://127.0.0.1:9090/api/thsamples")
	// addRoute(http.MethodGet, "/api/thsample/{{id}}", "http://127.0.0.1:9090/api/thsample/{{id}}?a={{id}}")

	fmt.Printf("Path nodes for GET,HEAD\n")
	printRoutes(gGetHeadPathMap, 0)
	fmt.Printf("Path nodes for POST,PUT\n")
	printRoutes(gPostPutPathMap, 0)

	return true
}