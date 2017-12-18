package pathmap

import (
	"fmt"
	"strings"
)

type ApiBindingInfo struct {
	Url string	
}

func NewApiBindingInfo(url string) *ApiBindingInfo {
	return &ApiBindingInfo{
		url,
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

// Global
var pathMap = NewPathNode("", "")


	
func GetPathNode(path string) (*PathNode, map[string]string) {
	parts := strings.Split(strings.Trim(path, "/ "), "/")
	currentNode := pathMap
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

func GetApiBindingInfo(path string) *ApiBindingInfo {
	pathNode, pathParamMap := GetPathNode(path)
	if pathNode != nil {
		url := GetUrl(pathNode, pathParamMap)
		
		return &ApiBindingInfo{
			Url: url,
		}
	} else {
		return nil
	}
}

func addRoute(path string, url string) {
	parts := strings.Split(strings.Trim(path, "/ "), "/")
	
	count := len(parts)
	currentNode := pathMap
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

func printRoutes(pathMap *PathNode, level int) {
	currentNode := pathMap
	for key, node := range currentNode.subNode {
		fmt.Printf("%s", strings.Repeat("    ", level))
		fmt.Printf("%s %q\n", key, node.abi.Url)
		printRoutes(node, level + 1)
	}
}

func Initialize() {
	pathMap.abi = nil
	pathMap.subNode = make(map[string]*PathNode)

	addRoute("/api/thsamples", "http://127.0.0.1:9090/api/thsamples")

	addRoute("/api/thsamples/debug", "http://127.0.0.1:9090/api/thsamples/debug")
	
	addRoute("/api/thsample/{{id}}", "http://127.0.0.1:9090/api/thsample/{{id}}?a={{id}}")

	addRoute("/api/thsample/{{id}}/info", "http://127.0.0.1:9090/api/thsample/{{id}}/info")

	printRoutes(pathMap, 0)
}