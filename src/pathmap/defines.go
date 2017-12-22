package pathmap

import (
	"middleware"
	"net/http"
)

const (
	AbiStatusUnknown	= 0
	AbiStatusEnabled	= 1
	AbiStatusDisabled	= 2
	AbiStatusDying		= 3
)

type ApiBindingInfo struct {
	BindId int32;
	GroupId int32;
	Url string;
	WarningLevel int32;
	LogLevel int32;
	CheckConfig int32;
	Status int32;
	HeadMiddleware *middleware.HeadMiddleware;
}

type PathNode struct {
	bindId 			int32
	pathNodeMap 	PathNodeMap
	pathParamName 	string
}

type PathNodeMap map[string]*PathNode


// Global
var gGetHeadPathMap = NewPathNode()

var gPostPutPathMap = NewPathNode()

var gApiBindingInfoMap = map[int32]*ApiBindingInfo {}

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