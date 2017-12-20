package pathmap

import (
	"errors"
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

const (
	AbiStatusUnknown	= 0
	AbiStatusEnabled	= 1
	AbiStatusDisabled	= 2
	AbiStatusDying		= 3
) 



type ApiBindingInfo struct {
	Url string;
	WarningLevel int32;
	LogLevel int32;
	CheckConfig int32;
	Status int32;
	counter *middleware.Counter;
}

var gApiBindingInfoMap = map[int32]*ApiBindingInfo {}


func NewApiBindingInfo(url string) *ApiBindingInfo {
	return &ApiBindingInfo{
		Url: url, 
		WarningLevel: 0, 
		LogLevel: 0, 
		CheckConfig: 0, 
		Status: AbiStatusEnabled,
		counter: middleware.NewCounter(),
	}
}

type PathNodeMap map[string]*PathNode

type PathNode struct {
	bindId 			int32
	pathNodeMap 	PathNodeMap
	pathParamName 	string
}

func (this *PathNode) FindApiBindingInfo() *ApiBindingInfo {
	return gApiBindingInfoMap[this.bindId]
}

func NewUrlPathNode(url string, pathParamName string, bindId int32) *PathNode {
	abi := NewApiBindingInfo(url)
	gApiBindingInfoMap[bindId] = abi
	return &PathNode {
		pathNodeMap: make(PathNodeMap),
		pathParamName: pathParamName,
		bindId: bindId,
	}
}

func NewPathNode() *PathNode {
	return &PathNode{
		pathNodeMap: make(PathNodeMap),
		pathParamName: "",
		bindId: 0,
	}
}

func (this *PathNode) Update(url string, pathParamName string, bindId int32) {
	// this.abi = NewApiBindingInfo(url)
	this.bindId = bindId
	this.pathParamName = pathParamName
}

func GetPathNode(method, path string) (*PathNode, map[string]string) {
	parts := strings.Split(strings.Trim(path, "/ "), "/")
	currentNode := GetPathMap(method)
	
	pathParamMap := make(map[string]string)
	for _, part := range parts {
		if node, ok := currentNode.pathNodeMap[part]; ok {
			currentNode = node
		} else {
			if node, ok := currentNode.pathNodeMap["$"]; ok {
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

func GetUrl(abi *ApiBindingInfo, pathParamMap map[string]string, query string) string {
	url := abi.Url
	for key, value := range pathParamMap {
		url = strings.Replace(url, fmt.Sprintf("{{%s}}", key), value, -1)
	}
	return strings.TrimRight(url, "? ") + "?" + query
}

func GetApiBindingInfo(method, path, query string) (*ApiBindingInfo, string) {
	pathNode, pathParamMap := GetPathNode(method, path)
	if pathNode != nil {
		abi := pathNode.FindApiBindingInfo()
		url := GetUrl(abi, pathParamMap, query)
		return abi, url
	} else {
		return nil, ""
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
		if node, ok := currentNode.pathNodeMap[part]; ok {
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
			currentNode.pathNodeMap[part] = newNode
			currentNode = newNode
		}
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

// TODO: Cache?
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

/**
 * Make a route forbidden or not
 */
func ChangeRouteStatus(method, path string, status int32) (int32, error) {
	pathNode, _ := GetPathNode(method, path)
	if pathNode != nil {
		abi := pathNode.FindApiBindingInfo()
		lastStatus := abi.Status
		abi.Status = status
		return lastStatus, nil
	}
	return AbiStatusUnknown, errors.New("Route not found")
}

/**
 * It's for Initialize the routes table and update the routes the table
 */
func updateRoutes(apiBindings []map[string]string) {
	for _, apiBinding := range apiBindings {
		fmt.Printf("%+v\n", apiBinding)
		methods := getHttpMethod(apiBinding["http_method"])
		for _, method := range methods {
			path := apiBinding["gateway_api"]
			ChangeRouteStatus(method, path, AbiStatusDying)
		}
	}
	// PrintRoutes()
	addRoutes(apiBindings)
}

/**
 * It's for Initialize the routes table and update the routes the table
 */
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
	PrintRoutes()
	// For test	
	// addRoute(http.MethodGet, "/api/thsamples", "http://127.0.0.1:9090/api/thsamples")
	// addRoute(http.MethodGet, "/api/thsample/{{id}}", "http://127.0.0.1:9090/api/thsample/{{id}}?a={{id}}")

	return true
}