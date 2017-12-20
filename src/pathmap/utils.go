package pathmap

import (
	"fmt"
	"strings"
)

func PrintRoutesNode(pathNode *PathNode, level int) {
	currentNode := pathNode
	for key, node := range currentNode.pathNodeMap {
		abi := node.FindApiBindingInfo()
		if abi != nil {
			fmt.Printf("|_%s[%s(%d)] %q %d\n", strings.Repeat("_", level * 4), key, node.bindId, abi.Url, abi.Status)
		} else {
			fmt.Printf("|_%s[%s]\n", strings.Repeat("_", level * 4), key)
		}
		
		PrintRoutesNode(node, level + 1)
	}
}

func PrintRoutes() {
	fmt.Printf("Path nodes for GET,HEAD\n")
	PrintRoutesNode(gGetHeadPathMap, 0)
	fmt.Printf("Path nodes for POST,PUT\n")
	PrintRoutesNode(gPostPutPathMap, 0)
}