package pathmap

import (
	"fmt"
	"middleware"
	"net/http"
	"strings"
)

// Global
var gGetHeadPathMap = NewPathNode("", "")

var gPostPutPathMap = NewPathNode("", "")

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
	subNode map[string]*PathNode
	abi *ApiBindingInfo
	pathParamName string
}

func NewPathNode(url string, pathParamName string) *PathNode {
	return &PathNode{
		subNode: make(map[string]*PathNode),
		abi: NewApiBindingInfo(url),
		pathParamName: pathParamName,
	}
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

func GetUrl(pathNode *PathNode, pathParamMap map[string]string) string {
	url := pathNode.abi.Url
	for key, value := range pathParamMap {
		url = strings.Replace(url, fmt.Sprintf("{{%s}}", key), value, -1)
	}
	return url
}

func GetApiBindingInfo(method, path, query string) *ApiBindingInfo {
	pathNode, pathParamMap := GetPathNode(method, path)
	if pathNode != nil {
		url := GetUrl(pathNode, pathParamMap)
		
		// TODO: ? how to add ?
		return NewApiBindingInfo(url + "?" + query)
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

func addRoute(method, path, url string) {
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
			currentNode = node
		} else {
			param := ""
			if index + 1 == count {
				param = url
			}
			newNode := NewPathNode(param, pathParam)
			currentNode.subNode[part] = newNode
			currentNode = newNode
		}
	}
}

func printRoutes(pathNode *PathNode, level int) {
	currentNode := pathNode
	for key, node := range currentNode.subNode {
		fmt.Printf("|_%s[%s] %s\n", strings.Repeat("____", level), key, node.abi.Url)
		printRoutes(node, level + 1)
	}
}

// TODO:
func getHttpMethod(httpMethod string) string {
	return "GET"
}

func addRoutes(apiBindings []map[string]string) {
	for _, apiBinding := range apiBindings {
		fmt.Printf("%+v\n", apiBinding)
		httpMethod := getHttpMethod(apiBinding["http_method"])
		addRoute(httpMethod, apiBinding["gateway_api"], apiBinding["service_api"])
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